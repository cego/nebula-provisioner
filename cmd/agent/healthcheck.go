package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/slackhq/nebula/cert"
	"github.com/spf13/cobra"
)

func init() {
	healthcheckCmd.Flags().Int("critical", 5, "Critical: Minimum number of days a certificate has to be valid")
	healthcheckCmd.Flags().Int("warning", 10, "Warning: Minimum number of days a certificate has to be valid")
}

var healthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Check certs and nebula is active",
	Run: func(cmd *cobra.Command, args []string) {
		outputPath := config.GetString("service.config_output", "/etc/nebula/")
		critDaysBeforeExpire, err := cmd.Flags().GetInt("critical")
		if err != nil {
			fmt.Printf("failed to get `--critical` parameter: %s\n", err)
			os.Exit(3)
		}
		warnDaysBeforeExpire, err := cmd.Flags().GetInt("warning")
		if err != nil {
			fmt.Printf("failed to get `--warning` parameter: %s\n", err)
			os.Exit(3)
		}

		_, err = os.Stat(outputPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("%s does not exists\n", outputPath)
			} else {
				fmt.Printf("%s does not exists error: %s\n", outputPath, err)
			}
			os.Exit(2)
		}

		_, err = os.Stat(filepath.Join(outputPath, NebulaCrtPath))
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("%s/%s does not exists\n", outputPath, NebulaCrtPath)
			} else {
				fmt.Printf("%s/%s does not exists error: %s\n", outputPath, NebulaCrtPath, err)
			}
			os.Exit(2)
		}

		_, err = os.Stat(filepath.Join(outputPath, NebulaKeyPath))
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("%s/%s does not exists\n", outputPath, NebulaKeyPath)
			} else {
				fmt.Printf("%s/%s does not exists error: %s\n", outputPath, NebulaKeyPath, err)
			}
			os.Exit(2)
		}

		bytes, err := os.ReadFile(filepath.Join(outputPath, NebulaCrtPath))
		if err != nil {
			fmt.Printf("failed to read file: %s err: %s\n", filepath.Join(outputPath, NebulaCrtPath), err)
			os.Exit(2)
		}
		crt, _, err := cert.UnmarshalNebulaCertificateFromPEM(bytes)
		if err != nil {
			fmt.Printf("failed to unmarshall nebula certificate err: %s\n", err)
			os.Exit(2)
		}

		if crt.Expired(time.Now()) {
			fmt.Println("nebula certificate is expired")
			os.Exit(2)
		}

		if crt.Expired(time.Now().AddDate(0, 0, critDaysBeforeExpire)) {
			fmt.Printf("nebula certificate will expire at %s\n", crt.Details.NotAfter)
			os.Exit(2)
		}

		if crt.Expired(time.Now().AddDate(0, 0, warnDaysBeforeExpire)) {
			fmt.Printf("nebula certificate will expire at %s\n", crt.Details.NotAfter)
			os.Exit(1)
		}

		checkCA(outputPath, warnDaysBeforeExpire, critDaysBeforeExpire)

		fmt.Println("OK")
	},
}

func checkCA(outputPath string, warnDaysBeforeExpire, critDaysBeforeExpire int) {
	_, err := os.Stat(filepath.Join(outputPath, NebulaCaPath))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s/%s does not exists\n", outputPath, NebulaCaPath)
		} else {
			fmt.Printf("%s/%s does not exists error: %s\n", outputPath, NebulaCaPath, err)
		}
		os.Exit(2)
	}

	bytes, err := os.ReadFile(filepath.Join(outputPath, NebulaCaPath))
	if err != nil {
		fmt.Printf("failed to read file: %s err: %s\n", filepath.Join(outputPath, NebulaCaPath), err)
		os.Exit(2)
	}

	var newestCA *cert.NebulaCertificate

	var ca *cert.NebulaCertificate
	for true {
		if len(bytes) <= 1 {
			break
		}
		ca, bytes, err = cert.UnmarshalNebulaCertificateFromPEM(bytes)
		if err != nil {
			fmt.Printf("failed to unmarshall nebula ca certificate err: %s\n", err)
			os.Exit(2)
		}
		if newestCA == nil {
			newestCA = ca
		} else if newestCA.Details.NotAfter.Before(ca.Details.NotAfter) {
			newestCA = ca
		}
	}

	if newestCA.Expired(time.Now()) {
		fmt.Println("nebula ca certificate is expired")
		os.Exit(2)
	}

	if newestCA.Expired(time.Now().AddDate(0, 0, critDaysBeforeExpire)) {
		fmt.Printf("nebula ca certificate will expire at %s\n", newestCA.Details.NotAfter)
		os.Exit(2)
	}

	if newestCA.Expired(time.Now().AddDate(0, 0, warnDaysBeforeExpire)) {
		fmt.Printf("nebula ca certificate will expire at %s\n", newestCA.Details.NotAfter)
		os.Exit(1)
	}

}
