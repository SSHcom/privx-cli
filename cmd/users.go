//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/spf13/cobra"
)

var (
	userID         string
	password       string
	query          string
	offset         int
	limit          int
	userQuery      []string
	userRoleGrant  []string
	userRoleRevoke []string
)

func init() {
	rootCmd.AddCommand(userListCmd)
	userListCmd.Flags().StringArrayVarP(&userQuery, "query", "q", []string{}, "query PrivX users with keyword")

	userListCmd.AddCommand(userCreateCmd)

	userListCmd.AddCommand(userUpdateCmd)
	userUpdateCmd.Flags().StringVar(&userID, "uid", "", "unique user id")
	userUpdateCmd.MarkFlagRequired("uid")

	userListCmd.AddCommand(userDeleteCmd)
	userDeleteCmd.Flags().StringVar(&userID, "uid", "", "unique user id")
	userDeleteCmd.MarkFlagRequired("uid")

	userListCmd.AddCommand(userUpdatePasswordCmd)
	userUpdatePasswordCmd.Flags().StringVar(&userID, "uid", "", "unique user id")
	userUpdatePasswordCmd.Flags().StringVar(&password, "password", "", "new password for local user")
	userUpdatePasswordCmd.MarkFlagRequired("uid")
	userUpdatePasswordCmd.MarkFlagRequired("password")

	userListCmd.AddCommand(userShowCmd)

	userListCmd.AddCommand(usersRolesCmd)
	usersRolesCmd.Flags().StringVar(&userID, "uid", "", "user unique id")
	usersRolesCmd.Flags().StringArrayVar(&userRoleGrant, "grant", []string{}, "grant role to user, requires role unique id.")
	usersRolesCmd.Flags().StringArrayVar(&userRoleRevoke, "revoke", []string{}, "revoke role from user, requires role unique id.")
	usersRolesCmd.MarkFlagRequired("uid")
}

//
//
var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new local user",
	Long:  `Create new local user`,
	Example: `
privx-cli users create [access flags] JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         userCreate,
}

func userCreate(cmd *cobra.Command, args []string) error {
	var newUser userstore.LocalUser
	api := userstore.New(curl())

	err := readJSON(args[0], &newUser)
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
var userUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update local user",
	Long:  `Update a local user inside`,
	Example: `
privx-cli users update [access flags] JSON-FILE --uid UID
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         userUpdate,
}

func userUpdate(cmd *cobra.Command, args []string) error {
	var updateUser userstore.LocalUser
	api := userstore.New(curl())

	err := readJSON(args[0], &updateUser)
	if err != nil {
		return err
	}
	fmt.Println(updateUser)
	err = api.UpdateLocalUser(userID, &updateUser)
	if err != nil {
		return err
	}

	return err
}

//
//
var userDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete local user",
	Long:  `Delete a local user`,
	Example: `
privx-cli users delete [access flags] --uid UID
	`,
	SilenceUsage: true,
	RunE:         userDelete,
}

func userDelete(cmd *cobra.Command, args []string) error {
	api := userstore.New(curl())

	err := api.DeleteLocalUser(userID)

	return err
}

//
//
var userUpdatePasswordCmd = &cobra.Command{
	Use:   "update-password",
	Short: "Update local user password",
	Long:  `Update a local users password`,
	Example: `
privx-cli users update-password [access flags] --uid UID --password NEW-PASSWORD
	`,
	SilenceUsage: true,
	RunE:         userUpdatePassword,
}

func userUpdatePassword(cmd *cobra.Command, args []string) error {
	newPassword := userstore.Password{
		Password: password,
	}
	api := userstore.New(curl())

	err := api.UpdateLocalUserPassword(userID, &newPassword)

	return err
}

//
//
var userListCmd = &cobra.Command{
	Use:   "users",
	Short: "PrivX users",
	Long:  `List and manage PrivX users`,
	Example: `
privx-cli users [access flags]
	`,
	SilenceUsage: true,
	RunE:         userList,
}

func userList(cmd *cobra.Command, args []string) error {
	api := rolestore.New(curl())
	users, err := api.SearchUsers(strings.Join(userQuery, ","), "")
	if err != nil {
		return err
	}

	return stdout(users)
}

//
//
var userShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Description about PrivX user",
	Long:  `Description about PrivX user`,
	Example: `
privx-cli users show [access flags] UID ...
	`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE:         userShow,
}

func userShow(cmd *cobra.Command, args []string) error {
	api := rolestore.New(curl())
	users := []rolestore.User{}

	for _, uid := range args {
		user, err := api.User(uid)
		if err != nil {
			return err
		}
		users = append(users, *user)
	}

	return stdout(users)
}

//
//
var usersRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "user roles",
	Long:  `list and manage user roles`,
	Example: `
privx-cli users roles [access flags] --uid UID
privx-cli users roles [access flags] --uid UID --grant role-uid
privx-cli users roles [access flags] --uid UID --revoke role-uid
	`,
	SilenceUsage: true,
	RunE:         userRoles,
}

func userRoles(cmd *cobra.Command, args []string) error {
	api := rolestore.New(curl())

	for _, role := range userRoleGrant {
		err := api.GrantUserRole(userID, role)
		if err != nil {
			return err
		}
	}

	for _, role := range userRoleRevoke {
		err := api.RevokeUserRole(userID, role)
		if err != nil {
			return err
		}
	}

	roles, err := api.UserRoles(userID)
	if err != nil {
		return err
	}
	return stdout(roles)
}

func readJSON(name string, object interface{}) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &object)
	if err != nil {
		return err
	}

	return nil
}
