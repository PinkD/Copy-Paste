package cpst

import (
	"database/sql"
	"runtime"
)

type dB struct {
	db *sql.DB
}

func newDB(url string) *dB {
	db := &dB{}
	_db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	_db.SetMaxOpenConns(runtime.NumCPU())
	db.db = _db
	_, err = _db.Exec(
		`CREATE TABLE IF NOT EXISTS "content" (
					code BIGINT PRIMARY KEY ,
					sha CHAR(40) NOT NULL,
					content TEXT NOT NULL
				)`)
	if err != nil {
		panic(err)
	}
	_, err = _db.Exec("CREATE INDEX IF NOT EXISTS sha_index ON content (sha)")
	if err != nil {
		panic(err)
	}
	return db
}
func (db *dB) ContainsContent(sha, content string) (code uint64, err error) {
	row, err := db.db.Query("SELECT code, content FROM content WHERE sha = $1", sha)
	if err != nil {
		return
	}
	defer row.Close()
	if row.Next() {
		var _content string
		err = row.Scan(&code, &_content)
		if err != nil {
			return
		}
		if len(content) == len(_content) && content == _content { //same content
			return //same content
		}
	}
	return 0, nil //no sha or sha collision
}

func (db *dB) SaveContent(data *contentData) (err error) {
	_, err = db.db.Exec("INSERT INTO content (code, sha, content) VALUES ($1, $2, $3)", data.Code, data.Sha, data.Content)
	return
}

func (db *dB) GetContent(code uint64) (content string, err error) {
	row, err := db.db.Query("SELECT content FROM content WHERE code = $1", code)
	if err != nil {
		return
	}
	defer row.Close()
	if row.Next() {
		err = row.Scan(&content)
	}
	return
}

func (db *dB) GetCount() (count uint64, err error) {
	row, err := db.db.Query("SELECT code FROM content ORDER BY code DESC LIMIT 1")
	if err != nil {
		return
	}
	defer row.Close()
	if row.Next() {
		err = row.Scan(&count)
	}
	return
}
