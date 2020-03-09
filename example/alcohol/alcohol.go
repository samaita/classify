package main

import (
	"log"

	"github.com/samaita/classify"
)

func main() {

	dataTrain, err := classify.ReadFileToMapString("train.csv")
	if err != nil {
		log.Fatal(err)
	}

	dataLibrary, err := classify.ReadFileToMapString("lib.csv")
	if err != nil {
		log.Fatal(err)
	}

	client := classify.Client{
		ClientID:     "example",
		Method:       classify.MethodNaiveBayes,
		TrainingData: dataTrain,
		LibraryData:  dataLibrary,
	}

	err = client.Init()
	if err != nil {
		log.Fatal(err)
	}

	usecase := make(map[string]bool)
	usecase["prost beer 100ml"] = true
	usecase["jaket hoodie sweater gambar jackdaniel's distro keren"] = false
	usecase["kaos heineken motif wash ombre"] = false

	for text, expected := range usecase {
		res := client.Classify(text)
		log.Println(text, res, "want:", expected)
	}
}
