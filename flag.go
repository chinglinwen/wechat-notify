package main

import (
	"github.com/namsral/flag"
)

var (
	addr = flag.String("a", ":8001", "listening address")
)
