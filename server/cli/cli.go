package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/longzekun/specter/server/certs"
	"github.com/longzekun/specter/server/transport"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {

	//	Operator
	operatorCmd.Flags().StringP(nameFlagStr, "n", "", "name of the operator")
	operatorCmd.Flags().StringP(hostFlagStr, "l", "", "multiplayer listener host")
	operatorCmd.Flags().Uint16P(portFlagStr, "p", 6784, "multiplayer listener port")
	operatorCmd.Flags().StringP(saveFlagStr, "s", "", "operator configuration file path")
	operatorCmd.Flags().StringP(outputFlagStr, "o", "file", "output file format( file\\stdout )")
	rootCmd.AddCommand(operatorCmd)

}

var rootCmd = &cobra.Command{
	Use:   "specter-server",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//	root cmd
		certs.SetupCAs()

		zap.S().Debugf("Starting Specter Server")

		transport.StartMtlsClientListener("0.0.0.0", 7777)

		for {
			time.Sleep(1 * time.Second)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
