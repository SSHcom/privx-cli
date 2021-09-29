//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/vault"
	"github.com/spf13/cobra"
)

type vaultOptions struct {
	secretName   string
	ownerID      string
	keywords     string
	sortkey      string
	sortdir      string
	vaultReadTo  []string
	vaultWriteTo []string
	limit        int
	offset       int
	filter       string
}

func init() {
	rootCmd.AddCommand(secretListCmd())
	rootCmd.AddCommand(userSecretListCmd())
}

//
//
func secretListCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "PrivX secrets",
		Long:  `List and manage PrivX secrets`,
		Example: `
	privx-cli secrets [access flags] --offset <OFFSET> --limit <LIMIT>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return secretList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")

	cmd.AddCommand(secretShowCmd())
	cmd.AddCommand(secretCreateCmd())
	cmd.AddCommand(vaultUpdateCmd())
	cmd.AddCommand(secretDeleteCmd())
	cmd.AddCommand(secretMetadataShowCmd())
	cmd.AddCommand(secretSearchCmd())
	cmd.AddCommand(secretSchemasShowCmd())

	return cmd
}

func secretList(options vaultOptions) error {
	api := vault.New(curl())

	secrets, err := api.Secrets(options.offset, options.limit)
	if err != nil {
		return err
	}

	return stdout(secrets)
}

//
//
func secretShowCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get a secret",
		Long:  `Get a secret. Secret Name's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli secrets show [access flags] --name <SECRET-NAME>,<SECRET-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return secretShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.secretName, "name", "", "secret name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func secretShow(options vaultOptions) error {
	api := vault.New(curl())
	secrets := []vault.Secret{}

	for _, name := range strings.Split(options.secretName, ",") {
		secret, err := api.Secret(name)
		if err != nil {
			return err
		}
		secrets = append(secrets, *secret)
	}

	return stdout(secrets)
}

//
//
func secretCreateCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new secret",
		Long:  `Create new secret`,
		Example: `
	privx-cli secrets create [access flags] --name <SECRET-NAME>
		--allow-read-to <ROLE-ID>
		--allow-write-to <ROLE-ID>
		...
		JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return secretCreate(args, options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.secretName, "name", "", "secret name")
	flags.StringArrayVar(&options.vaultReadTo, "allow-read-to", []string{}, "read by role ID")
	flags.StringArrayVar(&options.vaultWriteTo, "allow-write-to", []string{}, "write by role ID")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("read-by")
	cmd.MarkFlagRequired("write-by")

	return cmd
}

func secretCreate(args []string, options vaultOptions) error {
	secret, err := readJSON(args[0])
	if err != nil {
		return err
	}

	api := vault.New(curl())
	if err := api.CreateSecret(options.secretName, options.vaultReadTo,
		options.vaultWriteTo, secret); err != nil {
		return err
	}

	return stdout(secret)
}

//
//
func vaultUpdateCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update secret",
		Long:  `Update existing secret`,
		Example: `
	privx-cli secrets update [access flags] --name <SECRET-NAME> JSON-FILE

	privx-cli secrets update [access flags] --name <SECRET-NAME>
		--allow-read-to <ROLE-ID>
		--allow-write-to <ROLE-ID>
		...
		JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return secretUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.secretName, "name", "", "secret name")
	flags.StringArrayVar(&options.vaultReadTo, "allow-read-to", []string{}, "read by role ID")
	flags.StringArrayVar(&options.vaultWriteTo, "allow-write-to", []string{}, "write by role ID")
	cmd.MarkFlagRequired("name")

	return cmd
}

func secretUpdate(options vaultOptions, args []string) error {
	secret, err := readJSON(args[0])
	if err != nil {
		return err
	}

	api := vault.New(curl())
	bag, err := api.Secret(options.secretName)
	if err != nil {
		return err
	}

	if len(options.vaultReadTo) == 0 {
		for _, ref := range bag.AllowRead {
			options.vaultReadTo = append(options.vaultReadTo, ref.ID)
		}
	}

	if len(options.vaultWriteTo) == 0 {
		for _, ref := range bag.AllowWrite {
			options.vaultWriteTo = append(options.vaultWriteTo, ref.ID)
		}
	}

	if err := api.UpdateSecret(options.secretName, options.vaultReadTo,
		options.vaultWriteTo, secret); err != nil {
		return err
	}

	return stdout(secret)
}

//
//
func secretDeleteCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete secret",
		Long:  `Delete secret from PrivX Vault. Secret Name's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli vault delete [access flags] --name <SECRET-NAME>,<SECRET-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return secretDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.secretName, "name", "", "secret name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func secretDelete(options vaultOptions) error {
	api := vault.New(curl())

	for _, name := range strings.Split(options.secretName, ",") {
		err := api.DeleteSecret(name)
		if err != nil {
			return err
		} else {
			fmt.Println(name)
		}
	}

	return nil
}

//
//
func secretMetadataShowCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Get a secrets metadata",
		Long:  `Get a secrets metadata. Secret Name's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli secrets metadata [access flags] --name <SECRET-NAME>,<SECRET-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return secretMetadataShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.secretName, "name", "", "secret name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func secretMetadataShow(options vaultOptions) error {
	api := vault.New(curl())
	secrets := []vault.Secret{}

	for _, name := range strings.Split(options.secretName, ",") {
		secret, err := api.SecretMetadata(name)
		if err != nil {
			return err
		}
		secrets = append(secrets, *secret)
	}

	return stdout(secrets)
}

//
//
func secretSearchCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for secrets",
		Long:  `Search for secrets`,
		Example: `
	privx-cli secrets search [access flags] --offset <OFFSET> --keywords "<KEYWORD>,<KEYWORD>"
	privx-cli secrets search [access flags] --limit <LIMIT> --sortkey <SORTKEY>
	privx-cli secrets search [access flags] --limit <LIMIT> --filter "shared"
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return secretSearch(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort object by name, updated, or created.")
	flags.StringVar(&options.keywords, "keywords", "", "comma or space-separated string to search in secret's names")
	flags.StringVar(&options.filter, "filter", "", "filter defines what kind of secrets to search for; it accepts one value of the following list:"+
		"{'personal', 'shared', 'accessible', 'readable', 'writable', ''}")
	return cmd
}

func secretSearch(options vaultOptions) error {
	api := vault.New(curl())

	secrets, err := api.SearchSecrets(options.offset, options.limit, options.keywords,
		options.sortkey, strings.ToUpper(options.sortdir), options.filter)
	if err != nil {
		return err
	}

	return stdout(secrets)
}

//
//
func secretSchemasShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schemas",
		Short: "Returns the defined schemas",
		Long:  `Returns the defined schemas`,
		Example: `
	privx-cli secrets schemas [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return secretSchemasShow()
		},
	}

	return cmd
}

func secretSchemasShow() error {
	api := vault.New(curl())

	schemas, err := api.VaultSchemas()
	if err != nil {
		return err
	}

	return stdout(schemas)
}

func readJSON(name string) (secret interface{}, err error) {
	file, err := os.Open(name)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &secret)
	return
}

///// commands for user secrets

// list the user secrets
func userSecretListCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "user-secrets",
		Short: "PrivX user secrets",
		Long:  `List and manage PrivX secrets owned by a specific user`,
		Example: `
	privx-cli user-secrets [access flags] --offset <OFFSET> --limit <LIMIT>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userSecretList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.ownerID, "owner-id", "", "secret's owner ID")
	cmd.MarkFlagRequired("owner-id")

	cmd.AddCommand(userSecretShowCmd())
	cmd.AddCommand(userSecretCreateCmd())
	cmd.AddCommand(userSecretUpdateCmd())
	cmd.AddCommand(userSecretDeleteCmd())
	cmd.AddCommand(userSecretMetadataShowCmd())
	return cmd
}

func userSecretList(options vaultOptions) error {
	api := vault.New(curl())
	secretID := &vault.SecretID{
		OwnerID: options.ownerID,
		Name:    "",
	}
	secrets, err := api.UserSecrets(*secretID, options.offset, options.limit)
	if err != nil {
		return err
	}

	return stdout(secrets)
}

//

// show user secret details
func userSecretShowCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get a user secret",
		Long:  `Get a user secret. Secret Name's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli user-secrets show [access flags] --name <SECRET-NAME>,<SECRET-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userSecretShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.secretName, "name", "", "secret name")
	flags.StringVar(&options.ownerID, "owner-id", "", "secret's owner ID")
	cmd.MarkFlagRequired("owner-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func userSecretShow(options vaultOptions) error {
	api := vault.New(curl())
	secrets := []vault.Secret{}
	secretID := &vault.SecretID{
		OwnerID: options.ownerID,
		Name:    "",
	}
	for _, name := range strings.Split(options.secretName, ",") {
		secretID.Name = name
		secret, err := api.UserSecret(*secretID)
		if err != nil {
			return err
		}
		secrets = append(secrets, *secret)
	}

	return stdout(secrets)
}

//

// create user secret
func userSecretCreateCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new user secret",
		Long:  `Create new user secret`,
		Example: `
	privx-cli user-secrets create [access flags] --name <SECRET-NAME>
		--allow-read-to <ROLE-ID>
		--allow-write-to <ROLE-ID>
		...
		JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userSecretCreate(args, options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.secretName, "name", "", "secret name")
	flags.StringArrayVar(&options.vaultReadTo, "allow-read-to", []string{}, "read by role ID")
	flags.StringArrayVar(&options.vaultWriteTo, "allow-write-to", []string{}, "write by role ID")
	flags.StringVar(&options.ownerID, "owner-id", "", "secret's owner ID")
	cmd.MarkFlagRequired("owner-id")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("read-by")
	cmd.MarkFlagRequired("write-by")

	return cmd
}

func userSecretCreate(args []string, options vaultOptions) error {
	secret, err := readJSON(args[0])
	if err != nil {
		return err
	}

	secretID := &vault.SecretID{
		OwnerID: options.ownerID,
		Name:    options.secretName,
	}
	api := vault.New(curl())
	if err := api.CreateUserSecret(*secretID, options.vaultReadTo,
		options.vaultWriteTo, secret); err != nil {
		return err
	}

	return stdout(secret)
}

//

// update user secret
func userSecretUpdateCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update user secret",
		Long:  `Update existing user secret`,
		Example: `
	privx-cli user-secrets update [access flags] --owner-id <SECRET-NAME> --name <SECRET-NAME> JSON-FILE

	privx-cli user-secrets update [access flags] --owner-id <SECRET-NAME> --name <SECRET-NAME>
		--allow-read-to <ROLE-ID>
		--allow-write-to <ROLE-ID>
		...
		JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userSecretUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.secretName, "name", "", "secret name")
	flags.StringArrayVar(&options.vaultReadTo, "allow-read-to", []string{}, "read by role ID")
	flags.StringArrayVar(&options.vaultWriteTo, "allow-write-to", []string{}, "write by role ID")
	flags.StringVar(&options.ownerID, "owner-id", "", "secret's owner ID")
	cmd.MarkFlagRequired("owner-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func userSecretUpdate(options vaultOptions, args []string) error {
	secret, err := readJSON(args[0])
	if err != nil {
		return err
	}

	secretID := &vault.SecretID{
		OwnerID: options.ownerID,
		Name:    options.secretName,
	}
	api := vault.New(curl())
	bag, err := api.UserSecret(*secretID)
	if err != nil {
		return err
	}

	if len(options.vaultReadTo) == 0 {
		for _, ref := range bag.AllowRead {
			options.vaultReadTo = append(options.vaultReadTo, ref.ID)
		}
	}

	if len(options.vaultWriteTo) == 0 {
		for _, ref := range bag.AllowWrite {
			options.vaultWriteTo = append(options.vaultWriteTo, ref.ID)
		}
	}

	if err := api.UpdateUserSecret(*secretID, options.vaultReadTo,
		options.vaultWriteTo, secret); err != nil {
		return err
	}

	return stdout(secret)
}

//

//show user secret metadata
func userSecretMetadataShowCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Get a user secrets metadata",
		Long:  `Get a user secrets metadata. Secret Name's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli user-secrets metadata [access flags] --owner-id <SECRET-NAME> --name <SECRET-NAME>,<SECRET-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userSecretMetadataShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.secretName, "name", "", "secret name")
	flags.StringVar(&options.ownerID, "owner-id", "", "secret's owner ID")
	cmd.MarkFlagRequired("owner-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func userSecretMetadataShow(options vaultOptions) error {
	api := vault.New(curl())
	secrets := []vault.Secret{}

	secretID := &vault.SecretID{
		OwnerID: options.ownerID,
		Name:    "",
	}

	for _, name := range strings.Split(options.secretName, ",") {

		secretID.Name = name
		secret, err := api.UserSecretMetadata(*secretID)
		if err != nil {
			return err
		}
		secrets = append(secrets, *secret)
	}

	return stdout(secrets)
}

//

// delete user secret
func userSecretDeleteCmd() *cobra.Command {
	options := vaultOptions{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete secret",
		Long:  `Delete secret from PrivX Vault. Secret Name's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli user-secret delete [access flags] --name <SECRET-NAME>,<SECRET-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return userSecretDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.secretName, "name", "", "secret name")
	flags.StringVar(&options.ownerID, "owner-id", "", "secret's owner ID")
	cmd.MarkFlagRequired("owner-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func userSecretDelete(options vaultOptions) error {
	api := vault.New(curl())

	secretID := &vault.SecretID{
		OwnerID: options.ownerID,
		Name:    "",
	}

	for _, name := range strings.Split(options.secretName, ",") {
		secretID.Name = name
		err := api.DeleteUserSecret(*secretID)
		if err != nil {
			return err
		} else {
			fmt.Println(name)
		}
	}

	return nil
}

//
