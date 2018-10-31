package logger

import (
	"log"
	"net"
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip, port, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("userip: %q is not IP:port", r.RemoteAddr)
		}

		inner.ServeHTTP(w, r)
		log.Printf(
			"Request:\tIP:%s\tPORT:%s\t%s\t%s\t%s\t%s",
			ip,
			port,
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
