package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // 👈 load .env file
	if err != nil {
		log.Fatal(err)
	}

	env := os.Getenv("KEY_")

	fmt.Printf("Some KEY_ %s", env)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from input: %v\n", err)
	}
}

func checkDomain(domain string) {

	var hasMX, hasDMARC, hasSPF bool
	var spfRecord, dmarcRecord string

	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		printErr(err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		printErr(err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)

	if err != nil {
		printErr(err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
		}
	}

	fmt.Printf("%v, %v, %v, %v, %v, %v", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}

func printErr(err error) {
	log.Printf("Error: %v\n", err)
}
