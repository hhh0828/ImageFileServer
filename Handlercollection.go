package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/gommon/log"
)

type Project struct {
	Name      string `json:"name"`
	Shortdesc string `json:"shortdesc"`
	Longdesc  string `json:"longdesc"`
	Imgurl    string `json:"imgurl"`
}

type Imageurl struct {
	URL string
}

func Returnimagefile(w http.ResponseWriter, r *http.Request) {

	var img Imageurl
	json.NewDecoder(r.Body).Decode(&img)
	//get a 300x400 size image.
	http.ServeFile(w, r, img.URL)
	w.Header().Set("Content-Type", "application/json")

}

func Chartserving(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ChartsDir/index.yaml")
}

// 이미지 업로드시, 파일을 formfile형태로 받고 imagefile을 생성한다 생성한 파일을 PV에 저장한 후, 해당 위치를 API서버로 송신해준다.
// 그렇다면, 해당파일을 Resize와 PKey형식으로 고유한 파일로 만들기위한 인덱스 추가활동이 필요함.
/*
방법 1 check - 파일형식을 나눠서 이름_고유숫자 형식으로 저장.  >> 사용자가 _ 임의로 입력시 문제확률 100000%
방법 2 check - 해쉬함수를 이용하여 해쉬데이터를 비교하여 찾는방법
방법 3 check - 데이터베이스를 만들어 거기서 고유정보를 저장하여 비교하여 전달...
포인트 : 이미지파일을 저장할때, 이거를 식별 할 수있는... 그런정보로 저장해야함 보통 이름으로하지만, 중복된 이름이 두개 올라갈경우 문제가 생겨서..const
*/
//without image
func UploadProject(w http.ResponseWriter, r *http.Request) {

	var project Project

	project.Name = r.FormValue("name")
	project.Shortdesc = r.FormValue("shortdesc")
	project.Longdesc = r.FormValue("longdesc")

	w.Header().Set("Content-Type", "application/json")
	Webappserver := "http://172.17.0.3:8700/uploadproject"
	jsondata, err := json.Marshal(project)
	if err != nil {
		fmt.Println("error occured", err)
	}
	data, err := http.NewRequest("POST", Webappserver, bytes.NewBuffer(jsondata))
	if err != nil {
		fmt.Println("fatal error occured but can't stop server", err)
	}
	clicon, err := http.DefaultClient.Do(data)
	if err != nil {
		fmt.Println("error occured", clicon.StatusCode)
	}
	fmt.Println(clicon.StatusCode)

}

func ImageUpload(w http.ResponseWriter, r *http.Request) {
	//GET 이미지 from http.request.

	err := r.ParseMultipartForm(20)
	if err != nil {
		log.Error("error during parse data", err)
	}

	var project Project

	fmt.Println(r.FormFile("image"))

	imagefile, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	if imagefile == nil {
		UploadProject(w, r)
		return

	}
	defer imagefile.Close()

	//이름뿐인 파일 생성
	fmt.Println("here is working? if i didn't upload image on the post request? ")
	filePath := "./TempDir/" + handler.Filename
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	fmt.Println("the tempfile created", handler.Filename)
	defer dst.Close()

	//컨텐츠를 카피한다.
	_, err = io.Copy(dst, imagefile)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	//Image reszie UUID adding

	finalfile := ImageResizer(filePath)

	// multipart.fileheader 객체로부터 파일이름을 반환한 후 URL 반환
	fileUrl := fmt.Sprintf("https://www.hyunhoworld.site/files/%s", finalfile)
	w.Header().Set("Content-Type", "application/json")
	Webappserver := "http://172.17.0.3:8700/imageurlsaverequest"
	project.Name = r.FormValue("name")
	project.Shortdesc = r.FormValue("shortdesc")
	project.Longdesc = r.FormValue("longdesc")
	project.Imgurl = fileUrl
	jsonbytedata, err := json.Marshal(&project)
	if err != nil {
		fmt.Println("fatal error occured but can't stop server", err)
	}
	data, err := http.NewRequest("POST", Webappserver, bytes.NewBuffer(jsonbytedata))
	if err != nil {
		fmt.Println("fatal error occured but can't stop server", err)
	}

	serveranswer, err := http.DefaultClient.Do(data)
	if err != nil {
		fmt.Println("something wrong happened", err)
	}
	fmt.Println(serveranswer.StatusCode)

}

// ImageDelete API -
func ImageDelete(w http.ResponseWriter, r *http.Request) {

	var project Project
	json.NewDecoder(r.Body).Decode(&project)
	//https://www.hyunhoworld.site/files/8afbf97f-8824-4685-8cdc-4936986b67ca.png URL 형식
	fileuuid := "./ImageDir/" + strings.TrimLeft(project.Imgurl, "https://www.hyunhoworld.site/files/")
	os.Remove(fileuuid)

}
