package modules

// GOPRIVATE="git.ace/*"

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	icarus "git.ace/icarus/icarusclient/v5"
)

func GenerateAuthQuery() icarus.QueryPackage {
	/*
		name file creds.txt

		format:
		IP_ADDR PORT
		USERNAME PASSWORD
	*/

	filePath := "../resources/creds.txt"
	readFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	fileScanner.Scan()
	line := strings.SplitN(fileScanner.Text(), " ", 2)                               // split into IP addr and port
	query := icarus.NewQuery(strings.Trim(line[0], " "), strings.Trim(line[1], " ")) // query local icarus server

	fileScanner.Scan()
	line = strings.SplitN(fileScanner.Text(), " ", 2)                                      // split into username and password
	resp, ok := query.Authenticate(strings.Trim(line[0], " "), strings.Trim(line[1], " ")) // authenticate into Icarus
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
	}

	readFile.Close()

	query.ClearQueries()
	return query
}
