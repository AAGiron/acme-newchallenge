package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// Command line flags to configure the tests parameters
var (	
	IPserver   = flag.String("ipserver", "", "IP of the TLS Server")
	clientIP = flag.String("ipclient", "", "IP or name of the TLS Client")
	serverPort = flag.Int("port", 4433, "TLS Server port")
	handshakes = flag.Int("handshakes", 1, "Number of Handshakes desired")
	clientAuth = flag.Bool("clientauth", false, "Client authentication")
	pqtls      = flag.Bool("pqtls", false, "PQTLS")
	classic    = flag.Bool("classic", false, "TLS with classic algorithms")
	pebbleRootCA = flag.String("pebbleroot", "", "Path to the Pebble Root CA certificate")	
	serverCertPath = flag.String("cert", "", "Path to the server certificate")
	serverKeyPath = flag.String("key", "", "Path to the server private key")
	prequantumScenario = flag.Bool("prescenario", false, "By setting this flag to true, it will be simulated the pre-quantum scenario of the PQCACME proposal")
	postquantumScenario = flag.Bool("postscenario", false, "By setting this flag to true, it will be simulated the post-quantum scenario of the PQCACME proposal")
	truststorePath = flag.String("tspath", "", "Path to the client's trustore")
	ocspResponseFilePath = flag.String("ocspresponsepath", "", "Path to the file where the OCSP Response is written to.")
	syncClient = flag.Bool("syncclient", false, "The TLS client will wait for the Pebble initialization before initializing the TLS configuration.")
	measurementsDir = flag.String("measurements", "", "Path to the directory where the measurements are written to")
	kexAlgorithm = flag.String("kex", "Kyber512", "KEX algorithm to be used in the handshake.")
)

var (
	// Handshake KEX algorithms supported for tests
	hsKEXAlgorithms = map[string]tls.CurveID{
		"P256": tls.CurveP256, "P384": tls.CurveP384, "P521": tls.CurveP521,
		"Kyber512": tls.OQS_Kyber512, "P256_Kyber512": tls.P256_Kyber512,
		"Kyber768": tls.OQS_Kyber768, "P384_Kyber768": tls.P384_Kyber768,
		"Kyber1024": tls.OQS_Kyber1024, "P521_Kyber1024": tls.P521_Kyber1024,
		"LightSaber_KEM": tls.LightSaber_KEM, "P256_LightSaber_KEM": tls.P256_LightSaber_KEM,
		"Saber_KEM": tls.Saber_KEM, "P384_Saber_KEM": tls.P384_Saber_KEM,
		"FireSaber_KEM": tls.FireSaber_KEM, "P521_FireSaber_KEM": tls.P521_FireSaber_KEM,
		"NTRU_HPS_2048_509": tls.NTRU_HPS_2048_509, "P256_NTRU_HPS_2048_509": tls.P256_NTRU_HPS_2048_509,
		"NTRU_HPS_2048_677": tls.NTRU_HPS_2048_677, "P384_NTRU_HPS_2048_677": tls.P384_NTRU_HPS_2048_677,
		"NTRU_HPS_4096_821": tls.NTRU_HPS_4096_821, "P521_NTRU_HPS_4096_821": tls.P521_NTRU_HPS_4096_821,
		"NTRU_HPS_4096_1229": tls.NTRU_HPS_4096_1229, "P521_NTRU_HPS_4096_1229": tls.P521_NTRU_HPS_4096_1229,
		"NTRU_HRSS_701": tls.NTRU_HRSS_701, "P384_NTRU_HRSS_701": tls.P384_NTRU_HRSS_701,
		"NTRU_HRSS_1373": tls.NTRU_HRSS_1373, "P521_NTRU_HRSS_1373": tls.P521_NTRU_HRSS_1373,
		"P256_Classic-McEliece-348864": tls.P256_Classic_McEliece_348864,
	}
	
	// Message to be sent by the client after the handshake
	clientHSMsg = "hello, server"

	// Message to be sent by the server after the handshake
	serverHSMsg = "hello, client"	
)

func initClientAndAuth(k, kAuth string) (*tls.Config, error) {
	var clientConfig *tls.Config
	
	kexCurveID, err := nameToCurveID(k)
	if err != nil {
		return nil, err
	}		
	
	clientConfig = initClient(kexCurveID)

	return clientConfig, nil
}

func nameToCurveID(name string) (tls.CurveID, error) {
	curveID, prs := hsKEXAlgorithms[name]
	if !prs {
		fmt.Println("Algorithm not found. Available algorithms: ")
		for name, _ := range hsKEXAlgorithms {
			fmt.Println(name)
		}
		return 0, errors.New("ERROR: Algorithm not found")
	}
	return curveID, nil
}

func curveIDToName(cID tls.CurveID) (name string, e error) {
	for n, id := range hsKEXAlgorithms {
		if id == cID {
			return n, nil
		}
	}
	return "0", errors.New("ERROR: Algorithm not found")
}

// initServer initializes the server TLS configuration
func initServer(curveID tls.CurveID) *tls.Config {	
	cfg := &tls.Config{
		MinVersion:                 tls.VersionTLS10,
		MaxVersion:                 tls.VersionTLS13,
		InsecureSkipVerify:         false,
		SupportDelegatedCredential: false,
		CurvePreferences:           []tls.CurveID{curveID},
		PQTLSEnabled: *pqtls,
		OCSPResponseFilePath: *ocspResponseFilePath,
		ServerName: *IPserver,
	}		

	if *clientAuth {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
		cfg.ClientCAs = x509.NewCertPool()
	}
	
	serverCert, err := tls.LoadX509KeyPair(*serverCertPath, *serverKeyPath)
	if err != nil {
		panic(err)
	}
	
	cfg.Certificates = make([]tls.Certificate, 1)		
	cfg.Certificates[0] = serverCert

	return cfg
}

// initClient initializes the client TLS configuration
func initClient(curveID tls.CurveID) *tls.Config {	

	ccfg := &tls.Config{
		MinVersion:                 tls.VersionTLS10,
		MaxVersion:                 tls.VersionTLS13,
		InsecureSkipVerify:         false,
		SupportDelegatedCredential: false,
		CurvePreferences:           []tls.CurveID{curveID},
		RootCAs: x509.NewCertPool(),
		PQTLSEnabled: *pqtls,
	}

	if !*postquantumScenario {
		if *syncClient {
			// Wait for Pebble to be ready
			const expectedMessage = "pebble is ready"		

			server, err := net.Listen("tcp", *clientIP + ":9000")
			if err != nil {
				panic(err)
			}
			
			defer server.Close()	
			
			fmt.Println("Waiting")
			connection, err := server.Accept()
			if err != nil {
				panic(err)
			}			
			
			buffer := make([]byte, len([]byte(expectedMessage)))
			
			_, err = connection.Read(buffer)
			if err != nil {
				panic(err)
			}
			
			fmt.Println("READY")

			if string(buffer) != expectedMessage {
				panic("received message does not match expected message")			
			}
			connection.Close()

			// Pebble is ready, now we can proceed with the normal execution
		}		

		certPEMBlock, err := os.ReadFile(*pebbleRootCA)
		if err != nil {
			panic(err)
		}		
	
		certDERBlock, certPEMBlock := pem.Decode(certPEMBlock)
	
		pebbleRootCACert, err := x509.ParseCertificate(certDERBlock.Bytes)
		if err != nil {
			panic(err)
		}			
		ccfg.RootCAs.AddCert(pebbleRootCACert)
	}
	
	return ccfg
}

func newLocalListener(ip string, port string) net.Listener {
	ln, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		ln, err = net.Listen("tcp6", "[::1]:0")
	}
	if err != nil {
		log.Fatal(err)
	}
	return ln
}

// timingInfo contains the client and the server handshake timing measurements
type timingInfo struct {
	serverTimingInfo tls.CFEventTLS13ServerHandshakeTimingInfo
	clientTimingInfo tls.CFEventTLS13ClientHandshakeTimingInfo
}

func (ti *timingInfo) eventHandler(event tls.CFEvent) {
	switch e := event.(type) {
	case tls.CFEventTLS13ServerHandshakeTimingInfo:
		ti.serverTimingInfo = e
	case tls.CFEventTLS13ClientHandshakeTimingInfo:
		ti.clientTimingInfo = e
	}
}

// testConn performs the tests connections between the client and the server
// This function is called by both client and server
func testConn(clientMsg, serverMsg string, tlsConfig *tls.Config, peer string, ipserver string, port string) (timingState timingInfo, cconnState tls.ConnectionState, err error, success bool) {	
	tlsConfig.CFEventHandler = timingState.eventHandler
	
	if peer == "server" {

		handshakeSizes := make(map[string]uint32)
		
		var timingsFullProtocol []float64
		var timingsWriteServerHello []float64
		var timingsWriteCertVerify []float64
		serverAuthAlgorithmName := "unknown"
		
		buf := make([]byte, len(clientMsg))

		countConnections := 0

		ln := newLocalListener(ipserver, port)
		defer ln.Close()		
		
		for {
			serverConn, err := ln.Accept()
			if err != nil {
				fmt.Print(err)
				fmt.Print("error 1 %v", err)
			}
			server := tls.Server(serverConn, tlsConfig)
			if err := server.Handshake(); err != nil {
				fmt.Printf("Handshake error %v", err)
			}			

			//server read client hello
			n, err := server.Read(buf)
			if err != nil || n != len(clientMsg) {
				fmt.Print(err)
				fmt.Print("error 2 %v", err)
			}

			//server responds
			n, err = server.Write([]byte(serverMsg))
			if n != len(serverMsg) || err != nil {
				//error
				fmt.Print(err)
				fmt.Print("error 3 %v", err)
			}			

			countConnections++

			cconnState = server.ConnectionState()			

			if *pqtls || *classic {

				if (*pqtls && cconnState.DidPQTLS) || *classic {
										
					if *clientAuth {
						if !cconnState.DidClientAuthentication {
							fmt.Println("Server unsuccessful TLS with mutual authentication")
							continue
						}
					}

					timingsFullProtocol = append(timingsFullProtocol, float64(timingState.serverTimingInfo.FullProtocol)/float64(time.Millisecond))
					timingsWriteServerHello = append(timingsWriteServerHello, float64(timingState.serverTimingInfo.WriteServerHello)/float64(time.Millisecond))
					timingsWriteCertVerify = append(timingsWriteCertVerify, float64(timingState.serverTimingInfo.WriteCertificateVerify)/float64(time.Millisecond))

					if countConnections == *handshakes {
						// var kAuth string
						var err error

						kKEX, e := curveIDToName(tlsConfig.CurvePreferences[0])
						if e != nil {
							fmt.Print("4 %v", err)
						}

						// if *classic {							
						// 	priv, _ := tlsConfig.Certificates[0].PrivateKey.(*ecdsa.PrivateKey)
						// 	kAuth, err = sigIDToName(priv.PublicKey.Curve)							
						// } else {							
						// 	priv, _ := tlsConfig.Certificates[0].PrivateKey.(*liboqs_sig.PrivateKey)
						// 	kAuth, err = sigIDToName(priv.SigId)
						// }
						
						if err != nil {
							fmt.Print("5 %v", err)
						}

						handshakeSizes["ServerHello"] = cconnState.ServerHandshakeSizes.ServerHello
						handshakeSizes["EncryptedExtensions"] = cconnState.ServerHandshakeSizes.EncryptedExtensions
						handshakeSizes["Certificate"] = cconnState.ServerHandshakeSizes.Certificate
						handshakeSizes["CertificateRequest"] = cconnState.ServerHandshakeSizes.CertificateRequest
						handshakeSizes["CertificateVerify"] = cconnState.ServerHandshakeSizes.CertificateVerify
						handshakeSizes["Finished"] = cconnState.ServerHandshakeSizes.Finished

						
						tlsSaveCSVServer(timingsFullProtocol, timingsWriteServerHello, timingsWriteCertVerify, kKEX, serverAuthAlgorithmName, countConnections, handshakeSizes)
						countConnections = 0
						timingsFullProtocol = nil
						timingsWriteCertVerify = nil
						timingsWriteServerHello = nil
					}
				} else {
					fmt.Println("Server unsuccessful TLS")
					continue
				}
			}
		}
	}

	if peer == "client" {

		buf := make([]byte, len(serverMsg))

		client, err := tls.Dial("tcp", ipserver+":"+port, tlsConfig)
		if err != nil {
			fmt.Print(err)
		}
		defer client.Close()

		client.Write([]byte(clientMsg))

		_, err = client.Read(buf)		

		cconnState = client.ConnectionState()

		if *pqtls {
			if cconnState.DidPQTLS {

				if *clientAuth {

					if cconnState.DidClientAuthentication {
						log.Println("Client Success using PQTLS with mutual authentication")
					} else {
						log.Println("Client unsuccessful PQTLS with mutual authentication")
						return timingState, cconnState, nil, false
					}

				} else {
					log.Println("Client Success using PQTLS")
				}
			} else {
				log.Println("Client unsuccessful PQTLS")
				return timingState, cconnState, nil, false
			}
		} else if *classic {
			if *clientAuth {
				if cconnState.DidClientAuthentication {
					log.Println("Client Success using TLS with mutual authentication")
				} else {
					log.Println("Client unsuccessful TLS with mutual authentication")
					return timingState, cconnState, nil, false
				}			
			} else {
				log.Println("Client Success using TLS")
			}
		}
	}

	return timingState, cconnState, nil, true
}

func launchHTTPSServer(serverConfig *tls.Config, port string) {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	
	addr := ":"+ port

	server := &http.Server{
		Addr: addr, 
		Handler: nil,
		TLSConfig: serverConfig,
	}

	err := server.ListenAndServeTLS("", "")
	
	if err != nil {
			log.Fatal("ListenAndServe: ", err)
	}
}

func sliceContains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
  }
	return false
}

func getClassicAlgorithmName(pub interface{}) string {
	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		panic("classic algorithm not supported. Expected a *ecdsa.PublicKey")		
	}

	switch ecdsaPub.Curve {
	case elliptic.P256():
		return "ECDSA-P256"
	case elliptic.P384():
		return "ECDSA-P384"
	case elliptic.P521():
		return "ECDSA-P521"
	default:
		panic("unknown elliptic curve")		
	}
}
