//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/authorizer"
	"github.com/spf13/cobra"
)

type authorizerOptions struct {
	accessGroupID   string
	caID            string
	fileName        string
	trustedClientID string
	sortkey         string
	sortdir         string
	limit           int
	offset          int
}

func init() {
	rootCmd.AddCommand(authorizerListCmd())
}

func (m authorizerOptions) normalize_sortdir() string {
	return strings.ToUpper(m.sortdir)
}

//
//
func authorizerListCmd() *cobra.Command {
	options := authorizerOptions{}

	cmd := &cobra.Command{
		Use:          "authorizer",
		Short:        "List and manage authorizer root certificates",
		Long:         `List and manage authorizer root certificates`,
		SilenceUsage: true,
		Example: `
	privx-cli authorizer [access flags]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return authorizerList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.accessGroupID, "access-group-id", "", "access group ID filter")

	cmd.AddCommand(authorizerShowCmd())
	cmd.AddCommand(authorizerRevocationListCmd())
	cmd.AddCommand(targetHostCredentialShowCmd())
	cmd.AddCommand(deploymentScriptDownloadCmd())
	cmd.AddCommand(principalCommandScriptDownloadCmd())
	cmd.AddCommand(sslTrustAnchorShowCmd())
	cmd.AddCommand(extenderTrustAnchorShowCmd())
	cmd.AddCommand(certificateSearchCmd())
	cmd.AddCommand(getCertByIDCmd())
	cmd.AddCommand(certificateListCmd())

	return cmd
}

func authorizerList(options authorizerOptions) error {
	api := authorizer.New(curl())

	certificates, err := api.CACertificates(options.accessGroupID)
	if err != nil {
		return err
	}

	return stdout(certificates)
}

//
//
func authorizerShowCmd() *cobra.Command {
	options := authorizerOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get authorizer's root certificate",
		Long:  `Get authorizer's root certificate`,
		Example: `
	privx-cli authorizer show [access flags] --id <CA-ID> --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return authorizerShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.caID, "id", "", "ca ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func authorizerShow(options authorizerOptions) error {
	api := authorizer.New(curl())

	err := api.CACertificate(options.caID, options.fileName)
	if err != nil {
		return err
	}

	return nil
}

//
//
func authorizerRevocationListCmd() *cobra.Command {
	options := authorizerOptions{}

	cmd := &cobra.Command{
		Use:   "show-crl",
		Short: "Get authorizer CA's certificate revocation list",
		Long:  `Get authorizer CA's certificate revocation list`,
		Example: `
	privx-cli authorizer revocation-list [access flags] --id <CA-ID> --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return authorizerRevocationList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.caID, "id", "", "ca ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func authorizerRevocationList(options authorizerOptions) error {
	api := authorizer.New(curl())

	err := api.CertificateRevocationList(options.caID, options.fileName)
	if err != nil {
		return err
	}

	return nil
}

//
//
func targetHostCredentialShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "target-host-credentials",
		Short: "Get target host credentials for the user",
		Long:  `Get target host credentials for the user`,
		Example: `
	privx-cli authorizer target-host-credentials [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return targetHostCredential(args)
		},
	}

	return cmd
}

func targetHostCredential(args []string) error {
	var authorizationRequest authorizer.AuthorizationRequest
	api := authorizer.New(curl())

	err := decodeJSON(args[0], &authorizationRequest)
	if err != nil {
		return err
	}

	certificate, err := api.TargetHostCredentials(&authorizationRequest)
	if err != nil {
		return err
	}

	return stdout(certificate)
}

//
//
func deploymentScriptDownloadCmd() *cobra.Command {
	options := authorizerOptions{}

	cmd := &cobra.Command{
		Use:   "deployment-script",
		Short: "Get deployment script pre-configured for PrivX installation",
		Long:  `Get deployment script pre-configured for PrivX installation`,
		Example: `
	privx-cli authorizer deployment-script [access flags] --trusted-client-id <TRUSTED-CLIENT-ID> --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploymentScriptDownload(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.trustedClientID, "trusted-client-id", "", "trusted client ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	cmd.MarkFlagRequired("trusted-client-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func deploymentScriptDownload(options authorizerOptions) error {
	api := authorizer.New(curl())

	handler, err := api.DeployScriptDownloadHandle(options.trustedClientID)
	if err != nil {
		return err
	}

	err = api.DownloadDeployScript(options.trustedClientID, handler.SessionID, options.fileName)
	if err != nil {
		return err
	}

	return nil
}

//
//
func principalCommandScriptDownloadCmd() *cobra.Command {
	options := authorizerOptions{}

	cmd := &cobra.Command{
		Use:   "principal-cmd-script",
		Short: "Get the principals command script",
		Long:  `Get the principals command script`,
		Example: `
	privx-cli authorizer principal-cmd-script [access flags] --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return principalCommandScriptDownload(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.fileName, "name", "", "file name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func principalCommandScriptDownload(options authorizerOptions) error {
	api := authorizer.New(curl())

	err := api.DownloadPrincipalCommandScript(options.fileName)
	if err != nil {
		return err
	}

	return nil
}

//
//
func sslTrustAnchorShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ssl-trust-anchor",
		Short:        "Get the SSL trust anchor",
		Long:         `Get the SSL trust anchor`,
		SilenceUsage: true,
		Example: `
	privx-cli authorizer ssl-trust-anchor [access flags]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sslTrustAnchorShow()
		},
	}

	return cmd
}

func sslTrustAnchorShow() error {
	api := authorizer.New(curl())

	anchor, err := api.SSLTrustAnchor()
	if err != nil {
		return err
	}

	return stdout(anchor)
}

//
//
func extenderTrustAnchorShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "extender-trust-anchor",
		Short:        "Get the extender trust anchor",
		Long:         `Get the extender trust anchor`,
		SilenceUsage: true,
		Example: `
	privx-cli authorizer extender-trust-anchor [access flags]
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return extenderTrustAnchorShow()
		},
	}

	return cmd
}

func extenderTrustAnchorShow() error {
	api := authorizer.New(curl())

	anchor, err := api.ExtenderTrustAnchor()
	if err != nil {
		return err
	}

	return stdout(anchor)
}

//
//
func certificateSearchCmd() *cobra.Command {
	options := authorizerOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search certificates",
		Long:  `Search certificates`,
		Example: `
	privx-cli authorizer search [access flags] --offset <OFFSET> --sortkey <SORTKEY>
	privx-cli authorizer search [access flags] --limit <LIMIT> JSON-FILE
		`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return certificateSearch(options, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC")

	return cmd
}

func certificateSearch(options authorizerOptions, args []string) error {
	var searchObject authorizer.APICertificateSearch
	api := authorizer.New(curl())

	if len(args) == 1 {
		err := decodeJSON(args[0], &searchObject)
		if err != nil {
			return err
		}
	}

	cert, err := api.SearchCert(options.offset, options.limit, options.sortkey,
		options.normalize_sortdir(), &searchObject)
	if err != nil {
		return err
	}

	return stdout(cert)
}

//
func certificateListCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:          "cert-list",
		Short:        "Get all Certificates",
		Long:         `Get all Certificates`,
		SilenceUsage: true,
		Example: `
	privx-cli authorizer cert-list
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return certificateList()
		},
	}
	return cmd
}

func certificateList() error {
	api := authorizer.New(curl())

	certificates, err := api.GetAllCertificates()
	if err != nil {
		return err
	}

	return stdout(certificates)
}

//
//
func getCertByIDCmd() *cobra.Command {
	var ID string
	cmd := &cobra.Command{
		Use:          "get-cert",
		Short:        "Get Certificate by ID",
		Long:         `Get Certificate by ID`,
		SilenceUsage: true,
		Example: `
	privx-cli authorizer get-cert [access flags] -id <ID>
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCertByID(ID)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&ID, "id", "", "Certificate ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func getCertByID(ID string) error {
	api := authorizer.New(curl())

	cert, err := api.GetCertByID(ID)
	if err != nil {
		return err
	}

	return stdout(cert)
}
