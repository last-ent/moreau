package main

import (
	"database/sql"
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

	_ "github.com/mattn/go-sqlite3"
)

const bucketName = "www.last-ent.com"

var editorHtml []byte

type essay struct {
	Title     string
	Content   string
	ToPublish bool
}

var toDB chan essay

func init() {
	f, _ := os.Open("./editor.html")
	ed, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Unable to read html file:", err.Error())
	}

	editorHtml = []byte(ed)

	toDB = make(chan essay)

	dbSetup(toDB)
}

const crEssays = `
CREATE TABLE IF NOT EXISTS essays(
	title TEXT NOT NULL,
	content BLOB NOT NULL
);`

const crPublished = `
CREATE TABLE IF NOT EXISTS published(
	essay_id INTEGER
);`

const insEssay = `
INSERT INTO
	essays(title, content)
	VALUES($1, $2);
`

const insPub = `
INSERT INTO
	published(essay_id)
	VALUES($1);
`

const allPublished = `
SELECT e.rowid, *
FROM essays e
JOIN published p
ON p.essay_id = e.rowid;
`

const allEssays = `
SELECT e.rowid, *
FROM essays e
`

const getEssay = `
SELECT e.rowid, *
FROM essays e
WHERE e.rowid = $1
`

func handleDBCalls(db *sql.DB, ch chan essay) {

	for {
		newEssay := <-ch
		fmt.Println("New Essay:", newEssay)

		res, err := db.Exec(
			insEssay,
			newEssay.Title,
			newEssay.Content,
		)
		fmt.Printf("%#v\n%#v", res, err)
		eID, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		if !newEssay.ToPublish {
			continue
		}

		_, err = db.Exec(insPub, eID)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func dbSetup(ch chan essay) {
	db, err := sql.Open("sqlite3", "./db")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(crEssays); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(crPublished); err != nil {
		log.Fatal(err)
	}

	go handleDBCalls(db, ch)

}

func editorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		go uploadPost(r.Form)
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

func uploadPost(form url.Values) {
	title := form["title"][0]
	// form["content"] is a single string but better safe than sorry.
	content := strings.Join(form["content"], "\n")
	// upload(title, strings.NewReader(content))
	toPub := false
	if form["to_publish"][0] == "y" {
		toPub = true
	}

	toDB <- essay{
		Title:     title,
		ToPublish: toPub,
		Content:   content,
	}
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
