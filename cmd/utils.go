//
// Copyright (c) 2022 SSH Communications Security Inc.
//
// All rights reserved.
//

package cmd

import (
	"encoding/json"
	"io"
	"os"
)

func decodeJSON(name string, object interface{}) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &object)
	if err != nil {
		return err
	}

	return nil
}
