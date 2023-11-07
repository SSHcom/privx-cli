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

type collectorOptions struct {
	collectorID string
}

func init() {
	rootCmd.AddCommand(collectorListCmd())
}

//
//
func collectorListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collectors",
		Short: "List and manage logconf collectors",
		Long:  `List and manage logconf collectors`,
		Example: `
	privx-cli collectors [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return collectorList()
		},
	}

	cmd.AddCommand(collectorCreateCmd())
	cmd.AddCommand(collectorShowCmd())
	cmd.AddCommand(collectorUpdateCmd())
	cmd.AddCommand(collectorDeleteCmd())

	return cmd
}

func collectorList() error {
	api := rolestore.New(curl())

	collectors, err := api.LogconfCollectors()
	if err != nil {
		return err
	}

	return stdout(collectors)
}

//
//
func collectorCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create logconf collector",
		Long:  `Create logconf collector`,
		Example: `
	privx-cli collectors create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return collectorCreate(args)
		},
	}

	return cmd
}

func collectorCreate(args []string) error {
	var newCollector rolestore.LogconfCollector
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &newCollector)
	if err != nil {
		return err
	}

	id, err := api.CreateLogconfCollector(newCollector)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
func collectorShowCmd() *cobra.Command {
	options := collectorOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get logconf collector",
		Long:  `Get logconf collector`,
		Example: `
	privx-cli collectors show [access flags] --collector-id <COLLECTOR-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return collectorShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.collectorID, "collector-id", "", "collector ID")
	cmd.MarkFlagRequired("collector-id")

	return cmd
}

func collectorShow(options collectorOptions) error {
	api := rolestore.New(curl())

	collector, err := api.LogconfCollector(options.collectorID)
	if err != nil {
		return err
	}

	return stdout(collector)
}

//
//
func collectorUpdateCmd() *cobra.Command {
	options := collectorOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update logconf collector",
		Long:  `Update logconf collector`,
		Example: `
	privx-cli collectors update [access flags] --collector-id <COLLECTOR-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return collectorUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.collectorID, "collector-id", "", "collector ID")
	cmd.MarkFlagRequired("collector-id")

	return cmd
}

func collectorUpdate(options collectorOptions, args []string) error {
	var updateCollector rolestore.LogconfCollector
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &updateCollector)
	if err != nil {
		return err
	}

	err = api.UpdateLogconfCollector(options.collectorID, &updateCollector)
	if err != nil {
		return err
	}

	return nil
}

//
//
func collectorDeleteCmd() *cobra.Command {
	options := collectorOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete logconf collector",
		Long:  `Delete logconf collector. Collector ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli collectors delete [access flags] --collector-id <COLLECTOR-ID>,<COLLECTOR-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return collectorDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.collectorID, "collector-id", "", "collector ID")
	cmd.MarkFlagRequired("collector-id")

	return cmd
}

func collectorDelete(options collectorOptions) error {
	api := rolestore.New(curl())

	for _, id := range strings.Split(options.collectorID, ",") {
		err := api.DeleteLogconfCollector(id)
		if err != nil {
			return err
		} else {
			fmt.Println(id)
		}
	}

	return nil
}
