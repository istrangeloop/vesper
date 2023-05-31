package flashcards

import (
	"database/sql"
	"strings"
	"time"
	"vesper/utils"
)

type Card struct {
	id               int
	question, answer string
	//halflife_type      int
	//interval_in_days   int
	date time.Time
}
type Deck []Card

// TODO: make it not have to receive db
func (c Card) Add(db *sql.DB) {
	next := time.Now()
	next = next.Add(time.Hour * 24 * 7)
	insertDynStmt := `insert into "flashcards"("question", "answer", "date") values($1, $2, $3)`
	_, e := db.Exec(insertDynStmt, strings.TrimSpace(c.question), strings.TrimSpace(c.answer), next.UTC())
	utils.CheckError(e)
}

func (c Card) Delete(db *sql.DB) {
	deleteStmt := `delete from "flashcards" where id=$1`
	_, e := db.Exec(deleteStmt, c.id)
	utils.CheckError(e)
}

func (c Card) Update(db *sql.DB) {
	updateStmt := `update "flashcards" set "question"=$1, "answer"=$2, "date"=$3 where "id"=$4`
	date := time.Now().Add(time.Hour * 24 * 7)
	_, e := db.Exec(updateStmt, c.question, c.answer, date, c.id)
	utils.CheckError(e)

}

func (c Card) SendToFuture(db *sql.DB) {
	updateStmt := `update "flashcards" set "date"=$1 where "id"=$2`
	date := time.Now().Add(time.Hour * 24 * 7)
	_, e := db.Exec(updateStmt, date, c.id)
	utils.CheckError(e)
}

func Get(db *sql.DB, id int) Card {
	var question = Card{}
	selectStmt := `SELECT "question", "answer" FROM "flashcards" where "id"=$1`
	raw, e := db.Query(selectStmt, id)
	e = raw.Scan(&question.question, &question.answer)
	utils.CheckError(e)
	return question
}

func AllCards(db *sql.DB) []Card {
	var res = []Card{}
	rows, err := db.Query(`SELECT "id", "question", "answer", "date" FROM "flashcards"`)
	utils.CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var id int
		var question string
		var answer string
		var date time.Time

		err = rows.Scan(&id, &question, &answer, &date)
		utils.CheckError(err)
		res = append(res, Card{question: strings.TrimSpace(question), answer: strings.TrimSpace(answer), id: id, date: date})
	}

	utils.CheckError(err)
	return res
}

func DueCards(db *sql.DB) []Card {
	var res = []Card{}
	selectStmt := `SELECT "id", "question", "answer", "date" FROM "flashcards" where "date" < $1`
	rows, err := db.Query(selectStmt, time.Now())
	utils.CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var id int
		var question string
		var answer string
		var date time.Time

		err = rows.Scan(&id, &question, &answer, &date)
		utils.CheckError(err)
		res = append(res, Card{question: strings.TrimSpace(question), answer: strings.TrimSpace(answer), id: id, date: date})
	}

	utils.CheckError(err)
	return res
}
