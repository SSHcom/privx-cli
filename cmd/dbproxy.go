//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"github.com/SSHcom/privx-sdk-go/api/dbproxy"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dbProxyCmd())
}

func dbProxyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db-proxy",
		Short: "Db-proxy service commands",
		Long:  `Db-proxy service commands`,
		Example: `
	privx-cli db-proxy config [access flags]
		`,
		SilenceUsage: true,
	}

	cmd.AddCommand(dbproxyConfCmd())

	return cmd
}

func dbproxyConfCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Get DB proxy configuration",
		Long:  `Get DB proxy configuration. Includes info about ca_certificate and certificate_chain`,
		Example: `
	privx-cli db-proxy config [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return dbproxyConf()
		},
	}

	return cmd
}

func dbproxyConf() error {
	api := dbproxy.New(curl())
	config, err := api.DbProxyConf()
	if err != nil {
		return err
	}

	return stdout(config)
}
