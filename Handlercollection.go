package main

import (
	"encoding/json"
	"net/http"
)

type Imageurl struct {
	URL string
}

func Returnimagefile(w http.ResponseWriter, r *http.Request) {

	var img Imageurl
	json.NewDecoder(r.Body).Decode(&img)
	//get a 300x400 size image.
	http.ServeFile(w, r, img.URL)
	w.Header().Set("Content-Type", "Content-Type")

}

func Chartserving(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ChartsDir/index.yaml")
}
