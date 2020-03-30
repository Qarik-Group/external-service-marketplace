package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/starkandwayne/external-service-marketplace/tweed"
	"github.com/starkandwayne/external-service-marketplace/util"

	realtweed "github.com/tweedproject/tweed"
)

type Catalog struct {
	realtweed.Catalog
}

func (c *Catalog) Merge(prefix string, other realtweed.Catalog) {
	for _, svc := range other.Services {
		svc.ID = fmt.Sprintf("%s--%s", prefix, svc.ID)
		svc.Name = fmt.Sprintf("[%s] %s", prefix, svc.Name)
		c.Services = append(c.Services, svc)
	}
}

type API struct {
	Config *Config
	Bind   string
}

func (a *API) Catalog() (Catalog, error) {
	c := Catalog{}

	for _, broker := range a.Config.ServiceBrokers {
		cat := tweed.Catalog(broker.Username, broker.Password, broker.URL)
		c.Merge(broker.Prefix, cat)
	}

	return c, nil
}

func (a *API) Provision(prefix, service, plan string) (string, error) {
	fmt.Printf("PROVISIONING [%s][%s][%s]\n", prefix, service, plan)

	broker, found := a.Config.Broker(prefix)
	if !found {
		return "", fmt.Errorf("no such broker '%s'", prefix)
	}

	fmt.Printf("PROVISIONING against tweed at %s (u:%s, p:%s)\n", broker.URL, broker.Username, broker.Password)
	// This is where we dispatch off to the actual broker.

	// in James' ideal world, here's what we do
	/*
	instance, err := broker.Backend.Provision(service, plan)
	if err != nil {
		return "", err
	}

	return instance.ID, nil
	*/
	return "", nil
}

var config *Config
var url string

func findTweed(tweedname string) int {
	brokers := config.ServiceBrokers

	for i := 0; i < len(brokers); i++ {
		if tweedname == brokers[i].Prefix {
			return i
		}
	}
	return -1
}
func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Endpoint Hit")
}
func bindFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tweedIndex := findTweed(vars["tweed"])
	//instance := vars["instance"]
	//binding := vars["binding"]
	//nowait := vars["nowait"]
	username := config.ServiceBrokers[tweedIndex].Username
	password := config.ServiceBrokers[tweedIndex].Password
	url := config.ServiceBrokers[tweedIndex].URL
	//username, password, _ := r.BasicAuth()

	var bindCmd util.BindCommand
	err := json.NewDecoder(r.Body).Decode(&bindCmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error Reading Request Body"))
		return
	}

	res := tweed.Bind(username, password, url, bindCmd)
	if res.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(res.Error))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.Ref))

}

func unbindFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance := vars["instance"]
	binding := vars["binding"]
	tweedIndex := findTweed(vars["tweed"])
	//nowait := vars["nowait"]

	username := config.ServiceBrokers[tweedIndex].Username
	password := config.ServiceBrokers[tweedIndex].Password
	url := config.ServiceBrokers[tweedIndex].URL
	//username, password, _ := r.BasicAuth()

	instancebinding := make([]string, 2)
	instancebinding[0] = instance
	instancebinding[1] = binding

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error Reading Request Body"))
		return
	}

	var unbindCmd util.UnbindCommand
	err = json.Unmarshal(body, &unbindCmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Did you use the UnbindCommand struct in util directory when you created your request. Please use that format"))
		return
	}

	//unbindCmd.Args.InstanceBinding = instancebinding

	res := tweed.UnBind(username, password, url, unbindCmd)
	util.JSON(res)
	data, _ := json.Marshal(res)
	if res.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.Ref))

}

func provisionFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tweedIndex := findTweed(vars["tweed"])
	//instance := vars["instance"]
	service := vars["service"]
	plan := vars["plan"]
	//nowait := vars["nowait"]
	username := config.ServiceBrokers[tweedIndex].Username
	password := config.ServiceBrokers[tweedIndex].Password
	url := config.ServiceBrokers[tweedIndex].URL
	//username, password, _ := r.BasicAuth()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error Reading Request Body"))
		return
	}
	var provCmd util.ProvisionCommand
	s := make([]string, 2)
	s[0] = service
	s[1] = plan
	json.Unmarshal(body, &provCmd)

	//provCmd.Args.ServicePlan = s

	res := tweed.Provision(username, password, url, provCmd)
	if res.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error In Request"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.Ref))

}

func deprovisionFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instance := vars["instance"]
	tweedIndex := findTweed(vars["tweed"])
	//nowait := vars["nowait"]
	username := config.ServiceBrokers[tweedIndex].Username
	password := config.ServiceBrokers[tweedIndex].Password
	url := config.ServiceBrokers[tweedIndex].URL
	//username, password, _ := r.BasicAuth()

	s := make([]string, 1)
	s[0] = instance

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error Reading Request Body"))
		return
	}
	var deprovCmd util.DeprovisionCommand
	//deprovCmd.Args.InstanceIds = instance
	json.Unmarshal(body, &deprovCmd)

	res := tweed.DeProvision(username, password, url, deprovCmd)
	if res.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error In Request"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.Ref))

}

func (api API) Run() {
	config = api.Config
	//url = "http://10.128.32.138:31666"
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	//register broker
	r.HandleFunc("/register/{broker}", testResponse)
	//retrieve clouds
	r.HandleFunc("/clouds", testResponse)
	//retrieve a specific cloud
	r.HandleFunc("/clouds/{cloud}", testResponse)

	r.HandleFunc("/catalog", func(w http.ResponseWriter, r *http.Request) {
		cat, err := api.Catalog()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don't do it.
			return
		}

		b, err := json.Marshal(cat)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don't do it.
			return
		}

		w.WriteHeader(200)
		fmt.Fprintf(w, "%s\n", string(b))
	})

	//provision a service
	r.HandleFunc("/provision/{service}/{plan}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		parts := strings.SplitN(vars["service"], "--", 2)
		if len(parts) != 2 {
			panic("ahhhhhh") // FIXME don't do this its bad
		}

		prefix := parts[0]
		service := parts[1]
		plan := vars["plan"]

		inst, err := api.Provision(prefix, service, plan)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don"t do it.
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, "OK %s\n", inst) // FIXME - use JSON, give some info back
	})

	//get an instance
	r.HandleFunc("/instances/{instance}", testResponse)
	//deprovision an instance
	r.HandleFunc("/deprovision/{instance}", deprovisionFunction)
	//bind an instance
	r.HandleFunc("/bind/{instance}/{binding}/{nowait}", bindFunction)
	//retrieve binding
	r.HandleFunc("/getbinding/{instance}", testResponse)
	//unbind an instance
	r.HandleFunc("/unbind/{instance}", unbindFunction)

	//start the server
	log.Fatal(http.ListenAndServe(api.Bind, r))
}
