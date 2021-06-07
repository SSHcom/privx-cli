//
// Copyright (c) 2021 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"fmt"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/hoststore"
	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/spf13/cobra"
)

var (
	tagType string
)

func init() {
	rootCmd.AddCommand(tagListCmd)
	tagListCmd.Flags().IntVar(&offset, "offset", 0, "where to start fetching the items")
	tagListCmd.Flags().IntVar(&limit, "limit", 50, "number of items to return")
	tagListCmd.Flags().StringVar(&query, "query", "", "query string matches the tags")
	tagListCmd.Flags().StringVar(&sortdir, "sortdir", "", "sort direction, ASC or DESC (default ASC)")
	tagListCmd.Flags().StringVar(&tagType, "type", "", "choose the tag type, user or host")
	tagListCmd.MarkFlagRequired("type")
}

//
//
var tagListCmd = &cobra.Command{
	Use:   "tags",
	Short: "User | Host tags",
	Long:  `Get privx user or host tags`,
	Example: `
privx-cli tags [access flags] --type user
privx-cli tags [access flags] --type host --sortdir DESC
privx-cli tags [access flags] --type host --query TAG
privx-cli tags [access flags] --type user --offset OFFSET --limit LIMIT
	`,
	SilenceUsage: true,
	RunE:         tagList,
}

func tagList(cmd *cobra.Command, args []string) error {
	if tagType == "user" {
		userTags()
	} else if tagType == "host" {
		hostTags()
	} else {
		return fmt.Errorf("tag type does not exist: %s", tagType)
	}

	return nil
}

func userTags() error {
	api := userstore.New(curl())

	tags, err := api.LocalUserTags(offset, limit, strings.ToUpper(sortdir), query)
	if err != nil {
		return err
	}

	return stdout(tags)
}

func hostTags() error {
	api := hoststore.New(curl())

	tags, err := api.HostTags(offset, limit, strings.ToUpper(sortdir), query)
	if err != nil {
		return err
	}

	return stdout(tags)
}
