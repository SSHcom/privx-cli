//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"fmt"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/workflow"
	"github.com/spf13/cobra"
)

type workflowOptions struct {
	workflowID string
	limit      int
	offset     int
}

func init() {
	rootCmd.AddCommand(workflowListCmd())
}

//
//
func workflowListCmd() *cobra.Command {
	options := workflowOptions{}

	cmd := &cobra.Command{
		Use:   "workflows",
		Short: "List and manage workflows",
		Long:  `List and manage PrivX workflows`,
		Example: `
	privx-cli workflows [access flags] --offset <OFFSET> --limit <LIMIT>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workflowList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")

	cmd.AddCommand(workflowCreateCmd())
	cmd.AddCommand(workflowShowCmd())
	cmd.AddCommand(workflowDeleteCmd())
	cmd.AddCommand(workflowUpdateCmd())
	cmd.AddCommand(workflowSettingListCmd())
	cmd.AddCommand(workflowSettingsUpdateCmd())
	cmd.AddCommand(testEmailNotificationCmd())

	return cmd
}

func workflowList(options workflowOptions) error {
	api := workflow.New(curl())

	workflows, err := api.Workflows(options.offset, options.limit)
	if err != nil {
		return err
	}

	return stdout(workflows)
}

//
//
func workflowCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new workflow",
		Long:  `Create a new workflow`,
		Example: `
	privx-cli workflows create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workflowCreate(cmd, args)
		},
	}

	return cmd
}

func workflowCreate(cmd *cobra.Command, args []string) error {
	var newWorkflow workflow.Workflow
	api := workflow.New(curl())

	err := decodeJSON(args[0], &newWorkflow)
	if err != nil {
		return err
	}

	id, err := api.CreateWorkflow(&newWorkflow)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
func workflowShowCmd() *cobra.Command {
	options := workflowOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get workflow by ID",
		Long:  `Get workflow object by ID. Workflow ID's are separated by commas when using multiple values, see example.`,
		Example: `
	privx-cli workflows show [access flags] --id <WORKFLOW-ID>,<WORKFLOW-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workflowShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.workflowID, "id", "", "unique workflow ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func workflowShow(options workflowOptions) error {
	api := workflow.New(curl())
	workflows := []workflow.Workflow{}

	for _, id := range strings.Split(options.workflowID, ",") {
		result, err := api.Workflow(id)
		if err != nil {
			return err
		}
		workflows = append(workflows, *result)
	}

	return stdout(workflows)
}

//
//
func workflowDeleteCmd() *cobra.Command {
	options := workflowOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete workflow by ID",
		Long:  `Delete workflow by ID. Workflow ID's are separated by commas when using multiple values, see example.`,
		Example: `
	privx-cli workflows delete [access flags] --id <WORKFLOW-ID>,<WORKFLOW-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workflowDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.workflowID, "id", "", "unique workflow ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func workflowDelete(options workflowOptions) error {
	api := workflow.New(curl())

	for _, id := range strings.Split(options.workflowID, ",") {
		err := api.DeleteWorkflow(id)
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
func workflowUpdateCmd() *cobra.Command {
	options := workflowOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a workflow",
		Long:  `Update a workflow`,
		Example: `
	privx-cli workflows update [access flags] --id <WORKFLOW-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workflowUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.workflowID, "id", "", "unique workflow ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func workflowUpdate(options workflowOptions, args []string) error {
	var updateWorkflow workflow.Workflow
	api := workflow.New(curl())

	err := decodeJSON(args[0], &updateWorkflow)
	if err != nil {
		return err
	}

	err = api.UpdateWorkflow(options.workflowID, &updateWorkflow)
	if err != nil {
		return err
	}

	return nil
}

//
//
func workflowSettingListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Get workflow settings",
		Long:  `Get workflow settings`,
		Example: `
	privx-cli workflows settings [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workflowSettingList()
		},
	}

	return cmd
}

func workflowSettingList() error {
	api := workflow.New(curl())

	settings, err := api.Settings()
	if err != nil {
		return err
	}

	return stdout(settings)
}

//
//
func workflowSettingsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-settings",
		Short: "Update workflow settings",
		Long:  `Update workflow settings`,
		Example: `
	privx-cli workflows update-settings [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return workflowSettingsUpdate(cmd, args)
		},
	}

	return cmd
}

func workflowSettingsUpdate(cmd *cobra.Command, args []string) error {
	var updateSettings workflow.Settings
	api := workflow.New(curl())

	err := decodeJSON(args[0], &updateSettings)
	if err != nil {
		return err
	}

	err = api.UpdateSettings(&updateSettings)
	if err != nil {
		return err
	}

	return nil
}

//
//
func testEmailNotificationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testsmtp",
		Short: "Test the email settings",
		Long:  `Test the email settings`,
		Example: `
	privx-cli workflows testsmtp [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return testEmailNotification(cmd, args)
		},
	}

	return cmd
}

func testEmailNotification(cmd *cobra.Command, args []string) error {
	var smtp workflow.Settings
	api := workflow.New(curl())

	err := decodeJSON(args[0], &smtp)
	if err != nil {
		return err
	}

	testResult, err := api.TestEmailNotification(&smtp)
	if err != nil {
		return err
	}

	return stdout(testResult)
}
