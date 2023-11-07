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

type sourceOptions struct {
	sourceID string
}

func init() {
	rootCmd.AddCommand(sourceListCmd())
}

//
//
func sourceListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sources",
		Short: "List and manage user and host directories",
		Long:  `List and manage user and host directories`,
		Example: `
privx-cli sources [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sourceList()
		},
	}

	cmd.AddCommand(sourceCreateCmd())
	cmd.AddCommand(sourceShowCmd())
	cmd.AddCommand(sourceDeleteCmd())
	cmd.AddCommand(sourceUpdateCmd())
	cmd.AddCommand(sourceRefreshCmd())

	return cmd
}

func sourceList() error {
	api := rolestore.New(curl())

	sources, err := api.Sources()
	if err != nil {
		return err
	}

	return stdout(sources)
}

//
//
func sourceCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new source",
		Long:  `Create a new source`,
		Example: `
	privx-cli sources create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sourceCreate(cmd, args)
		},
	}

	return cmd
}

func sourceCreate(cmd *cobra.Command, args []string) error {
	var newSource rolestore.Source
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &newSource)
	if err != nil {
		return err
	}

	id, err := api.CreateSource(newSource)
	if err != nil {
		return err
	}

	return stdout(id)
}

//
//
func sourceShowCmd() *cobra.Command {
	options := sourceOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get source by ID",
		Long:  `Get source by ID`,
		Example: `
	privx-cli sources show [access flags] --id <SOURCE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sourceShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.sourceID, "id", "", "source ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func sourceShow(options sourceOptions) error {
	api := rolestore.New(curl())

	source, err := api.Source(options.sourceID)
	if err != nil {
		return err
	}

	return stdout(source)
}

//
//
func sourceDeleteCmd() *cobra.Command {
	options := sourceOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete source",
		Long:  `Delete source. Source ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli sources delete [access flags] --id <SOURCE-ID>,<SOURCE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sourceDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.sourceID, "id", "", "unique source ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func sourceDelete(options sourceOptions) error {
	api := rolestore.New(curl())

	for _, id := range strings.Split(options.sourceID, ",") {
		err := api.DeleteSource(id)
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
func sourceUpdateCmd() *cobra.Command {
	options := sourceOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update source",
		Long:  `Update source`,
		Example: `
	privx-cli sources update [access flags] --id <SOURCE-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sourceUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.sourceID, "id", "", "unique source ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func sourceUpdate(options sourceOptions, args []string) error {
	var updateSource rolestore.Source
	api := rolestore.New(curl())

	err := decodeJSON(args[0], &updateSource)
	if err != nil {
		return err
	}

	err = api.UpdateSource(options.sourceID, &updateSource)
	if err != nil {
		return err
	}

	return nil
}

//
//
func sourceRefreshCmd() *cobra.Command {
	options := sourceOptions{}

	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh Source",
		Long:  `Refresh Source. Source ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli sources refresh [access flags] --id <SOURCE-ID>,<SOURCE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sourceRefresh(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.sourceID, "id", "", "source ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func sourceRefresh(options sourceOptions) error {
	api := rolestore.New(curl())

	err := api.RefreshSources(strings.Split(options.sourceID, ","))
	if err != nil {
		return err
	}

	return nil
}
