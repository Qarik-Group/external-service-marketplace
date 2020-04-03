package main

import (
	"github.com/starkandwayne/external-service-marketplace/tweed"
)

type UniqueServices struct {
	services map[string]tweed.Service
}

func newUniqueServices() *UniqueServices {
	us := UniqueServices{}
	us.services = make(map[string]tweed.Service)
	return &us
}

//Populates map with unique Services found in the config passed
func EsmCatalog(config Config) *UniqueServices {
	us := newUniqueServices()
	for i := 0; i < len(config.ServiceBrokers); i++ {
		broker := config.ServiceBrokers[i]
		cat := tweed.Catalog(broker.Username, broker.Password, broker.URL)
		for j := 0; j < len(cat.Services); j++ {
			if v, ok := us.services[cat.Services[j].ID]; !ok {
				us.services[cat.Services[j].ID] = v
			}
		}
	}
	return us
}

func PrintUniqueServices(us UniqueServices) {
	for k, _ := range us.services {
		tweed.JSON(us.services[k])
	}
}
