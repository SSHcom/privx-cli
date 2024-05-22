//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"fmt"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/spf13/cobra"
)

type roleOptions struct {
	roleID    string
	roleName  string
	tokenCode string
	ttl       int
}

func init() {
	rootCmd.AddCommand(roleListCmd())
	rootCmd.AddCommand(idendityProvidersListCmd())
}

func roleListCmd() *cobra.Command {
	options := searchOptions{}

	cmd := &cobra.Command{
		Use:   "roles",
		Short: "List and manage PrivX roles",
		Long:  `List and manage PrivX roles`,
		Example: `
	privx-cli roles [access flags] --offset offset --limit limit
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return roleList(options.offset, options.limit, options.sortdir, options.sortkey)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "Offset where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "Number of items to return")
	flags.StringVar(&options.sortdir, "sortdir", "ASC", "Sort direction, ASC or DESC")
	flags.StringVar(&options.sortkey, "sortkey", "name", "Sort by specific object property")

	cmd.AddCommand(roleCreateCmd())
	cmd.AddCommand(roleShowCmd())
	cmd.AddCommand(roleDeleteCmd())
	cmd.AddCommand(roleUpdateCmd())
	cmd.AddCommand(rolesMemberListCmd())
	cmd.AddCommand(roleResolveCmd())
	cmd.AddCommand(awsTokenShowCmd())

	return cmd
}

func roleList(offset, limit int, sortdir, sortkey string) error {
	api := rolestore.New(curl())

	roles, err := api.Roles(offset, limit, sortkey, sortdir)
	if err != nil {
		return err
	}

	return stdout(roles)
}

func roleCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new role",
		Long:  `Create new role`,
		Example: `
	privx-cli roles create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return roleCreate(args)
		},
	}

	return cmd
}

func roleCreate(args []string) error {
	var newRole rolestore.Role
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &newRole)
	if err != nil {
		return err
	}

	id, err := api.CreateRole(newRole)
	if err != nil {
		return err
	}

	return stdout(id)
}

func roleShowCmd() *cobra.Command {
	options := roleOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get role by ID",
		Long:  `Get role by ID`,
		Example: `
	privx-cli roles show [access flags] --id <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return roleShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleID, "id", "", "role ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func roleShow(options roleOptions) error {
	api := rolestore.New(curl())

	role, err := api.Role(options.roleID)
	if err != nil {
		return err
	}

	return stdout(role)
}

func roleDeleteCmd() *cobra.Command {
	options := roleOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete role",
		Long:  `Delete role. Role ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli roles delete [access flags] --id <ROLE-ID>,<ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return roleDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleID, "id", "", "role ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func roleDelete(options roleOptions) error {
	api := rolestore.New(curl())

	for _, id := range strings.Split(options.roleID, ",") {
		err := api.DeleteRole(id)
		if err != nil {
			return err
		} else {
			fmt.Println(id)
		}
	}

	return nil
}

func roleUpdateCmd() *cobra.Command {
	options := roleOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update role",
		Long:  `Update role`,
		Example: `
	privx-cli roles update [access flags] --id <ROLE-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return roleUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleID, "id", "", "role ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func roleUpdate(options roleOptions, args []string) error {
	var updateRole rolestore.Role
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &updateRole)
	if err != nil {
		return err
	}

	err = api.UpdateRole(options.roleID, &updateRole)
	if err != nil {
		return err
	}

	return nil
}

func rolesMemberListCmd() *cobra.Command {
	options := roleOptions{}

	cmd := &cobra.Command{
		Use:   "members",
		Short: "Get members of PrivX role",
		Long:  `Get members of PrivX role`,
		Example: `
	privx-cli roles members [access flags] UID ...
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return roleMemberList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleID, "id", "", "role ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func roleMemberList(options roleOptions) error {
	api := rolestore.New(curl())
	members := []rolestore.User{}

	for _, role := range strings.Split(options.roleID, ",") {
		member, err := api.GetRoleMembers(role)
		if err != nil {
			return err
		}
		members = append(members, member...)
	}

	return stdout(members)
}

func roleResolveCmd() *cobra.Command {
	options := roleOptions{}

	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve role names and return corresponding ID's",
		Long:  `Resolve role names and return corresponding ID's. Role name's separated by commas when using multiple values, see example`,
		Example: `
	privx-cli roles resolve [access flags] --name <ROLE-NAME>,<ROLE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return roleResolve(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleName, "name", "", "role name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func roleResolve(options roleOptions) error {
	api := rolestore.New(curl())

	id, err := api.ResolveRoles(strings.Split(options.roleName, ","))
	if err != nil {
		return err
	}

	return stdout(id)
}

func awsTokenShowCmd() *cobra.Command {
	options := roleOptions{}

	cmd := &cobra.Command{
		Use:   "aws-token",
		Short: "Get an AWS token for a role",
		Long: `Get an AWS token for a role. Return 403 on an initial request if the AWS role has multi-factor authentication enabled.
Subsequent request must contain MFA as a query parameter. Return 403 if the user does not have the role.`,
		Example: `
	privx-cli roles aws-token [access flags] --id <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return awsTokenShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.roleID, "id", "", "role ID")
	flags.StringVar(&options.tokenCode, "mfa", "", "multi-factor-authentication code")
	flags.IntVar(&options.ttl, "ttl", 50, "max time validity for the token")
	cmd.MarkFlagRequired("id")

	return cmd
}

func awsTokenShow(options roleOptions) error {
	api := rolestore.New(curl())

	token, err := api.AWSToken(options.roleID, options.tokenCode, options.ttl)
	if err != nil {
		return err
	}

	return stdout(token)
}

//
// Idendity providers
//

func idendityProvidersListCmd() *cobra.Command {
	options := searchOptions{}
	cmd := &cobra.Command{
		Use:   "idendity",
		Short: "List all idendity providers",
		Long:  `List all idendity providers`,
		Example: `
	privx-cli idendity [access flags] --offset offset --limit limit
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return idendityProvidersList(options.offset, options.limit)
		},
	}
	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "Offset where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "Number of items to return")

	cmd.AddCommand(idendityCreateCmd())
	cmd.AddCommand(idendityShowCmd())
	cmd.AddCommand(idendityDeleteCmd())
	cmd.AddCommand(idendityUpdateCmd())
	cmd.AddCommand(idenditySearchCmd())

	return cmd
}

func idendityProvidersList(offset, limit int) error {
	api := rolestore.New(curl())

	idendity, err := api.GetAllIdendityProviders(offset, limit)
	if err != nil {
		return err
	}

	return stdout(idendity.Items)
}

func idendityCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new idendity provider",
		Long:  `Create new idendity provider`,
		Example: `
	privx-cli idendity create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return idendityCreate(args)
		},
	}

	return cmd
}

func idendityCreate(args []string) error {
	var newIDProvider rolestore.IdentityProvider
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &newIDProvider)
	if err != nil {
		return err
	}

	id, err := api.CreateIdendityProvider(newIDProvider)
	if err != nil {
		return err
	}

	return stdout(id)
}

func idendityShowCmd() *cobra.Command {
	var ID string

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get idendity provider by ID",
		Long:  `Get idendity provider by ID`,
		Example: `
	privx-cli idendity show [access flags] --id <IDENDITY-PROVIDER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return idendityShow(ID)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&ID, "id", "", "Idendity provider ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func idendityShow(ID string) error {
	api := rolestore.New(curl())

	IDProvider, err := api.GetIdendityProviderByID(ID)
	if err != nil {
		return err
	}

	return stdout(IDProvider)
}

func idendityDeleteCmd() *cobra.Command {
	var IDs string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete idendity providers",
		Long:  `Delete idendity providers. Idendity provider's ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli idendity delete [access flags] --id <IDENDITY-PROVIDER-ID>,<IDENDITY-PROVIDER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return idendityDelete(IDs)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&IDs, "id", "", "role ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func idendityDelete(IDs string) error {
	api := rolestore.New(curl())

	for _, id := range strings.Split(IDs, ",") {
		err := api.DeleteIdendityProviderByID(id)
		if err != nil {
			return err
		} else {
			fmt.Println(id)
		}
	}

	return nil
}

func idendityUpdateCmd() *cobra.Command {
	var ID string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update idendity provider",
		Long:  `Update idendity provider`,
		Example: `
	privx-cli idendity update [access flags] --id <IDENDITY-PROVIDER-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return idendityUpdate(ID, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&ID, "id", "", "role ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func idendityUpdate(ID string, args []string) error {
	var updateIDProvider rolestore.IdentityProvider
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &updateIDProvider)
	if err != nil {
		return err
	}

	err = api.UpdateIdendityProvider(updateIDProvider, ID)
	if err != nil {
		return err
	}

	return nil
}

func idenditySearchCmd() *cobra.Command {
	searchParams := rolestore.Params{}
	var keywords string
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search idendity provider",
		Long:  `Search idendity provider`,
		Example: `
	privx-cli idendity search [access flags] --offset <OFFSET> --limit <LIMIT> --sortkey <SORTKEY> --sortdir <SORTDIR>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return idenditySearch(searchParams, keywords)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&searchParams.Sortkey, "sortkey", "", "Sort by specific object property")
	flags.StringVar(&searchParams.Sortdir, "sortdir", "ASC", "Sort direction, ASC or DESC")
	flags.StringVar(&keywords, "keywords", "", "comma or space separated list of search keywords")
	flags.IntVar(&searchParams.Limit, "limit", 50, "Number of items to return")
	flags.IntVar(&searchParams.Offset, "offset", 0, "Offset where to start fetching the items")

	return cmd
}

func idenditySearch(params rolestore.Params, keywords string) error {
	api := rolestore.New(curl())

	result, err := api.SearchIdendityProviders(
		params.Offset,
		params.Limit,
		params.Sortkey,
		strings.ToUpper(params.Sortdir),
		keywords)
	if err != nil {
		return err
	}

	return stdout(result)
}

//
//
