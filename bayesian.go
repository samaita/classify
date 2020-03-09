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
	LibraryData            map[string]string
	TrainingData           map[string]string
	SkipSanitizeOnGenerate bool
}

type Source struct {
	DB       *sql.DB
	Filepath string
}

func (b *Bayesian) init() error {
	var err error

	b.initRegex()

	err = b.trainData()
	if err != nil {
		return err
	}

	logger(fmt.Sprintf("Load Library: %d", len(b.LibraryData)))
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

func (b *Bayesian) trainData() error {
	var err error
	var truePositive, trueNegative []string

	if len(b.TrainingData) > 0 {
		for k, v := range b.TrainingData {
			switch v {
			case scoreTruePositive:
				truePositive = append(truePositive, k)
			case scoreTrueNegative:
				trueNegative = append(trueNegative, k)
			}
		}
	}

	b.Classifier = bayesian.NewClassifierTfIdf(trueClass, falseClass) // Create a classifier with TF-IDF support.
	b.Classifier.Learn(b.sanitizeInput(truePositive), trueClass)
	b.Classifier.Learn(b.sanitizeInput(trueNegative), falseClass)
	b.Classifier.ConvertTermsFreqToTfIdf() // required

	logger(fmt.Sprintf("Train Data - True Positive: %d, True Negative: %d", len(truePositive), len(trueNegative)))
	return err
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
