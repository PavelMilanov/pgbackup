package db

import (
	"fmt"
	"testing"
)

var token = Token{}

func TestCreate(t *testing.T) {
	token.Generate()
	fmt.Println(token.Hash)
}

func TestValidate(t *testing.T) {
	token.Generate()
	token.Validate()
}
