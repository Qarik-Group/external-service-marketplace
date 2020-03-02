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

func GetPassword() string {
	pw, err := os.LookupEnv("TWEED_PASSWORD")
	if !err {
		log.Fatal("Set your TWEED_PASSWROD env variable before using cli")
	}
	return pw
}

func JSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{(error)} failed to marshal JSON:\n")
		fmt.Fprintf(os.Stderr, "        @R{%s}\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", string(b))
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
