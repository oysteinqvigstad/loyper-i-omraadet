package web

import (
	"api/internal/dto"
	"log"
	"net/http"
)

func SetupRoutes(port string, db *dto.TrailsDB) *http.ServeMux {
	mux := http.ServeMux{}
	mux.Handle("/near/", getNearByTrails(db))

	domainNamePort := "https://localhost:" + port
	log.Println("Starting webservice")
	log.Println(domainNamePort + "/near/")

	return &mux
}
