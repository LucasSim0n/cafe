package cafe

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	app := NewServer()
	if len(app.children) != 0 {
		t.Errorf("Expected no children, got %d", len(app.children))
	}
	if len(app.routes) != 0 {
		t.Errorf("Expected no routes, got %d", len(app.routes))
	}
	if len(app.middlewares) != 0 {
		t.Errorf("Expected no global middlewares, got %d", len(app.middlewares))
	}
	if app.parent != nil {
		t.Error("Expected no parent")
	}
	if app.prefix != "" {
		t.Errorf("Expected no prefix, got %s", app.prefix)
	}
}

func TestApp_Get(t *testing.T) {
	app := NewServer()
	var called bool
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}
	app.Get("/test", handler)

	if len(app.routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(app.routes))
	}
	if app.routes[0].path != "/test" {
		t.Errorf("Expected path /test, got %s", app.routes[0].path)
	}
	if app.routes[0].method != "GET" {
		t.Errorf("Expected method GET, got %s", app.routes[0].method)
	}

	req := httptest.NewRequest("GET", "/test/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	if !called {
		t.Error("Handler was not called")
	}
}

func TestApp_Post(t *testing.T) {
	app := NewServer()
	var called bool
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	}
	app.Post("/test", handler)
	req := httptest.NewRequest("POST", "/test/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status Created, got %d", rr.Code)
	}
	if !called {
		t.Error("Handler was not called")
	}
}

func TestApp_Put(t *testing.T) {
	app := NewServer()
	var called bool
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusAccepted)
	}
	app.Put("/test", handler)
	req := httptest.NewRequest("PUT", "/test/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("Expected status Accepted, got %d", rr.Code)
	}
	if !called {
		t.Error("Handler was not called")
	}
}

func TestApp_Delete(t *testing.T) {
	app := NewServer()
	var called bool
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}
	app.Delete("/test", handler)
	req := httptest.NewRequest("DELETE", "/test/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %d", rr.Code)
	}
	if !called {
		t.Error("Handler was not called")
	}
}

func TestApp_UseRouter(t *testing.T) {
	app := NewServer()
	var called bool
	rtr := app.Group("/api")
	rtr.Get("/sub", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Write([]byte("sub-route"))
	})

	if len(app.children) != 1 {
		t.Errorf("Expected 1 mounted router, got %d", len(app.children))
	}
	if app.children[0].prefix != "/api" {
		t.Errorf("Expected mounted router path /api, got %s", app.children[0].prefix)
	}

	req := httptest.NewRequest("GET", "/api/sub/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	if !called {
		t.Error("Sub-router handler was not called")
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "sub-route" {
		t.Errorf("Expected body 'sub-route', got '%s'", string(body))
	}
}

func TestApp_UseRouter_DuplicatePath(t *testing.T) {
	app := NewServer()
	app.Group("/api")
	app.Group("/api")

	if len(app.children) != 1 {
		t.Errorf("Expected 1 mounted router due to duplicate path, got %d", len(app.children))
	}
}

func TestNewRouter(t *testing.T) {
	rtr := NewRouter("")
	if len(rtr.routes) != 0 {
		t.Errorf("Expected no routes, got %d", len(rtr.routes))
	}
	if len(rtr.children) != 0 {
		t.Errorf("Expected no routers, got %d", len(rtr.children))
	}
	if len(rtr.middlewares) != 0 {
		t.Errorf("Expected no router middlewares, got %d", len(rtr.middlewares))
	}
}

func TestRouter_Get(t *testing.T) {
	rtr := NewRouter("")
	handler := func(w http.ResponseWriter, r *http.Request) {}
	rtr.Get("/item", handler)

	if len(rtr.routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(rtr.routes))
	}
	if rtr.routes[0].path != "/item" {
		t.Errorf("Expected path /item, got %s", rtr.routes[0].path)
	}
	if rtr.routes[0].method != "GET" {
		t.Errorf("Expected method GET, got %s", rtr.routes[0].method)
	}
}

func TestRouter_UseRouter(t *testing.T) {
	parentRouter := NewRouter("/main")
	childRouter := parentRouter.Group("/sub")
	var called bool
	childRouter.Get("/nested", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Write([]byte("nested-route"))
	})

	if len(parentRouter.children) != 1 {
		t.Errorf("Expected 1 mounted router, got %d", len(parentRouter.children))
	}
	if parentRouter.children[0].prefix != "/sub" {
		t.Errorf("Expected mounted router path /sub, got %s", parentRouter.children[0].prefix)
	}

	allRoutes := parentRouter.compileRoutes()
	found := false
	for _, r := range allRoutes {
		if r.path == "/sub/nested" && r.method == "GET" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected '/sub/nested' route not found in flattened routes")
	}

	app := NewServer()
	app.UseRouter(parentRouter)
	req := httptest.NewRequest("GET", "/main/sub/nested/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	if !called {
		t.Error("Nested router handler was not called")
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "nested-route" {
		t.Errorf("Expected body 'nested-route', got '%s'", string(body))
	}
}

func TestAddRoute_Duplicate(t *testing.T) {
	routes := []*Route{}
	handler1 := func(w http.ResponseWriter, r *http.Request) {}
	handler2 := func(w http.ResponseWriter, r *http.Request) {}

	routes = addRoute(routes, "/test", "GET", handler1)
	routes = addRoute(routes, "/test", "GET", handler2) // Should not add a duplicate

	if len(routes) != 1 {
		t.Errorf("Expected 1 route after adding a duplicate, got %d", len(routes))
	}
	if routes[0].handler == nil {
		t.Errorf("Expected handler to be non-nil")
	}
}

func TestRouters_Order(t *testing.T) {
	app := NewServer()
	var order []string

	app.Get("/root", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "root")
	})

	router1 := app.Group("/api1")
	router1.Get("/sub1", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "sub1")
	})

	router2 := app.Group("/api2")
	router2.Get("/sub2", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "sub2")
	})

	// Test root route
	reqRoot := httptest.NewRequest("GET", "/root/", nil)
	rrRoot := httptest.NewRecorder()
	app.ServeHTTP(rrRoot, reqRoot)
	if rrRoot.Code != http.StatusOK {
		t.Errorf("Expected status OK for root, got %d", rrRoot.Code)
	}

	// Test api1/sub1 route
	reqSub1 := httptest.NewRequest("GET", "/api1/sub1/", nil)
	rrSub1 := httptest.NewRecorder()
	app.ServeHTTP(rrSub1, reqSub1)
	if rrSub1.Code != http.StatusOK {
		t.Errorf("Expected status OK for /api1/sub1, got %d", rrSub1.Code)
	}

	// Test api2/sub2 route
	reqSub2 := httptest.NewRequest("GET", "/api2/sub2/", nil)
	rrSub2 := httptest.NewRecorder()
	app.ServeHTTP(rrSub2, reqSub2)
	if rrSub2.Code != http.StatusOK {
		t.Errorf("Expected status OK for /api2/sub2, got %d", rrSub2.Code)
	}
}

func TestRouter_GetRoutes_NestedRouters(t *testing.T) {
	r1 := NewRouter("")
	r2 := r1.Group("/admin")
	r3 := r2.Group("/sub")
	r4 := r1.Group("/admin")

	r1.Get("/users", func(w http.ResponseWriter, r *http.Request) {})
	r2.Get("/profile", func(w http.ResponseWriter, r *http.Request) {})
	r3.Get("/details", func(w http.ResponseWriter, r *http.Request) {})
	r4.Get("/settings", func(w http.ResponseWriter, r *http.Request) {})
	r1.Get("/billing", func(w http.ResponseWriter, r *http.Request) {})

	allRoutes := r1.compileRoutes()

	expectedPaths := map[string]bool{
		"/users":             false,
		"/admin/profile":     false,
		"/admin/sub/details": false,
		"/admin/settings":    false,
		"/billing":           false,
	}

	for _, r := range allRoutes {
		if _, ok := expectedPaths[r.path]; ok {
			expectedPaths[r.path] = true
		}
	}

	for path, found := range expectedPaths {
		if !found {
			t.Errorf("Expected route '%s' not found in flattened routes", path)
		}
	}

	if len(allRoutes) != len(expectedPaths) {
		t.Errorf("Expected %d routes, got %d", len(expectedPaths), len(allRoutes))
		var rString string
		for _, r := range allRoutes {
			rString += r.path + "!\n"
		}
		t.Logf("Found routes: %+v", rString)
	}
}

func TestApp_Use(t *testing.T) {
	app := NewServer()
	mw := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)
		}
	}
	app.Use(mw)

	if len(app.middlewares) != 1 {
		t.Errorf("Expected 1 global middleware, got %d", len(app.middlewares))
	}
}

func TestRouter_Use(t *testing.T) {
	rtr := NewRouter("")
	mw := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)
		}
	}
	rtr.Use(mw)

	if len(rtr.middlewares) != 1 {
		t.Errorf("Expected 1 router middleware, got %d", len(rtr.middlewares))
	}
}

func TestSetUpMiddlewares(t *testing.T) {
	var callOrder []string
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "finalHandler")
		w.Write([]byte("ok"))
	})

	mw1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "mw1")
			next(w, r)
		}
	}
	mw2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "mw2")
			next(w, r)
		}
	}

	mws := []Middleware{mw1, mw2}
	chainedHandler := chain(finalHandler, mws)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	chainedHandler.ServeHTTP(rr, req)

	expectedOrder := []string{"mw1", "mw2", "finalHandler"}
	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Expected call order length %d, got %d. Order: %v", len(expectedOrder), len(callOrder), callOrder)
	}
	for i, expected := range expectedOrder {
		if callOrder[i] != expected {
			t.Errorf("At index %d, expected '%s', got '%s'. Full order: %v", i, expected, callOrder[i], callOrder)
		}
	}
	if rr.Body.String() != "ok" {
		t.Errorf("Expected body 'ok', got '%s'", rr.Body.String())
	}
}

func TestApp_MiddlewareExecution(t *testing.T) {
	app := NewServer()
	var callOrder []string

	app.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "globalMw1")
			next(w, r)
		}
	})
	app.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "globalMw2")
			next(w, r)
		}
	})

	app.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "finalHandler")
		w.Write([]byte("route handled"))
	})

	req := httptest.NewRequest("GET", "/test/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	expectedOrder := []string{"globalMw1", "globalMw2", "finalHandler"}
	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Expected call order length %d, got %d. Order: %v", len(expectedOrder), len(callOrder), callOrder)
	}
	for i, expected := range expectedOrder {
		if callOrder[i] != expected {
			t.Errorf("At index %d, expected '%s', got '%s'. Full order: %v", i, expected, callOrder[i], callOrder)
		}
	}
	if rr.Body.String() != "route handled" {
		t.Errorf("Expected body 'route handled', got '%s'", rr.Body.String())
	}
}

func TestRouter_MiddlewareExecution(t *testing.T) {
	app := NewServer()
	var callOrder []string

	rtr := app.Group("/api")
	rtr.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "routerMw1")
			next(w, r)
		}
	})
	rtr.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "routerMw2")
			next(w, r)
		}
	})

	rtr.Get("/subroute", func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "finalSubHandler")
		w.Write([]byte("subroute handled"))
	})

	req := httptest.NewRequest("GET", "/api/subroute/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	expectedOrder := []string{"routerMw1", "routerMw2", "finalSubHandler"}
	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Expected call order length %d, got %d. Order: %v", len(expectedOrder), len(callOrder), callOrder)
	}
	for i, expected := range expectedOrder {
		if callOrder[i] != expected {
			t.Errorf("At index %d, expected '%s', got '%s'. Full order: %v", i, expected, callOrder[i], callOrder)
		}
	}
	if rr.Body.String() != "subroute handled" {
		t.Errorf("Expected body 'subroute handled', got '%s'", rr.Body.String())
	}
}

func TestAppAndRouter_MiddlewareCombinedExecution(t *testing.T) {
	app := NewServer()
	var callOrder []string

	app.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "globalAppMw")
			next(w, r)
		}
	})

	rtr := app.Group("/data")
	rtr.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "routerMw")
			next(w, r)
		}
	})

	rtr.Get("/item", func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "finalItemHandler")
		w.Write([]byte("item handled"))
	})

	req := httptest.NewRequest("GET", "/data/item/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	expectedOrder := []string{"globalAppMw", "routerMw", "finalItemHandler"}
	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Expected call order length %d, got %d. Order: %v", len(expectedOrder), len(callOrder), callOrder)
	}
	for i, expected := range expectedOrder {
		if callOrder[i] != expected {
			t.Errorf("At index %d, expected '%s', got '%s'. Full order: %v", i, expected, callOrder[i], callOrder)
		}
	}
	if rr.Body.String() != "item handled" {
		t.Errorf("Expected body 'item handled', got '%s'", rr.Body.String())
	}
}

func TestDeepNestedMiddlewareCombinedExecution(t *testing.T) {
	app := NewServer()
	var callOrder []string

	app.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "globalAppMw")
			next(w, r)
		}
	})

	router1 := app.Group("/api")
	router1.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "router1Mw")
			next(w, r)
		}
	})

	router2 := router1.Group("/nested")
	router2.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "router2Mw")
			next(w, r)
		}
	})
	router2.Get("/resource", func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "finalResourceHandler")
		w.Write([]byte("resource handled"))
	})

	req := httptest.NewRequest("GET", "/api/nested/resource/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	expectedOrder := []string{"globalAppMw", "router1Mw", "router2Mw", "finalResourceHandler"}

	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Expected call order length %d, got %d. Order: %v", len(expectedOrder), len(callOrder), callOrder)
	}
	for i, expected := range expectedOrder {
		if callOrder[i] != expected {
			t.Errorf("At index %d, expected '%s', got '%s'. Full order: %v", i, expected, callOrder[i], callOrder)
		}
	}
	if rr.Body.String() != "resource handled" {
		t.Errorf("Expected body 'resource handled', got '%s'", rr.Body.String())
	}
}

func TestMiddlewareCanHaltExecution(t *testing.T) {
	app := NewServer()
	var callOrder []string

	// Middleware que detiene la ejecución
	app.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "haltingMw")
			w.WriteHeader(http.StatusUnauthorized) // Detiene la ejecución
			// No llama a next(w, r)
		}
	})

	app.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "finalHandler") // No debería ser llamada
		w.Write([]byte("protected content"))
	})

	req := httptest.NewRequest("GET", "/protected/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	expectedOrder := []string{"haltingMw"}
	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Expected call order length %d, got %d. Order: %v", len(expectedOrder), len(callOrder), callOrder)
	}
	for i, expected := range expectedOrder {
		if callOrder[i] != expected {
			t.Errorf("At index %d, expected '%s', got '%s'. Full order: %v", i, expected, callOrder[i], callOrder)
		}
	}
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized, got %d", rr.Code)
	}
	if rr.Body.String() != "" { // Body debería estar vacío si el handler final no se ejecuta
		t.Errorf("Expected empty body, got '%s'", rr.Body.String())
	}
}

func TestMiddlewareCanModifyRequest(t *testing.T) {
	app := NewServer()
	var modifiedValue string

	app.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Modificar el contexto de la solicitud o un encabezado, por ejemplo
			ctx := r.Context()
			newCtx := context.WithValue(ctx, "user_id", "123")
			r = r.WithContext(newCtx)
			next(w, r)
		}
	})

	app.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value("user_id")
		if val != nil {
			modifiedValue = val.(string)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(modifiedValue))
	})

	req := httptest.NewRequest("GET", "/info/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	if modifiedValue != "123" {
		t.Errorf("Expected modifiedValue '123', got '%s'", modifiedValue)
	}
	if rr.Body.String() != "123" {
		t.Errorf("Expected body '123', got '%s'", rr.Body.String())
	}
}

// --- NUEVOS TESTS PARA PARÁMETROS DE RUTA ---

func TestApp_PathParameter(t *testing.T) {
	app := NewServer()
	expectedID := "123"
	app.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id != expectedID {
			t.Errorf("Expected ID '%s', got '%s'", expectedID, id)
		}
		w.Write([]byte("User ID: " + id))
	})

	req := httptest.NewRequest("GET", "/users/123/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "User ID: 123" {
		t.Errorf("Expected body 'User ID: 123', got '%s'", string(body))
	}
}

func TestApp_MultiplePathParameters(t *testing.T) {
	app := NewServer()
	expectedUserID := "user456"
	expectedBookID := "book789"

	app.Get("/users/{userID}/books/{bookID}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userID")
		bookID := r.PathValue("bookID")

		if userID != expectedUserID {
			t.Errorf("Expected UserID '%s', got '%s'", expectedUserID, userID)
		}
		if bookID != expectedBookID {
			t.Errorf("Expected BookID '%s', got '%s'", expectedBookID, bookID)
		}
		w.Write(fmt.Appendf([]byte{}, "User: %s, Book: %s", userID, bookID))
	})

	req := httptest.NewRequest("GET", "/users/user456/books/book789/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "User: user456, Book: book789" {
		t.Errorf("Expected body 'User: user456, Book: book789', got '%s'", string(body))
	}
}

func TestRouter_PathParameter(t *testing.T) {
	app := NewServer()
	userRouter := app.Group("/api/users")
	expectedUserID := "987"

	userRouter.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id != expectedUserID {
			t.Errorf("Expected ID '%s', got '%s'", expectedUserID, id)
		}
		w.Write([]byte("Router User ID: " + id))
	})

	req := httptest.NewRequest("GET", "/api/users/987/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "Router User ID: 987" {
		t.Errorf("Expected body 'Router User ID: 987', got '%s'", string(body))
	}
}

func TestRouter_NestedPathParameters(t *testing.T) {
	app := NewServer()

	expectedProductID := "prodA"
	expectedReviewID := "revB"

	mainRouter := app.Group("/store")
	subRouter := mainRouter.Group("/products/{productID}")

	subRouter.Get("/reviews/{reviewID}", func(w http.ResponseWriter, r *http.Request) {
		productID := r.PathValue("productID") // Vendrá del router padre
		reviewID := r.PathValue("reviewID")

		if productID != expectedProductID {
			t.Errorf("Expected ProductID '%s', got '%s'", expectedProductID, productID)
		}
		if reviewID != expectedReviewID {
			t.Errorf("Expected ReviewID '%s', got '%s'", expectedReviewID, reviewID)
		}
		w.Write(fmt.Appendf([]byte{}, "Product: %s, Review: %s", productID, reviewID))
	})

	req := httptest.NewRequest("GET", "/store/products/prodA/reviews/revB/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "Product: prodA, Review: revB" {
		t.Errorf("Expected body 'Product: prodA, Review: revB', got '%s'", string(body))
	}
}

func TestPathParameter_WithMiddleware(t *testing.T) {
	app := NewServer()
	var middlewareCalled bool
	expectedID := "42"

	// Middleware que puede leer el contexto (aunque r.PathValue no está en el contexto en sí,
	// se accede directamente del request)
	app.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			if r.PathValue("id") != "" && r.PathValue("id") != expectedID {
				t.Errorf("Middleware: Expected ID '%s', got '%s' from PathValue", expectedID, r.PathValue("id"))
			}
			next(w, r)
		}
	})

	app.Get("/items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id != expectedID {
			t.Errorf("Handler: Expected ID '%s', got '%s'", expectedID, id)
		}
		w.Write([]byte("Item ID: " + id))
	})

	req := httptest.NewRequest("GET", "/items/42/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	if !middlewareCalled {
		t.Error("Middleware was not called")
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "Item ID: 42" {
		t.Errorf("Expected body 'Item ID: 42', got '%s'", string(body))
	}
}

// Test para asegurar que el comportamiento de trailing slash sigue funcionando con path parameters
func TestApp_PathParameter_WithTrailingSlash(t *testing.T) {
	app := NewServer()
	expectedID := "55"
	app.Get("/widgets/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id != expectedID {
			t.Errorf("Expected ID '%s', got '%s'", expectedID, id)
		}
		w.Write([]byte("Widget ID: " + id))
	})

	// Probar con trailing slash
	req := httptest.NewRequest("GET", "/widgets/55/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK for trailing slash, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "Widget ID: 55" {
		t.Errorf("Expected body 'Widget ID: 55', got '%s'", string(body))
	}

	// Probar sin trailing slash (para asegurar que ambas funcionan)
	req = httptest.NewRequest("GET", "/widgets/55", nil)
	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status OK for no trailing slash, got %d", rr.Code)
	}
}
