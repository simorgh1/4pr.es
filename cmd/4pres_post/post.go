package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/url"

	"github.com/eawsy/aws-lambda-go/service/lambda/runtime"
	"github.com/inge4pres/4pr.es/shortener"
)

var domain = "4pr.es"

func handle(evt json.RawMessage, ctx *runtime.Context) (interface{}, error) {
	log.Println("Received event: ", string(evt))
	var values map[string]string
	if err := json.Unmarshal(evt, &values); err != nil {
		return nil, err
	}
	decoded, err := url.QueryUnescape(values["url"])
	if err != nil {
		log.Println("Decode URL err: ", err)
		return nil, err
	}
	surl, err := shortener.SaveShortUrl(decoded, shortener.GetDyndbTable())
	if err != nil {
		log.Println("Save Short Url err: ", err)
		return nil, err
	}
	buf := new(bytes.Buffer)
	resp := template.New("postresp")
	resp, err = resp.Parse(shortener.PostHtml)
	if err != nil {
		return nil, err
	}
	pr := &shortener.PostResp{
		Url: surl,
	}
	resp.Execute(buf, pr)
	return string(buf.Bytes()), nil
}

func init() {
	runtime.HandleFunc(handle)
}

func main() {}