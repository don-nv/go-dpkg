package dprom

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var serveHTTP = promhttp.Handler().ServeHTTP

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveHTTP(w, r)
}
