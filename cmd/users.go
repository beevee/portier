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
