package controller

import (
	"github.com/rancher/machine-controller/controller/machine"
	"github.com/rancher/machine-controller/controller/machine_driver"
	"github.com/rancher/types/config"
)

func Register(management *config.ManagementContext) {
	machine.Register(management)
	machine_driver.Register(management)
}
