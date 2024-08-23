package main

import "net/http"

func FileserverHandler() *http.ServeMux {

	mux := http.NewServeMux()
	//mux.HandleFunc("/"+imgsrc, Returnimagefile)

	Fileserver := http.FileServer(http.Dir("./ImageDir"))
	mux.Handle("/", http.StripPrefix("/", Fileserver))

	return mux
}

/*
요청은 http://localhost:8770/ImageDir/"image.png"로온다.
localhost:8770/image.png

*/
