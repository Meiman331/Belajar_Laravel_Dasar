package main

import (
	"fmt"
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	statusCode int
	writer     http.ResponseWriter
}

func (lrw *loggingResponseWriter) Header() http.Header {
	return lrw.writer.Header()
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.writer.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}
	return lrw.writer.Write(b)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{writer: w}
		log.Printf("[LOG] %s %s", r.Method, r.URL.String())
		next.ServeHTTP(lrw, r)
		if lrw.statusCode == 0 {
			lrw.statusCode = http.StatusOK
		}

	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Halaman Utama")
	})

	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Tentang Kami")
	})

	wrapped := loggingMiddleware(mux)

	addr := ":8080"
	log.Printf("Server berjalan di %s", addr)
	if err := http.ListenAndServe(addr, wrapped); err != nil {
		log.Fatal(err)
	}
}
