package middleware

import (
	"fmt"
	"net/http"
	"runtime"
)

func WithRecover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				fmt.Printf("recovering from err %v\n %s", err, buf)

				w.Header().Set("Content-Type", "application/json")
				_, writeErr := w.Write([]byte(`{"error":"server panic!"}`))
				if writeErr != nil {
					fmt.Println(writeErr)
				}
			}
		}()

		h.ServeHTTP(w, r)
	})
}
