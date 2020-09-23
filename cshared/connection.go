package main

import (
	"github.com/tkmn0/sylph"
)

type Connection interface {
	OnTransport(func(transport sylph.Transport))
}
