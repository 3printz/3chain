package main

import (
	"strings"
)

func executeContract(z string) {
	println("contract request ... " + z)

	senz := parse(z)

	if senz.Attr["type"] == "PREQ" {
		// handle purchase order
		// matching logic
		// TODO check with calling biz api
		var rz string
		if strings.EqualFold(senz.Attr["location"], "singapore") {
			// with amc1
			zid := "8c43a1e0-794f-11e8-8c3a-2f9c177c5396"
			rz = respSenz(senz.Attr["uid"], "YES", "3ops", zid)
		} else if strings.EqualFold(senz.Attr["location"], "malaysia") {
			// with amc2
			zid := "937ed420-794f-11e8-8c3a-2f9c177c5396"
			rz = respSenz(senz.Attr["uid"], "YES", "3ops", zid)
		} else {
			// not match
			rz = respSenz(senz.Attr["uid"], "NO", "3ops", "zid")
		}
		kmsg := Kmsg{
			Topic: "opsresp",
			Msg:   rz,
		}
		kchan <- kmsg

		return
	}

	if senz.Attr["type"] == "PORD" {
		// handle purchase order
		// call oem to get design
		//notifyPorder(senz)

		return
	}

	if senz.Attr["type"] == "PRINT" {
		// handle print request
		// call amc to print
		notifyPrnt()

		return
	}
}
