package main

import (
	"fmt"
	"net/http"
)

func checkIntervalServerError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func redirectToIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", 301)
}

func ifErrorIn(err error, what string) {
	if err != nil {
		fmt.Println("error in ", what)
		panic(err)
	}
}