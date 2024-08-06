package cmd

import (
	"os"

	"github.com/jiftle/sckproxy/internal/proxy"
	"github.com/spf13/cobra"
)

var (
	addr string
)

var rootCmd = &cobra.Command{
	Use:   "sckproxy",
	Short: "socket5 proxy",
	Long:  `socket5 proxy written by golang.`,
	Run: func(cmd *cobra.Command, args []string) {
		// proxy.StartTcpProxy(addr)
		proxy.StartHttpProxy(addr)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&addr, "listen", "l", ":1080", "listen addr. eg. :1080")
}
