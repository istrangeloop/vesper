// you tell me a lot of teminders that you don't want to forget, so I'll keep them safe
// and show it to you on the screen.

// please use this to say nice things to yourself c:

package reminders

import (
	"arthur/utils"
	"database/sql"
)

type Reminder struct {
	text  string
	topic string
}

// TODO: make it not have to receive db
func AddReminder(db *sql.DB, Reminder string, topic string) {
	insertDynStmt := `insert into "reminders"("text", "topic") values($1, $2)`
	_, e := db.Exec(insertDynStmt, Reminder, topic)
	utils.CheckError(e)
}

func DeleteReminder(db *sql.DB, id int) {
	deleteStmt := `delete from "reminders" where id=$1`
	_, e := db.Exec(deleteStmt, 1)
	utils.CheckError(e)
}

func UpdateReminder(db *sql.DB, id int, Reminder string, topic string) {
	updateStmt := `update "reminders" set "text"=$1, "topic"=$2 where "id"=$3`
	_, e := db.Exec(updateStmt, Reminder, topic, id)
	utils.CheckError(e)

}

func ShowReminder(db *sql.DB, id int) {
	selectStmt := `SELECT "text", "topic" FROM "reminders" where "id"=$1`
	_, e := db.Exec(selectStmt, id)
	utils.CheckError(e)

}

func ShowReminders(db *sql.DB) []string {
	var res = []string{}
	rows, err := db.Query(`SELECT "text", "topic" FROM "reminders"`)
	utils.CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var text string
		var topic string

		err = rows.Scan(&text, &topic)
		utils.CheckError(err)
		res = append(res, text)
		//fmt.Println(text, topic)
	}

	utils.CheckError(err)
	return res
}
