package main

import (
    "encoding/json"
    "os"
    "fmt"
)

type Configuration struct {
    InstagramClientId string
    HashTags []string
}

func ReadConfiguration(filename string) Configuration {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
	  log.Fatal(fmt.Sprintf("Error parsing configuration: %v", err))
	}
	return configuration
}
