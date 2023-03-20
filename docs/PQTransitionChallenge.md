# ACME PQ-Transition Challenge details

We implemented a new ACME challenge (called "PQ-transition challenge"),
allowing speed-ups near to 4x compared to the HTTP-01 challenge. In
order to use the new challenge, we specified the following CLI flags:

-   Go-LEGO (client):

    -   `genCSRatLoadTest`: this boolean flag enables the generation of
        a new Certificate-Signing Request (CSR) for each client thread
        during a load test, if true.

    -   `newchallenge`: this required boolean flag tells LEGO to execute
        the issuance by means of the new challenge, if true.

-   Go-Pebble (server):

    -   `perMessageTimingCSVPath`: Path to the CSV file where the timing
        metrics for each message are written to.

    -   `memoryCSVPath`: Path to the CSV file where the memory
        measurements are written to.

    -   `newchallenge`: By setting this flag to true, the ACME pq-order/
        endpoint will be activated.

    -   `pqorderroot`: a Root CA and Issuer CA will be created for
        pq-order/ so this flag specifies which PQC algorithm will be
        used (root CA).

    -   `pqorderissuer`: this flag specifies which PQC algorithm will be
        used for the issuer CA.

The '/pq-order' is the entry point in Pebble for the new challenge
proposal. Clients can request to that challenge, given that they have an
account created and possess a previously-issued certificate (and
corresponding key). The code below shows how Pebble implementation
starts the new challenge. Note the TLS mutual authentication requirement
in the configuration (`tlsCfg` variable). Without it, Pebble denies
issuance. Below, an excerpt of relevant parts of the code.

``` {tabsize="4"}
if *newchallenge{
    //sets defaults if config not present
    if c.Pebble.PQOrderTLSPortAddress == 0 {
        c.Pebble.PQOrderTLSPortAddress = 10001
    }
    //(...)
    // In our tests we always suppose that it will have at most two chains: one classical 
    // and one post-quantum.
    // So if it is needed one of them, so they should be the first one in their respective 
    // field (classicChains and PQChains)
    if pqChain[0] == "" {
        rootCertX509 := caImpl.GetRootCert(0).Cert
        caCertPool.AddCert(rootCertX509)
    } else {
        rootCertX509 := caImpl.PQChains[0].Root.Cert.Cert
        caCertPool.AddCert(rootCertX509)
    }

    tlsCfg := &tls.Config {
        PQTLSEnabled: true,			
        InsecureSkipVerify: false,
        ClientAuth: tls.RequireAndVerifyClientCert, //mandatory Client Auth		
        ClientCAs:	caCertPool,
    }
    
    if *pqOrderRoot == "" || *pqOrderIssuer == ""{
        panic("If new challenge you need --pqOrderRoot and --pqOrderIssuer algorithms")
    }
    pqOrderChain := []string{*pqOrderRoot, *pqOrderIssuer, *pqOrderIssuer}

    // PQCACME Modification: Creates a second CA chain (PQC one for the new challenge)
    caImpl.AddPQChain(chainLength, *dirToSaveRoot, pqOrderChain, *hybrid)

    // Grabs WFE instance
    grabbedWFE := &nc.NewChallengeWFE{WebFrontEndImpl: wfeImpl}

    // Starts pq-order endpoint in a different TLS config (requires client auth).
    go func() {
        http.HandleFunc(string(wfe.NewChallengePath), grabbedWFE.HandlePQOrder)
        err := http.ListenAndServeTLSWithConfig(
            c.Pebble.PQOrderListenAddress,				
            c.Pebble.Certificate,
            c.Pebble.PrivateKey,
            nil, //default handler
            tlsCfg,
        )
        cmd.FailOnError(err, "Calling ListenAndServeTLS() for /pq-order")
        }()
    logger.Printf("ACME %s%s endpoint is activated",c.Pebble.PQOrderListenAddress,
            string(wfe.NewChallengePath))
    
}
```


In Pebble, `/pq-order` handling function and the main implementation of
the new challenge is placed in a new package. The variables,
structs, and functions are:

-   `func storePQOrder()`: store a PQOrder in the DB. Follows `wfe.go`
    and `memory` `store.go` original functions.

-   `func HandlePQOrder()`: Entry point to our new challenge work
    properly.



The core function of the new challenge is `HandlePQOrder()`. Below, we
list the relevant code of this function.

``` {tabsize="4"}
//func HandlePQOrder(...) {
    //0. First thing is to replace the Classic CA (Root and Intermediates) by PQC ones	
	w.Ca.PQCACME = true

	//1. Parse request	
    //(...)
	postData, prob := grabbedWFE.VerifyPOST(req, grabbedWFE.LookupJWK)
	if prob != nil { //(Error...) }

	//2. There is no order (yet), so go straight parsing and processing CSR 	
	//(...)
	csrBytes, err := base64.RawURLEncoding.DecodeString(finalizeMessage.CSR)
	if err != nil { //(Error...) }
	parsedCSR, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil { //(Error...) }
    //(...)	
	
	//3. check if the TLS-layer client certificate domain name 
	// matches the domain asked in the CSR	
	if ! (req.TLS.PeerCertificates[0].Subject.CommonName == csrDNSs[0]) {
		//(Error...)
	}
    //(...)
	//4. Store a new order directly
	existingOrder := storePQOrder(rw,orderID,csrDNSs,csrIPs,existingAcct.ID,parsedCSR)
	if existingOrder == nil { //(Error...) }

	//5. Issue the certificate (CompleteOrder(existingOrder))
	existingOrder.Status = acme.StatusValid
	issuePQCert(existingOrder)
	
	//6. and create URL for download.
	orderReq := grabbedWFE.OrderForDisplay(existingOrder, req)
	err = grabbedWFE.WriteJSONResponse(rw, http.StatusOK, orderReq)
//}
```

Lastly, changing the CAImpl struct in the ca.go file was also required.
The name of certificate chains are ClassicChains, and we added two
different fields: PQACME and PQChains. PQACME is a boolean flag that
uses the post-quantum issuance chain when it is set by true; PQChains
saves all post-quantum chains. Both were added since the new challenge
requires changing between classical to post-quantum issuance chains
without losing any of them.

    type CAImpl struct {
       log              *log.Logger
       db               *db.MemoryStore
       ocspResponderURL string


       // ClassicChains contains all classical chains.
       ClassicChains []*chain
       // PQCACME is set to true when a post-quantum chain is used
       PQCACME bool
       // PQChains contains all post-quantum chains.
       PQChains []*chain
       certValidityPeriod uint
    }

Secondly, we created `pq_ca.go` to place every function related to post-quantum scenarios. 
For instance, a new function `AddPQChain()` adds a post-quantum issuance 
chain on `PQChains` and it enables generating a new certificate to an 
authenticated ACME client. The `hybrid` flag tells if a concatenation of public keys is needed.
We concatenate a classical and a PQC one. Further details of the PQC integration is described
in the paper (and also implemented accross the repositories, such as `go-std` and `go-jose`).

    // PQCACME Modification: Adds `pqChain` as `ca.PQChains`
    // It also suits for newchallenge when new post-quantum chain is required on the demand
    func (ca *CAImpl) AddPQChain(chainLength int, dirToSaveRoot string, pqChain []string, 
    hybrid bool) {
    	ca.PQCACME = true
    	//(...)
    	ca.PQChains = make([]*chain, 1)
    	for i := 0; i < len(ca.PQChains); i++ {
    		ca.PQChains[i] = ca.newPqChain(intermediateKey, intermediateSubject, 
                subjectKeyID, chainLength, dirToSaveRoot, pqChain, hybrid)
    	}
    }

Note that PQACME field is set to true in order to force the instance of
`CAImpl` using the new post-quantum issuance chain in the certificate's
generation. There is a POST example in the paper, showing how the client
can interact with the PQ-Transition Challenge. 
