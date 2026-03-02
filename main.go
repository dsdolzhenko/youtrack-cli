package main

import (
	"github.com/dsdolzhenko/youtrack-cli/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version)
}
