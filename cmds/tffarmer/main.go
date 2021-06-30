package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"os"

	"github.com/threefoldtech/tfexplorer"
	"github.com/threefoldtech/tfexplorer/client"
	"github.com/urfave/cli"
)

var (
	db     client.Directory
	userid = &tfexplorer.UserIdentity{}
)

func main() {
	app := cli.NewApp()
	app.Usage = "Create and manage a Threefold farm"
	app.Version = "0.0.1"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "enable debug logging",
		},
		cli.StringFlag{
			Name:  "seed",
			Usage: "seed filename",
			Value: "user.seed",
		},
		cli.StringFlag{
			Name:   "explorer, e",
			Usage:  "URL of the Explorer",
			Value:  "https://explorer.grid.tf/api/v1",
			EnvVar: "BCDB_URL",
		},
	}

	app.Before = func(c *cli.Context) error {
		var err error
		debug := c.Bool("debug")
		if !debug {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		url := c.String("explorer")
		// check seed file in ~/.config/tffarmer/user.seed
		err = userid.Load(c.String("seed"))
		if err != nil {
			defaultSeedLocation, err := getSeedPath()
			if err != nil {
				return err
			}
			err = userid.Load(defaultSeedLocation)
			if err != nil {
				// Seed not found in default location ~/.config/tffarmer.seed
				generateSeedOption := "n"
				fmt.Print("Seed not provided. Do you want to generate a new seed? [y/n] ")
				fmt.Scanln(&generateSeedOption)
				if generateSeedOption == "n" {
					return err
				}
				userid, err = generateNewUser(c, url, defaultSeedLocation)
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("User seed found at %s and will be used\n", defaultSeedLocation)
			}

		}

		cl, err := client.NewClient(url, userid)
		if err != nil {
			return errors.Wrap(err, "failed to create client to explorer")
		}

		db = cl.Directory

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "farm",
			Usage: "Manage and create farms",
			Subcommands: []cli.Command{
				{
					Name:      "register",
					Usage:     "register a new farm",
					Category:  "identity",
					ArgsUsage: "farm_name",
					Flags: []cli.Flag{
						cli.StringSliceFlag{
							Name:     "addresses",
							Usage:    "wallet address",
							Required: true,
						},
						cli.StringSliceFlag{
							Name:     "email",
							Usage:    "email address of the farmer. It is used to send communication to the farmer and for the minting",
							Required: true,
						},
						cli.StringSliceFlag{
							Name:     "iyo_organization",
							Usage:    "the It'sYouOnline organization used by your farm in v1",
							Required: false,
						},
					},
					Action: registerFarm,
				},
				{
					Name:     "update",
					Usage:    "update an existing farm",
					Category: "identity",
					Flags: []cli.Flag{
						cli.Int64Flag{
							Name:     "id",
							Usage:    "farm ID",
							Required: true,
						},
						cli.StringSliceFlag{
							Name:     "addresses",
							Usage:    "wallet address. the format is 'asset:address: e.g: 'TFT:GBUPOYJ7I4D4TYSFXPJNLSATHCCF2QDDQCIIIXBG7CV7S2U36UMAQENV'",
							Required: false,
						},
						cli.StringSliceFlag{
							Name:     "email",
							Usage:    "email address of the farmer. It is used to send communication to the farmer and for the minting",
							Required: false,
						},
						cli.StringSliceFlag{
							Name:     "iyo_organization",
							Usage:    "the It'sYouOnline organization used by your farm in v1",
							Required: false,
						},
					},
					Action: updateFarm,
				},
				{
					Name:     "addip",
					Usage:    "add Ip addresses to farm",
					Category: "identity",
					Flags: []cli.Flag{
						cli.Int64Flag{
							Name:     "id",
							Usage:    "farm ID",
							Required: true,
						},
						cli.StringFlag{
							Name:     "address",
							Usage:    "IP address",
							Required: true,
						},
						cli.StringSliceFlag{
							Name:     "gateway",
							Usage:    "gateway address",
							Required: true,
						},
					},
					Action: addIP,
				},
				{
					Name:     "deleteip",
					Usage:    "delete Ip addresses from farm",
					Category: "identity",
					Flags: []cli.Flag{
						cli.Int64Flag{
							Name:     "id",
							Usage:    "farm ID",
							Required: true,
						},
						cli.StringFlag{
							Name:     "address",
							Usage:    "IP address",
							Required: true,
						},
					},
					Action: deleteIP,
				},
				{
					Name:     "list",
					Usage:    "list farms",
					Category: "identity",
					Flags:    []cli.Flag{},
					Action:   listFarms,
				},
			},
		},
		{
			Name:  "network",
			Usage: "Manage network of a farm and hand out allocation to the grid",
			Subcommands: []cli.Command{
				{
					Name:     "configure-public",
					Category: "network",
					Usage: `configure the public interface of a node.
You can specify multime time the ip and gw flag to configure multiple IP on the public interface`,
					ArgsUsage: "node ID",
					Flags: []cli.Flag{
						cli.StringSliceFlag{
							Name:  "ip",
							Usage: "ip address to set to the exit interface",
						},
						cli.StringSliceFlag{
							Name:  "gw",
							Usage: "gw address to set to the exit interface",
						},
						cli.StringFlag{
							Name:  "iface",
							Usage: "name of the interface to use as public interface",
						},
					},
					Action: configPublic,
				},
			},
		},
		{
			Name:  "nodes",
			Usage: "Manage nodes from a farm",
			Subcommands: []cli.Command{
				{
					Name:     "free",
					Category: "nodes",
					Usage:    "mark some nodes as free to use",
					Flags: []cli.Flag{
						cli.StringSliceFlag{
							Name:  "nodes",
							Usage: "node IDs. can be specified multiple time",
						},
						cli.BoolFlag{
							Name:  "free",
							Usage: "if set, the node is marked free, it not the node is mark not free",
						},
					},
					Action: markFree,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
