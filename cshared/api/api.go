package api

import (
	"github.com/tkmn0/sylph"
)

var clients map[uintptr]*sylph.Client
var servers map[uintptr]*sylph.Server
