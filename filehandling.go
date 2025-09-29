package main

import (
	"fmt"
	"mime"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/widget"
)

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			fmt.Printf("failed to get file %s: %s\n", path, err)
			return false
		}
	}
	return true
}

func getConfigPath() string {
	userConfigPath, err := app.DataDir()
	if err != nil {
		fmt.Printf("failed to get config dir: %s\n", err)
		os.Exit(2)
	}

	return userConfigPath + "/numbercounter"
}

func getOrCreateConfigDir() string {
	configPath := getConfigPath()
	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(configPath, 0766)
			if err != nil {
				fmt.Printf("failed to create config dir: %s\n", err)
				os.Exit(3)
			}
		} else {
			fmt.Printf("failed to get config dir: %s\n", err)
			os.Exit(4)
		}
	}

	return configPath
}

func writeToCache(configPath, folderPath string) {
	cachePath := fmt.Sprintf(
		"%s/cache",
		configPath,
	)

	err := os.WriteFile(cachePath, []byte(folderPath), 0766)
	if err != nil {
		fmt.Printf("Failed to write to cache: %s\n", err)
		return
	}
}

func readFromCache(configPath string) error {
	cachePath := fmt.Sprintf(
		"%s/cache",
		configPath,
	)

	if !fileExist(cachePath) {
		return nil
	}

	file, err := os.ReadFile(cachePath)
	if err != nil {
		fmt.Printf("failed to read from cache %v\n", err)
		return fmt.Errorf("Failed to read cache: %w", err)
	}

	folderPath = string(file)
	readFolderForCounters()

	return nil
}

func readFolderForCounters() {
	listOfFiles, err := os.ReadDir(folderPath)
	if err != nil {
		fmt.Printf("listOfFiles ERR %#v\n", err.Error())
	}
	counters = []*Counter{}
	for _, v := range listOfFiles {
		if v.IsDir() {
			continue
		}
		name := v.Name()
		splitName := strings.Split(name, ".")
		extension := mime.TypeByExtension("." + splitName[len(splitName)-1])
		if !strings.Contains(extension, "text/") {
			continue
		}

		data, err := os.ReadFile(fmt.Sprintf("%s/%s", folderPath, name))

		if err != nil {
			fmt.Printf("failed to read file %s %#v\n", name, err)
		}
		value, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err != nil {
			fmt.Printf("failed to read from: %s, value: %s, err: %#v\n", folderPath + "/" + name, string(data), err)
		}

		counters = append(counters, &Counter{
			value:           value,
			fileName:        &name,
			incrementButton: new(widget.Clickable),
			decrementButton: new(widget.Clickable),
		})
	}
}

func printFileName(filename string) string {
	return fmt.Sprintf(
		"%s: ",
		strings.Replace(
			strings.ReplaceAll(filename, "_", " "),
			".txt",
			"",
			1,
		),
	)

}
