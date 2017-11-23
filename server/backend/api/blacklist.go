package main

import (
	"bufio"
	"os"
	"path"
)

var (
	blacklistFiles = map[string]string{
		"fromCampus": "blacklist/from/campus.txt",
		"fromClass":  "blacklist/from/class.txt",
		"fromPeriod": "blacklist/from/period.txt",
		"toCampus":   "blacklist/to/campus.txt",
		"toClass":    "blacklist/to/class.txt",
		"toPeriod":   "blacklist/to/period.txt",
	}
)

// loadBlacklists loads the blacklists for students' transfer.
func loadBlacklists() (map[string][]string, error) {
	blacklists := map[string][]string{}

	for k, v := range blacklistFiles {
		fileName := path.Join(currentDir, v)
		f, err := os.Open(fileName)
		if err != nil {
			return blacklists, err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			blacklists[k] = append(blacklists[k], scanner.Text())
		}

		if err = scanner.Err(); err != nil {
			return blacklists, err
		}

	}
	return blacklists, nil
}
