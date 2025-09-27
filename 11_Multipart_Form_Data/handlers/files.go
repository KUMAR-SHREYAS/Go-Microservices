package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"product-images/files"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type Files struct {
	log   hclog.Logger
	store files.Storage
}

func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{log: l, store: s}
}

func (f *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024)
	if err != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected multipart form data", http.StatusBadRequest)
		return
	}
	id, idErr := strconv.Atoi(r.FormValue("id"))
	f.log.Info("Process form for id", "id", id)

	if idErr != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected expected integer id", http.StatusBadRequest)
		return
	}

	ff, mh, err := r.FormFile("file")
	if err != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected file", http.StatusBadRequest)
		return
	}
	f.saveFile(r.FormValue("id"), mh.Filename, rw, ff)
}

func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", fn)
	f.saveFile(id, fn, rw, r.Body)
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r io.ReadCloser) { // use io.ReadCloser that implementts read() and close() thats needed for multipart.File data type
	f.log.Info("Save file for product", "id", id, "path", path)
	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r) // src - r.Body , destination - fp
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
