package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
)

// tlsClientResultsInfo contains handshake information and measurements statistics from the client
type tlsClientResultsInfo struct {
	kexName                 string
	authName                string
	avgTotalTime            float64
	avgProcessServerHello   float64
	avgWriteClientHello     float64
	stdevTotalTime          float64
	stdevProcessServerHello float64
	stdevWriteClientHello   float64
}

// TLSClientResultsInfo contains handshake information and measurements statistics from the server
type tlsServerResultsInfo struct {
	kexName              string
	authName             string
	avgTotalTime         float64
	avgWriteCertVerify   float64
	stdevTotalTime       float64
	stdevWriteCertVerify float64
}

// tlsComputeStats calculates the average and standard deviation of the following client measurements: 
// Full TLS handshake time, time to process the ServerHello message, and time to write the ClientHello message
func tlsComputeStats(timingsFullProtocol []float64, timingsProcessServerHello []float64, timingsWriteClientHello []float64, hs int) (r tlsClientResultsInfo) {

	//counts
	var countTotalTime float64
	var countProcessServerHello float64
	var countWriteClientHello float64

	//Average
	countTotalTime, countProcessServerHello, countWriteClientHello = 0, 0, 0
	for i := 0; i < hs; i++ {
		countTotalTime += timingsFullProtocol[i]
		countProcessServerHello += timingsProcessServerHello[i]
		countWriteClientHello += timingsWriteClientHello[i]
	}

	r.avgTotalTime = (countTotalTime) / float64(hs)
	r.avgProcessServerHello = (countProcessServerHello) / float64(hs)
	r.avgWriteClientHello = (countWriteClientHello) / float64(hs)

	//stdev
	for i := 0; i < hs; i++ {
		r.stdevTotalTime += math.Pow(float64(timingsFullProtocol[i])-r.avgTotalTime, 2)
		r.stdevProcessServerHello += math.Pow(float64(timingsProcessServerHello[i])-r.avgProcessServerHello, 2)
		r.stdevWriteClientHello += math.Pow(float64(timingsWriteClientHello[i])-r.avgWriteClientHello, 2)
	}
	r.stdevTotalTime = math.Sqrt(r.stdevTotalTime / float64(hs))
	r.stdevProcessServerHello = math.Sqrt(r.stdevProcessServerHello / float64(hs))
	r.stdevWriteClientHello = math.Sqrt(r.stdevWriteClientHello / float64(hs))

	return r
}

// tlsPrintStatistics prints to stdout the client measurements' average and standard deviation.
func tlsPrintStatistics(results []tlsClientResultsInfo) {
	//header
	fmt.Printf("%-47s | ", "TestName")
	fmt.Printf("%-20s | ", "AvgClientTotalTime")
	fmt.Printf("%-20s | ", "StdevClientTotalTime")
	fmt.Printf("%-20s | ", "AvgWrtCHelloTime")
	fmt.Printf("%-20s | ", "StdevWrtCHelloTime")
	fmt.Printf("%-20s | ", "AvgPrSHelloTime")
	fmt.Printf("%-20s  ", "StdevPrSHelloTime")

	for _, r := range results {
		//content
		fmt.Println()
		fmt.Printf("%23s %23s |", r.kexName, r.authName)

		fmt.Printf(" %-20f |", r.avgTotalTime)
		fmt.Printf(" %-20f |", r.stdevTotalTime)
		fmt.Printf(" %-20f |", r.avgWriteClientHello)
		fmt.Printf(" %-20f |", r.stdevWriteClientHello)
		fmt.Printf(" %-20f |", r.avgProcessServerHello)
		fmt.Printf(" %-20f ", r.stdevProcessServerHello)
	}
}

// tlsInitCSVClient creates the CSV files which will hold the client's measurments data.
// Additionally it adds a header in this CSV files.
func tlsInitCSVClient() {
	csvFile, err := os.Create(*measurementsDir + "/tls-client.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)

	header := []string{"KEXAlgo", "authAlgo", "timingFullProtocol", "timingProcessServerHello", "timingWriteClientHello"}

	csvwriter.Write(header)
	csvwriter.Flush()
	csvFile.Close()

	csvFile, err = os.Create(*measurementsDir + "/tls-client-sizes.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter = csv.NewWriter(csvFile)

	header = []string{"KEXAlgo", "authAlgo", "ClientHello", "Certificate", "CertificateVerify", "Finished", "Total"}

	csvwriter.Write(header)
	csvwriter.Flush()
	csvFile.Close()
}

// tlsSaveCSVClient writes the client's measurements data to CSV files.
func tlsSaveCSVClient(timingsFullProtocol []float64, timingsProcessServerHello []float64, timingsWriteClientHello []float64, name, authName string, hs int, sizes map[string]uint32) {
	csvFile, err := os.OpenFile(*measurementsDir + "/tls-client.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for i := 0; i < hs; i++ {
		arrayStr := []string{name, authName, fmt.Sprintf("%f", timingsFullProtocol[i]),
			fmt.Sprintf("%f", timingsProcessServerHello[i]),
			fmt.Sprintf("%f", timingsWriteClientHello[i])}

		if err := csvwriter.Write(arrayStr); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		csvwriter.Flush()
	}
	csvFile.Close()

	csvFile, err = os.OpenFile(*measurementsDir + "/tls-client-sizes.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	csvwriter = csv.NewWriter(csvFile)

	totalSizes := sizes["ClientHello"] + sizes["Certificate"] + sizes["CertificateVerify"] + sizes["Finished"]

	arrayStr := []string{name, authName,
		fmt.Sprintf("%d", sizes["ClientHello"]),		
		fmt.Sprintf("%d", sizes["Certificate"]),
		fmt.Sprintf("%d", sizes["CertificateVerify"]),
		fmt.Sprintf("%d", sizes["Finished"]),
		fmt.Sprintf("%d", totalSizes),
	}

	if err := csvwriter.Write(arrayStr); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	csvwriter.Flush()
	csvFile.Close()
}

// tlsInitCSVServer creates the CSV files which will hold the server's measurments data.
// Additionally it adds a header in this CSV files.
func tlsInitCSVServer() {
	csvFile, err := os.Create(*measurementsDir + "/tls-server.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)

	header := []string{"KEXAlgo", "authAlgo", "timingFullProtocol", "timingWriteServerHello", "timingWriteCertVerify"}

	csvwriter.Write(header)
	csvwriter.Flush()
	csvFile.Close()

	
	csvFile, err = os.Create(*measurementsDir + "/tls-server-sizes.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter = csv.NewWriter(csvFile)

	header = []string{"KEXAlgo", "authAlgo", "ServerHello", "EncryptedExtensions", "Certificate", "CertificateRequest", "CertificateVerify", "Finished", "Total"}

	csvwriter.Write(header)
	csvwriter.Flush()
	csvFile.Close()
}

// tlsSaveCSVServer writes the server's measurements data to CSV files.
func tlsSaveCSVServer(timingsFullProtocol []float64, timingsWriteServerHello []float64, timingsWriteCertVerify []float64, name string, authName string, hs int, sizes map[string]uint32) {
	csvFile, err := os.OpenFile(*measurementsDir + "/tls-server.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for i := 0; i < hs; i++ {
		arrayStr := []string{name, authName, fmt.Sprintf("%f", timingsFullProtocol[i]),
			fmt.Sprintf("%f", timingsWriteServerHello[i]),
			fmt.Sprintf("%f", timingsWriteCertVerify[i]),
		}

		if err := csvwriter.Write(arrayStr); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		csvwriter.Flush()
	}
	csvFile.Close()

	csvFile, err = os.OpenFile(*measurementsDir + "/tls-server-sizes.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	csvwriter = csv.NewWriter(csvFile)

	totalSizes := sizes["ServerHello"] + sizes["EncryptedExtensions"] + sizes["Certificate"] + sizes["CertificateRequest"] + sizes["CertificateVerify"] + sizes["Finished"]

	arrayStr := []string{name, authName, 
		fmt.Sprintf("%d", sizes["ServerHello"]),
		fmt.Sprintf("%d", sizes["EncryptedExtensions"]),
		fmt.Sprintf("%d", sizes["Certificate"]),
		fmt.Sprintf("%d", sizes["CertificateRequest"]),
		fmt.Sprintf("%d", sizes["CertificateVerify"]),
		fmt.Sprintf("%d", sizes["Finished"]),
		fmt.Sprintf("%d", totalSizes),
	}

	if err := csvwriter.Write(arrayStr); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	csvwriter.Flush()
	csvFile.Close()

}
