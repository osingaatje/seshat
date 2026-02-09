package main

import (
	"github.com/osingaatje/seshat/driver"
	"github.com/osingaatje/seshat/ui/cli"
	"os"
)

func main() {
	os.Exit(driver.Bootstrap(cli.RunCli))
}
