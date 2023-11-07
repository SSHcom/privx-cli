//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/connectionmanager"
	"github.com/spf13/cobra"
)

type connectionOptions struct {
	channID    string
	fileID     string
	connID     string
	roleID     string
	userID     string
	hostID     string
	fileName   string
	sortkey    string
	sortdir    string
	format     string
	filter     string
	offset     int
	limit      int
	fuzzyCount bool
	force      bool
}

type uebaOptions struct {
	datasetID                 string
	logs                      bool
	bin_count                 int
	set_active_after_training bool
}

func init() {
	rootCmd.AddCommand(connectionListCmd())
	rootCmd.AddCommand(uebaCmd())

}

func connectionListCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "connections",
		Short: "List and manage connections",
		Long:  `List and manage connections`,
		Example: `
	privx-cli connections [access flags] --offset <OFFSET> --sortkey <SORTKEY>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionList(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")
	flags.BoolVarP(&options.fuzzyCount, "fuzzycount", "", false, "return a fuzzy total count instead of exact total count")

	cmd.AddCommand(connectionSearchCmd())
	cmd.AddCommand(connectionShowCmd())
	cmd.AddCommand(storedFileDownloadCmd())
	cmd.AddCommand(trailLogDownloadCmd())
	cmd.AddCommand(accessRoleListCmd())
	cmd.AddCommand(connectionAccessRoleGrantCmd())
	cmd.AddCommand(connectionAccessRoleRevokeCmd())
	cmd.AddCommand(connectionTerminateCmd())

	return cmd
}

func connectionList(options connectionOptions) error {
	api := connectionmanager.New(curl())

	conn, err := api.Connections(options.offset, options.limit,
		options.sortkey, options.sortdir, options.fuzzyCount)
	if err != nil {
		return err
	}

	return stdout(conn)
}

func connectionSearchCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for connections",
		Long:  `Search for connections`,
		Example: `
	privx-cli connections search [access flags] --offset <OFFSET> --sortkey <SORTKEY>
	privx-cli connections search [access flags] --limit <LIMIT> JSON-FILE
		`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionSearch(options, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")

	return cmd
}

func connectionSearch(options connectionOptions, args []string) error {
	var searchObject connectionmanager.ConnectionSearch
	api := connectionmanager.New(curl())

	if len(args) == 1 {
		err := decodeJSON(args[0], &searchObject)
		if err != nil {
			return err
		}
	}

	conn, err := api.SearchConnections(options.offset, options.limit, options.sortdir,
		options.sortkey, options.fuzzyCount, searchObject)
	if err != nil {
		return err
	}

	return stdout(conn)
}

func connectionShowCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get connection by ID",
		Long:  `Get connection by ID. Connection ID's are separated by commas when using multiple values, see example`,
		Example: `
	privx-cli connections show [access flags] --conn-id <CONN-ID>,<CONN-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	cmd.MarkFlagRequired("conn-id")

	return cmd
}

func connectionShow(options connectionOptions) error {
	api := connectionmanager.New(curl())
	conns := []connectionmanager.Connection{}

	for _, id := range strings.Split(options.connID, ",") {
		conn, err := api.Connection(id)
		if err != nil {
			return err
		}
		conns = append(conns, *conn)
	}

	return stdout(conns)
}

func storedFileDownloadCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "download-file",
		Short: "Download trail stored file",
		Long:  `Download trail stored file`,
		Example: `
	privx-cli connections download-file [access flags] --conn-id <CONN-ID> --file-id <FILE-ID> --channel-id <CHANNEL-ID> --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return storedFileDownload(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	flags.StringVar(&options.channID, "channel-id", "", "channel ID")
	flags.StringVar(&options.fileID, "file-id", "", "file ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	cmd.MarkFlagRequired("conn-id")
	cmd.MarkFlagRequired("channel-id")
	cmd.MarkFlagRequired("file-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func storedFileDownload(options connectionOptions) error {
	api := connectionmanager.New(curl())

	sessionID, err := api.CreateSessionIDFileDownload(options.connID, options.channID, options.fileID)
	if err != nil {
		return err
	}

	err = api.DownloadStoredFile(options.connID, options.channID, options.fileID,
		sessionID, options.fileName)
	if err != nil {
		return err
	}

	return nil
}

func trailLogDownloadCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "download-log",
		Short: "Download trail log",
		Long:  `Download trail log`,
		Example: `
	privx-cli connections download-log [access flags] --conn-id <CONN-ID> --channel-id <CHANNEL-ID> --sid <SESSION-ID> --name <FILE-NAME>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return trailLogDownload(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	flags.StringVar(&options.channID, "channel-id", "", "channel ID")
	flags.StringVar(&options.fileName, "name", "", "file name")
	flags.StringVar(&options.format, "format", "", "trail log format, json or hex")
	flags.StringVar(&options.filter, "filter", "", "trail log event filter")
	cmd.MarkFlagRequired("conn-id")
	cmd.MarkFlagRequired("channel-id")
	cmd.MarkFlagRequired("name")

	return cmd
}

func trailLogDownload(options connectionOptions) error {
	api := connectionmanager.New(curl())

	sessionID, err := api.CreateSessionIDTrailLog(options.connID, options.channID)
	if err != nil {
		return err
	}

	err = api.DownloadTrailLog(options.connID, options.channID, sessionID,
		options.format, options.filter, options.fileName)
	if err != nil {
		return err
	}

	return nil
}

func accessRoleListCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "access-roles",
		Short: "List access roles for a connection",
		Long:  `List access roles for a connection`,
		Example: `
	privx-cli connections access-roles [access flags] --conn-id <CONN-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return accessRoleList(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	cmd.MarkFlagRequired("conn-id")

	return cmd
}

func accessRoleList(options connectionOptions) error {
	api := connectionmanager.New(curl())

	roles, err := api.AccessRoles(options.connID)
	if err != nil {
		return err
	}

	return stdout(roles)
}

func connectionAccessRoleGrantCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "grant-access-role",
		Short: "Add the role to a list of roles that can explicitly access this connection data",
		Long:  `Add the role to a list of roles that can explicitly access this connection data`,
		Example: `
	privx-cli connections grant-role [access flags] --conn-id <CONN-ID> --role-id <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionAccessRoleGrant(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	flags.StringVar(&options.roleID, "role-id", "", "role ID")
	cmd.MarkFlagRequired("conn-id")
	cmd.MarkFlagRequired("role-id")

	return cmd
}

func connectionAccessRoleGrant(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.GrantAccessRoleToConnection(options.connID, options.roleID)
	if err != nil {
		return err
	}

	return nil
}

func connectionAccessRoleRevokeCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "revoke-access-role",
		Short: "Remove the role from a list of roles that can explicitly access this connection data",
		Long:  `Remove the role from a list of roles that can explicitly access this connection data`,
		Example: `
	privx-cli connections revoke-access-role [access flags] --conn-id <CONN-ID> --role-id <ROLE-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionAccessRoleRevoke(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "connection ID")
	flags.StringVar(&options.roleID, "role-id", "", "role ID")
	flags.BoolVarP(&options.force, "force", "f", false, "force command")
	cmd.MarkFlagRequired("role-id")

	return cmd
}

func connectionAccessRoleRevoke(options connectionOptions) error {
	api := connectionmanager.New(curl())

	if options.connID != "" {
		err := api.RevokeAccessRoleFromConnection(options.connID, options.roleID)
		if err != nil {
			return err
		}
	} else {
		if !options.force {
			fmt.Fprintln(os.Stderr, "Error: this action will revoke data access rights from this role to all connections.\nUse --force | -f flag to revoke data access to all connections or use --conn-id to revoke data access to a specific connection")
			os.Exit(1)
		} else {
			err := api.RevokeAccessRoleFromAllConnections(options.roleID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func connectionTerminateCmd() *cobra.Command {
	options := connectionOptions{}

	cmd := &cobra.Command{
		Use:   "terminate",
		Short: "Terminate connection by ID",
		Long:  `Terminate connection by ID`,
		Example: `
	privx-cli connections terminate [access flags] --conn-id <CONN-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectionTerminate(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.connID, "conn-id", "", "terminate connection by ID")
	flags.StringVar(&options.hostID, "by-target", "", "terminate connection by host ID")
	flags.StringVar(&options.userID, "by-user", "", "terminate connection by user ID")

	return cmd
}

func connectionTerminate(options connectionOptions) error {
	if (options == connectionOptions{}) {
		fmt.Println("Specify at least one flag for the termination type of the connection: --conn-id, --by-target or --by-user")
	} else if options.hostID != "" {
		terminateConnectionByTargerHost(options)
	} else if options.userID != "" {
		terminateConnectionByUser(options)
	} else if options.connID != "" {
		terminateConnectionByConnection(options)
	}

	return nil
}

func terminateConnectionByConnection(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.TerminateConnection(options.connID)
	if err != nil {
		return err
	}

	return nil
}

func terminateConnectionByTargerHost(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.TerminateConnectionsByTargetHost(options.hostID)
	if err != nil {
		return err
	}

	return nil
}

func terminateConnectionByUser(options connectionOptions) error {
	api := connectionmanager.New(curl())

	err := api.TerminateConnectionsByUser(options.userID)
	if err != nil {
		return err
	}

	return nil
}

func uebaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ueba",
		Short: "Ueba Related Commands",
		Long:  `Ueba Related Commands`,
		Example: `
	privx-cli connections ueba [access flags]
		`,
		SilenceUsage: true,
	}

	cmd.AddCommand(uebaConfigCmd())
	cmd.AddCommand(uebaAnomalySettingsCmd())
	cmd.AddCommand(uebaStartAnalyzingCmd())
	cmd.AddCommand(uebaStopAnalyzingCmd())
	cmd.AddCommand(uebaScriptDownloadCmd())
	cmd.AddCommand(uebaDatasetListCmd())
	cmd.AddCommand(uebaStatusCmd())
	cmd.AddCommand(uebaInternalStatusCmd())
	return cmd
}

func uebaConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "List and manage Ueba Configurations",
		Long:  `List and manage Ueba Configurations`,
		Example: `
	privx-cli ueba config [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaConfig()
		},
	}

	cmd.AddCommand(uebaConfigSetCmd())

	return cmd
}

func uebaConfig() error {
	api := connectionmanager.New(curl())

	configs, err := api.UebaConfigurations()
	if err != nil {
		return err
	}

	return stdout(configs)
}

func uebaConfigSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set ueba config",
		Long:  `Set ueba config. Requires path to json file with config passed as args`,
		Example: `
	privx-cli ueba config set [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaConfigSet(args)
		},
	}

	return cmd
}

func uebaConfigSet(args []string) error {
	var newConfig connectionmanager.UebaConfigurations
	api := connectionmanager.New(curl())

	err := decodeJSON(args[0], &newConfig)
	if err != nil {
		return err
	}

	return api.SetUebaConfigurations(&newConfig)
}

func uebaAnomalySettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "anomaly-settings",
		Short: "Get and create Ueba Anomaly Settings",
		Long:  `Get and create Ueba Anomaly Settings`,
		Example: `
	privx-cli ueba anomaly-settings [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaAnomalySettings(args)
		},
	}

	cmd.AddCommand(uebaAnomalySettingsCreateCmd())

	return cmd
}

func uebaAnomalySettings(args []string) error {
	api := connectionmanager.New(curl())

	settings, err := api.UebaAnomalySettings()
	if err != nil {
		return err
	}

	return stdout(settings)
}

func uebaAnomalySettingsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Ueba Anomaly Settings",
		Long:  `Create Ueba Anomaly Settings`,
		Example: `
	privx-cli ueba anomaly-settings create [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaAnomalySettingsCreate(args)
		},
	}

	return cmd
}

func uebaAnomalySettingsCreate(args []string) error {
	var newSettings connectionmanager.UebaAnomalySettings
	api := connectionmanager.New(curl())
	err := decodeJSON(args[0], &newSettings)
	if err != nil {
		return err
	}

	return api.CreateAnomalySettings(newSettings)
}

func uebaStartAnalyzingCmd() *cobra.Command {
	options := uebaOptions{}
	cmd := &cobra.Command{
		Use:   "start-analysis",
		Short: "Start Ueba analyzing",
		Long:  `Start analyzing connections with a saved dataset. Fails if training not done, has not finished or failed.`,
		Example: `
	privx-cli ueba start-analysis [access flags] --id <DATASET-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaStartAnalyzing(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.datasetID, "id", "", "dataset ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func uebaStartAnalyzing(options uebaOptions) error {
	api := connectionmanager.New(curl())
	return api.StartAnalyzing(options.datasetID)
}

func uebaStopAnalyzingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-analysis",
		Short: "Stop Ueba analyzing",
		Long:  `Stop analyzing connections with a saved dataset. Fails if training not done, has not finished or failed.`,
		Example: `
	privx-cli ueba stop-analysis [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaStopAnalyzing()
		},
	}

	return cmd
}

func uebaStopAnalyzing() error {
	api := connectionmanager.New(curl())
	return api.StopAnalyzing()
}

func uebaScriptDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-script",
		Short: "Download Ueba setup script.",
		Long:  `Download Ueba setup script.`,
		Example: `
	privx-cli ueba download-script [access flags]
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaScriptDownload()
		},
	}

	return cmd
}

func uebaScriptDownload() error {
	api := connectionmanager.New(curl())
	sessionIdstruct, err := api.CreateIdForUebaScript()
	if err != nil {
		return err
	}
	return api.DownloadUebaScript(sessionIdstruct.ID)

}

func uebaDatasetListCmd() *cobra.Command {
	options := uebaOptions{}
	cmd := &cobra.Command{
		Use:   "datasets",
		Short: "List and manage Ueba Datasets",
		Long:  `List and manage Ueba Datasets`,
		Example: `
	privx-cli ueba datasets [access flags]
	privx-cli ueba datasets [access flags] --logs false --bin-count 50
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaDatasets(options)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.bin_count, "bin-count", 50, "how many bins from training history")
	flags.BoolVarP(&options.logs, "logs", "l", false, "add pandas and tensorflow log prints")

	cmd.AddCommand(uebaDatasetCreateCmd())
	cmd.AddCommand(uebaDatasetShowCmd())
	cmd.AddCommand(uebaDatasetUpdateCmd())
	cmd.AddCommand(uebaDatasetDeleteCmd())
	cmd.AddCommand(uebaDatasetTrainCmd())
	cmd.AddCommand(uebaConnectionCountsCmd())

	return cmd
}

func uebaDatasets(options uebaOptions) error {
	api := connectionmanager.New(curl())

	datasets, err := api.UebaDatasets(options.logs, options.bin_count)
	if err != nil {
		return err
	}

	return stdout(datasets)
}

func uebaDatasetCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new dataset definition.",
		Long:  `Create new dataset definition.`,
		Example: `
	privx-cli ueba datasets create [access flags] JSON-FILE
			`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaDatasetCreate(args)
		},
	}

	return cmd
}

func uebaDatasetCreate(args []string) error {
	var datasetBodyParam connectionmanager.DatasetBodyParam
	api := connectionmanager.New(curl())

	err := decodeJSON(args[0], &datasetBodyParam)
	if err != nil {
		return err
	}

	datasetID, err := api.CreateUebaDataset(datasetBodyParam)
	if err != nil {
		return err
	}

	return stdout(datasetID)
}

func uebaDatasetShowCmd() *cobra.Command {
	options := uebaOptions{}
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get dataset by ID",
		Long:  `Get dataset by ID, possibility to filter training history`,
		Example: `
	privx-cli ueba datasets show [access flags] --id <DATASET-ID>
	privx-cli ueba datasets show [access flags] --id <DATASET-ID> --logs false --bin-count 50
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaDatasetShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.datasetID, "id", "", "dataset ID")
	flags.IntVar(&options.bin_count, "bin-count", 50, "how many bins from training history")
	flags.BoolVarP(&options.logs, "logs", "l", false, "add pandas and tensorflow log prints")
	cmd.MarkFlagRequired("id")

	return cmd
}

func uebaDatasetShow(options uebaOptions) error {
	api := connectionmanager.New(curl())

	dataset, err := api.UebaDataset(options.logs, options.bin_count, options.datasetID)
	if err != nil {
		return err
	}

	return stdout(dataset)
}

func uebaDatasetUpdateCmd() *cobra.Command {
	options := uebaOptions{}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update dataset",
		Long:  `Update dataset. Note this will cause backend to empty training history and delete trained weights in ueba machine.`,
		Example: `
	privx-cli ueba datasets update [access flags] --id <DATASET-ID> JSON-FILE
			`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaDatasetUpdate(options, args)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.datasetID, "id", "", "dataset ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func uebaDatasetUpdate(options uebaOptions, args []string) error {
	var datasetBodyParam connectionmanager.DatasetBodyParam
	api := connectionmanager.New(curl())

	err := decodeJSON(args[0], &datasetBodyParam)
	if err != nil {
		return err
	}

	return api.UpdateUebaDataset(datasetBodyParam, options.datasetID)
}

func uebaDatasetDeleteCmd() *cobra.Command {
	options := uebaOptions{}
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete dataset",
		Long:  `Delete dataset`,
		Example: `
	privx-cli ueba datasets delete [access flags] --id <DATASET-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaDatasetDelete(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.datasetID, "id", "", "dataset ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func uebaDatasetDelete(options uebaOptions) error {
	api := connectionmanager.New(curl())

	err := api.DeleteUebaDataset(options.datasetID)
	if err != nil {
		return err
	}

	return err
}

func uebaDatasetTrainCmd() *cobra.Command {
	options := uebaOptions{}
	cmd := &cobra.Command{
		Use:   "train",
		Short: "Train dataset",
		Long:  `Train dataset`,
		Example: `
	privx-cli ueba datasets train [access flags] --id <DATASET-ID> --set-active <SET_ACTIVE>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaDatasetTrain(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.datasetID, "id", "", "dataset ID")
	flags.BoolVarP(&options.set_active_after_training, "set-active", "a", false, "dataset ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func uebaDatasetTrain(options uebaOptions) error {
	api := connectionmanager.New(curl())

	connectionCount, err := api.TrainUebaDataset(options.datasetID, options.set_active_after_training)
	if err != nil {
		return err
	}

	return stdout(connectionCount)
}

func uebaConnectionCountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connection-count",
		Short: "Get number of connections for dataset with given parameters.",
		Long:  `Get number of connections for dataset with given parameters. All connections, if json empty in body.`,
		Example: `
	privx-cli ueba datasets connection-count [access flags] JSON-FILE
		`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaConnectionCounts(args)
		},
	}

	return cmd
}

func uebaConnectionCounts(args []string) error {
	var timeRange connectionmanager.TimeRange

	api := connectionmanager.New(curl())

	if len(args) == 1 {
		err := decodeJSON(args[0], &timeRange)
		if err != nil {
			return err
		}
	}
	count, err := api.ConnectionCounts(timeRange)
	if err != nil {
		return err
	}

	return stdout(count)
}

func uebaStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get Ueba microservice Status",
		Long:  `Get Ueba microservice Status`,
		Example: `
	privx-cli ueba status [access flags]
			`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaStatus()
		},
	}

	return cmd
}

func uebaStatus() error {
	api := connectionmanager.New(curl())
	status, err := api.UebaStatus()
	if err != nil {
		return err
	}

	return stdout(status)
}

func uebaInternalStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "internal-status",
		Short: "Get Ueba microservice intenal status",
		Long:  `Get Ueba microservice internal status`,
		Example: `
	privx-cli ueba internal-status [access flags]
			`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uebaInternalStatus()
		},
	}

	return cmd
}

func uebaInternalStatus() error {
	var status connectionmanager.UebaInternalStatus
	api := connectionmanager.New(curl())
	status, err := api.UebaInternalStatus()
	if err != nil {
		return err
	}

	return stdout(status)
}
