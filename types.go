package main

import "database/sql"

type Song struct {
	Id int64
	Song_name string
	File_path string
}

type User struct {
	Id int64
	First_name string
	Last_name string
	Selected_song int64
	Selected_song_nullInt64 sql.NullInt64
	Song_name string
}

type IndexPageContext struct {
	Songs []Song
	Users []User
}

type upDateUserContext struct {
	Songs []Song
	User User
}