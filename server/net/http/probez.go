package http

import (
	"fmt"
	"net/http"
)

type ProbezHandler struct {
	Check func() bool
	Name  string
}

func (h *ProbezHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !h.Check() {
		http.Error(w, h.Name, http.StatusTeapot)
	}
	fmt.Fprint(w, "OK-"+h.Name)
}
