package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var err error

func indexHandler(w http.ResponseWriter, r *http.Request) {
	song_rows, songs_err := db.Query("SELECT * FROM songs ORDER BY id ASC")
	checkIntervalServerError(songs_err, w)
	user_rows, users_err := db.Query("SELECT * FROM users ORDER BY id ASC")
	checkIntervalServerError(users_err, w)

	var songIdToName = make(map[int64]string)
	var songs []Song
	var users []User
	var song Song
	var user User

	for song_rows.Next() {
		err = song_rows.Scan(&song.Id, &song.Song_name, &song.File_path)
		checkIntervalServerError(err, w)
		songIdToName[song.Id] = song.Song_name
		songs = append(songs, song)
	}
	for user_rows.Next() {
		err = user_rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Selected_song_nullInt64)
		if user.Selected_song_nullInt64.Valid {
			user.Song_name = songIdToName[user.Selected_song_nullInt64.Int64]
		} else {
			user.Song_name = "NULL"
		}
		checkIntervalServerError(err, w)
		users = append(users, user)
	}

	context := IndexPageContext{songs, users}
	t, err := template.ParseFiles("templates/index.html")
	checkIntervalServerError(err, w)
	err = t.Execute(w, context)
	checkIntervalServerError(err, w)
}

func serveSongCreationPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		redirectToIndex(w, r)
	}
	t, err := template.ParseFiles("templates/create_song.html")
	ifErrorIn(err, "parsing template file")
	err = t.Execute(w, nil)
	ifErrorIn(err, "sering the page")
}

func createHandler_song(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectToIndex(w, r)
	}
	var song Song
	song.Song_name = r.FormValue("song_name")
	song.File_path = r.FormValue("file_path")

	fmt.Println("Song to be created:")
	fmt.Println("song name: ", song.Song_name, "file path: ", song.File_path)

	_, err := db.Exec("INSERT INTO songs(song_name, file_path) VALUES ($1, $2)", song.Song_name, song.File_path)
	ifErrorIn(err, "insertion")
	redirectToIndex(w, r)
}

func serveUserCreationPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		redirectToIndex(w, r)
	}
	t, err := template.ParseFiles("templates/create_user.html")
	ifErrorIn(err, "parsing template file")

	song_rows, songs_err := db.Query("SELECT * FROM songs ORDER BY id ASC")
	ifErrorIn(songs_err, "quering for songs data")
	var songs []Song
	var song Song

	for song_rows.Next() {
		err = song_rows.Scan(&song.Id, &song.Song_name, &song.File_path)
		checkIntervalServerError(err, w)
		songs = append(songs, song)
	}

	err = t.Execute(w, songs)
	ifErrorIn(err, "serving the page")
}


func createHandler_user(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectToIndex(w, r)
	}
	var user User
	user.First_name = r.FormValue("first_name")
	user.Last_name = r.FormValue("last_name")
	user.Selected_song, _ = strconv.ParseInt(r.FormValue("selected_song"), 10, 64)

	fmt.Println("User to be created:")
	fmt.Println("first_name: ", user.First_name, "last_name: ", user.Last_name, "selected_song_id: ", user.Selected_song)
	stmt := "INSERT INTO users (first_name, last_name, selected_song) VALUES ('" + user.First_name + "', '" + user.Last_name + "', "
	if user.Selected_song != -1 {
		stmt += fmt.Sprint(user.Selected_song) + ")"
	} else {
		stmt += "NULL)"
	}
	fmt.Println(stmt)
	_, err := db.Exec(stmt)
	ifErrorIn(err, "insertion")
	redirectToIndex(w, r)
}

func serveSongUpdatePageHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "GET" {
		redirectToIndex(w, r)
	}
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	song_rows, songs_err := db.Query("SELECT * FROM songs WHERE id=$1", id)
	ifErrorIn(songs_err, "quering songs data")
	var song Song
	for song_rows.Next() {
		err = song_rows.Scan(&song.Id, &song.Song_name, &song.File_path)
		checkIntervalServerError(err, w)
	}
	t, err := template.ParseFiles("templates/update_song.html")
	ifErrorIn(err, "parsing template file")
	err = t.Execute(w, song)
	ifErrorIn(err, "serving the page")
}

func updateHandler_song(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectToIndex(w, r)
	}
	var song Song

	song.Id, _ = strconv.ParseInt(r.FormValue("id"), 10, 64)
	song.Song_name = r.FormValue("song_name")
	song.File_path = r.FormValue("file_path")

	fmt.Println("Updated Song: ")
	fmt.Println("song name: ", song.Song_name, "file path: ", song.File_path)

	stmt := "UPDATE songs SET song_name = '" + song.Song_name + "', file_path = '" + song.File_path + "' WHERE id = " + fmt.Sprint(song.Id)
	fmt.Println(stmt)

	_, err := db.Exec("UPDATE songs SET song_name=$1, file_path=$2 WHERE id=$3", song.Song_name, song.File_path, song.Id)
	ifErrorIn(err, "execution")
	redirectToIndex(w, r)
}

func serveUserUpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		redirectToIndex(w, r)
	}
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	song_rows, songs_err := db.Query("SELECT * FROM songs ORDER BY id ASC")
	checkIntervalServerError(songs_err, w)
	user_rows, users_err := db.Query("SELECT * FROM users WHERE id=$1", id)
	checkIntervalServerError(users_err, w)

	var songs []Song
	var song Song
	var user User

	for song_rows.Next() {
		err = song_rows.Scan(&song.Id, &song.Song_name, &song.File_path)
		checkIntervalServerError(err, w)
		songs = append(songs, song)
	}
	for user_rows.Next() {
		err = user_rows.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Selected_song_nullInt64)
		checkIntervalServerError(err, w)
	}

	fmt.Println(user.First_name, user.Last_name, user.Selected_song_nullInt64.Valid)

	if user.Selected_song_nullInt64.Valid {
		user.Selected_song = user.Selected_song_nullInt64.Int64
	} else {
		user.Selected_song = -1
	}

	context := upDateUserContext{songs, user}
	t, err := template.ParseFiles("templates/update_user.html")
	checkIntervalServerError(err, w)
	err = t.Execute(w, context)
	checkIntervalServerError(err, w)
}

func updateHandler_user(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		redirectToIndex(w, r)
	}
	var user User

	user.Id, _ = strconv.ParseInt(r.FormValue("id"), 10, 64)
	user.First_name = r.FormValue("first_name")
	user.Last_name = r.FormValue("last_name")
	user.Selected_song, _ = strconv.ParseInt(r.FormValue("selected_song"), 10, 64)

	var stmt string
	if user.Selected_song != -1 {
		stmt = "UPDATE users SET first_name = '" + user.First_name + "', last_name = '" + user.Last_name + "', selected_song = '" + fmt.Sprint(user.Selected_song) + "' WHERE id = " + fmt.Sprint(user.Id)
	} else {
		stmt = "UPDATE users SET first_name = '" + user.First_name + "', last_name = '" + user.Last_name + "', selected_song = NULL" + " WHERE id = " + fmt.Sprint(user.Id)
	}
	fmt.Println(stmt)

	_, err := db.Exec(stmt)
	ifErrorIn(err, "updating uesr")
	redirectToIndex(w, r)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Println("aha")
		redirectToIndex(w, r)
	}
	var table_name = r.FormValue("type") + "s"
	var id, _ = strconv.ParseInt(r.FormValue("id"), 10, 64)

	stmt := "DELETE FROM "+ table_name + " WHERE id = " + fmt.Sprint(id)
	fmt.Println("id of the table entry to be deleted: ", id)
	fmt.Println(stmt)

		_, err := db.Exec(stmt)
	ifErrorIn(err, "deletion")
	redirectToIndex(w, r)
}

