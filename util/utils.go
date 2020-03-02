package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func TweedURL() string {
	url, here := os.LookupEnv("TWEED_URL")
	if !here {
		var hostStr = "kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type==" + "ExternalIP" + ")].address"
		var portStr = "ubectl -n ${NAMESPACE:-tweed} get service tweed -o jsonpath='{.spec.ports[?(@.name==" + "tweed" + ")].nodePort}'"
		hostOut, err := exec.Command(hostStr).Output()
		if err != nil {
			log.Fatal("Error getting Tweed host from kubectl", err)
		}
		portOut, err := exec.Command(portStr).Output()
		if err != nil {
			log.Fatal("Error getting Tweed port from kubectl", err)
		}
		url = "http://" + string(hostOut) + ":" + string(portOut)
	}
	return url
}

func GetUserName() string {
	un, err := os.LookupEnv("TWEED_USERNAME")
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
