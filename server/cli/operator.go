package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/longzekun/specter/server/certs"
	"github.com/longzekun/specter/server/console"
	"github.com/spf13/cobra"
)

const (
	nameFlagStr   = "name"
	hostFlagStr   = "lhost"
	portFlagStr   = "lport"
	saveFlagStr   = "save"
	outputFlagStr = "output"
)

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Generate operator configuration files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString(nameFlagStr)
		if err != nil {
			fmt.Printf("failed to parse flags %v: %v\n", nameFlagStr, err)
			return
		}
		if name == "" {
			fmt.Printf("failed to parse flags %v: empty flag name\n", nameFlagStr)
			return
		}

		host, err := cmd.Flags().GetString(hostFlagStr)
		if err != nil {
			fmt.Printf("failed to parse flags %v: %v\n", hostFlagStr, err)
			return
		}
		if host == "" {
			fmt.Printf("failed to parse flags %v: empty flag host\n", hostFlagStr)
			return
		}

		port, err := cmd.Flags().GetUint16(portFlagStr)
		if err != nil {
			fmt.Printf("failed to parse flags %v: %v\n", portFlagStr, err)
			return
		}

		if port == 0 {
			fmt.Printf("failed to parse flags %v: empty flag port\n", portFlagStr)
			return
		}

		save, err := cmd.Flags().GetString(saveFlagStr)
		if err != nil {
			fmt.Printf("failed to parse flags %v: %v\n", saveFlagStr, err)
			return
		}

		output, err := cmd.Flags().GetString(outputFlagStr)
		if err != nil {
			fmt.Printf("failed to parse flags %v: %v\n", outputFlagStr, err)
			return
		}

		if output == "" {
			output = "file"
		} else if output != "file" {
			if output != "stdout" {
				output = "file"
			}
		}

		wd, _ := os.Getwd()
		if output == "file" {
			if save == "" {
				fmt.Println("The output format is a file, and a save path must be specified.")
				return
			} else {
				if _, err := os.Stat(filepath.Join(wd, save)); err == nil {
					fmt.Println("output file already exists. ")
					return
				} else if !os.IsNotExist(err) {
					fmt.Println("output file already exists.")
					return
				}
			}
		}

		certs.SetupCAs()

		clientConfig, err := console.NewOperatorClientConfig(name, host, port)
		if err != nil {
			fmt.Printf("failed to create operator client config: %v\n", err)
			return
		}

		if output != "file" {
			fmt.Println(string(clientConfig))
			return
		}

		savePath := filepath.Join(wd, save)
		f, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Printf("failed to open save file: %v\n", err)
			return
		}
		defer f.Close()

		_, err = f.Write(clientConfig)
		if err != nil {
			fmt.Printf("failed to write save file: %v\n", err)
		}
		return

	},
}
