package printers

import (
	"embed"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/everettraven/synkr/pkg/plugins"
)

//go:embed web.html
var webpage embed.FS

type Web struct{}

func (w *Web) Print(results ...plugins.SourceResult) error {
	return serve(results...)
}

func serve(results ...plugins.SourceResult) error {
	data, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("marshalling results to json: %w", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		view, err := webpage.ReadFile("web.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("could not render web view: %s", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(view)
	})

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	})

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return fmt.Errorf("creating listener: %w", err)
	}
	fmt.Println("Web UI is available at", listener.Addr().String())

	return http.Serve(listener, mux)
}
