package main

import (
	httptransport "darkstalker/pkg/protocols/http"
	"darkstalker/pkg/services"
	"darkstalker/pkg/services/defaultService"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	_ = godotenv.Load()

	var svc services.Service
	{
		svc = defaultService.NewDefaultService(os.Getenv("DARKSTALKER_VK_TOKEN"))
	}

	var h http.Handler
	{
		h = httptransport.MakeHTTPHandler("v1", svc)
		h = corsMiddleware(h)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), h))
}

// Allows CORS for given handler.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		next.ServeHTTP(w, r)
	})
}
