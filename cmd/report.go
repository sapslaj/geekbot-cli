package cmd

import (
	"fmt"
	"strings"

	"github.com/sapslaj/geekbot-cli/internal/flows"
	"github.com/sapslaj/geekbot-cli/internal/geekbotclient"

	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report",
	Run: func(cmd *cobra.Command, args []string) {
		standupName, err := cmd.Flags().GetString("standup")
		flows.Fuck(err)

		geekbot := geekbotclient.NewGeekbotClient(&geekbotclient.GeekbotClient{})
		standup := flows.SelectStandup(geekbot, standupName)

		filename := flows.CreateStandupFile(standup)
		flows.OpenEditor(filename)

		responses := flows.ParseResponses(filename, standup)

		report, err := geekbot.QuestionAnswersToJson(responses)
		flows.Fuck(err)

		var confirm string
		fmt.Println(string(report))
		fmt.Println("Confirm and send? [y/n]: ")
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm)[0] == 'y' {
			fmt.Println(geekbot.CreateReport(responses))
		} else {
			fmt.Println("aborting.")
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringP("standup", "s", "", "Standup name to report")
}
