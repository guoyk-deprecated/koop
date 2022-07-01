package main

import (
	"errors"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func exit(err *error) {
	if *err != nil {
		log.Println("exited with error:", (*err).Error())
		os.Exit(1)
	} else {
		log.Println("exited")
	}
}

func main() {
	var err error
	defer exit(&err)

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "no-update",
			Usage:   "no update if already existed",
			EnvVars: []string{"KOOP_NO_UPDATE"},
		},
		&cli.BoolFlag{
			Name:    "zero-replicas",
			Usage:   "create workloads with zero replicas",
			EnvVars: []string{"KOOP_ZERO_REPLICAS"},
		},
	}
	app.Usage = "file based kubernetes operation tool"
	app.Commands = append(app.Commands, &cli.Command{
		Name:        "pull",
		Description: "pull resources from existing cluster",
		Action: func(c *cli.Context) error {
			if c.NArg() != 4 {
				return errors.New("invalid number of arguments")
			}
			isNoUpdate = c.Bool("no-update")
			isZeroReplicas = c.Bool("zero-replicas")
			return commandPull(c.Context, c.Args().Get(0), c.Args().Get(1), c.Args().Get(2), c.Args().Get(3))
		},
	})
	app.Commands = append(app.Commands, &cli.Command{
		Name:        "push",
		Description: "push resources to existing cluster",
		Action: func(c *cli.Context) error {
			if c.NArg() != 4 {
				return errors.New("invalid number of arguments")
			}
			isNoUpdate = c.Bool("no-update")
			isZeroReplicas = c.Bool("zero-replicas")
			return commandPush(c.Context, c.Args().Get(0), c.Args().Get(1), c.Args().Get(2), c.Args().Get(3))
		},
	})
	err = app.Run(os.Args)
}
