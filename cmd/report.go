package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sapslaj/geekbot-cli/internal/geekbotclient"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

func fuck(err error) {
	if err != nil {
		panic(err)
	}
}

func slugify(s string) string {
	r := regexp.MustCompile(`(?m)[!@#$%^&\*\(\)\[\]\{\};:\,\./<>\?\|\x60~=_+ ]`)
	return r.ReplaceAllString(strings.ToLower(s), "-")
}

func qr(question string) string {
	// qr = question render, ie convert template variables to plain English
	question = strings.ReplaceAll(question, "{last_report_date}", "last report date")
	return question
}

func getMdSection(header, contents string) (string, error) {
	start := strings.Index(contents, header) + len(header)
	if start < 0 {
		return "", errors.New("Can't find section `" + header + "`")
	}
	end := strings.Index(contents[start:], "\n#")
	if end < 0 {
		return contents[start:], nil
	} else {
		return contents[start : start+end], nil
	}
}

func selectStandup(c *geekbotclient.GeekbotClient, name string) *geekbotclient.Standup {
	standups, err := c.StandupList()
	fuck(err)
	if len(standups) == 0 {
		fuck(errors.New("no standups to report to"))
	}
	if len(name) == 0 {
		if len(standups) > 1 {
			idx, err := fuzzyfinder.Find(
				standups,
				func(i int) string {
					return standups[i].Name
				},
			)
			fuck(err)
			return standups[idx]
		} else {
			return standups[0]
		}
	} else {
		for _, standup := range standups {
			if standup.Name == name {
				return standup
			}
		}
		fuck(errors.New("standup not found: `" + name + "`"))
	}
	return nil
}

func createStandupFile(standup *geekbotclient.Standup) string {
	userHome, err := os.UserHomeDir()
	fuck(err)
	standupsDir := filepath.Join(userHome, ".geekbot", "standups")
	err = os.MkdirAll(standupsDir, os.ModePerm)
	fuck(err)

	today := time.Now().Format("2006-01-02")
	filename := filepath.Join(standupsDir, slugify(today+"-"+standup.Name)+".md")

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(filename)
		fuck(err)
		defer file.Close()
		defer file.Sync()
		writer := bufio.NewWriter(file)
		defer writer.Flush()

		writer.WriteString("# " + standup.Name + " (" + today + ")\n\n")
		for _, question := range standup.Questions {
			writer.WriteString("## " + qr(question.Text) + "\n\n")
		}
	}

	return filename
}

func openEditor(filename string) {
	editor, set := os.LookupEnv("GEEKBOT_EDITOR")
	if !set {
		editor, set = os.LookupEnv("EDITOR")
		if !set {
			editor = "vi"
		}
	}
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	fuck(cmd.Run())
}

func parseResponses(filename string, standup *geekbotclient.Standup) []*geekbotclient.QuestionAnswer {
	b, err := ioutil.ReadFile(filename)
	fuck(err)
	contents := string(b)

	if len(contents) < 4 {
		fuck(errors.New("looks like the standup file is kinda empty :/"))
	}

	var responses []*geekbotclient.QuestionAnswer
	for _, question := range standup.Questions {
		response, err := getMdSection("## "+qr(question.Text), contents)
		fuck(err)
		responses = append(responses, &geekbotclient.QuestionAnswer{
			Standup:  standup,
			Question: question,
			Answer:   strings.Trim(response, "\n"),
		})
	}
	return responses
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a report",
	Run: func(cmd *cobra.Command, args []string) {
		standupName, err := cmd.Flags().GetString("standup")
		fuck(err)

		geekbot := geekbotclient.NewGeekbotClient(&geekbotclient.GeekbotClient{})
		standup := selectStandup(geekbot, standupName)

		filename := createStandupFile(standup)
		openEditor(filename)

		responses := parseResponses(filename, standup)

		report, err := geekbot.QuestionAnswersToJson(responses)
		fuck(err)

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

	reportCmd.Flags().String("standup", "", "Standup name to report")
}
