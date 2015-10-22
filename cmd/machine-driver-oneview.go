package main

import (
	"github.com/HewlettPackard/docker-machine-oneview/oneview"
	"github.com/docker/machine/libmachine/drivers/plugin"
)

func main() {
	plugin.RegisterDriver(oneview.NewDriver("", ""))
}
