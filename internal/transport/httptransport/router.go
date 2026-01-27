// // package httptransport

// // // router placeholder

// // import (
// // 	"net/http"
// // )

// // func NewRouter(h *UserHandler) http.Handler {

// //		mux := http.NewServeMux()
// //		mux.HandleFunc("/users/register/", h.Create)
// //		return mux
// //	}
// package httptransport

// import (
// 	"net/http"

// 	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/observability/logger"
// 	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/transport/middleware"
// )

// // NewRouter wires all HTTP routes and middleware.
// // This is the SINGLE place where middleware ordering is defined.
// func NewRouter(
// 	h *UserHandler,
// 	log logger.Logger,
// 	env string,
// 	serviceName string,
// ) http.Handler {

// 	mux := http.NewServeMux()

// 	// ------------------------------------------------
// 	// Centralized error handler (wraps handlers)
// 	// ------------------------------------------------
// 	errHandler := middleware.NewErrorHandler(log)

// 	// ------------------------------------------------
// 	// Routes
// 	// ------------------------------------------------
// 	// Wrap regular handler to work with error handler middleware
// 	createHandler := func(w http.ResponseWriter, r *http.Request) error {
// 		h.Create(w, r)
// 		return nil // Regular handlers don't return errors
// 	}
// 	mux.Handle(
// 		"/users/register/",
// 		errHandler.Handle(createHandler),
// 	)

// 	// ------------------------------------------------
// 	// Global middleware order (IMPORTANT)
// 	// ------------------------------------------------
// 	var handler http.Handler = mux

// 	// 1️⃣ OpenTelemetry tracing (FIRST)
// 	handler = middleware.Tracing(serviceName)(handler)

// 	// 2️⃣ Panic recovery
// 	handler = middleware.PanicRecover(handler)

// 	// 3️⃣ Request logging (INFO / DEBUG)
// 	handler = middleware.RequestLogger(log, env)(handler)

// 	return handler
// }

package httptransport

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/observability/logger"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/transport/middleware"
)

func NewRouter(
	h *UserHandler,
	log logger.Logger,
	env string,
	serviceName string,
) http.Handler {

	r := chi.NewRouter()

	// ------------------------------------------------
	// Centralized error handler
	// ------------------------------------------------
	errHandler := middleware.NewErrorHandler(log)

	// ------------------------------------------------
	// Global middleware (order matters)
	// ------------------------------------------------
	r.Use(middleware.Tracing(serviceName))     // 1️⃣ Start trace
	r.Use(middleware.PanicRecover(errHandler)) // 2️⃣ Recover panics
	r.Use(middleware.RequestLogger(log, env))  // 3️⃣ Log lifecycle

	// ------------------------------------------------
	// Routes
	// ------------------------------------------------
	r.Route("/users", func(r chi.Router) {
		r.Post("/register", errHandler.Handle(h.Create))
	})

	return r
}
