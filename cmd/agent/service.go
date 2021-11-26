package main

import (
	"context"
	"github.com/mitchellh/go-ps"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/yaml.v3"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Service keeping nebula configuration up to date",
	Run: func(cmd *cobra.Command, args []string) {
		agent, err := NewClient(l, config)
		if err != nil {
			l.WithError(err).Error("failed to create client")
			os.Exit(1)
		}
		defer agent.Close()

		templatePath := config.GetString("service.config_template", "nebula.yml.template")
		outputPath := config.GetString("service.config_output", "/etc/nebula/")

		templatePath = filepath.Join(configDir, templatePath)
		if ok, _ := fileExists(templatePath); !ok {
			l.Errorf("file not found: %s", templatePath)
			os.Exit(1)
		}

		interval := config.GetDuration("service.interval", 10*time.Minute)

		ticker := time.NewTicker(interval)
		sighup := make(chan os.Signal, 1)
		signal.Notify(sighup, syscall.SIGHUP)

		quit := make(chan struct{})
		go func() {
			run(agent, templatePath, outputPath)
			for {
				select {
				case <-sighup:
					l.Debug("received HUP signal, triggering update")
					run(agent, templatePath, outputPath)
				case <-ticker.C:
					l.Debug("received tick, triggering update")
					run(agent, templatePath, outputPath)
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()

		serviceShutdownBlock()
	},
}

func init() {
}

func serviceShutdownBlock() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT)

	select {
	case rawSig := <-sigChan:
		sig := rawSig.String()
		l.WithField("signal", sig).Info("Caught signal, shutting down")
	}
}

func run(agent *agentClient, configTemplatePath, outputPath string) {
	l.Infoln("Generate config dir")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, err := agent.client.GetEnrollStatus(ctx, &emptypb.Empty{})
	if err != nil {
		l.WithError(err).Errorf("failed to get enrollment status")
		return
	}

	if !status.IsEnrolled {
		l.Info("agent is not enrolled yet")
		return
	}

	var nebulaConfig map[interface{}]interface{}
	bs, err := ioutil.ReadFile(configTemplatePath)
	if err != nil {
		l.WithError(err).Errorf("failed to read template file")
		return
	}
	if err := yaml.Unmarshal(bs, &nebulaConfig); err != nil {
		l.WithError(err).Errorf("error when reading contents of %s", configTemplatePath)
		return
	}

	_, err = os.Stat(outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(outputPath, 0700)
			if err != nil {
				l.WithError(err).Errorf("failed to create directory `%s`\n", outputPath)
				return
			}
		} else {
			l.WithError(err).Errorf("failed to create directory `%s`\n", outputPath)
			return
		}
	}

	caBytes := make([]byte, 0)
	for i := range status.CertificateAuthorities {
		caBytes = append(caBytes, status.CertificateAuthorities[i].PublicKeyPEM...)
		caBytes = append(caBytes, '\n')
	}

	caFilePath := filepath.Join(outputPath, "ca.crt")
	if err := ioutil.WriteFile(caFilePath, caBytes, 0600); err != nil {
		l.WithError(err).Errorf("failed to write ca.crt to `%s`\n", caFilePath)
		return
	}

	if _, ok := nebulaConfig["pki"]; !ok {
		nebulaConfig["pki"] = make(map[string]interface{})
	}
	pki := nebulaConfig["pki"].(map[string]interface{})
	pki["ca"], err = filepath.Abs(caFilePath)
	if err != nil {
		l.WithError(err).Errorf("failed to get absolut path for `%s`\n", caFilePath)
		return
	}
	if ok, _ := fileExists(pki["ca"].(string)); !ok {
		l.WithError(err).Errorf("file not found `%s`\n", pki["ca"].(string))
		return
	}

	crtFilePath := filepath.Join(outputPath, "nebula.crt")
	if err := ioutil.WriteFile(crtFilePath, []byte(status.SignedPEM), 0600); err != nil {
		l.WithError(err).Errorf("failed to write nebula.crt to `%s`\n", crtFilePath)
		return
	}
	pki["cert"], err = filepath.Abs(crtFilePath)
	if err != nil {
		l.WithError(err).Errorf("failed to get absolut path for `%s`\n", crtFilePath)
		return
	}
	if ok, _ := fileExists(pki["cert"].(string)); !ok {
		l.WithError(err).Errorf("file not found `%s`\n", pki["cert"].(string))
		return
	}

	keyPath := resolvePath(AgentNebulaKeyPath)
	if ok, _ := fileExists(keyPath); !ok {
		l.Errorf("nebula key not exists: %s", keyPath)
		return
	}

	keyFileSource, err := ioutil.ReadFile(keyPath)
	if err != nil {
		l.WithError(err).Errorf("error open file: %s", keyPath)
		return
	}

	keyPathDst := filepath.Join(outputPath, "nebula.key")
	if ok, info := fileExists(keyPathDst); ok && info.Mode() != 0600 {
		l.Warnf("wrong permission on key %s, fixing...", keyPathDst)
		err = os.Chmod(keyPathDst, 0600)
		if err != nil {
			l.WithError(err).Errorf("failed to fix permission on %s", keyPathDst)
		}
	}

	if err := ioutil.WriteFile(keyPathDst, keyFileSource, 0600); err != nil {
		l.WithError(err).Errorf("failed to write nebula.key to `%s`\n", keyPathDst)
		return
	}

	pki["key"], err = filepath.Abs(keyPathDst)
	if err != nil {
		l.WithError(err).Errorf("failed to get absolut path for `%s`", keyPathDst)
		return
	}
	if ok, _ := fileExists(pki["key"].(string)); !ok {
		l.WithError(err).Errorf("file not found `%s`\n", pki["key"].(string))
		return
	}

	blockedFingerprints := make([]string, 0)

	for _, crl := range status.CertificateRevocationList {
		blockedFingerprints = append(blockedFingerprints, crl.Fingerprints...)
	}

	pki["blocklist"] = blockedFingerprints

	bs, err = yaml.Marshal(nebulaConfig)
	if err != nil {
		l.WithError(err).Error("failed to generate config")
		return
	}
	configFilePath := filepath.Join(outputPath, "config.yml")
	if err := ioutil.WriteFile(configFilePath, bs, 0600); err != nil {
		l.WithError(err).Errorf("failed to write config file to `%s`\n", configFilePath)
		return
	}

	reloadNebula()
}

func reloadNebula() {
	processes, err := ps.Processes()
	if err != nil {
		l.WithError(err).Error("failed to get os processes")
		return
	}

	for i := range processes {
		if processes[i].Executable() == "nebula" {
			l.Infof("nebula pid is %d", processes[i].Pid())
			err = syscall.Kill(processes[i].Pid(), syscall.SIGHUP)
			if err != nil {
				l.WithError(err).Errorf("failure when trying to trigger nebula reload")
			}
			break
		}
	}
}
