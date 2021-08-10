package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	remote, err := url.Parse("https://210.86.179.233:6443")
	if err != nil {
		panic(err)
	}
	// TLSNextProto = map[string]func(authority string, c *tls.Conn) http.RoundTripper{}
	// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// http.DefaultTransport.(*http.Transport).TLSNextProto = map[string]func(authority string, c *tls.Conn) http.RoundTripper{}
	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {

			log.Println(r.URL)
			// log.Println(r.Header)
			r.Host = remote.Host

			transport := http.DefaultTransport.(interface {
				Clone() *http.Transport
			}).Clone()
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			disableHTTP2 := isSPDY(r)
			if disableHTTP2 {
				log.Println(transport.TLSNextProto)
				transport.TLSNextProto = map[string]func(authority string, c *tls.Conn) http.RoundTripper{}
				log.Println(transport.TLSNextProto)
			}
			p.Transport = transport
			p.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServeTLS(":8080", "fullchain.pem", "privkey.pem", nil)
	if err != nil {
		panic(err)
	}

}

func isSPDY(r *http.Request) bool {
	isSPDY := strings.HasPrefix(strings.ToLower(r.Header.Get("Upgrade")), "spdy/")
	if isSPDY {
		log.Println("it a SPDY!")
	}

	return isSPDY
}
