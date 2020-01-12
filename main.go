package main

import (
	"os"
	"strings"

	"github.com/imander/bootcampspot/cmd"
	"github.com/imander/bootcampspot/config"
)

func main() {
	args := os.Args
	prompt := len(args) > 1 && strings.HasPrefix(args[1], "conf")
	config.Load(!prompt)
	cmd.Execute()
}
