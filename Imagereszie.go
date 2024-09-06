package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

// please put the file path.
// openfile > decodefile > resizefile > createemptyfile > encodefileto created one > return file.

// TempDir의 이미지파일을 받으면,파일 리사이징 과 동시에 UUID 추가함 고유식별자를 만들어줌
// 이미지 자체를 객체화해서 다루어도 되지만, 필요없음 간단하기때문...
func ImageResizer(str string) string {
	//오픈 하고
	file, err := os.Open(str)
	if err != nil {
		log.Fatal("fatal error occured", err)
	}
	defer file.Close()
	//디코드 하고고
	decodedfile, _, _ := image.Decode(file)
	if err != nil {
		log.Fatal("fatal error occured", err)
	}
	//이미지 리사이즈 해주고
	resizedimg := resize.Resize(300, 400, decodedfile, resize.Lanczos3)

	//User upload > NginX > API > PostgreDB > URL 작성 https://hyunhoworld.site/files/"파일제목_PK.PNG"

	//Proxy Server Part >> 프록시 개념 더 공부해보기
	//NginX에 location 블록 추가 > /files/ 파일서버 포트 8085 Completed.

	//프론트엔드 파일업로드 구현,
	//업로드 파일 데이터베이스에 URL 형태로 저장하는 api 작성,
	//업로드된 파일 fileserver api를 이용하여 보내기 코드작성,
	//fileserver api 코드 작성 > get image from wepAPI서버 and resize image and save it on mounted volume//통신프로토콜 http 통신템플릿 json
	//API서버에서 이미지를 받을 때 Json으로 받고 받은 데이터를 디코드하지않고 바로 Fileserver로 토스해도됨.

	//받은 이미지, http로 Private 네트워크의 fileserver로 전달, 전달된 이미지는 이미지 핸들러를 통해 300x400사이즈로 변경 후
	//파일 제목에 _PK 추가, 같이온 PK아이디 참조하고, 파일서버내 ImageDir 폴더에 저장.

	//Option 파일서버에 이미지가 많이 있는경우를 대비해서, 해당 이미지 데이터를 분산 저장할수있고, 이미지 폴더를 리스트화 하여,
	//검색하여 찾는기능 추가 또는 이미지를 디코드하여 DB에 저장 또는 바이트형태로 문서로 저장

	//inject a Identity of image and resize it to proper size showing on page.

	//이름 새로만들고
	newname := uuid.NewString()
	//경로 지정해주고
	resizedfilepath := fmt.Sprintf("./ImageDir/%s.png", newname)

	//해당경로에 파일을만들고
	resiezedimgfile, _ := os.Create(resizedfilepath)
	defer resiezedimgfile.Close()
	newresizedfilepath := fmt.Sprintf("ImageDir/%s.png", newname)
	//리사이즈된 이미지 엔코드 해주고
	png.Encode(resiezedimgfile, resizedimg)

	os.Remove(str)

	//끝 Desired State > ImageDir에 UUID.png파일 이있어야함
	fmt.Println("the image that you has input has been encoded with new size png file.")
	return newresizedfilepath
}
