package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
)

type Senzie struct {
    name        string
	out         chan string
    quit        chan bool
    tik         chan string
    reader      *bufio.Reader
    writer      *bufio.Writer
    conn        *net.TCPConn
}

type Senz struct {
    Msg         string
    Uid          string
    Ztype       string
    Sender      string
    Receiver    string
    Attr        map[string]string
    Digsig      string
}

func main() {
    // first init key pair
    setUpKeys()

    // init cassandra session
    initCStarSession()

    // address
    tcpAddr, err := net.ResolveTCPAddr("tcp4", config.switchHost + ":" + config.switchPort)
    if err != nil {
        fmt.Println("Error address:", err.Error())
        os.Exit(1)
    }

    // tcp connect
    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

    // close on app closes
    defer conn.Close()

    fmt.Println("connected to switch")

    // create senzie
    senzie := &Senzie {
        name: config.senzieName,
        out: make(chan string),
        quit: make(chan bool),
        tik: make(chan string),
        reader: bufio.NewReader(conn),
        writer: bufio.NewWriter(conn),
        conn: conn,
    }
    registering(senzie)

    // close session
    clearCStarSession()
}

func registering(senzie *Senzie) {
    // send reg
    z := regSenz()
    senzie.writer.WriteString(z + ";")
    senzie.writer.Flush()

    // listen for reg status
    msg, err := senzie.reader.ReadString(';')
    if err != nil {
        fmt.Println("Error reading: ", err.Error())

        senzie.conn.Close()
        os.Exit(1)
    }

    // parse senz
    // check reg status
    senz := parse(msg)
    if(senz.Attr["status"] == "REG_DONE" || senz.Attr["status"] == "REG_ALR") {
        println("reg done...")
        // start reading and writing
        go writing(senzie)
        reading(senzie)
    } else {
        // close and exit
        senzie.conn.Close()
        os.Exit(1)
    }
}

func reading(senzie *Senzie) {
    READER:
    for {
        // read data
        msg, err := senzie.reader.ReadString(';')
        if err != nil {
            fmt.Println("Error reading: ", err.Error())

            senzie.quit <- true
            break READER
        }

        // not handle TAK, TIK, TUK
        if (msg == "TAK;") {
            // when connect, we recive TAK
            continue READER
        } else if(msg == "TIK;") {
            // send TIK
            senzie.tik <- "TUK;"
            continue READER
        } else if(msg == "TUK;") {
            continue READER
        }

        println("---- " + msg)

        // parse and handle
        senz := parse(msg)
        go handling(senzie, &senz)
    }
}

func writing(senzie *Senzie)  {
    // write
    WRITER:
    for {
        select {
        case <- senzie.quit:
            println("quiting/write -- ")
            break WRITER
        case senz := <-senzie.out:
            println("writing -- ")
            println(senz)
            // send
            senzie.writer.WriteString(senz + ";")
            senzie.writer.Flush()
        case tik := <- senzie.tik:
            println("ticking -- " )
            senzie.writer.WriteString(tik)
            senzie.writer.Flush()
        }
    }
}

func handling(senzie *Senzie, senz *Senz) {
    // frist send AWA back
    senzie.out <- awaSenz(senz.Attr["uid"])

    if(senz.Ztype == "SHARE") {
        // we only handle share cheques
        if cId, ok := senz.Attr["cid"]; !ok {
            // this means new cheque
            // create cheque
            cheque := senzToCheque(senz)
            cheque.Id = uuid()
            createCheque(cheque)

            // create trans
            trans := senzToTrans(senz)
            trans.ChequeId = cheque.Id
            trans.Status = "TRANSFER"
            createTrans(trans)

            // TODO send status back to fromAcc

            // forward cheque to toAcc
            senzie.out <- forwardChequeSenz(cheque, senz.Sender, senz.Attr["to"], uid())
        } else {
            // this mean already transfered cheque
            // check for double spend
            if(isDoubleSpend(senz.Sender, senz.Attr["to"], cId)) {
                // TODO send error status back
            } else {
                // TODO create trans
                trans := senzToTrans(senz)
                // TODO trans.ChequeId = cId
                trans.Status = "DEPOSIT"
                createTrans(trans)
            }
        }
    }
}