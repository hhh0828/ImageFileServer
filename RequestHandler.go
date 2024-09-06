package main

import "net/http"

func FileserverHandler() *http.ServeMux {

	mux := http.NewServeMux()
	//mux.HandleFunc("/"+imgsrc, Returnimagefile) // will be deprecated

	//외부로부터 www.hyunhoworld.site/files/imageupload. 너무헷갈령 도메인들어가니까 ㅡ.ㅡ
	mux.HandleFunc("/imageupload", ImageUpload)

	//정적 파일서버 URL www.hyunhoworld.site/files/"imgfile" 을 받으면 해당 imgfile파일을 반환한다.
	Fileserver := http.FileServer(http.Dir("./ImageDir"))
	mux.Handle("/", http.StripPrefix("/", Fileserver))

	return mux
}

/*
요청은 http://localhost:8770/ImageDir/"image.png"로온다.
localhost:8770/image.png

*/
