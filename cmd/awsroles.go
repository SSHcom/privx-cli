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

type awsRoleOptions struct {
	awsRoleID string
	refresh   bool
}

func init() {
	rootCmd.AddCommand(awsRoleListCmd())
}

//
//
func awsRoleListCmd() *cobra.Command {
	options := awsRoleOptions{}

	cmd := &cobra.Command{
		Use:   "aws-roles",
		Short: "List and manage AWS role links",
		Long:  `List and manage AWS role links`,
		Example: `
	privx-cli aws-roles [access flags] --refresh=true
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return awsRoleList(options)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&options.refresh, "refresh", false, "refresh the AWS roles from AWS directories before fetching")

	cmd.AddCommand(awsRoleShowCmd())
	cmd.AddCommand(awsRoleDeleteCmd())
	cmd.AddCommand(awsRoleUpdateCmd())
	cmd.AddCommand(linkedRoleListCmd())

	return cmd
}

func awsRoleList(options awsRoleOptions) error {
	api := rolestore.New(curl())

	roles, err := api.AWSRoleLinks(options.refresh)
	if err != nil {
		return err
	}

	return stdout(roles)
}

//
//
func awsRoleShowCmd() *cobra.Command {
	options := awsRoleOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get AWS role by ID",
		Long:  `Get AWS role by ID`,
		Example: `
	privx-cli aws-roles show [access flags] --id <AWS-ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return awsRoleShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.awsRoleID, "id", "", "AWS role ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func awsRoleShow(options awsRoleOptions) error {
	api := rolestore.New(curl())

	role, err := api.AWSRoleLink(options.awsRoleID)
	if err != nil {
		return err
	}

	return stdout(role)
}

//
//
func awsRoleDeleteCmd() *cobra.Command {
	options := awsRoleOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete cached AWS role and its mappings",
		Long: `Delete cached AWS role and its mappings. Does not affect the AWS service, if the role still exists on AWS, it will re-appear on the next role scan.
Host ID's are separated by commas when using multiple values, see example.`,
		Example: `
	privx-cli aws-roles delete [access flags] --id <AWS-ROLE-ID>,<AWS-ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return awsRoleDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.awsRoleID, "id", "", "AWS role ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func awsRoleDelete(options awsRoleOptions) error {
	api := rolestore.New(curl())

	for _, id := range strings.Split(options.awsRoleID, ",") {
		err := api.DeleteAWSRoleLInk(id)
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
func awsRoleUpdateCmd() *cobra.Command {
	options := awsRoleOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a AWS role granting PrivX roles",
		Long:  `Update a AWS role granting PrivX roles`,
		Example: `
	privx-cli aws-roles update [access flags] --id <AWS-ROLE-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return awsRoleUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.awsRoleID, "id", "", "AWS role ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func awsRoleUpdate(options awsRoleOptions, args []string) error {
	var updateHost []rolestore.RoleRef
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &updateHost)
	if err != nil {
		return err
	}

	err = api.UpdateAWSRoleLink(options.awsRoleID, updateHost)
	if err != nil {
		return err
	}

	return nil
}

//
//
func linkedRoleListCmd() *cobra.Command {
	options := awsRoleOptions{}

	cmd := &cobra.Command{
		Use:   "linked-roles",
		Short: "Get AWS role granting PrivX roles",
		Long:  `Get AWS role granting PrivX roles`,
		Example: `
	privx-cli aws-roles linked-roles [access flags] --id <AWS-ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return linkedRoleList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.awsRoleID, "id", "", "AWS role ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func linkedRoleList(options awsRoleOptions) error {
	api := rolestore.New(curl())

	roles, err := api.LinkedRoles(options.awsRoleID)
	if err != nil {
		return err
	}

	return stdout(roles)
}
