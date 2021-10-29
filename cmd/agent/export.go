package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export data for use by nebula",
	Run: func(cmd *cobra.Command, args []string) {
		agent, err := NewClient(l, config)
		if err != nil {
			fmt.Printf("failed to create client: %s", err)
			os.Exit(1)
		}
		defer agent.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := agent.client.GetEnrollStatus(ctx, &emptypb.Empty{})
		if err != nil {
			l.WithError(err).Fatalln("failed to get enrollment status")
		}

		if res.IsEnrolled {
			outCA, err := cmd.Flags().GetString("out-ca")
			if err != nil {
				l.WithError(err).Fatalln("failed to get out-ca")
			}
			if outCA != "" {
				var caPEM []byte

				for _, ca := range res.CertificateAuthorities {
					caPEM = append(caPEM, ca.PublicKeyPEM...)
				}

				err = ioutil.WriteFile(outCA, caPEM, 0600)
				if err != nil {
					l.WithError(err).Errorf("error while writing out-ca")
				}
			}

			outCert, err := cmd.Flags().GetString("out-pub")
			if err != nil {
				l.WithError(err).Fatalln("failed to get out-pub")
			}
			if outCert != "" {
				err = ioutil.WriteFile(outCert, []byte(res.SignedPEM), 0600)
				if err != nil {
					l.WithError(err).Errorf("error while writing out-pub")
				}
			}

		} else if res.IsEnrollmentRequested {
			l.Errorln("Agent has requested to be enrolled")
		} else {
			l.Errorln("Agent enrollment not started")
		}
	},
}

func init() {
	exportCmd.Flags().String("out-ca", "", "Path to export ca to")
	exportCmd.Flags().String("out-pub", "", "Path to export signed pub client to")
}
