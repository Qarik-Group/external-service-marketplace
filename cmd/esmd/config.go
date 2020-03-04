package main

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
	return nil, nil
}
