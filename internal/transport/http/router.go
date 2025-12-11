package httptransport

// router placeholder

import (
	"net/http"
)

func NewRouter(h *UserHandler) http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/users/register/", h.Create)
	return mux
}
