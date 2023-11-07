package cmd

import (
	authApi "github.com/SSHcom/privx-sdk-go/api/auth"
	"github.com/spf13/cobra"
)

type idpclientOptions struct {
	idpID string
}

func init() {
	rootCmd.AddCommand(idpClientsCmd())
}

//
//
func idpClientsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "idp-clients",
		Short: "IDP Clients Related Commands",
		Long:  `IDP Clients Related Commands`,
		Example: `
	privx-cli connections idp-clients [access flags]
		`,
		SilenceUsage: true,
	}

	cmd.AddCommand(idpClientCreateCmd())
	cmd.AddCommand(idpClientUpdateCmd())
	cmd.AddCommand(idpClientShowCmd())
	cmd.AddCommand(idpClientDeleteCmd())
	cmd.AddCommand(idpClientRegenerateIDCmd())

	return cmd
}

//
//
func idpClientCreateCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create PrivX IDP client",
		Long:  `Create PrivX IDP client`,
		Example: `
	privx-cli idp-clients create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return idpClientCreate(args)
		},
	}

	return cmd
}

func idpClientCreate(args []string) error {
	var idpClient authApi.IDPClient
	api := authApi.New(curl())

	err := decodeJSON(args[0], &idpClient)
	if err != nil {
		return err
	}
	idpId, err := api.CreateIdpClient(&idpClient)
	if err != nil {
		return err
	}

	return stdout(idpId)
}

//
//
func idpClientUpdateCmd() *cobra.Command {
	var options idpclientOptions

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update PrivX IDP client",
		Long:  `Update PrivX IDP client`,
		Example: `
	privx-cli idp-clients update [access flags] --id <IDP-ID> JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return idpClientUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.idpID, "id", "", "IDP ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func idpClientUpdate(options idpclientOptions, args []string) error {
	var idpClient authApi.IDPClient
	api := authApi.New(curl())

	err := decodeJSON(args[0], &idpClient)
	if err != nil {
		return err
	}
	return api.UpdateIdpClient(&idpClient, options.idpID)
}

//
//
func idpClientShowCmd() *cobra.Command {
	var options idpclientOptions

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get PrivX IDP client by ID",
		Long:  `Get PrivX IDP client by ID`,
		Example: `
	privx-cli idp-clients show [access flags] --id <IDP-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return idpClientShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.idpID, "id", "", "IDP ID")
	cmd.MarkFlagRequired("id")
	return cmd
}

func idpClientShow(options idpclientOptions) error {
	api := authApi.New(curl())

	client, err := api.IdpClient(options.idpID)
	if err != nil {
		return err
	}

	return stdout(client)
}

//
//
func idpClientDeleteCmd() *cobra.Command {
	var options idpclientOptions

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete PrivX IDP client by ID",
		Long:  `Delete PrivX IDP client by ID`,
		Example: `
	privx-cli idp-clients delete [access flags] --id <IDP-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return idpClientDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.idpID, "id", "", "IDP ID")
	cmd.MarkFlagRequired("id")
	return cmd
}

func idpClientDelete(options idpclientOptions) error {
	api := authApi.New(curl())
	return api.DeleteIdpClient(options.idpID)
}

//
//
func idpClientRegenerateIDCmd() *cobra.Command {
	var options idpclientOptions

	cmd := &cobra.Command{
		Use:   "regenerate",
		Short: "Regenerate client_id and client_secret for OIDC IDP client",
		Long:  `Regenerate client_id and client_secret for OIDC IDP client`,
		Example: `
	privx-cli idp-clients regenerate [access flags] --id <IDP-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return idpClientRegenerateID(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.idpID, "id", "", "IDP ID")
	cmd.MarkFlagRequired("id")
	return cmd
}

func idpClientRegenerateID(options idpclientOptions) error {
	api := authApi.New(curl())

	clientConfig, err := api.RegenerateIdpClientConfig(options.idpID)
	if err != nil {
		return err
	}

	return stdout(clientConfig)
}
