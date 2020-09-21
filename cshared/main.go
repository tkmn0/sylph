package main

import (
	"C"

	"github.com/tkmn0/sylph"
)

var client *sylph.Client

//export Initialize
func Initialize() {
	client = sylph.NewClient()
}

func Connect(address string, port uint) {

}

func main() {}
