package main

import (
	"github.com/osingaatje/seshat/src/driver"
	"github.com/osingaatje/seshat/src/ui/cli"
	"os"
)

func main() {
	os.Exit(driver.Bootstrap(cli.RunCli))
}
