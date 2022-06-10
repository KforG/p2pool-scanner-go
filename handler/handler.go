package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KforG/p2pool-scanner-go/config"
	"github.com/KforG/p2pool-scanner-go/logging"
	"github.com/KforG/p2pool-scanner-go/scanner"

	"github.com/go-chi/chi"
)

func Router(n *scanner.Nodes) {
	logging.Infof("Setting up router..\n")

	r := chi.NewRouter()
	//r.Use(middleware.Logger) // Useful for debugging purposes. Doesn't need to clog log in production.

	r.Get("/nodes", getNodes(n))

	logging.Infof("Listening on port %s", config.Active.WebPort)
	http.ListenAndServe(fmt.Sprintf(":"+config.Active.WebPort), r)
}

func getNodes(n *scanner.Nodes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(n)
	}
}
