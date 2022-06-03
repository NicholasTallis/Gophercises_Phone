package main

import (
	"database/sql"
	"log"
	"regexp"
	"strings"

	"github.com/golift/imessage"
	_ "github.com/golift/imessage"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "nicholas"
	password = "password"
	dbname   = "gophercises_phone"
)

func main() {
	println("Test")
}
func messageRoutine(s *imessage.Messages, done chan imessage.Incoming, end *bool) {
	for msg := range done { // wait here for messages to come in.
		if len(msg.Text) < 60 {
			log.Println("id:", msg.RowID, "from:", msg.From, "attachment?", msg.File, "msg:", msg.Text)
		} else {
			log.Println("id:", msg.RowID, "from:", msg.From, "length:", len(msg.Text))
		}
		if strings.HasPrefix(msg.Text, "Help") {
			// Reply to any incoming message that has the word "Help" as the first word.
			s.Send(imessage.Outgoing{Text: "no help for you", To: msg.From})
		}
		if strings.HasPrefix(msg.Text, "test") {
			s.Send(imessage.Outgoing{Text: "Auto Reply", To: msg.From})
		}
		if strings.HasPrefix(msg.Text, "Test") {
			s.Send(imessage.Outgoing{Text: "Auto Reply 2", To: msg.From})
		}
	}

}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}

func displayStudents(db *sql.DB) {
	row, err := db.Query("SELECT * FROM student ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var code string
		var name string
		var program string
		row.Scan(&id, &code, &name, &program)
		log.Println("Student: ", code, " ", name, " ", program)
	}
}

// func normalize(phone string) string {
// 	var buf bytes.Buffer
// 	for _, ch := range phone {
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}
// 	}
// 	return buf.String()
// }
