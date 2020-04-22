package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/tweedproject/tweed/api"

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
	Config *util.Config
	Bind   string
}

func (a *API) Catalog() (realtweed.Catalog, error) {
	c := Catalog{}

	for _, broker := range a.Config.ServiceBrokers {
		cat := tweed.Connect(a.Config).SingleCatalog(broker.URL)
		c.Merge(broker.Prefix, cat)
	}
	//loop over catalogs and condense
	/*cats := tweed.Connect(a.Config).Catalog()
	for _, catalog := range cats {
		c.Services = append(c.Services, catalog.Services)
	}*/

	return c.Catalog, nil
}

func (a *API) Provision(cmd util.ProvisionCommand, prefix string) (api.ProvisionResponse, error) {
	fmt.Printf("PROVISIONING [%s][%s][%s]\n", prefix, cmd.Service, cmd.Plan)
	broker, found := a.Config.Broker(prefix)
	var nothing api.ProvisionResponse
	if !found {
		return nothing, fmt.Errorf("no such broker '%s'", prefix)
	}
	fmt.Println("Provision From API")
	util.JSON(cmd)
	fmt.Printf("PROVISIONING against tweed at %s (u:%s, p:%s)\n", broker.URL, broker.Username, broker.Password)
	// This is where we dispatch off to the actual broker.
	t := tweed.Connect(a.Config)
	provInst := t.Provision(broker.URL, cmd)
	// in James' ideal world, here's what we do
	/*
		instance, err := broker.Backend.Provision(service, plan)
		if err != nil {
			return "", err
		}
		return instance.ID, nil
	*/
	return provInst, nil
}

func (a *API) Deprovision(prefix, instance string) (api.DeprovisionResponse, error) {
	fmt.Printf("DEPROVISIONING [%s][%s]\n", prefix, instance)

	broker, found := a.Config.Broker(prefix)
	var nothing api.DeprovisionResponse
	if !found {
		return nothing, fmt.Errorf("no such broker '%s'", prefix)
	}

	fmt.Printf("DEPROVISIONING against tweed at %s (u:%s, p:%s)\n", broker.URL, broker.Username, broker.Password)
	// This is where we dispatch off to the actual broker.
	t := tweed.Connect(a.Config)
	deprovInst := t.DeProvision(broker.URL, instance)

	// in James' ideal world, here's what we do
	/*
		instance, err := broker.Backend.Provision(service, plan)
		if err != nil {
			return "", err
		}
		return instance.ID, nil
	*/
	return deprovInst, nil
}
func (a *API) BindSVC(prefix string, instance string) (api.BindResponse, error) {
	fmt.Printf("Binding [%s][%s]\n", prefix, instance)

	broker, found := a.Config.Broker(prefix)
	var nothing api.BindResponse
	if !found {
		return nothing, fmt.Errorf("no such broker '%s'", prefix)
	}

	fmt.Printf("Binding against tweed at %s (u:%s, p:%s)\n", broker.URL, broker.Username, broker.Password)
	t := tweed.Connect(a.Config)
	bindInst := t.Bind(broker.URL, instance)

	return bindInst, nil
}
func (a *API) Unbind(prefix, instance string, binding string) (api.UnbindResponse, error) {
	fmt.Printf("Unbinding [%s][%s]\n", prefix, instance)

	broker, found := a.Config.Broker(prefix)
	var nothing api.UnbindResponse
	if !found {
		return nothing, fmt.Errorf("no such broker '%s'", prefix)
	}

	fmt.Printf("Unbinding against tweed at %s (u:%s, p:%s)\n", broker.URL, broker.Username, broker.Password)

	t := tweed.Connect(a.Config)
	unbindInst := t.UnBind(broker.URL, instance, prefix)

	return unbindInst, nil
}

var config *util.Config
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error Reading Request Body"))
			return
		}
		var provCmd util.ProvisionCommand
		json.Unmarshal(body, &provCmd)
		provCmd.Service = service
		provCmd.Plan = plan //not sure if this works

		inst, err := api.Provision(provCmd, prefix)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don"t do it.
			return
		} else if inst.Error != "" {
			w.Write([]byte(inst.Ref))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(inst.Ref))
		//json.NewEncoder(w).Encode(inst)
		fmt.Fprintf(w, "OK %s\n", inst) // FIXME - use JSON, give some info back
	})

	//get an instance
	r.HandleFunc("/instances/{instance}", testResponse)
	//deprovision an instance
	r.HandleFunc("/deprovision/{prefix}/{instance}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		inst, err := api.Deprovision(vars["prefix"], vars["instance"])
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don"t do it.
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(inst)
		//fmt.Fprintf(w, "OK %s\n", inst) // FIXME - use JSON, give some info back
	})
	//bind an instance
	r.HandleFunc("/bind/{prefix}/{instance}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		inst, err := api.BindSVC(vars["prefix"], vars["instance"])
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don"t do it.
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(inst)
		//fmt.Fprintf(w, "OK %s\n", inst) // FIXME - use JSON, give some info back
	})
	//retrieve binding
	r.HandleFunc("/getbinding/{instance}", testResponse)
	//unbind an instance
	r.HandleFunc("/unbind/{prefix}/{instance}/{binding}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		inst, err := api.Unbind(vars["prefix"], vars["instance"], vars["binding"])
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don"t do it.
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(inst)
		//fmt.Fprintf(w, "OK %s\n", inst) // FIXME - use JSON, give some info back
	})

	//start the server
	log.Fatal(http.ListenAndServe(":"+api.Bind, r))
}
