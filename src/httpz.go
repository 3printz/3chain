package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type NotifyPo struct {
	Oem    string
	Amc    string
	Msg    string
	Status string
}

type NotifyAmc struct {
	PO_ID  string
	AMC_ID string
}

type NotifyOem struct {
	PO_ID  string
	OEM_ID string
}

type NotifyPnt struct {
	Oem    string
	Amc    string
	Msg    string
	Status string
}

func notifyPord() {
	// json req
	obj := NotifyPo{
		Oem:    "oem1",
		Amc:    "amc1",
		Msg:    "match done",
		Status: "SUCCESS",
	}
	j, _ := json.Marshal(obj)
	notify(j, apiConfig.poApi)
}

func notifyPorder(senz Senz) {
	// notify amc
	obj1 := NotifyAmc{
		PO_ID:  senz.Attr["poid"],
		AMC_ID: senz.Attr["amcid"],
	}
	j, _ := json.Marshal(obj1)
	notify(j, senz.Attr["amcapi"])

	// notify oem
	obj2 := NotifyOem{
		PO_ID:  senz.Attr["poid"],
		OEM_ID: senz.Attr["oemid"],
	}
	j, _ = json.Marshal(obj2)
	notify(j, senz.Attr["oemapi"])
}

func notifyPrnt() {
	// json req
	obj := NotifyPnt{
		Oem:    "oem1",
		Amc:    "amc1",
		Msg:    "match done",
		Status: "SUCCESS",
	}
	j, _ := json.Marshal(obj)
	notify(j, apiConfig.prntApi)
}

func notify(j []byte, uri string) {
	log.Printf("INFO: send request to, %s with %s", uri, string(j))

	// new request
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: fail request, %s", err.Error())
		return
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	// status
	if resp.StatusCode != 200 {
		log.Printf("ERROR: fail request, status: %s response: %s", resp.StatusCode, string(b))
		return
	}
}
