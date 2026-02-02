package main

import (
	"os"
	"seshat/driver"
	"seshat/ui/cli"
)

func main() {
	os.Exit(driver.Bootstrap(cli.RunCli))
}
