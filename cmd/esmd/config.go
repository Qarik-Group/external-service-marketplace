package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServiceBrokers []struct {
		Prefix     string `yaml:"prefix"`
		URL        string `yaml:"url"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		SkipVerify bool   `yaml:"skip-verify"`
	} `yaml:"service-brokers"`

	Clouds []struct {
		ID   string `yaml:"id"`
		Name string `yaml:"name"`
		Type string `yaml:"type"`

		// figure out: how to specify creds for CF / K8s
	} `yaml:"clouds"`
}

func ReadConfig(path string) (*Config, error) {
	var c *Config
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c, err
}
