//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/spf13/cobra"
)

type UserSearchOptions struct {
	limit   int
	offset  int
	sortkey string
	sortdir string
	keyword string
	source  string
}
type userOptions struct {
	userID         string
	enable         bool
	disable        bool
	reset          bool
	sources        []string
	keywords       []string
	userRoleGrant  []string
	userRoleRevoke []string
	userIDs        []string
	search         UserSearchOptions
}

func init() {
	rootCmd.AddCommand(userListCmd())
}

//
//

func userListCmd() *cobra.Command {

	options := userOptions{}

	cmd := &cobra.Command{
		Use:   "users",
		Short: "List and manage users",
		Long:  `List and manage users`,
		Example: `
	privx-cli users [access flags] --keywords <KEYWORD>,<KEYWORD>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.search.keyword, "keywords", "", "comma or space-separated string to search in secret's names")
	flags.IntVar(&options.search.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.search.limit, "limit", 50, "max number of items to return")
	flags.StringVar(&options.search.sortdir, "sortdir", "ASC", "sort direction, ASC or DESC (default ASC)")
	flags.StringVar(&options.search.sortkey, "sortkey", "", "sort object by property: source, email, principal, full_name.")
	flags.StringArrayVar(&options.userIDs, "userid", []string{}, "list of users IDs.")

	cmd.AddCommand(userShowCmd())
	cmd.AddCommand(userSettingShowCmd())
	cmd.AddCommand(userSettingsUpdateCmd())
	cmd.AddCommand(usersRolesCmd())
	cmd.AddCommand(userMFACmd())
	cmd.AddCommand(externalUserSearchCmd())

	return cmd
}
func usersValidateSortDir(sortdir string) error {
	sortdirAllowedValues := []string{"ASC", "DESC"}

	for _, a := range sortdirAllowedValues {
		if a == strings.ToUpper(sortdir) {
			return nil
		}
	}

	return fmt.Errorf("sortdir field must be one of these values %q", sortdirAllowedValues)
}
func usersValidateSortKey(sortkey string) error {
	sortkeyAllowedValues := []string{"name", "updated", "created"}

	for _, a := range sortkeyAllowedValues {
		if a == strings.ToLower(sortkey) {
			return nil
		}
	}

	return fmt.Errorf("filter field must be one of these values %q", sortkeyAllowedValues)
}

func userList(options userOptions) error {
	err := usersValidateSortDir(options.search.sortdir)
	if err != nil {
		return err
	}
	err = usersValidateSortKey(options.search.sortkey)
	if err != nil {
		return err
	}
	api := rolestore.New(curl())

	searchBody := rolestore.UserSearchObject{
		Keywords: options.search.keyword,
		Source:   options.search.source,
		UserIDs:  options.userIDs,
	}

	users, err := api.SearchUsers(options.search.offset, options.search.limit, strings.ToLower(options.search.sortkey), strings.ToUpper(options.search.sortdir), searchBody)
	if err != nil {
		return err
	}

	return stdout(users)
}

//
//
func userShowCmd() *cobra.Command {
	options := userOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Description about PrivX user",
		Long:  `Description about PrivX user. User ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli users show [access flags] --id <USER-ID>,<USER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "user ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func userShow(options userOptions) error {
	api := rolestore.New(curl())
	users := []rolestore.User{}

	for _, id := range strings.Split(options.userID, ",") {
		user, err := api.User(id)
		if err != nil {
			return err
		}
		users = append(users, *user)
	}

	return stdout(users)
}

//
//
func userSettingShowCmd() *cobra.Command {
	options := userOptions{}

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Show specific user settings",
		Long:  `Show specific user settings.`,
		Example: `
	privx-cli users settings [access flags] --id <USER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userSettingShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "user ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func userSettingShow(options userOptions) error {
	api := rolestore.New(curl())

	settings, err := api.UserSettings(options.userID)
	if err != nil {
		return err
	}

	return stdout(settings)
}

//
//
func userSettingsUpdateCmd() *cobra.Command {
	options := userOptions{}

	cmd := &cobra.Command{
		Use:   "update-settings",
		Short: "Update specific user's settings",
		Long:  `Update specific user's settings`,
		Example: `
	privx-cli users update-settings [access flags] --id <USER-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userSettingsUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "user ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func userSettingsUpdate(options userOptions, args []string) error {
	var updateSettings *json.RawMessage
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &updateSettings)
	if err != nil {
		return err
	}

	err = api.UpdateUserSettings(updateSettings, options.userID)
	if err != nil {
		return err
	}

	return nil
}

//
//
func usersRolesCmd() *cobra.Command {
	options := userOptions{}

	cmd := &cobra.Command{
		Use:   "roles",
		Short: "Show and manage specific user roles",
		Long:  `Show and manage specific user roles`,
		Example: `
	privx-cli users roles [access flags] --id <USER-ID>
	privx-cli users roles [access flags] --id <USER-ID> --grant <ROLE-ID>
	privx-cli users roles [access flags] --id <USER-ID> --revoke <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userRoles(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "user ID")
	flags.StringArrayVar(&options.userRoleGrant, "grant", []string{}, "grant role to user, requires role unique id.")
	flags.StringArrayVar(&options.userRoleRevoke, "revoke", []string{}, "revoke role from user, requires role unique id.")
	cmd.MarkFlagRequired("id")

	return cmd
}

func userRoles(options userOptions) error {
	api := rolestore.New(curl())

	for _, role := range options.userRoleGrant {
		err := api.GrantUserRole(options.userID, role)
		if err != nil {
			return err
		}
	}

	for _, role := range options.userRoleRevoke {
		err := api.RevokeUserRole(options.userID, role)
		if err != nil {
			return err
		}
	}

	roles, err := api.UserRoles(options.userID)
	if err != nil {
		return err
	}
	return stdout(roles)
}

//
//
func userMFACmd() *cobra.Command {
	options := userOptions{}

	cmd := &cobra.Command{
		Use:   "mfa",
		Short: "Enable, disable or reset multifactor authentication",
		Long:  `Enable, disable or reset multifactor authentication. User ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli users mfa [access flags] --id <USER-ID>,<USER-ID> --enable
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userMFA(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "user ID")
	flags.BoolVarP(&options.enable, "enable", "e", false, "turn on multifactor authentication")
	flags.BoolVarP(&options.disable, "disable", "d", false, "turn off multifactor authentication")
	flags.BoolVarP(&options.reset, "reset", "r", false, "reset multifactor authentication")
	cmd.MarkFlagRequired("id")

	return cmd
}

func userMFA(options userOptions) error {
	if options.enable {
		enableMFA(options)
	} else if options.disable {
		disableMFA(options)
	} else if options.reset {
		resetMFA(options)
	} else {
		fmt.Fprintln(os.Stderr, "Error: you have to specify one of the following flag: --enable, --disable or --reset")
		os.Exit(1)
	}

	return nil
}

func enableMFA(options userOptions) error {
	api := rolestore.New(curl())

	err := api.EnableMFA(strings.Split(options.userID, ","))
	if err != nil {
		return err
	}

	return nil
}

func disableMFA(options userOptions) error {
	api := rolestore.New(curl())

	err := api.DisableMFA(strings.Split(options.userID, ","))
	if err != nil {
		return err
	}

	return nil
}

func resetMFA(options userOptions) error {
	api := rolestore.New(curl())

	err := api.ResetMFA(strings.Split(options.userID, ","))
	if err != nil {
		return err
	}

	return nil
}

//
//
func externalUserSearchCmd() *cobra.Command {
	options := userOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search external users",
		Long:  `Search external users`,
		Example: `
	privx-cli users search-external-users [access flags] --keywords <KEYWORD>,<KEYWORD>
	privx-cli users search-external-users [access flags] --keywords <KEYWORD> --sources <SOURCE>,<SOURCE>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return externalUserSearch(options)
		},
	}

	flags := cmd.Flags()
	flags.StringArrayVarP(&options.keywords, "keywords", "", []string{}, "search keywords")
	flags.StringArrayVarP(&options.sources, "sources", "", []string{}, "the source ID where to search the user from")

	return cmd
}

func externalUserSearch(options userOptions) error {
	api := rolestore.New(curl())
	users, err := api.SearchUsersExternal(strings.Join(options.keywords, ","),
		strings.Join(options.sources, ","))
	if err != nil {
		return err
	}

	return stdout(users)
}
