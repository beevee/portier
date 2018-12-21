// Copyright © 2018 Alexey Kirpichnikov <alex.kirp@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"github.com/beevee/portier/yandex"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var role string

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users (enable|disable)",
	Short: "Batch process corporate users",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		if err := cobra.OnlyValidArgs(cmd, args); err != nil {
			return err
		}
		return nil
	},
	ValidArgs: []string{"enable", "disable"},
	Run: func(cmd *cobra.Command, args []string) {
		api := &yandex.API{}
		cmd.Print("initializing Yandex API\n")
		err := api.Init(
			viper.GetString("session_id"),
			viper.GetString("client_id"))
		if err != nil {
			cmd.Printf("failed to init Yandex API: %s\n", err)
			return
		}

		users, err := api.GetUsersByRole(role)
		if err != nil {
			cmd.Printf("failed to fetch users with role %s: %s\n", role, err)
			return
		}
		cmd.Printf("fetched %d users with role %s\n", len(users), role)

		for _, user := range users {
			switch args[0] {
			case "enable":
				cmd.Printf("enabling user %s (%s)\n", user.ID, user.Name)
				if !user.IsActive {
					if err := api.EnableUser(user); err != nil {
						cmd.Printf("failed to enable user %s (%s): %s\n", user.ID, user.Name, err)
					} else {
						cmd.Printf("successfully enabled user %s (%s)\n", user.ID, user.Name)
					}
				} else {
					cmd.Printf("user %s (%s) is already active\n", user.ID, user.Name)
				}
			case "disable":
				cmd.Printf("disabling user %s (%s)\n", user.ID, user.Name)
				if user.IsActive {
					if err := api.DisableUser(user); err != nil {
						cmd.Printf("failed to disable user %s (%s): %s\n", user.ID, user.Name, err)
					} else {
						cmd.Printf("successfully disabled user %s (%s)\n", user.ID, user.Name)
					}
				} else {
					cmd.Printf("user %s (%s) is already inactive\n", user.ID, user.Name)
				}
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)

	usersCmd.PersistentFlags().StringVar(&role, "role", "Кирпичников", "operate on users in this role only")
}
