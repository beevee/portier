package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beevee/portier/yandex"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "portier"
	app.HelpName = "portier"
	app.Usage = "provides convenience functions for corporate Yandex.Taxi accounts"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:    "users",
			Aliases: []string{"u"},
			Usage:   "enable or disable orders from application for user role",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "role",
					Usage: "operate on users in this role only",
					Value: "Кирпичников", // this is to prevent accidental runs of batch operations on all users
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					return cli.NewExitError("please specify: enable or disable orders", 1)
				}

				api := &yandex.API{}
				log.Print("initializing Yandex API\n")
				err := api.Init(c.Parent().String("sessionid"), c.Parent().String("clientid"))
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("failed to init Yandex API: %s\n", err), 1)
				}

				users, err := api.GetUsersByRole(c.String("role"))
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("failed to fetch users with role %s: %s\n", c.String("role"), err), 1)
				}
				log.Printf("fetched %d users with role %s\n", len(users), c.String("role"))

				for _, user := range users {
					switch c.Args().First() {
					case "enable":
						log.Printf("enabling user %s (%s)\n", user.ID, user.Name)
						if !user.IsActive {
							if err := api.EnableUser(user); err != nil {
								log.Printf("failed to enable user %s (%s): %s\n", user.ID, user.Name, err)
							} else {
								log.Printf("successfully enabled user %s (%s)\n", user.ID, user.Name)
							}
						} else {
							log.Printf("user %s (%s) is already active\n", user.ID, user.Name)
						}
					case "disable":
						log.Printf("disabling user %s (%s)\n", user.ID, user.Name)
						if user.IsActive {
							if err := api.DisableUser(user); err != nil {
								log.Printf("failed to disable user %s (%s): %s\n", user.ID, user.Name, err)
							} else {
								log.Printf("successfully disabled user %s (%s)\n", user.ID, user.Name)
							}
						} else {
							log.Printf("user %s (%s) is already inactive\n", user.ID, user.Name)
						}
					default:
						return cli.NewExitError(fmt.Sprintf("we can only enable and disable, and you said: %s\n", c.Args().First()), 1)
					}
				}

				return nil
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "sessionid, s",
			Usage:    "value of Yandex Session_id cookie",
			EnvVar:   "SESSION_ID",
			Required: true,
		},
		cli.StringFlag{
			Name:     "clientid, c",
			Usage:    "Yandex client id",
			EnvVar:   "CLIENT_ID",
			Required: true,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
