package main

import (
	"encoding/csv"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	inputFile := flag.String("input", "", "The input file.")

	flag.Parse()

	err := read(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
}

func read(inputFilePath string) error {
	file, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}

	//noinspection GoUnhandledErrorResult
	defer file.Close()

	var pattern strings.Builder
	uniqueClassNames := map[string]bool{}
	// ignore
	uniqueClassNames["com.intellij.TestAll"] = true

	reader := csv.NewReader(file)
	classCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if classCount == 0 && record[0] != "1" {
			continue
		}

		methodName := strings.TrimPrefix(record[1], "com.intellij.tests.BootstrapTests: ")
		className := methodName[0:strings.LastIndexByte(methodName, '.')]

		if uniqueClassNames[className] {
			continue
		}

		classCount++
		if classCount > 1 {
			pattern.WriteString("||")
		}

		uniqueClassNames[className] = true
		pattern.WriteString(className)
	}

	err = ioutil.WriteFile(inputFilePath+".pattern.txt", []byte(pattern.String()), 0644)
	if err != nil {
		return err
	}
	return nil
}
