//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"fmt"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/spf13/cobra"
)

type localUserOptions struct {
	userID   string
	userName string
	password string
	offset   int
	limit    int
}

func init() {
	rootCmd.AddCommand(localUserListCmd())
}

//
//
func localUserListCmd() *cobra.Command {
	options := localUserOptions{}

	cmd := &cobra.Command{
		Use:   "local-users",
		Short: "List and manage local users",
		Long:  `List and manage local users`,
		Example: `
	privx-cli local-users [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return localUserList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.userName, "name", "", "username of the user")
	flags.StringVar(&options.userID, "id", "", "ID of the user")

	cmd.AddCommand(localUserShowCmd())
	cmd.AddCommand(localUserCreateCmd())
	cmd.AddCommand(localUserUpdateCmd())
	cmd.AddCommand(localUserDeleteCmd())
	cmd.AddCommand(localUserUpdatePasswordCmd())

	return cmd
}

func localUserList(options localUserOptions) error {
	api := userstore.New(curl())

	users, err := api.LocalUsers(options.offset, options.limit,
		options.userID, options.userName)
	if err != nil {
		return err
	}

	return stdout(users)
}

//
//
func localUserShowCmd() *cobra.Command {
	options := localUserOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get local user by ID",
		Long:  `Get local user by ID. User ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli local-users show [access flags] --id <USER-ID>,<USER-ID>
		`,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return localUserShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "user ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func localUserShow(options localUserOptions) error {
	api := userstore.New(curl())
	users := []userstore.LocalUser{}

	for _, id := range strings.Split(options.userID, ",") {
		user, err := api.LocalUser(id)
		if err != nil {
			return err
		}
		users = append(users, *user)
	}

	return stdout(users)
}

//
//
func localUserCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new local user",
		Long:  `Create new local user`,
		Example: `
	privx-cli local-users create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return localUserCreate(args)
		},
	}

	return cmd
}

func localUserCreate(args []string) error {
	var newUser userstore.LocalUser
	api := userstore.New(curl())

	err := decodeJSON(args[0], &newUser)
	if err != nil {
		return err
	}

	uid, err := api.CreateLocalUser(newUser)
	if err != nil {
		return err
	}

	return stdout(uid)
}

//
//
func localUserUpdateCmd() *cobra.Command {
	options := localUserOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update local user",
		Long:  `Update local user`,
		Example: `
	privx-cli local-users update [access flags] --id <USER-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return localUserUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "unique user ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func localUserUpdate(options localUserOptions, args []string) error {
	var updateUser userstore.LocalUser
	api := userstore.New(curl())

	err := decodeJSON(args[0], &updateUser)
	if err != nil {
		return err
	}

	err = api.UpdateLocalUser(options.userID, &updateUser)
	if err != nil {
		return err
	}

	return nil
}

//
//
func localUserDeleteCmd() *cobra.Command {
	options := localUserOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete local user",
		Long:  `Delete local user. User ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli local-users delete [access flags] --id <USER-ID>,<USER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return localUserDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "unique user ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func localUserDelete(options localUserOptions) error {
	api := userstore.New(curl())

	for _, id := range strings.Split(options.userID, ",") {
		err := api.DeleteLocalUser(id)
		if err != nil {
			return err
		} else {
			fmt.Println(id)
		}
	}

	return nil
}

//
//
func localUserUpdatePasswordCmd() *cobra.Command {
	options := localUserOptions{}

	cmd := &cobra.Command{
		Use:   "update-password",
		Short: "Update local user password",
		Long:  `Update local user password`,
		Example: `
	privx-cli local-users update-password [access flags] --id <USER-ID> --password <NEW-PASSWORD>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return localUserUpdatePassword(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "unique user id")
	flags.StringVar(&options.password, "password", "", "new password for local user")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("password")

	return cmd
}

func localUserUpdatePassword(options localUserOptions) error {
	newPassword := userstore.Password{
		Password: options.password,
	}
	api := userstore.New(curl())

	err := api.UpdateLocalUserPassword(options.userID, &newPassword)
	if err != nil {
		return err
	}

	return nil
}
