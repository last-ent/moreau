package data

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
SELECT e.rowid, e.title, e.content, p.rowid
FROM essays e
JOIN published p
ON p.essay_id = e.rowid;
`

const allEssays = `
SELECT e.rowid, e.title, e.content
FROM essays e
`

const getEssay = `
SELECT e.rowid, e.title, e.content
FROM essays e
WHERE e.rowid = $1
`
