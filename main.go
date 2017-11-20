package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

const bucketName = "www.last-ent.com"

var editorHtml []byte

func init() {
	f, _ := os.Open("./editor.html")
	ed, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Unable to read html file:", err.Error())
	}

	editorHtml = []byte(ed)
}

func editorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		go uploadPost(r.Form)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	default:
		w.Write(editorHtml)
	}
}

func uploadPost(form url.Values) {
	title := form["title"][0]
	// form["content"] is a single string but better safe than sorry.
	content := strings.Join(form["content"], "\n")
	upload(title, strings.NewReader(content))
}

func main() {
	http.HandleFunc("/", editorHandler)
	http.HandleFunc("/editor", editorHandler)

	fmt.Println("Starting server at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func upload(object string, content *strings.Reader) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	wc := client.Bucket(bucketName).Object(object).NewWriter(ctx)
	var d int64
	if d, err = io.Copy(wc, content); err != nil {
		panic(err)
	}
	if err := wc.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done uploading", object, d)
}

func makeBucketPublic(bucket string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	acl := client.Bucket(bucket).DefaultObjectACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		log.Fatal(err)
	}
}
