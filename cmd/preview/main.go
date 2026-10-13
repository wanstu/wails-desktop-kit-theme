package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"

	theme "github.com/wanstu/wails-desktop-kit-theme"
)

//go:embed all:web
var previewFS embed.FS

type packView struct {
	Name           string `json:"name"`
	DisplayName    string `json:"displayName"`
	Description    string `json:"description"`
	StylesheetPath string `json:"stylesheetPath"`
	GradientLight  string `json:"gradientLight"`
	GradientDark   string `json:"gradientDark"`
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "preview listen address")
	flag.Parse()

	app, err := fs.Sub(previewFS, "web")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/packs", func(w http.ResponseWriter, _ *http.Request) {
		packs := theme.Packs()
		manifest, err := theme.BuildRemoteManifest("")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		colors := make(map[string]theme.RemoteManifestPack)
		for _, item := range manifest.Packs {
			colors[item.Name] = item
		}
		result := make([]packView, 0, len(packs))
		for _, pack := range packs {
			path, _ := theme.StylesheetPath(pack.Name)
			result = append(result, packView{
				Name:           pack.Name,
				DisplayName:    pack.DisplayName,
				Description:    pack.Description,
				StylesheetPath: path,
				GradientLight:  colors[pack.Name].GradientLight,
				GradientDark:   colors[pack.Name].GradientDark,
			})
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("GET /desktopkit-theme/manifest.json", func(w http.ResponseWriter, _ *http.Request) {
		manifest, err := theme.BuildRemoteManifest("")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(manifest)
	})
	mux.Handle("/", http.FileServer(http.FS(theme.MountWithKit(app))))

	log.Printf("Theme Gallery: http://%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
