package utils

import "net/http"

type HookedResponseWriter struct {
	http.ResponseWriter
	Got404 bool
}

func (w *HookedResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *HookedResponseWriter) WriteHeader(status int) {
	if status == http.StatusNotFound {
		w.Got404 = true
	} else {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *HookedResponseWriter) Write(p []byte) (int, error) {
	if w.Got404 {
		return len(p), nil
	}

	return w.ResponseWriter.Write(p)
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (w *StatusRecorder) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *StatusRecorder) WriteHeader(code int) {
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}
