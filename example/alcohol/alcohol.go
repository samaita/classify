package main

import (
	"log"

	"github.com/samaita/classify"
)

func main() {
	client := classify.Client{
		ClientID: "example",
		Method:   classify.NaiveBayes,
	}
	client.TrainingSrc.Filepath = "train.csv"
	err := client.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(client.Classify("prost beer 100ml"))
}
