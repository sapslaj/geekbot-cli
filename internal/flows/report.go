package flows

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/sapslaj/geekbot-cli/internal/geekbotclient"
)

func Qr(question string) string {
	// qr = question render, ie convert template variables to plain English
	question = strings.ReplaceAll(question, "{last_report_date}", "last report date")
	return question
}

func GetMdSection(header, contents string) (string, error) {
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

func SelectStandup(c *geekbotclient.GeekbotClient, name string) *geekbotclient.Standup {
	standups, err := c.StandupList()
	Fuck(err)
	if len(standups) == 0 {
		Fuck(errors.New("no standups to report to"))
	}
	if len(name) == 0 {
		if len(standups) > 1 {
			idx, err := fuzzyfinder.Find(
				standups,
				func(i int) string {
					return standups[i].Name
				},
			)
			Fuck(err)
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
		Fuck(errors.New("standup not found: `" + name + "`"))
	}
	return nil
}

func CreateStandupFile(standup *geekbotclient.Standup) string {
	userHome, err := os.UserHomeDir()
	Fuck(err)
	standupsDir := filepath.Join(userHome, ".geekbot", "standups")
	err = os.MkdirAll(standupsDir, os.ModePerm)
	Fuck(err)

	today := time.Now().Format("2006-01-02")
	filename := filepath.Join(standupsDir, Slugify(today+"-"+standup.Name)+".md")

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(filename)
		Fuck(err)
		defer file.Close()
		defer file.Sync()
		writer := bufio.NewWriter(file)
		defer writer.Flush()

		writer.WriteString("# " + standup.Name + " (" + today + ")\n\n")
		for _, question := range standup.Questions {
			writer.WriteString("## " + Qr(question.Text) + "\n\n")
		}
	}

	return filename
}

func OpenEditor(filename string) {
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
	Fuck(cmd.Run())
}

func ParseResponses(filename string, standup *geekbotclient.Standup) []*geekbotclient.QuestionAnswer {
	b, err := ioutil.ReadFile(filename)
	Fuck(err)
	contents := string(b)

	if len(contents) < 4 {
		Fuck(errors.New("looks like the standup file is kinda empty :/"))
	}

	var responses []*geekbotclient.QuestionAnswer
	for _, question := range standup.Questions {
		response, err := GetMdSection("## "+Qr(question.Text), contents)
		Fuck(err)
		responses = append(responses, &geekbotclient.QuestionAnswer{
			Standup:  standup,
			Question: question,
			Answer:   strings.Trim(response, "\n"),
		})
	}
	return responses
}
