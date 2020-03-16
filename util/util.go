package util

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func GetUsername() string {
	return os.Getenv("TWEED_USERNAME")
}

func GetPassword() string {
	return os.Getenv("TWEED_PASSWORD")
}

func GetTweedUrl() string {
	url := os.Getenv("TWEED_URL")
	strings.Trim(url, "http://")
	return url
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
