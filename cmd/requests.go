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

type requestOptions struct {
	requestID string
	filter    string
	sortkey   string
	sortdir   string
	limit     int
	offset    int
}

func init() {
	rootCmd.AddCommand(requestListCmd())
}

//
//
func requestListCmd() *cobra.Command {
	options := requestOptions{}

	cmd := &cobra.Command{
		Use:   "requests",
		Short: "List and manage request queues",
		Long:  `List and manage the request queue for the user`,
		Example: `
	privx-cli requests [access flags] --offset <OFFSET> --limit <LIMIT> --filter <FILTER>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return requestList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.filter, "filter", "", "filter request items")

	cmd.AddCommand(requestCreateCmd())
	cmd.AddCommand(requestShowCmd())
	cmd.AddCommand(requestDeleteCmd())
	cmd.AddCommand(requestHandlingCmd())
	cmd.AddCommand(requestSearchCmd())

	return cmd
}

func requestList(options requestOptions) error {
	api := workflow.New(curl())

	requests, err := api.Requests(options.offset, options.limit, options.filter)
	if err != nil {
		return err
	}

	return stdout(requests)
}

//
//
func requestCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create request",
		Long:  `Add a workflow to the request queue`,
		Example: `
	privx-cli requests create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return requestCreate(cmd, args)
		},
	}

	return cmd
}

func requestCreate(cmd *cobra.Command, args []string) error {
	var newRequest workflow.Request
	api := workflow.New(curl())

	err := decodeJSON(args[0], &newRequest)
	if err != nil {
		return err
	}

	id, err := api.CreateRequest(&newRequest)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
func requestShowCmd() *cobra.Command {
	options := requestOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get request by ID",
		Long:  `Get a request object by ID. Request ID's are separated by commas when using multiple values, see example.`,
		Example: `
	privx-cli requests show [access flags] --id <REQUEST-ID>,<REQUEST-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return requestShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.requestID, "id", "", "request ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func requestShow(options requestOptions) error {
	api := workflow.New(curl())
	requests := []workflow.Request{}

	for _, id := range strings.Split(options.requestID, ",") {
		request, err := api.Request(id)
		if err != nil {
			return err
		}
		requests = append(requests, *request)
	}

	return stdout(requests)
}

func requestRevokeTargetRoleCmd() *cobra.Command {
	options := requestOptions{}

	cmd := &cobra.Command{
		Use:   "Revoke",
		Short: "Revoke role",
		Long:  `Revoke role by Request ID. Request ID's are separated by commas when using multiple values, see example.`,
		Example: `
	privx-cli requests delete [access flags] --id <REQUEST-ID>,<REQUEST-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return requestRevokeTargetRole(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.requestID, "id", "", "request ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func requestRevokeTargetRole(options requestOptions) error {
	api := workflow.New(curl())

	for _, id := range strings.Split(options.requestID, ",") {
		err := api.RevokeTargetRole(id)
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
func requestDeleteCmd() *cobra.Command {
	options := requestOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete request",
		Long:  `Delete request item by ID. Request ID's are separated by commas when using multiple values, see example.`,
		Example: `
	privx-cli requests delete [access flags] --id <REQUEST-ID>,<REQUEST-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return requestDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.requestID, "id", "", "request ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func requestDelete(options requestOptions) error {
	api := workflow.New(curl())

	for _, id := range strings.Split(options.requestID, ",") {
		err := api.DeleteRequest(id)
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
func requestHandlingCmd() *cobra.Command {
	options := requestOptions{}

	cmd := &cobra.Command{
		Use:   "handle-request",
		Short: "Update a request in queue",
		Long:  `Update a request in queue. Only users with matching role are permitted to change the status of a step requiring such role.`,
		Example: `
	privx-cli requests request-decision [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return requestHandling(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.requestID, "id", "", "unique workflow ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func requestHandling(options requestOptions, args []string) error {
	var request workflow.Decision
	api := workflow.New(curl())

	err := decodeJSON(args[0], &request)
	if err != nil {
		return err
	}

	err = api.MakeDecisionOnRequest(options.requestID, request)
	if err != nil {
		return err
	}

	return nil
}

//
//
func requestSearchCmd() *cobra.Command {
	options := requestOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search workflow requests",
		Long:  `Search workflow requests`,
		Example: `
	privx-cli requests search [access flags] --offset <OFFSET> --limit <LIMIT> --filter <FILTER>
	privx-cli requests search [access flags] JSON-FILE
		`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return requestSearch(options, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.filter, "filter", "", "filter request items(requests, active_requests, approvals, etc.)")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC")

	return cmd
}

func requestSearch(options requestOptions, args []string) error {
	var searchObject workflow.Search
	api := workflow.New(curl())

	if len(args) == 1 {
		err := decodeJSON(args[0], &searchObject)
		if err != nil {
			return err
		}
	}

	requests, err := api.SearchRequests(options.offset, options.limit, strings.ToUpper(options.sortdir),
		options.sortkey, strings.ToUpper(options.filter), &searchObject)
	if err != nil {
		return err
	}

	return stdout(requests)
}
