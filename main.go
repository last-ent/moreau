package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/last-ent/globr/data"
	"github.com/last-ent/globr/gcp"
	_ "github.com/mattn/go-sqlite3"
)

const bucketName = "www.last-ent.com"

var editorHtml []byte
var gpcAPI *gcp.GcpAPI

func init() {
	f, _ := os.Open("./editor.html")
	ed, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Unable to read html file:", err.Error())
	}

	editorHtml = []byte(ed)

	toDB := make(chan data.Essay)

	data.NewDBInstance(toDB)
	gpcAPI = gcp.NewGcpAPI(bucketName, toDB)
}

func editorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		go gpcAPI.UploadPost(r.Form)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	case http.MethodGet:
		w.Write(editorHtml)

	default:
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}

func main() {
	http.HandleFunc("/", editorHandler)
	http.HandleFunc("/editor", editorHandler)

	fmt.Println("Starting server at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
