package gcp

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/last-ent/globr/data"
)

type GcpAPI struct {
	BucketName string
	ToDB       chan<- data.Essay
	Callback   <-chan data.PubRow
}

var gcpAPI *GcpAPI

func NewGcpAPI(bucketName string, dbCh chan<- data.Essay, callback <-chan data.PubRow) *GcpAPI {
	gcpAPI = &GcpAPI{
		ToDB:       dbCh,
		Callback:   callback,
		BucketName: bucketName,
	}
	go gcpAPI.upload()

	return gcpAPI
}

func (gcp *GcpAPI) upload() {
	for essay := range gcp.Callback {
		object := fmt.Sprintf("%d", essay.PubID)

		UploadToGCP(object, essay.Content)
	}
}

func UploadToGCP(object string, c string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
	content := strings.NewReader(c)

	wc := client.Bucket(gcpAPI.BucketName).Object(object).NewWriter(ctx)
	var d int64
	if d, err = io.Copy(wc, content); err != nil {
		panic(err)
	}
	if err := wc.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done uploading", object, d)

}

func (gcp *GcpAPI) UploadPost(form url.Values) {
	title := form["title"][0]
	// form["content"] is a single string but better safe than sorry.
	content := strings.Join(form["content"], "\n")
	// upload(title, strings.NewReader(content))
	toPub := false
	if form["to_publish"][0] == "y" {
		toPub = true
	}

	gcp.ToDB <- data.Essay{
		Title:   title,
		Content: content,

		ToPublish: toPub,
	}
}

func (gcp *GcpAPI) makeBucketPublic() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	acl := client.Bucket(gcp.BucketName).DefaultObjectACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		log.Fatal(err)
	}
}
