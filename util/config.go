package util

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ServiceBroker struct {
	Prefix     string `yaml:"prefix"`
	URL        string `yaml:"url"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	SkipVerify bool   `yaml:"skip-verify"`
}

type Config struct {
	ServiceBrokers []ServiceBroker `yaml:"service-brokers"`

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
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}

	return c, err
}

func (c Config) Broker(prefix string) (ServiceBroker, bool) {
	for _, broker := range c.ServiceBrokers {
		if broker.Prefix == prefix {
			return broker, true
		}
	}
	return ServiceBroker{}, false
}
