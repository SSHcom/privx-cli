//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/spf13/cobra"
)

type authorizedkeyOptions struct {
	userID  string
	keyID   string
	sortkey string
	sortdir string
	limit   int
	offset  int
}

func init() {
	rootCmd.AddCommand(authorizedkeyListCmd())
}

//
//
func authorizedkeyListCmd() *cobra.Command {
	options := authorizedkeyOptions{}

	cmd := &cobra.Command{
		Use:   "authorized-keys",
		Short: "List and manage authorized keys",
		Long:  `List and manage authorized keys`,
		Example: `
	privx-cli authorized-keys [access flags] --offset <OFFSET> --sortkey <SORTKEY>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return authorizedkeyList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort object by name, updated, or created.")

	cmd.AddCommand(authorizedkeyShowCmd())
	cmd.AddCommand(authorizedkeyCreateCmd())
	cmd.AddCommand(authorizedkeyUpdateCmd())
	cmd.AddCommand(authorizedkeyDeleteCmd())
	cmd.AddCommand(authorizedkeyResolveCmd())

	return cmd
}

func authorizedkeyList(options authorizedkeyOptions) error {
	api := rolestore.New(curl())

	keys, err := api.AllAuthorizedKeys(options.offset, options.limit,
		strings.ToUpper(options.sortdir), options.sortkey)
	if err != nil {
		return err
	}

	return stdout(keys)
}

//
//
func authorizedkeyShowCmd() *cobra.Command {
	options := authorizedkeyOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get user's authorized keys",
		Long:  `Get user's authorized keys`,
		Example: `
	privx-cli authorized-keys show [access flags] --user-id <USER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return authorizedkeyShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "user-id", "", "user ID")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

func authorizedkeyShow(options authorizedkeyOptions) error {
	api := rolestore.New(curl())

	key, err := api.AuthorizedKeys(options.userID)
	if err != nil {
		return err
	}

	return stdout(key)
}

//
//
func authorizedkeyCreateCmd() *cobra.Command {
	options := authorizedkeyOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new authorized key for user",
		Long:  `Create new authorized key for user`,
		Example: `
	privx-cli authorized-keys create [access flags] --user-id <USER-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return authorizedkeyCreate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "user-id", "", "user ID")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

func authorizedkeyCreate(options authorizedkeyOptions, args []string) error {
	var newKey rolestore.AuthorizedKey
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &newKey)
	if err != nil {
		return err
	}

	id, err := api.CreateAuthorizedKey(newKey, options.userID)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
func authorizedkeyUpdateCmd() *cobra.Command {
	options := authorizedkeyOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update authorized key for user",
		Long:  `Update authorized key for user`,
		Example: `
	privx-cli authorized-keys update [access flags] --id <KEY-ID> --user-id <USER-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return authorizedkeyUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "user-id", "", "user ID")
	flags.StringVar(&options.keyID, "id", "", "key ID")
	cmd.MarkFlagRequired("user-id")
	cmd.MarkFlagRequired("id")

	return cmd
}

func authorizedkeyUpdate(options authorizedkeyOptions, args []string) error {
	var updateKey rolestore.AuthorizedKey
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &updateKey)
	if err != nil {
		return err
	}

	err = api.UpdateAuthorizedKey(&updateKey, options.userID, options.keyID)
	if err != nil {
		return err
	}

	return nil
}

//
//
func authorizedkeyDeleteCmd() *cobra.Command {
	options := authorizedkeyOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete user's authorized key",
		Long:  `Delete user's authorized key`,
		Example: `
	privx-cli authorized-keys delete [access flags] --id <KEY-ID> --user-id <USER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return authorizedkeyDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "user-id", "", "user ID")
	flags.StringVar(&options.keyID, "id", "", "key ID")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("user-id")

	return cmd
}

func authorizedkeyDelete(options authorizedkeyOptions) error {
	api := rolestore.New(curl())

	err := api.DeleteAuthorizedKey(options.userID, options.keyID)
	if err != nil {
		return err
	}

	return nil
}

//
//
func authorizedkeyResolveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve authorized keys",
		Long:  `Resolve authorized keys`,
		Example: `
	privx-cli authorized-keys resolve [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return authorizedkeyResolve(args)
		},
	}

	return cmd
}

func authorizedkeyResolve(args []string) error {
	var resolveKey rolestore.ResolveAuthorizedKey
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &resolveKey)
	if err != nil {
		return err
	}

	key, err := api.ResolveAuthorizedKey(resolveKey)
	if err != nil {
		return err
	}

	return stdout(key)
}
