package main

import ( //init
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
)

const (
	host = "192.168.2.163"
	port = 5432
	user = "jiahonghe"
	password = "asdfsfasdf"
	dbname = "gocrud"
)

var (
	db *sql.DB //concurrency / not global variable
)

//grourotine //channel
//return json instead of html

func main() { //args

	//connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close() // multiple defer/ custom defer
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	//route
	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/create/song", serveSongCreationPageHandler) //handlerfunc, how is it implemented
	http.HandleFunc("/create_song", createHandler_song)
	http.HandleFunc("/create/user", serveUserCreationPageHandler)
	http.HandleFunc("/create_user", createHandler_user)

	http.HandleFunc("/update/song", serveSongUpdatePageHandler)
	http.HandleFunc("/update_song", updateHandler_song)
	http.HandleFunc("/update/user", serveUserUpdatePageHandler)
	http.HandleFunc("/update_user", updateHandler_user)

	http.HandleFunc("/delete", deleteHandler)

	fmt.Println("Server is running at localhost:8080")
	http.ListenAndServe(":8080", nil)
}