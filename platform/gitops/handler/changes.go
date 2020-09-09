package handler

import (
	"fmt"
	"path/filepath"
)

// determineChanges takes a slice of commmits and returns the status of the services, e.g. "test" =>
// created, "test/api" => created, "foo" => modified.
func determineChanges(commits []commit) map[string]changeType {
	// services contains the directories and the status, e.g. "test" => created. It can also contain
	// sub-directories, e.g. "foo/handler" or "foo/api"
	result := make(map[string]changeType)

	// check for addition / deletion of main.go files which indicates a service was created or deleted
	for _, commit := range commits {
		for _, file := range commit.Added {
			fmt.Println("file", file)

			if filepath.Base(file) == "main.go" {
				dir := filepath.Dir(file)
				if _, ok := result[dir]; !ok {
					result[dir] = created
				}
			}
		}

		for _, file := range commit.Removed {
			if filepath.Base(file) == "main.go" {
				dir := filepath.Dir(file)
				if _, ok := result[dir]; !ok {
					result[dir] = deleted
				}
			}
		}
	}

	return result
}
