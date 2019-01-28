package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type virtualmachine struct {
	Action    string `json:"action"`
	VMName    string `json:"vmname"`
	IPAddress string `json:"ipaddress"`
}

func main() {
	path, format := parseFlags()
	virtualmachines := collectvm()

	var output io.Writer

	if path != "" {
		f, err := os.Create(path)
		handleError(err)
		defer f.Close()
		output = f

	} else {
		output = os.Stdout
	}

	if format == "json" {
		date, err := json.MarshalIndent(virtualmachines, "", " ")
		handleError(err)
		output.Write(date)

	} else if format == "csv" {
		output.Write([]byte("Action,VMName,IPAddress\n"))
		writer := csv.NewWriter(output)
		for _, virtualmachine := range virtualmachines {
			err := writer.Write([]string{virtualmachine.Action, virtualmachine.VMName, virtualmachine.IPAddress})
			handleError(err)
		}
		writer.Flush()

	}

}

func parseFlags() (path, format string) {
	flag.StringVar(&path, "path", "", "The path to export file.")
	flag.StringVar(&format, "format", "json", "The output format fo rthe virtual machines information. Available options are 'json' and 'csv'.")
	flag.Parse()

	format = strings.ToLower(format)
	if format != "csv" && format != "json" {
		fmt.Println("Error: Invalid format. Use 'json' or 'csv' instead.")
		flag.Usage()
		os.Exit(1)
	}
	return
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func collectvm() (virtualmachines []virtualmachine) {
	f, err := os.Open("vmlist.csv")
	handleError(err)
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ','
	lines, err := reader.ReadAll()
	handleError(err)

	for _, line := range lines {
		handleError(err)

		virtualmachine := virtualmachine{
			Action:    line[0],
			VMName:    line[1],
			IPAddress: line[2],
		}
		virtualmachines = append(virtualmachines, virtualmachine)

	}

	return
}
