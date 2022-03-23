package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zhongshuwen/zswchain-go/ecc"
)

type SimpleConverter func(inputData []byte, fromFormat string, toFormat string) ([]byte, int, error)

var SIMPLE_CONVERT_FILE_REGISTRY map[string]map[string]SimpleConverter = map[string]map[string]SimpleConverter{
	"pubpem": {
		"pubzswkey": func(inputData []byte, fromFormat string, toFormat string) ([]byte, int, error) {
			data, err := ecc.SM2PemToZSWPublicKeyString(inputData)

			return []byte(data), 1, err
		},
	},
}

func ConvertFile(inputData []byte, fromFormat string, toFormat string, inputFilePath string, outputFilePath string) ([]byte, error) {
	if fromMap, ok := SIMPLE_CONVERT_FILE_REGISTRY[fromFormat]; ok {
		if convertor, nxt := fromMap[toFormat]; nxt {
			if inputFilePath != "" {
				inputFileData, err := os.ReadFile(inputFilePath)
				if err != nil {
					return nil, err
				}
				inputData = inputFileData
			}
			outputData, _, err := convertor(inputData, fromFormat, toFormat)
			if err != nil {
				return nil, err
			}
			if outputFilePath != "" {
				f, err := os.Create(outputFilePath)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
				_, err2 := f.Write(outputData)
				if err2 != nil {
					return nil, err2
				}
			} else {
				fmt.Print(string(outputData))
				return outputData, nil
			}

		}
	}
	return nil, fmt.Errorf("there is currently no converter defined to convert between %s and %s", fromFormat, toFormat)
}
