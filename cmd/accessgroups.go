//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/authorizer"
	"github.com/spf13/cobra"
)

type accessGroupOptions struct {
	accessGroupID string
	caID          string
	sortkey       string
	sortdir       string
	offset        int
	limit         int
}

func (m accessGroupOptions) normalize_sortdir() string {
	return strings.ToUpper(m.sortdir)
}

func init() {
	rootCmd.AddCommand(accessGroupListCmd())
}

func accessGroupListCmd() *cobra.Command {
	options := accessGroupOptions{}

	cmd := &cobra.Command{
		Use:   "access-groups",
		Short: "List and manage access groups",
		Long:  `List and manage access groups`,
		Example: `
	privx-cli access-groups [access flags] --limit <LIMIT> --offset <OFFSET>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return accessGroupList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC")

	cmd.AddCommand(accessGroupCreateCmd())
	cmd.AddCommand(accessGroupSearchCmd())
	cmd.AddCommand(accessGroupShowCmd())
	cmd.AddCommand(accessGroupUpdateCmd())
	cmd.AddCommand(renewCAKeyCmd())
	cmd.AddCommand(revokeCAKeyCmd())

	return cmd
}

func accessGroupList(options accessGroupOptions) error {
	api := authorizer.New(curl())

	groups, err := api.AccessGroups(options.offset, options.limit,
		options.sortkey, options.normalize_sortdir())
	if err != nil {
		return err
	}

	return stdout(groups)
}

func accessGroupCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create access group",
		Long:  `Create access group`,
		Example: `
	privx-cli access-groups create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return accessGroupCreate(cmd, args)
		},
	}

	return cmd
}

func accessGroupCreate(cmd *cobra.Command, args []string) error {
	var newAccessGroup authorizer.AccessGroup
	api := authorizer.New(curl())

	err := decodeJSON(args[0], &newAccessGroup)
	if err != nil {
		return err
	}

	id, err := api.CreateAccessGroup(&newAccessGroup)
	if err != nil {
		return err
	}

	return stdout(id)
}

func accessGroupSearchCmd() *cobra.Command {
	options := accessGroupOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search access groups",
		Long:  `Search access groups`,
		Example: `
	privx-cli access-groups search [access flags] --offset <OFFSET> --sortkey <SORTKEY>
	privx-cli access-groups search [access flags] --limit <LIMIT> JSON-FILE
		`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return accessGroupSearch(options, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC")

	return cmd
}

func accessGroupSearch(options accessGroupOptions, args []string) error {
	var searchObject authorizer.SearchParams
	api := authorizer.New(curl())

	if len(args) == 1 {
		err := decodeJSON(args[0], &searchObject)
		if err != nil {
			return err
		}
	}

	hosts, err := api.SearchAccessGroup(options.offset, options.limit, options.sortkey,
		options.normalize_sortdir(), &searchObject)
	if err != nil {
		return err
	}

	return stdout(hosts)
}

func accessGroupShowCmd() *cobra.Command {
	options := accessGroupOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get access group by ID",
		Long:  `Get access group by ID. Access group ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli access-groups show [access flags] --id <ACCESS-GROUP-ID>,<ACCESS-GROUP-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return accessGroupShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.accessGroupID, "id", "", "access group ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func accessGroupShow(options accessGroupOptions) error {
	api := authorizer.New(curl())

	group, err := api.AccessGroup(options.accessGroupID)
	if err != nil {
		return err
	}

	return stdout(group)
}

func accessGroupUpdateCmd() *cobra.Command {
	options := accessGroupOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update access group",
		Long:  `Update access group`,
		Example: `
	privx-cli access-groups update [access flags] --id <ACCESS-GROUP-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return accessGroupUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.accessGroupID, "id", "", "access group ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func accessGroupUpdate(options accessGroupOptions, args []string) error {
	var updateAccessGroup authorizer.AccessGroup
	api := authorizer.New(curl())

	err := decodeJSON(args[0], &updateAccessGroup)
	if err != nil {
		return err
	}

	err = api.UpdateAccessGroup(options.accessGroupID, &updateAccessGroup)
	if err != nil {
		return err
	}

	return nil
}

func renewCAKeyCmd() *cobra.Command {
	options := accessGroupOptions{}

	cmd := &cobra.Command{
		Use:   "renew-ca",
		Short: "Renew CA key",
		Long:  `Renew CA key for a given access group`,
		Example: `
	privx-cli access-groups renew-ca [access flags] --id <ACCESS-GROUP-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return renewCAKey(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.accessGroupID, "id", "", "access group ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func renewCAKey(options accessGroupOptions) error {
	api := authorizer.New(curl())

	id, err := api.CreateAccessGroupsIdCas(options.accessGroupID)
	if err != nil {
		return err
	}

	return stdout(id)
}

func revokeCAKeyCmd() *cobra.Command {
	options := accessGroupOptions{}

	cmd := &cobra.Command{
		Use:   "revoke-ca",
		Short: "Revoke CA key",
		Long:  `Revoke CA key for a given access group`,
		Example: `
	privx-cli access-groups revoke-ca [access flags] --id <ACCESS-GROUP-ID> --ca-id <CA-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return revokeCAKey(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.accessGroupID, "id", "", "access group ID")
	flags.StringVar(&options.caID, "ca-id", "", "CA ID")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("ca-id")

	return cmd
}

func revokeCAKey(options accessGroupOptions) error {
	api := authorizer.New(curl())

	err := api.DeleteAccessGroupsIdCas(options.accessGroupID, options.caID)
	if err != nil {
		return err
	}

	return nil
}
