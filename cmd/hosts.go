//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	apiConfig "github.com/SSHcom/privx-sdk-go/api/config"
	"github.com/SSHcom/privx-sdk-go/api/hoststore"
	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/spf13/cobra"
)

type hostOptions struct {
	hostID         string
	filter         string
	sortkey        string
	sortdir        string
	deployStatus   bool
	disabledStatus bool
	limit          int
	offset         int
}

func init() {
	rootCmd.AddCommand(hostListCmd())
}

//
//
func hostListCmd() *cobra.Command {
	options := hostOptions{}

	cmd := &cobra.Command{
		Use:   "hosts",
		Short: "List and manage PrivX hosts",
		Long:  `List and manage PrivX hosts`,
		Example: `
	privx-cli hosts [access flags] --offset <OFFSET> --sortkey <SORTKEY>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hostList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort object by name, updated, or created.")
	flags.StringVar(&options.filter, "filter", "", "filter hosts, possible values: accessible or configured")

	cmd.AddCommand(hostSearchCmd())
	cmd.AddCommand(hostCreateCmd())
	cmd.AddCommand(hostShowCmd())
	cmd.AddCommand(hostUpdateCmd())
	cmd.AddCommand(hostDeleteCmd())
	cmd.AddCommand(hostResolveCmd())
	cmd.AddCommand(hostDeployableCmd())
	cmd.AddCommand(hostDisableCmd())
	cmd.AddCommand(hostSettingListCmd())
	cmd.AddCommand(hostsDeployCmd())

	return cmd
}

func hostList(options hostOptions) error {
	api := hoststore.New(curl())

	hosts, err := api.Hosts(options.offset, options.limit, options.sortkey,
		strings.ToUpper(options.sortdir), options.filter)
	if err != nil {
		return err
	}

	return stdout(hosts)
}

//
//
func hostSearchCmd() *cobra.Command {
	options := hostOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search hosts",
		Long:  `Search hosts`,
		Example: `
	privx-cli hosts search [access flags] --offset <OFFSET> --sortkey <SORTKEY>
	privx-cli hosts search [access flags] --limit <LIMIT> JSON-FILE
		`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hostSearch(options, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.filter, "filter", "", "filter hosts, possible values: accessible or configured")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC")

	return cmd
}

func hostSearch(options hostOptions, args []string) error {
	var searchObject hoststore.HostSearchObject
	api := hoststore.New(curl())

	if len(args) == 1 {
		err := decodeJSON(args[0], &searchObject)
		if err != nil {
			return err
		}
	}

	hosts, err := api.SearchHost(options.sortkey, strings.ToUpper(options.sortdir), options.filter,
		options.offset, options.limit, &searchObject)
	if err != nil {
		return err
	}

	return stdout(hosts)
}

//
//
func hostCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new host",
		Long:  `Create new host`,
		Example: `
	privx-cli hosts create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hostCreate(cmd, args)
		},
	}

	return cmd
}

func hostCreate(cmd *cobra.Command, args []string) error {
	var newHost hoststore.Host
	api := hoststore.New(curl())

	err := decodeJSON(args[0], &newHost)
	if err != nil {
		return err
	}

	id, err := api.CreateHost(newHost)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
func hostShowCmd() *cobra.Command {
	options := hostOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get host by ID",
		Long:  `Get host by ID. Host ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli hosts show [access flags] --id <HOST-ID>,<HOST-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hostShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.hostID, "id", "", "host ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func hostShow(options hostOptions) error {
	api := hoststore.New(curl())
	hosts := []hoststore.Host{}

	for _, id := range strings.Split(options.hostID, ",") {
		host, err := api.Host(id)
		if err != nil {
			return err
		}
		hosts = append(hosts, *host)
	}

	return stdout(hosts)
}

//
//
func hostUpdateCmd() *cobra.Command {
	options := hostOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update host",
		Long:  `Update host`,
		Example: `
	privx-cli hosts update [access flags] --id <HOST-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hostUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.hostID, "id", "", "unique host ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func hostUpdate(options hostOptions, args []string) error {
	var updateHost hoststore.Host
	api := hoststore.New(curl())

	err := decodeJSON(args[0], &updateHost)
	if err != nil {
		return err
	}

	err = api.UpdateHost(options.hostID, &updateHost)
	if err != nil {
		return err
	}

	return nil
}

//
//
func hostDeleteCmd() *cobra.Command {
	options := hostOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete host",
		Long:  `Delete host. Host ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli hosts delete [access flags] --id <HOST-ID>,<HOST-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hostDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.hostID, "id", "", "unique host ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func hostDelete(options hostOptions) error {
	api := hoststore.New(curl())

	for _, id := range strings.Split(options.hostID, ",") {
		err := api.DeleteHost(id)
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
func hostResolveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve host",
		Long:  `Resolve service and address to a single host`,
		Example: `
	privx-cli hosts resolve [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hostResolve(cmd, args)
		},
	}

	return cmd
}

func hostResolve(cmd *cobra.Command, args []string) error {
	var service hoststore.Service
	api := hoststore.New(curl())

	err := decodeJSON(args[0], &service)
	if err != nil {
		return err
	}

	host, err := api.ResolveHost(service)
	if err != nil {
		return err
	}

	return stdout(host)
}

//
//
func hostDeployableCmd() *cobra.Command {
	options := hostOptions{}

	cmd := &cobra.Command{
		Use:   "deployable",
		Short: "Set a host to be depoyable or undeployable",
		Long:  `Set a host to be depoyable or undeployable. Host ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli hosts deployable [access flags] --id <HOST-ID>,<HOST-ID> --status=true
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hostDeployable(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.hostID, "id", "", "host ID")
	flags.BoolVar(&options.deployStatus, "status", false, "host deploy status")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("status")

	return cmd
}

func hostDeployable(options hostOptions) error {
	api := hoststore.New(curl())

	for _, id := range strings.Split(options.hostID, ",") {
		err := api.UpdateDeployStatus(id, options.deployStatus)
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
func hostDisableCmd() *cobra.Command {
	options := hostOptions{}

	cmd := &cobra.Command{
		Use:   "disabled",
		Short: "Enable/disable host",
		Long:  `Enable(false)/disable(true) host. Host ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli hosts disabled [access flags] --id <HOST-ID>,<HOST-ID> --status=true
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hostDisable(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.hostID, "id", "", "host ID")
	flags.BoolVar(&options.disabledStatus, "status", false, "host disabled status")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("status")

	return cmd
}

func hostDisable(options hostOptions) error {
	api := hoststore.New(curl())

	for _, id := range strings.Split(options.hostID, ",") {
		err := api.UpdateDisabledHostStatus(id, options.disabledStatus)
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
func hostSettingListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Get the default service options",
		Long:  `Get the default service options`,
		Example: `
	privx-cli hosts settings [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hostSettingList()
		},
	}

	return cmd
}

func hostSettingList() error {
	api := hoststore.New(curl())

	settings, err := api.ServiceOptions()
	if err != nil {
		return err
	}

	return stdout(settings)
}

//
//
func hostsDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Creates target hosts deployment config",
		Long:  `Creates target hosts deployment config`,
		Example: `
	privx-cli hosts deploy [access flags] <NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return hostDeploy(args)
		},
	}

	return cmd
}

func hostDeploy(args []string) error {
	if len(args) < 1 {
		return errors.New("requires name of deployment configuration as an argument")
	}
	name := args[0]

	curl := curl()
	store := userstore.New(curl)

	seq, err := store.TrustedClients()
	if err != nil {
		return err
	}

	cli := findClientID(seq, name)
	if cli == "" {
		cli, err = store.CreateTrustedClient(
			userstore.HostProvisioning(name),
		)
		if err != nil {
			return err
		}
	}

	conf := apiConfig.New(curl)
	file, err := conf.ConfigDeploy(cli)
	if err != nil {
		return err
	}

	os.Stdout.Write(file)
	return nil
}

func findClientID(seq []userstore.TrustedClient, name string) string {
	for _, cli := range seq {
		if cli.Name == name {
			return cli.ID
		}
	}
	return ""
}
