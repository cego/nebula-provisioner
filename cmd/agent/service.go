package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
	"github.com/pkg/errors"
	conf "github.com/slackhq/nebula/config"

	"github.com/mitchellh/go-ps"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/yaml.v3"
)

var sighupCh = make(chan os.Signal, 1)
var advertiseAddr = ""
var advertiseAddrMutex = &sync.RWMutex{}
var nebulaPort = 4242

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

		nebulaPort = getNebulaPort(templatePath)

		interval := config.GetDuration("service.interval", 10*time.Minute)

		ticker := time.NewTicker(interval)
		tickerPM := time.NewTicker(1 * time.Minute)

		signal.Notify(sighupCh, syscall.SIGHUP)

		quit := make(chan struct{})
		go func() {
			if config.GetBool("service.port_mapping.enabled", false) {
				go natPMP(config)
			}
			run(agent, templatePath, outputPath)
			for {
				select {
				case <-sighupCh:
					l.Debug("received HUP signal, triggering update")
					run(agent, templatePath, outputPath)
				case <-ticker.C:
					l.Debug("received tick, triggering update")
					run(agent, templatePath, outputPath)
				case <-tickerPM.C:
					if config.GetBool("service.port_mapping.enabled", false) {
						l.Debug("received tick for port mapping")
						go natPMP(config)
					}
				case <-quit:
					ticker.Stop()
					tickerPM.Stop()
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

	nebulaPort = getNebulaPort(configTemplatePath)

	nebulaConfig, err := readNebulaTemplate(configTemplatePath)
	if err != nil {
		l.WithError(err).Errorf("failed to read template")
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

	crtFilePath := filepath.Join(outputPath, NebulaCrtPath)
	if err := ioutil.WriteFile(crtFilePath, []byte(status.SignedPEM), 0600); err != nil {
		l.WithError(err).Errorf("failed to write %s to `%s`\n", NebulaCrtPath, crtFilePath)
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

	keyPathDst := filepath.Join(outputPath, NebulaKeyPath)
	if ok, info := fileExists(keyPathDst); ok && info.Mode() != 0600 {
		l.Warnf("wrong permission on key %s, fixing...", keyPathDst)
		err = os.Chmod(keyPathDst, 0600)
		if err != nil {
			l.WithError(err).Errorf("failed to fix permission on %s", keyPathDst)
		}
	}

	if err := ioutil.WriteFile(keyPathDst, keyFileSource, 0600); err != nil {
		l.WithError(err).Errorf("failed to write %s to `%s`\n", NebulaKeyPath, keyPathDst)
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

	if _, ok := nebulaConfig["lighthouse"]; !ok {
		nebulaConfig["lighthouse"] = make(map[string]interface{})
	}
	lh := nebulaConfig["lighthouse"].(map[string]interface{})

	advertiseAddrMutex.RLock()
	if advertiseAddr != "" {
		if _, ok := lh["advertise_addrs"]; ok {
			switch x := lh["advertise_addrs"].(type) {
			case []string:
				lh["advertise_addrs"] = append(x, advertiseAddr)
			}
		} else {
			lh["advertise_addrs"] = []string{advertiseAddr}
		}
	}
	advertiseAddrMutex.RUnlock()

	bs, err := yaml.Marshal(nebulaConfig)
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

func readNebulaTemplate(configTemplatePath string) (map[interface{}]interface{}, error) {
	var nebulaConfig map[interface{}]interface{}
	bs, err := ioutil.ReadFile(configTemplatePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read template file")
	}
	if err := yaml.Unmarshal(bs, &nebulaConfig); err != nil {
		return nil, errors.Wrapf(err, "error when reading contents of %s", configTemplatePath)
	}
	return nebulaConfig, nil
}

func getNebulaPort(configTemplatePath string) int {
	nebulaConfig, err := readNebulaTemplate(configTemplatePath)
	if err != nil {
		return 4242
	}
	if _, ok := nebulaConfig["listen"]; ok {
		listen := nebulaConfig["listen"].(map[string]interface{})
		if _, ok := listen["port"]; ok {
			switch x := listen["port"].(type) {
			case string:
				port, err := strconv.Atoi(x)
				if err != nil {
					return 4242
				}
				return port
			}
		}
	}
	return 4242
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

func natPMP(config *conf.C) {
	var gatewayIP net.IP
	var err error

	if config.IsSet("service.port_mapping.gateway") {
		gw := config.GetString("service.port_mapping.gateway", "")
		gatewayIP = net.ParseIP(gw)
	} else {
		gatewayIP, err = gateway.DiscoverGateway()
		if err != nil {
			return
		}
	}

	client := natpmp.NewClient(gatewayIP)

	addr, err := client.GetExternalAddress()
	if err != nil {
		return
	}

	portMapping, err := client.AddPortMapping("udp", nebulaPort, 0, 120)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	externalIP := net.IP(addr.ExternalIPAddress[:])

	advertiseAddrNew := fmt.Sprintf("%s:%d", externalIP.String(), portMapping.MappedExternalPort)

	advertiseAddrMutex.Lock()
	defer advertiseAddrMutex.Unlock()

	if advertiseAddr != advertiseAddrNew {
		l.Debugf("advertiseAddr changed from %s -> %s", advertiseAddr, advertiseAddrNew)
		advertiseAddr = advertiseAddrNew
		sighupCh <- syscall.SIGHUP
	}
}
