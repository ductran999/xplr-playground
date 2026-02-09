package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	targetURL, _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
		req.Header.Set("X-Custom-Header", "Golang-PoC-Power")
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()

		newBody := strings.ReplaceAll(string(body), "Backend", "BACKEND (PROXIED)")
		resp.Body = io.NopCloser(bytes.NewBufferString(newBody))
		resp.ContentLength = int64(len(newBody))
		resp.Header.Set("Content-Length", fmt.Sprint(len(newBody)))
		resp.Header.Set("X-Proxy-Handled", "true")

		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("[Proxy] Error: %v\n", err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Proxy Error: Backend is unreachable."))
	}

	fmt.Println("Proxy Server is listening on :8081 -> Forward to :8080")
	log.Fatal(http.ListenAndServe(":8081", proxy))
}
