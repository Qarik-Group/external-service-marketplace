package main

import (
	"github.com/starkandwayne/external-service-marketplace/tweed"
)

type UniqueServices struct {
	services map[string]tweed.Service
}

//Constructor for UniqueServices map
func NewUniqueServices() *UniqueServices {
	us := UniqueServices{}
	us.services = make(map[string]tweed.Service)
	return &us
}

//Populates map with unique Services found in the config passed
func (us UniqueServices) EsmCatalog(config Config) {
	for i := 0; i < len(config.ServiceBrokers); i++ {
		broker := config.ServiceBrokers[i]
		cat := tweed.Catalog(broker.Username, broker.Password, broker.URL)
		for j := 0; j < len(cat.Services); j++ {
			if v, ok := us.services[cat.Services[j].ID]; !ok {
				us.services[cat.Services[j].ID] = v
			}
		}
	}
	for _, val := range us.services {
		tweed.JSON(val)
	}
}
