package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration of nebula",
	Run: func(cmd *cobra.Command, args []string) {
		agent, err := NewClient(l, config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer agent.Close()

		l.Infoln("Generate config dir")

		configTemplatePath, err := cmd.Flags().GetString("template")
		if err != nil {
			l.WithError(err).Fatalln("failed to get template")
		}

		var nebulaConfig map[interface{}]interface{}
		bs, err := ioutil.ReadFile(configTemplatePath)
		if err != nil {
			panic(err)
		}
		if err := yaml.Unmarshal(bs, &nebulaConfig); err != nil {
			panic(err)
		}

		if _, ok := nebulaConfig["pki"]; !ok {
			nebulaConfig["pki"] = make(map[string]interface{})
		}
		pki := nebulaConfig["pki"].(map[string]interface{})

		ca, err := cmd.Flags().GetString("ca")
		if err != nil {
			l.WithError(err).Fatalln("failed to get ca")
		}
		pki["ca"], err = filepath.Abs(ca)
		if err != nil {
			l.WithError(err).Fatalf("failed to get absolut path for `%s`\n", ca)
		}
		if ok, _ := fileExists(pki["ca"].(string)); !ok {
			l.WithError(err).Fatalf("file not found `%s`\n", pki["ca"].(string))
		}

		crt, err := cmd.Flags().GetString("crt")
		if err != nil {
			l.WithError(err).Fatalln("failed to get crt")
		}
		pki["cert"], err = filepath.Abs(crt)
		if err != nil {
			l.WithError(err).Fatalf("failed to get absolut path for `%s`\n", crt)
		}
		if ok, _ := fileExists(pki["cert"].(string)); !ok {
			l.WithError(err).Fatalf("file not found `%s`\n", pki["cert"].(string))
		}

		key, err := cmd.Flags().GetString("key")
		if err != nil {
			l.WithError(err).Fatalln("failed to get key")
		}
		pki["key"], err = filepath.Abs(key)
		if err != nil {
			l.WithError(err).Fatalf("failed to get absolut path for `%s`\n", key)
		}
		if ok, _ := fileExists(pki["key"].(string)); !ok {
			l.WithError(err).Fatalf("file not found `%s`\n", pki["key"].(string))
		}

		configPath, err := cmd.Flags().GetString("output-config-dir")
		if err != nil {
			l.WithError(err).Fatalln("failed to get output-config-dir")
		}

		_, err = os.Stat(configPath)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(configPath, 0700)
				if err != nil {
					l.WithError(err).Fatalf("failed to create directory `%s`\n", configPath)
				}
			} else {
				l.WithError(err).Fatalf("failed to create directory `%s`\n", configPath)
			}
		}

		bs, err = yaml.Marshal(nebulaConfig)
		if err != nil {
			panic(err)
		}
		configFilePath := filepath.Join(configPath, "config.yml")
		if err := ioutil.WriteFile(configFilePath, bs, 0600); err != nil {
			l.WithError(err).Fatalf("failed to write config file to `%s`\n", configFilePath)
		}
	},
}

func init() {
	configCmd.Flags().String("ca", AgentNebulaCaPath, "Path to ca.crt")
	configCmd.Flags().String("crt", AgentNebulaCrtPath, "Path to signed pub client .crt")
	configCmd.Flags().String("key", AgentNebulaKeyPath, "Path to client .key")
	configCmd.Flags().String("template", "nebula-config.yml", "Path to nebula config template file")
	configCmd.Flags().String("output-config-dir", "nebula", "Path to output generated nebula config")
}
