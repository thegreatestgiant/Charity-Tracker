package middleware

import "net/http"

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	// Notice: we do NOT call r.ResponseWriter.WriteHeader here yet
}

// This wraps any http.Handler and catches 404s
func SpaFallback(fs http.Handler, fallback string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		fs.ServeHTTP(rec, r)

		if rec.status == http.StatusNotFound {
			// The file didn't exist — serve index.html instead
			r.URL.Path = "/"
			fs.ServeHTTP(w, r)
			return
		}

		// File existed — write the captured status and we're done
		w.WriteHeader(rec.status)
	})
}
