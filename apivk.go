package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type RawResponse struct {
	Response []json.RawMessage `json:"response"`
}

func Bar() {
	fmt.Println("Bar")
}

type Response struct {
	Id          int          `json:"id"`
	FromId      int          `json:"from_id"`
	ToId        int          `json:"to_id"`
	Date        int          `json:"date"`
	PostType    string       `json:"post_type"`
	Text        string       `json:"text"`
	SignerId    int          `json:"signer_id"`
	IsPinned    int8         `json:"is_pinned"`
	Attachment  Attachment   `json:"attachment"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Type  string `json:"type"`
	Photo *Photo `json:"photo"`
	Link  *Link  `json:"link"`
}

type Photo struct {
	Pid        int    `json:"pid"`
	Aid        int    `json:"aid"`
	OwnerId    int    `json:"owner_id"`
	UserId     int    `json:"user_id"`
	Src        string `json:"src"`
	SrcBig     string `json:"src_big"`
	SrcSmall   string `json:"src_small"`
	SrcXbig    string `json:"src_xbig"`
	SrcXxbig   string `json:"src_xxbig"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Text       string `json:"text"`
	Created    int    `json:"created"`
	Access_key string `json:"access_key"`
}

type Link struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageSrc    string `json:"image_src"`
	ImageBig    string `json:"image_big"`
}

func tst() {
	raw := []byte(`{"response":[2468,{"id":5513,"from_id":-58014516,"to_id":-58014516,"date":1453516206,"post_type":"post","text":"Мир вам и ближним вашим!<br><br>Оправдываться - грех. Прп. авва Исайя поучал: \"Венец добродетелей - любовь; венец страстей - оправдание грехов своих\". Вместо оправданий будем говорить - \"прости\".<br><br>Братия Троицкой Селенгинской обители","signer_id":235194773,"is_pinned":1,"attachment":{"type":"photo","photo":{"pid":404037784,"aid":-7,"owner_id":-58014516,"user_id":100,"src":"http:\/\/cs631329.vk.me\/v631329773\/eefe\/8qL1uNcH-Kg.jpg","src_big":"http:\/\/cs631329.vk.me\/v631329773\/eeff\/IvmP9O8aFLc.jpg","src_small":"http:\/\/cs631329.vk.me\/v631329773\/eefd\/B6ADu5ntll0.jpg","src_xbig":"http:\/\/cs631329.vk.me\/v631329773\/ef00\/bfMgaBRpBcI.jpg","src_xxbig":"http:\/\/cs631329.vk.me\/v631329773\/ef01\/qcrVhaq--Ic.jpg","width":1000,"height":666,"text":"","created":1453470652,"access_key":"053e6253631cae5038"}},"attachments":[{"type":"photo","photo":{"pid":404037784,"aid":-7,"owner_id":-58014516,"user_id":100,"src":"http:\/\/cs631329.vk.me\/v631329773\/eefe\/8qL1uNcH-Kg.jpg","src_big":"http:\/\/cs631329.vk.me\/v631329773\/eeff\/IvmP9O8aFLc.jpg","src_small":"http:\/\/cs631329.vk.me\/v631329773\/eefd\/B6ADu5ntll0.jpg","src_xbig":"http:\/\/cs631329.vk.me\/v631329773\/ef00\/bfMgaBRpBcI.jpg","src_xxbig":"http:\/\/cs631329.vk.me\/v631329773\/ef01\/qcrVhaq--Ic.jpg","width":1000,"height":666,"text":"","created":1453470652,"access_key":"053e6253631cae5038"}},{"type":"link","link":{"url":"http:\/\/selenginskii-monastery.cerkov.ru\/sms-rassylka-monastyrya-pouchenie-dnya\/pouchenie-dnya-373\/","title":" » Поучение дня","description":"","image_src":"http:\/\/cs631418.vk.me\/v631418773\/c773\/OMWEVcI7cno.jpg","image_big":"http:\/\/cs631418.vk.me\/v631418773\/c775\/PbCrn5A8iuA.jpg"}}],"comments":{"count":0},"likes":{"count":77},"reposts":{"count":9}}]}`)

	var raw_res RawResponse
	err := json.Unmarshal(raw, &raw_res)
	if err != nil {
		log.Fatal("Error parsing json: ", err)
	}

	var res Response
	err = json.Unmarshal(raw_res.Response[1], &res)
	if err != nil {
		log.Fatal("Error parsing json: ", err)
	}

	log.Printf("%+v", res)
	log.Println(res.Attachments[0].Photo.Pid)
	log.Println(res.Attachments[1].Link.Title)
}
