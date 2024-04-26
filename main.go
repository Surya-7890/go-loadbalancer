package loadbalancer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var (
	ProxyList   = make(map[string][]*httputil.ReverseProxy)
	recentProxy = make(map[string]uint8)
)

/* @param hosts example "http://localhost:7000"  */
/* @param path example "/api"  */
func CreateNewProxy(path string, hosts ...string) {
	for _, host := range hosts {
		newProxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: strings.Split(host, ":")[0],
			Host:   strings.Split(host, "://")[1],
		})
		ProxyList[path] = append(ProxyList[path], newProxy)
	}
}

func handleProxy(w http.ResponseWriter, r *http.Request, path string) {
	remainder := uint8(len(ProxyList[path]))
	ProxyList[path][(recentProxy[path]+1)%remainder].ServeHTTP(w, r)
	defer func() {
		recentProxy[path] = (recentProxy[path] + 1) % remainder
	}()
}

func StartServer() {
	for k := range ProxyList {
		path := k
		http.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
			handleProxy(w, r, path)
		})
	}

	http.ListenAndServe(":8000", nil)
}
