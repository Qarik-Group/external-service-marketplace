package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func TweedURL() string {
	url, here := os.LookupEnv("TWEED_URL")
	fmt.Printf(url)
	if !here {
		log.Fatal("There was no TWEED_URL found on ur local environment")
	}
	return url
}

func GetUserName() string {
	un, err := os.LookupEnv("TWEED_USERNAME")
	fmt.Printf(un)
	if !err {
		log.Fatal("Set your TWEED_USERNAME env variable before using cli")
	}
	return un
}

func ReadResponse(r *http.Response) []byte {
	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Tried reading the response from the user and it aint good\n" + err.Error())
	}
	return bodyBytes
}

func MakeBody(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatal("Error converting struct to []byte in MakeBody\n struct:" + err.Error())
	}
	return b
}
