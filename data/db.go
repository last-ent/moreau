package data

import (
	"database/sql"
	"fmt"
	"log"
)

type Essay struct {
	Title     string
	Content   string
	ToPublish bool
}

type PubRow struct {
	PubID   int64
	EssayID int64

	Title   string
	Content string
}

type DBInstance struct {
	DB          *sql.DB
	GcpCallback chan<- PubRow
}

func newDBInstance(db *sql.DB, callback chan<- PubRow) *DBInstance {
	if _, err := db.Exec(crEssays); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(crPublished); err != nil {
		log.Fatal(err)
	}

	return &DBInstance{DB: db, GcpCallback: callback}
}

func NewDBInstance(ch chan Essay, callback chan<- PubRow) *DBInstance {
	db, err := sql.Open("sqlite3", "./db")
	if err != nil {
		log.Fatal(err)
	}
	dbi := newDBInstance(db, callback)
	go dbi.handleDBCalls(ch)
	return dbi
}

func (dbi *DBInstance) handleDBCalls(ch chan Essay) {
	for newEssay := range ch {
		fmt.Println("New Essay:", newEssay)
		eID := dbi.InsertEssay(newEssay.Title, newEssay.Content)

		if newEssay.ToPublish {
			pID := dbi.InsertPub(eID)
			dbi.GcpCallback <- PubRow{
				PubID:   pID,
				Title:   newEssay.Title,
				Content: newEssay.Content,
			}
		}
		fmt.Printf("%#v\n", dbi.GetAllPublished())
	}
}
