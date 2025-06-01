package main

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/TerraQuest-Studios/skinsdb-go/help"
	"github.com/TerraQuest-Studios/skinsdb-go/pages"
	"github.com/TerraQuest-Studios/skinsdb-go/skins"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.NotFound(templ.Handler(pages.Error(404), templ.WithStatus(http.StatusNotFound)).ServeHTTP)
	r.MethodNotAllowed(templ.Handler(pages.Error(405), templ.WithStatus(http.StatusNotFound)).ServeHTTP)
	r.Get("/", templ.Handler(pages.Index()).ServeHTTP)
	r.Get("/help", templ.Handler(pages.Help()).ServeHTTP)
	r.Get("/help/{page}", func(w http.ResponseWriter, req *http.Request) {
		page := chi.URLParam(req, "page")

		if !help.IsValidPage(page) {
			templ.Handler(pages.Error(404), templ.WithStatus(http.StatusNotFound)).ServeHTTP(w, req)
			return
		}

		templ.Handler(pages.HelpItem(page)).ServeHTTP(w, req)
	})
	r.Get("/skinsdata/skins/{skinimage}", func(w http.ResponseWriter, req *http.Request) {
		skinImage := chi.URLParam(req, "skinimage")

		if !skins.IsValidSkinImage(skinImage) {
			templ.Handler(pages.Error(404), templ.WithStatus(http.StatusNotFound)).ServeHTTP(w, req)
			return
		}

		skinData := skins.GetSkinImage(skinImage)
		if skinData == nil {
			templ.Handler(pages.Error(403), templ.WithStatus(http.StatusNotFound)).ServeHTTP(w, req)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write(skinData)
	})
	//"/api/v1/content?per_page={per_page}&page={page}"
	r.Get("/api/v1/content", func(w http.ResponseWriter, req *http.Request) {
		perPage := req.URL.Query().Get("per_page")
		page := req.URL.Query().Get("page")

		//convert to ints
		perPageInt, err := strconv.Atoi(perPage)
		if err != nil || perPageInt <= 0 {
			perPageInt = 10
		}
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt <= 0 {
			pageInt = 1
		}

		//limt to 1000
		if perPageInt > 1000 {
			perPageInt = 1000
		}
		if pageInt > 1000 {
			pageInt = 1000
		}

		response := skins.MakeSkinsApiResponse(perPageInt, pageInt)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	})
	http.ListenAndServe(":3420", r)
}
