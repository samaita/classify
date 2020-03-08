package classify

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

func generateID() string {
	return uuid.New().String()
}

func errMsg(err string) error {
	return fmt.Errorf("%v %v", errClassify, err)
}

func logger(txt string) {
	log.Printf("%v %v", errClassify, txt)
}
