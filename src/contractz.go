package main

func executeContract(z string) {
	println("contract request ... " + z)

	senz := parse(z)

	if senz.Attr["type"] == "PREQ" {
		// handle purchase order
		// TODO check weather given item/design exists
		rz := respSenz(senz.Attr["uid"], "DONE", "3ops")
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
		notifyPord()

		return
	}

	if senz.Attr["type"] == "PRINT" {
		// handle print request
		// call amc to print
		notifyPrnt()

		return
	}
}
