package data

import "log"

func (dbi *DBInstance) InsertPub(eID int64) int64 {
	ips, err := dbi.DB.Prepare(insPub)
	if err != nil {
		log.Fatal(err)
	}

	res, err := ips.Exec(eID)
	if err != nil {
		log.Fatal(err)
	}

	newID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	if err := ips.Close(); err != nil {
		log.Fatal(err)
	}

	return newID
}

func (dbi *DBInstance) InsertEssay(title string, content string) int64 {
	ies, err := dbi.DB.Prepare(insEssay)
	if err != nil {
		log.Fatal(err)
	}

	res, err := ies.Exec(title, content)
	if err != nil {
		log.Fatal(err)
	}

	newID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	if err := ies.Close(); err != nil {
		log.Fatal(err)
	}
	return newID
}

func (dbi *DBInstance) GetAllPublished() []*PubRow {
	aps, err := dbi.DB.Prepare(allPublished)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := aps.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	pubRows := []*PubRow{}
	for rows.Next() {
		var rowID, pubID int64
		var title, content string

		if err := rows.Scan(&rowID, &title, &content, &pubID); err != nil {
			log.Fatal(err)
		}

		pubRows = append(pubRows, &PubRow{
			PubID:   pubID,
			EssayID: rowID,

			Title:   title,
			Content: content,
		})
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return pubRows
}

func (dbi *DBInstance) GetEssay(eID int64) *PubRow {
	ges, err := dbi.DB.Prepare(getEssay)
	if err != nil {
		log.Fatal(err)
	}

	var rowID, pubID int64
	var title, content string

	err = ges.QueryRow(eID).Scan(&rowID, &title, &content, &pubID)
	if err != nil {
		log.Fatal(err)
	}

	return &PubRow{
		PubID:   pubID,
		EssayID: rowID,

		Title:   title,
		Content: content,
	}
}

func (dbi *DBInstance) GetAllEssays() []*PubRow {
	aps, err := dbi.DB.Prepare(allEssays)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := aps.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	essayRows := []*PubRow{}
	for rows.Next() {
		var rowID int64
		var title, content string

		if err := rows.Scan(&rowID, &title, &content); err != nil {
			log.Fatal(err)
		}

		essayRows = append(essayRows, &PubRow{
			EssayID: rowID,

			Title:   title,
			Content: content,
		})
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return essayRows
}
