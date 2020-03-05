package main

import (
	"fmt"
	xhttp "github.com/cbrgm/rbaclister/http"
	"github.com/ghodss/yaml"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type appConf struct {
	ApiServerAddr string
	Output string
}

var (
	appConfig = appConf{}
)

func main() {
	app := cli.NewApp()
	app.Name = "rbaclister"

	app.Action = action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "apiserverUrl",
			Usage: "kubernetes apiserver url",
			Destination: &appConfig.ApiServerAddr,
			Value: "127.0.0.1:8001",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "output directory for generated files",
			Destination: &appConfig.Output,
			Value: filepath.Join("./output"),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}

func action(c *cli.Context) error {
	var remoteUrl = fmt.Sprintf("http://%s", appConfig.ApiServerAddr)
	var baseDir = appConfig.Output

	client := xhttp.NewClient(remoteUrl)

	resourceLists, err := client.GetApiResourceLists()
	if err != nil {
		return err
	}

	for _, list := range resourceLists {
		dirName := cleanupString(list.GroupVersion)
		err := makeDirectory(baseDir, dirName)
		if err != nil {
			return err
		}

		var rules = list.AsRules()
		for _, rule := range rules {
			b, err := yaml.Marshal(rule)
			if err != nil {
				return err
			}

			err = writeToFile(filepath.Join(baseDir, dirName), rule.Resources[0], b)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// makeDirectory creates a new directory at the given base path with dirname.
// If the directory already exists, no new directory will be created.
func makeDirectory(base, dirName string) error {
	path := filepath.Join(base, dirName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// writeToFile writes byte data to the given file with filename at dirName
func writeToFile(dirName string, fileName string, data []byte) error {
	filepath := filepath.Join(dirName, cleanupString(fileName) + ".yaml")
	err := ioutil.WriteFile(filepath, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func cleanupString(str string) string {
	return strings.ReplaceAll(str, "/", "_")
}
