package main

func executeContract(z string) {
	println("contract request ... " + z)

	//senz := parse(z)

	// save event (request received)
	//t := eventTrans("orderzreq", "chainz", "Contract request received")
	//createTrans(t)

	rz := "DATA #print yes #qunt 323 @orderzresp ^chainz digisg"

	// save event (response send)
	//t = eventTrans("chainz", "orderzresp", "Forward contract request")
	//createTrans(t)

	// TODO execute contract function
	kmsg := Kmsg{
		Topic: "opsresp",
		Msg:   rz,
	}
	kchan <- kmsg
}
