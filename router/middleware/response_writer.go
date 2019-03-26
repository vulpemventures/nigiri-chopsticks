package middleware

import (
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about response.
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher

	Status() int
	Body() []byte
	Written() bool
	Before(func(ResponseWriter))
}

type beforeFunc func(ResponseWriter)

func newResponseWriter(rw http.ResponseWriter) ResponseWriter {
	nrw := &responseWriter{
		ResponseWriter: rw,
	}

	if _, ok := rw.(http.CloseNotifier); ok {
		return &responseWriterCloseNotifer{nrw}
	}

	return nrw
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	body        []byte
	size        int
	beforeFuncs []beforeFunc
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.status = s
	rw.callBefore()
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	rw.body = b
	return size, err
}

// Status returns the status code of the response or 0 if the response has not been written
func (rw *responseWriter) Status() int {
	return rw.status
}

// Body returns the response body
func (rw *responseWriter) Body() []byte {
	return rw.body
}

// Written returns whether or not the ResponseWriter has been written.
func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

// Before allows for a function to be called before the ResponseWriter has been written to. This is
// useful for setting headers or any other operations that must happen before a response has been written.
func (rw *responseWriter) Before(before func(ResponseWriter)) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

func (rw *responseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i](rw)
	}
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		if !rw.Written() {
			// The status will be StatusOK if WriteHeader has not been called yet
			rw.WriteHeader(http.StatusOK)
		}
		flusher.Flush()
	}
}

type responseWriterCloseNotifer struct {
	*responseWriter
}

func (rw *responseWriterCloseNotifer) CloseNotify() <-chan bool {
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
