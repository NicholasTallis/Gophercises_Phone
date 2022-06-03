package main

import (
	"database/sql"
	"log"
	"os"
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

	println("Beginning")
	iChatDBLocation := "/Users/nicholas/Library/Messages/chat.db"
	//iChatDBLocation := "/Users/nicholas/Desktop/chat.db"
	c := &imessage.Config{
		SQLPath:   iChatDBLocation,                               // Set this correctly
		QueueSize: 10,                                            // 10-20 is fine. If your server is super busy, tune this.
		Retries:   3,                                             // run the applescript up to this many times to send a message. 3 works well.
		DebugLog:  log.New(os.Stdout, "[DEBUG] ", log.LstdFlags), // Log debug messages.
		ErrorLog:  log.New(os.Stderr, "[ERROR] ", log.LstdFlags), // Log errors.
	}
	s, err := imessage.Init(c)
	checkErr(err)

	done := make(chan imessage.Incoming) // Make a channel to receive incoming messages.
	s.IncomingChan(".*", done)           // Bind to all incoming messages.
	err = s.Start()                      // Start outgoing and incoming message go routines.
	checkErr(err)
	log.Print("waiting for msgs")

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

	//s.Send(imessage.Outgoing{Text: "no help for you", To: "+12038327601"})
	//println("Should have sent")
	/*
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
		must(phonedb.Reset("postgres", psqlInfo, dbname))

		psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
		must(phonedb.Migrate("postgres", psqlInfo))

		db, err := phonedb.Open("postgres", psqlInfo)
		must(err)
		defer db.Close()

		err = db.Seed()
		must(err)

		phones, err := db.AllPhones()
		must(err)
		for _, p := range phones {
			fmt.Printf("Working on... %+v\n", p)
			number := normalize(p.Number)
			if number != p.Number {
				fmt.Println("Updating or removing...", number)
				existing, err := db.FindPhone(number)
				must(err)
				if existing != nil {
					must(db.DeletePhone(p.ID))
				} else {
					p.Number = number
					must(db.UpdatePhone(&p))
				}
			} else {
				fmt.Println("No changes required")
			}
		}
	*/
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
