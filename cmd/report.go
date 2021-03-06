package cmd

import (
	"fmt"
	"os"

	"github.com/sapslaj/geekbot-cli/internal/flows"
	"github.com/sapslaj/geekbot-cli/internal/geekbotclient"

	"github.com/spf13/cobra"
)

var (
	standupName     string
	quickMode       bool
	interactiveMode bool
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report",
	Run: func(cmd *cobra.Command, args []string) {
		if quickMode && interactiveMode {
			fmt.Fprintln(os.Stderr, "Uh you can't have --quick and --interactive at the same time! Pick one or none!")
			os.Exit(1)
		}
		if quickMode {
			geekbot, standup := preflow()
			flows.RunQuickFlow(geekbot, standup)
		} else if interactiveMode {
			fmt.Fprintln(os.Stderr, "haven't build interactive mode yet sorry")
			os.Exit(501)
		} else {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "Provide report commands or use `-q` or `-i`. See `-h` for help.")
				os.Exit(1)
			}
			geekbot, standup := preflow()
			flows.RunSteppedFlow(geekbot, standup, args)
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringVarP(&standupName, "standup", "s", "", "Standup name to report")
	reportCmd.Flags().BoolVarP(&quickMode, "quick", "q", false, "Create, open, edit, send all in one action")
	reportCmd.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "Semi-guided standup process")
}

func preflow() (*geekbotclient.GeekbotClient, *geekbotclient.Standup) {
	geekbot := geekbotclient.NewGeekbotClient(&geekbotclient.GeekbotClient{})
	standup := flows.SelectStandup(geekbot, standupName)
	return geekbot, standup
}
