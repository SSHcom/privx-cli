//
// Copyright (c) 2023 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"errors"

	authApi "github.com/SSHcom/privx-sdk-go/api/auth"
	"github.com/spf13/cobra"
)

type sessionStorageOptions struct {
	userID    string
	sourceID  string
	sessionID string
	sortkey   string
	sortdir   string
	offset    int
	limit     int
}

func init() {
	rootCmd.AddCommand(sessionStorageCmd())
}

//
//
func sessionStorageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sessions",
		Short: "List and manage sessions",
		Long:  `List and manage sessions. Only applicable for user sessions, and not other types of sessions.`,
		Example: `
	privx-cli sessions [access flags]
		`,
		SilenceUsage: true,
	}

	cmd.AddCommand(sessionsShowCmd())
	cmd.AddCommand(sessionsSearchCmd())
	cmd.AddCommand(terminateSessionCmd())
	cmd.AddCommand(terminateUserSessionsCmd())

	return cmd
}

//
//
func sessionsShowCmd() *cobra.Command {
	options := sessionStorageOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get sessions by userID or sourceID",
		Long:  `Get sessions by userID or sourceID. Provide either the sourceID or userID`,
		Example: `
	privx-cli sessions show [access flags] --user-id <USER-ID>
	privx-cli sessions show [access flags] --source-id <SOURCE-ID>
		`,
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate the flags, either user-id or source-id is required
			if options.sourceID == "" && options.userID == "" {
				return errors.New("either user-id or source-id is required")
			} else if options.sourceID != "" && options.userID != "" {
				return errors.New("only one of user-id or source-id is allowed")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return sessionsShow(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "user-id", "", "User ID")
	flags.StringVar(&options.sourceID, "source-id", "", "Source ID")
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "expires", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "ASC", "sort direction, ASC or DESC")

	return cmd
}

func sessionsShow(options sessionStorageOptions) error {
	api := authApi.New(curl())
	if options.userID != "" {
		userSessions, err := api.UserSessions(options.offset, options.limit,
			options.sortkey, options.sortdir, options.userID)
		if err != nil {
			return err
		}
		return stdout(userSessions)
	} else if options.sourceID != "" {
		sourceSessions, err := api.SourceSessions(options.offset, options.limit,
			options.sortkey, options.sortdir, options.sourceID)
		if err != nil {
			return err
		}
		return stdout(sourceSessions)
	}
	return nil
}

//
//
func sessionsSearchCmd() *cobra.Command {
	options := sessionStorageOptions{}

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for sessions",
		Long:  `Search for sessions. Search with keywords, userid or type`,
		Example: `
	privx-cli sessions search [access flags] JSON-FILE
		`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return sessionsSearch(options, args)
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&options.offset, "offset", 0, "where to start fetching the items")
	flags.IntVar(&options.limit, "limit", 50, "number of items to return")
	flags.StringVar(&options.sortkey, "sortkey", "expires", "sort by specific object property")
	flags.StringVar(&options.sortdir, "sortdir", "ASC", "sort direction, ASC or DESC")

	return cmd
}

func sessionsSearch(options sessionStorageOptions, args []string) error {
	var searchParams authApi.SearchParams
	api := authApi.New(curl())

	err := decodeJSON(args[0], &searchParams)
	if err != nil {
		return err
	}

	sessions, err := api.SearchSessions(options.offset, options.limit,
		options.sortkey, options.sortdir, &searchParams)
	if err != nil {
		return err
	}

	return stdout(sessions)
}

//
//
func terminateSessionCmd() *cobra.Command {
	options := sessionStorageOptions{}

	cmd := &cobra.Command{
		Use:   "terminate",
		Short: "Terminate single session by ID",
		Long:  `Terminate single session by sessionID`,
		Example: `
	privx-cli sessions terminate [access flags] --id <SESSION-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return terminateSession(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.sessionID, "id", "", "Session ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func terminateSession(options sessionStorageOptions) error {
	api := authApi.New(curl())

	return api.TerminateSession(options.sessionID)
}

//
//
func terminateUserSessionsCmd() *cobra.Command {
	options := sessionStorageOptions{}

	cmd := &cobra.Command{
		Use:   "terminate-all",
		Short: "Terminate all sessions for a user",
		Long:  `Terminate all sessions for a user by userID`,
		Example: `
	privx-cli sessions terminate-all [access flags] --id <USER-ID>
		`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return terminateUserSessions(options)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.userID, "id", "", "User ID")
	cmd.MarkFlagRequired("id")

	return cmd
}

func terminateUserSessions(options sessionStorageOptions) error {
	api := authApi.New(curl())

	return api.TerminateUserSessions(options.userID)
}
