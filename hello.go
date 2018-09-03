package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func main() {
	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		var t = template.Must(template.ParseFiles("flacupload.html"))
		if err := t.Execute(w, nil); err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		apikey, err := ioutil.ReadFile("apikey")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var json = `{
		  "config" : {
		    "languageCode" : "ja-JP",
		    "maxAlternatives" : 10
		  },
		  "audio" : {
		    "content" : "` + base64.StdEncoding.EncodeToString(bufbody.Bytes()) + `"
		  }
		}`
		jsonStr := []byte(json)
		ctx := appengine.NewContext(r)
		client := urlfetch.Client(ctx)
		url := "https://speech.googleapis.com/v1/speech:recognize?key=" + string(apikey)
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintln(w, "response Body:", string(body))
	}
}
