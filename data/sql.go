package data

type essay struct {
	Title     string
	Content   string
	ToPublish bool
}

type pubRow struct {
	PubID   int64
	EssayID int64

	Title   string
	Content string
}

var toDB chan essay
