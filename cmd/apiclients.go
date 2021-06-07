//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/spf13/cobra"
)

var (
	clientID       string
	apiClientRoles string
)

func init() {
	rootCmd.AddCommand(apiClientListCmd)

	apiClientListCmd.AddCommand(apiClientCreateCmd)
	apiClientCreateCmd.Flags().StringVar(&name, "name", "", "API client name")
	apiClientCreateCmd.Flags().StringVar(&apiClientRoles, "roles", "", "list of roles possessed by the API client")

	apiClientListCmd.AddCommand(apiClientShowCmd)
	apiClientShowCmd.Flags().StringVar(&clientID, "id", "", "unique API client id")
	apiClientShowCmd.MarkFlagRequired("id")

	apiClientListCmd.AddCommand(apiClientDeleteCmd)
	apiClientDeleteCmd.Flags().StringVar(&clientID, "id", "", "unique API client id")
	apiClientDeleteCmd.MarkFlagRequired("id")

	apiClientListCmd.AddCommand(apiClientUpdateCmd)
	apiClientUpdateCmd.Flags().StringVar(&clientID, "id", "", "unique API client id")
	apiClientUpdateCmd.MarkFlagRequired("id")
}

//
//
var apiClientListCmd = &cobra.Command{
	Use:   "api-clients",
	Short: "Get API clients",
	Long:  `Get all API clients`,
	Example: `
privx-cli api-clients [access flags]
	`,
	SilenceUsage: true,
	RunE:         apiClientList,
}

func apiClientList(cmd *cobra.Command, args []string) error {
	api := userstore.New(curl())
	result, err := api.APIClients()
	if err != nil {
		return err
	}

	return stdout(result)
}

//
//
var apiClientCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new API client",
	Long:  `Create new API client`,
	Example: `
privx-cli api-clients create [access flags] --name NAME --roles ROLE-ID,ROLE-ID
	`,
	SilenceUsage: true,
	RunE:         apiClientCreate,
}

func apiClientCreate(cmd *cobra.Command, args []string) error {
	api := userstore.New(curl())

	id, err := api.CreateAPIClient(name, strings.Split(apiClientRoles, ","))
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
var apiClientShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Get API client by ID",
	Long:  `Get API client by ID`,
	Example: `
privx-cli api-clients show [access flags] --id API-CLIENT-ID
	`,
	SilenceUsage: true,
	RunE:         apiClientShow,
}

func apiClientShow(cmd *cobra.Command, args []string) error {
	api := userstore.New(curl())

	result, err := api.APIClient(clientID)
	if err != nil {
		return err
	}

	return stdout(result)
}

//
//
var apiClientDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete API client",
	Long:  `Delete a API client`,
	Example: `
privx-cli api-clients delete [access flags] --id API-CLIENT-ID
	`,
	SilenceUsage: true,
	RunE:         apiClientDelete,
}

func apiClientDelete(cmd *cobra.Command, args []string) error {
	api := userstore.New(curl())

	err := api.DeleteAPIClient(clientID)

	return err
}

//
//
var apiClientUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update API client",
	Long:  `Update an existing API client`,
	Example: `
privx-cli users local api-clients update [access flags] --id API-CLIENT-ID JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         apiClientUpdate,
}

func apiClientUpdate(cmd *cobra.Command, args []string) error {
	var apiClient userstore.APIClient
	api := userstore.New(curl())

	err := readJSON(args[0], &apiClient)
	if err != nil {
		return err
	}

	err = api.UpdateAPIClient(clientID, &apiClient)

	return err
}
