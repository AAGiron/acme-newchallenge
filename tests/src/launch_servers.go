package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"sync"
)

var wg sync.WaitGroup

// startServerHybrid launches a goroutine that will run a TLS server in the specified ipserver:port address
func startServerHybrid(clientMsg, serverMsg string, serverConfig *tls.Config, ipserver string, port string) {		
	go testConn(clientMsg, serverMsg, serverConfig, "server", ipserver, port)
}

func launchServer() {

	fmt.Println("Starting servers...")

	flag.Parse()

	port := 4433

	tlsInitCSVServer()
	
	strport := fmt.Sprintf("%d", port)

	kexCurveID, err := nameToCurveID(*kexAlgorithm)
	if err != nil {
		log.Fatal(err)
	}						
	
	serverConfig := initServer(kexCurveID)

	wg.Add(1)
	//start

	fmt.Println("\n" + "Starting " + *kexAlgorithm + " TLS server at " + *IPserver + ":" + strport + "...")

	startServerHybrid(clientHSMsg, serverHSMsg, serverConfig, *IPserver, strport)
	
	wg.Wait() //endless but required
}

func main() {
	launchServer()
}
