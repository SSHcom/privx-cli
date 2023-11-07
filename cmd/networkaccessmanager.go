//
// Copyright (c) 2022 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/networkaccessmanager"
	"github.com/spf13/cobra"
)

type searchOptions struct {
	offset   int
	limit    int
	sortkey  string
	sortdir  string
	name     string
	ID       string
	filter   string
	keywords string
}

func init() {
	rootCmd.AddCommand(networkListCmd())
}

//
//
func networkListCmd() *cobra.Command {
	options := searchOptions{}

	cmd := &cobra.Command{
		Use:   "nam",
		Short: "get network targets",
		Long:  `get Network targets`,
		Example: `
	privx-cli nam [access flags] --offset 0 --limit 50 --sortkey <SORTKEY> --sortdir <SORTDIR> --name <NAME>,<NAME> --id <ID> 
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return networkList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "id", "sort object by id, name, comment.., .")
	flags.StringVar(&options.sortdir, "sortdir", "ASC", "sort direction, ASC or DESC (default ASC)")
	flags.StringVar(&options.name, "name", "", "comma or space-separated string to search in secret's names")
	flags.StringVar(&options.ID, "id", "", "comma or space-separated string to search in secret's names")

	cmd.AddCommand(networkAccessManagerStatusCmd())
	cmd.AddCommand(createNetworkCmd())
	cmd.AddCommand(searchNetworkCmd())
	cmd.AddCommand(getNetworkByIDCmd())
	cmd.AddCommand(updateNetworkCmd())
	cmd.AddCommand(deleteNetworkByIDCmd())
	cmd.AddCommand(disableNetworkByIDCmd())

	return cmd
}

func networkList(options searchOptions) error {
	api := networkaccessmanager.New(curl())

	networks, err := api.GetNetworkTargets(
		options.offset,
		options.limit,
		options.sortkey,
		strings.ToUpper(options.sortdir),
		options.name,
		options.ID)
	if err != nil {
		return err
	}

	return stdout(networks)
}

//
//
//
func networkAccessManagerStatusCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get the status of the microservice",
		Long:  `Get the status of the Network Access Manager microservice`,
		Example: `
	privx-cli nam status
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return networkAccessManagerStatus()
		},
	}

	return cmd
}

func networkAccessManagerStatus() error {
	api := networkaccessmanager.New(curl())

	status, err := api.NetworkAccessManagerStatus()
	if err != nil {
		return err
	}

	return stdout(status)
}

//
//
//
func createNetworkCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new network",
		Long:  `Create a new network`,
		Example: `
	privx-cli nam create JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createNetwork(args[0])
		},
	}

	return cmd
}

func createNetwork(args string) error {
	var network networkaccessmanager.Item

	api := networkaccessmanager.New(curl())

	err := decodeJSON(args, &network)
	if err != nil {
		return err
	}

	stdout(network)
	id, err := api.CreateNetworkTargets(network)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
//
func searchNetworkCmd() *cobra.Command {
	options := searchOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search network targets",
		Long:  `Search Network targets`,
		Example: `
	privx-cli nam search[access flags] --offset 0 --limit 50 --sortkey <SORTKEY> --sortdir <SORTDIR> --keywords <KEYWORDS>,<KEYWORDS> --filter <FILTER> 
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return searchNetwork(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort object by id, name, comment.., .")
	flags.StringVar(&options.sortdir, "sortdir", "ASC", "sort direction, ASC or DESC (default ASC)")
	flags.StringVar(&options.filter, "filter", "", "comma or space-separated string to search in secret's names")
	flags.StringVar(&options.keywords, "keywords", "", "search keywords")

	return cmd
}

func searchNetwork(options searchOptions) error {

	api := networkaccessmanager.New(curl())

	result, err := api.SearchNetworkTargets(
		options.offset,
		options.limit,
		options.sortkey,
		strings.ToUpper(options.sortdir),
		options.filter,
		options.keywords)
	if err != nil {
		return err
	}

	return stdout(result)
}

//
//
//
func getNetworkByIDCmd() *cobra.Command {
	var networkID string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a network by ID",
		Long:  `Get a network by ID`,
		Example: `
	privx-cli nam get --id <NETWORK-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getNetworkByID(networkID)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&networkID, "id", "", "id of the target network")
	cmd.MarkFlagRequired("id")

	return cmd
}

func getNetworkByID(networkID string) error {
	api := networkaccessmanager.New(curl())

	result, err := api.GetNetworkTargetByID(networkID)
	if err != nil {
		return err
	}

	return stdout(result)
}

//
//
//
func updateNetworkCmd() *cobra.Command {
	var networkID string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a specific netword target",
		Long:  `Update a specific netword target`,
		Example: `
	privx-cli nam update --id <NETWORK-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateNetwork(args[0], networkID)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&networkID, "id", "", "id of the target network")
	cmd.MarkFlagRequired("id")

	return cmd
}

func updateNetwork(args string, targetID string) error {
	var network networkaccessmanager.Item

	api := networkaccessmanager.New(curl())

	err := decodeJSON(args, &network)
	if err != nil {
		return err
	}

	err = api.UpdateNetworkTarget(&network, targetID)
	if err != nil {
		return err
	}

	return nil
}

//
//
//
func deleteNetworkByIDCmd() *cobra.Command {
	var networkID string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a network by ID",
		Long:  `Delete a network by ID`,
		Example: `
	privx-cli nam delete --id <NETWORK-ID>
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteNetworkByID(networkID)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&networkID, "id", "", "id of the target network to delete")
	cmd.MarkFlagRequired("id")

	return cmd
}

func deleteNetworkByID(networkID string) error {
	api := networkaccessmanager.New(curl())

	err := api.DeleteNetworkTargetByID(networkID)
	if err != nil {
		return err
	}

	return nil
}

//
//
//
func disableNetworkByIDCmd() *cobra.Command {
	var networkID string
	var disableState bool
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "disable or enable a network by ID",
		Long:  `disable or enable a network by ID`,
		Example: `
	privx-cli nam delete --id <NETWORK-ID> --disable <true or false>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return disableNetworkByID(disableState, networkID)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&networkID, "id", "", "id of the target network to delete")
	flags.BoolVar(&disableState, "disable", true, "disable true or false")
	cmd.MarkFlagRequired("id")

	return cmd
}

func disableNetworkByID(disableState bool, networkID string) error {
	api := networkaccessmanager.New(curl())

	err := api.DisableNetworkTargetByID(disableState, networkID)
	if err != nil {
		return err
	}

	return nil
}
