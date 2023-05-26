package main

import (
	"flag"
	"fmt"	
	"time"
	"crypto/tls"
)


func main_client() {
	flag.Parse()

	port := 4433
	
	handshakeSizes := make(map[string]uint32)
	var cconnState tls.ConnectionState

	tlsInitCSVClient()

	// struct for the metrics
	var algoResults tlsClientResultsInfo

	// list of structs
	var algoResultsList []tlsClientResultsInfo

	strport := fmt.Sprintf("%d", port)

	kexCurveID, err := nameToCurveID(*kexAlgorithm)
	if err != nil {
		panic(err)
	}			
	
	clientConfig := initClient(kexCurveID)								

	fmt.Printf("Starting TLS Handshakes: KEX Algorithm: %s\n", *kexAlgorithm)

	var timingsFullProtocol []float64
	var timingsProcessServerHello []float64
	var timingsWriteClientHello []float64
	var serverAuthAlgorithmName string

	for i := 0; i < *handshakes; i++ {
		var timingState timingInfo
		var err error
		var success bool
		
		timingState, cconnState, err, success = testConn(clientHSMsg, serverHSMsg, clientConfig, "client", *IPserver, strport)
		if err != nil || success == false{
			//log.Fatal(err)
			i--
			continue
		}
		timingsFullProtocol = append(timingsFullProtocol, float64(timingState.clientTimingInfo.FullProtocol)/float64(time.Millisecond))
		timingsProcessServerHello = append(timingsProcessServerHello, float64(timingState.clientTimingInfo.ProcessServerHello)/float64(time.Millisecond))
		timingsWriteClientHello = append(timingsWriteClientHello, float64(timingState.clientTimingInfo.WriteClientHello)/float64(time.Millisecond))			
	}

	handshakeSizes["ClientHello"] = cconnState.ClientHandshakeSizes.ClientHello							
	handshakeSizes["Certificate"] = cconnState.ClientHandshakeSizes.Certificate
	handshakeSizes["CertificateVerify"] = cconnState.ClientHandshakeSizes.CertificateVerify
	handshakeSizes["Finished"] = cconnState.ClientHandshakeSizes.Finished

	// TODO: Check if it is necessary to change this.
	serverAuthAlgorithmName = "unknown"

	//save results first
	tlsSaveCSVClient(timingsFullProtocol, timingsProcessServerHello, timingsWriteClientHello, *kexAlgorithm, serverAuthAlgorithmName, *handshakes, handshakeSizes)

	algoResults = tlsComputeStats(timingsFullProtocol, timingsProcessServerHello, timingsWriteClientHello, *handshakes)
	algoResults.kexName = *kexAlgorithm
	algoResults.authName = serverAuthAlgorithmName

	algoResultsList = append(algoResultsList, algoResults)		

	
	tlsPrintStatistics(algoResultsList)

	fmt.Println("End of test.")
}
