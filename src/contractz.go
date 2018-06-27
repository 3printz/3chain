package main

func executeContract(contract string) {
	println("executing... " + contract)

	// TODO save event (processing event)

	//z := "DATA #status done #qnt 323 " +

	// TODO execute contract function
	kmsg := Kmsg{
		Topic: "orderzresp",
		Msg:   contract,
	}
	kchan <- kmsg
}
