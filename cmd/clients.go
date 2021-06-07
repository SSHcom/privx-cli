//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/spf13/cobra"
)

var (
	trustedClientID string
)

func init() {
	rootCmd.AddCommand(clientListCmd)

	clientListCmd.AddCommand(clientCreateCmd)

	clientListCmd.AddCommand(clientShowCmd)
	clientShowCmd.Flags().StringVar(&trustedClientID, "id", "", "unique trusted client id")
	clientShowCmd.MarkFlagRequired("id")

	clientListCmd.AddCommand(clientDeleteCmd)
	clientDeleteCmd.Flags().StringVar(&trustedClientID, "id", "", "unique trusted client id")
	clientDeleteCmd.MarkFlagRequired("id")

	clientListCmd.AddCommand(clientUpdateCmd)
	clientUpdateCmd.Flags().StringVar(&trustedClientID, "id", "", "unique trusted client id")
	clientUpdateCmd.MarkFlagRequired("id")
}

//
//
var clientListCmd = &cobra.Command{
	Use:   "clients",
	Short: "Get trusted clients",
	Long:  `Get trusted clients`,
	Example: `
privx-cli clients [access flags]
	`,
	SilenceUsage: true,
	RunE:         trustedClients,
}

func trustedClients(cmd *cobra.Command, args []string) error {
	api := userstore.New(curl())

	trustedClients, err := api.TrustedClients()
	if err != nil {
		return err
	}

	return stdout(trustedClients)
}

//
//
var clientCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new trusted-client",
	Long:  `Create new trusted client`,
	Example: `
privx-cli clients create [access flags] JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         clientCreate,
}

func clientCreate(cmd *cobra.Command, args []string) error {
	var trustedClient userstore.TrustedClient
	api := userstore.New(curl())

	err := readJSON(args[0], &trustedClient)
	if err != nil {
		return err
	}

	id, err := api.CreateTrustedClient(trustedClient)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
var clientShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Get trusted client by ID",
	Long:  `Get trusted client by ID`,
	Example: `
privx-cli clients show [access flags] --id TRUSTED-CLIENT-ID
	`,
	SilenceUsage: true,
	RunE:         clientShow,
}

func clientShow(cmd *cobra.Command, args []string) error {
	api := userstore.New(curl())

	trustedClient, err := api.TrustedClient(trustedClientID)
	if err != nil {
		return err
	}

	return stdout(trustedClient)
}

//
//
var clientDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete trusted client",
	Long:  `Delete a trusted client`,
	Example: `
privx-cli clients delete [access flags] --id TRUSTED-CLIENT-ID
	`,
	SilenceUsage: true,
	RunE:         clientDelete,
}

func clientDelete(cmd *cobra.Command, args []string) error {
	api := userstore.New(curl())

	err := api.DeleteTrustedClient(trustedClientID)

	return err
}

//
//
var clientUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update trusted client",
	Long:  `Update an existing trusted client`,
	Example: `
privx-cli clients update [access flags] --id TRUSTED-CLIENT-ID JSON-FILE
	`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         clientUpdate,
}

func clientUpdate(cmd *cobra.Command, args []string) error {
	var trustedClient userstore.TrustedClient
	api := userstore.New(curl())

	err := readJSON(args[0], &trustedClient)
	if err != nil {
		return err
	}

	err = api.UpdateTrustedClient(trustedClientID, &trustedClient)

	return err
}
