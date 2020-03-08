package classify

import (
	"fmt"
	"time"
)

const (
	NaiveBayes = "naive-bayes" // method naive bayes

	errClassify                           = "[classify]"
	errClientInitClientIDorMethodRequired = "Client ID and method required"
	errClientInvalidMethod                = "Invalid method"
)

type Client struct {
	ClientID    string
	Method      string
	Model       Classifier
	TrainingSrc Source
}

type Classifier interface {
	Classify(input string) Classification
}

type Classification struct {
	Likely bool
	Score  float64
}

func (c *Client) Init() error {
	t := time.Now()

	if c.ClientID == "" || c.Method == "" {
		return errMsg(errClientInitClientIDorMethodRequired)
	}

	switch c.Method {
	case NaiveBayes:
		model, err := c.initNaiveBayes()
		if err != nil {
			return errMsg(err.Error())
		}
		c.Model = &model
	default:
		return errMsg(errClientInvalidMethod)
	}

	logger(fmt.Sprintf("Time Elapsed on init: %v", time.Since(t)))
	return nil
}

func (c *Client) initNaiveBayes() (Bayesian, error) {
	var err error

	b := Bayesian{
		ModelID:     generateID(),
		TrainingSrc: c.TrainingSrc,
	}

	err = b.init()
	return b, err
}

func (c *Client) Classify(input string) Classification {
	switch c.Method {
	case NaiveBayes:
		return c.Model.Classify(input)
	default:
		return Classification{
			Likely: false,
			Score:  -999,
		}
	}
}
