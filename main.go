//go:generate go run generator/main.go

package main

import (
	"os"

	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/rancher/machine-controller/controller/machine"
	machineDriver "github.com/rancher/machine-controller/controller/machine_driver"
)

var (
	GITCOMMIT = "HEAD"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "Kube config for accessing kubernetes cluster",
			EnvVar: "KUBECONFIG",
		},
		cli.BoolFlag{
			Name: "debug",
			Usage: "Enable debug log",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return run(c.String("config"))
	}

	app.ExitErrHandler = func(c *cli.Context, err error) {
		logrus.Fatal(err)
	}

	app.Run(os.Args)
}

func run(kubeConfigFile string) error {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
	if err != nil {
		return err
	}

	management, err := config.NewManagementContext(*kubeConfig)
	if err != nil {
		return err
	}

	machine.Register(management)
	machineDriver.Register(management)

	return management.StartAndWait()
}
