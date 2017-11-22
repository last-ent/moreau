package main

import (
	"bytes"
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
var gcpAPI *gcp.GcpAPI
var dbi *data.DBInstance

func init() {
	f, _ := os.Open("./editor.html")
	ed, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Unable to read html file:", err.Error())
	}

	editorHtml = []byte(ed)

	toDB := make(chan data.Essay)
	callback := make(chan data.PubRow)

	dbi = data.NewDBInstance(toDB, callback)
	gcpAPI = gcp.NewGcpAPI(bucketName, toDB, callback)
}

func editorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		go gcpAPI.UploadPost(r.Form)
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

const trow = `
<li>
	<a href="/%d"><b>%s</b></a>
</li>`

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer

	for _, essay := range dbi.GetAllPublished() {
		fmt.Println(".")
		buffer.WriteString(fmt.Sprintf(trow, essay.PubID, essay.Title))
	}

	var body bytes.Buffer
	body.WriteString(`<html><body>`)
	body.WriteString(`<table>`)
	body.WriteString(buffer.String())
	body.WriteString(`</table>`)
	body.WriteString(`</body></html>`)

	gcp.UploadToGCP("index.html", body.String())
	// http.Redirect(w, r, "/", http.StatusMovedPermanently)
	w.Write([]byte("asdfasdfadfdas"))
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/editor", editorHandler)

	// http.HandleFunc("/", editorHandler)
	fmt.Println("Starting server at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
