package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func compressionHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := GzipResponseWriter{Writer: gz, ResponseWriter: w}
		handler(gzr, r)
	}
}
