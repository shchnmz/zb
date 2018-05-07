package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/northbright/pathhelper"
	"github.com/shchnmz/zb"
)

var (
	currentDir, configFile string
)

func main() {
	var (
		err       error
		blacklist *zb.Blacklist
	)

	defer func() {
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	if blacklist, err = loadBlacklistFromTXT(); err != nil {
		return
	}

	fmt.Printf("Generating blacklist.json ...\n")

	if err = writeJSON(blacklist); err != nil {
		return
	}

	fmt.Printf("Done.\n")
}

// init initializes path variables.
func init() {
	currentDir, _ = pathhelper.GetCurrentExecDir()
	configFile = path.Join(currentDir, "config.json")
}

// load blacklist data from TXT files.
func loadBlacklistFromTXT() (*zb.Blacklist, error) {
	var blacklist = zb.Blacklist{map[string][]string{}}

	for t, _ := range zb.BlacklistTypes {
		fileName := fmt.Sprintf("%v.txt", t)
		p := path.Join(currentDir, "blacklist", fileName)

		// Skip no-existent blacklist txt
		if _, err := os.Stat(p); os.IsNotExist(err) {
			continue
		}

		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			// Trim line breaker and spaces.
			line := strings.Trim(scanner.Text(), "\n")
			line = strings.Trim(line, " ")

			if line != "" {
				blacklist.List[t] = append(blacklist.List[t], scanner.Text())
			}
		}

		if err = scanner.Err(); err != nil {
			return nil, err
		}
	}
	return &blacklist, nil
}

func writeJSON(blacklist *zb.Blacklist) error {
	p := path.Join(currentDir, "blacklist.json")

	f, err := os.Create(p)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)
	// Pretty print JSON: set ident to 4 spaces.
	encoder.SetIndent("", "    ")
	if err = encoder.Encode(blacklist); err != nil {
		return err
	}
	return nil
}
