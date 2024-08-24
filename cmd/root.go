package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/jiftle/sckproxy/internal/proxy"
	"github.com/spf13/cobra"
)

var (
	addr string
	mode string
)

var rootCmd = &cobra.Command{
	Use:   "sckproxy",
	Short: "socket5 proxy",
	Long:  `socket5 proxy written by golang.`,
	Run: func(cmd *cobra.Command, args []string) {
		// g.Log().SetFlags(g.Log().GetFlags() | glog.F_FILE_SHORT) //log contain filename and linenum

		if strings.EqualFold(mode, "socket") {
			proxy.StartSocket5Proxy(addr)
		} else if strings.EqualFold(mode, "http") {
			proxy.StartHttpProxy(addr)
		} else {
			g.Log().Warningf(context.Background(), "don't support mode: %v", mode)
			fmt.Printf("don't support mode: %v\n", mode)
		}
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
	rootCmd.Flags().StringVarP(&mode, "mode", "m", "socket", "socket,http")
}
