package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"path"
	"strings"
)

//go:embed dashboard/dist/*
var dashboardFiles embed.FS

func DashboardServer() error {
	staticFiles, assetCount, err := dashboardAssets()
	if err != nil {
		return fmt.Errorf("load embedded dashboard assets: %w", err)
	}

	listener, err := net.Listen("tcp", DashboardHTTPAddr)
	if err != nil {
		return fmt.Errorf("start dashboard server on %s: %w", DashboardHTTPAddr, err)
	}

	startupLog("frontend status", fmt.Sprintf("embedded dashboard assets ready (%d files)", assetCount))
	startupLog("dashboard status", "serving dashboard on "+DashboardHTTPAddr)
	fmt.Println("Listening to Web Dashboard at " + DashboardHTTPAddr)

	httpServer := &http.Server{Handler: dashboardHandler(staticFiles)}
	go func() {
		if serveErr := httpServer.Serve(listener); serveErr != nil && serveErr != http.ErrServerClosed {
			log.Printf("dashboard server stopped: %v", serveErr)
		}
	}()

	return nil
}

func dashboardAssets() (fs.FS, int, error) {
	staticFiles, err := fs.Sub(dashboardFiles, "dashboard/dist")
	if err != nil {
		return nil, 0, err
	}

	count := 0
	if err := fs.WalkDir(staticFiles, ".", func(current string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() {
			count++
		}
		return nil
	}); err != nil {
		return nil, 0, err
	}

	if _, err := fs.Stat(staticFiles, "index.html"); err != nil {
		return nil, 0, err
	}

	return staticFiles, count, nil
}

func dashboardHandler(staticFiles fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(staticFiles))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		if cleanPath == "." || cleanPath == "/" {
			http.ServeFileFS(w, r, staticFiles, "index.html")
			return
		}

		if _, err := fs.Stat(staticFiles, cleanPath); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		if path.Ext(cleanPath) != "" {
			http.NotFound(w, r)
			return
		}

		http.ServeFileFS(w, r, staticFiles, "index.html")
	})
}
