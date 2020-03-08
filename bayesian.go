package classify

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	bayesian "github.com/jbrukh/bayesian"
)

const (
	trueClass  bayesian.Class = "True"
	falseClass bayesian.Class = "False"

	tagSkipWord = "-"

	scoreTruePositive = "1"
	scoreTrueNegative = "-1"

	errNoLibrary = "No Library being used"

	pathLibraryCsv = "lib.csv"
)

var regexSanitize *regexp.Regexp

type Bayesian struct {
	Classifier             *bayesian.Classifier
	ModelID                string
	TrainingData           map[string][]string
	TrainingSrc            Source
	SkipSanitizeOnGenerate bool

	LibrarySrc  Source
	LibraryData map[string]string
}

type Source struct {
	DB       *sql.DB
	Filepath string
}

func (b *Bayesian) init() error {
	var err error

	err = b.loadTrainingSrc()
	if err != nil {
		return err
	}

	err = b.loadLibrarySrc()
	if err != nil {
		return err
	}

	b.initRegex()
	b.generateTrainingData()
	return nil
}

func (b *Bayesian) Classify(input string) Classification {
	var clearInput []string
	var likely int
	var result Classification
	var scores []float64

	clearInput = b.sanitizeInput([]string{input})
	if len(clearInput) <= 0 {
		return result
	}

	scores, likely, _ = b.Classifier.ProbScores(clearInput)
	return Classification{
		Likely: likely == 0,
		Score:  scores[0],
	}
}

func (b *Bayesian) loadTrainingSrc() error {
	var err error
	var truePositive, trueNegative []string

	dataTrainRaw := make(map[string]string)
	dataTrain := make(map[string][]string)

	if b.TrainingSrc.Filepath != "" {
		dataTrainRaw, err = readFileToMapString(b.TrainingSrc.Filepath)
		if err != nil {
			return err
		}

		for k, v := range dataTrainRaw {
			switch v {
			case scoreTruePositive:
				truePositive = append(truePositive, k)
			case scoreTrueNegative:
				trueNegative = append(trueNegative, k)
			}
		}
	}

	dataTrain["true"] = truePositive
	dataTrain["false"] = trueNegative

	logger(fmt.Sprintf("Load True Positive: %d, True Negative: %d", len(truePositive), len(trueNegative)))
	b.TrainingData = dataTrain

	return err
}

func (b *Bayesian) loadLibrarySrc() error {
	var err error
	dataLib := make(map[string]string)

	b.LibrarySrc.Filepath = pathLibraryCsv
	if b.LibrarySrc.Filepath != "" {
		dataLib, err = readFileToMapString(b.LibrarySrc.Filepath)
		if err != nil {
			return err
		}
	} else {
		logger(errNoLibrary)
		return nil
	}

	b.LibraryData = dataLib
	logger(fmt.Sprintf("Load Library: %d", len(dataLib)))

	return err
}

func (b *Bayesian) generateTrainingData() {
	dataTrue := b.TrainingData["true"]
	dataFalse := b.TrainingData["false"]

	if !b.SkipSanitizeOnGenerate {
		dataTrue = b.sanitizeInput(dataTrue)
		dataFalse = b.sanitizeInput(dataFalse)
	}

	b.Classifier = bayesian.NewClassifierTfIdf(trueClass, falseClass) // Create a classifier with TF-IDF support.
	b.Classifier.Learn(dataTrue, trueClass)
	b.Classifier.Learn(dataFalse, falseClass)
	b.Classifier.ConvertTermsFreqToTfIdf() // required
}

func (b *Bayesian) sanitizeInput(input []string) []string {
	var result []string
	var res string

	for _, txt := range input {
		res = strings.Trim(regexSanitize.ReplaceAllString(strings.ToLower(txt), " "), " ")
		for _, w := range strings.Split(res, " ") {
			if b.LibraryData[w] == tagSkipWord || len(w) < 2 {
				continue
			}
			if b.LibraryData[w] != "" {
				result = append(result, b.LibraryData[w])
			} else {
				result = append(result, w)
			}
		}
	}
	return result
}

func (b *Bayesian) initRegex() {
	regexSanitize = regexp.MustCompile("[^a-z]+") // only allow character
}
