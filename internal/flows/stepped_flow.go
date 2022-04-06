package flows

import (
	"fmt"
	"os"

	"github.com/sapslaj/geekbot-cli/internal/geekbotclient"
)

var ValidSteps = []string{"start", "edit", "check", "send"}

type SteppedFlow struct {
	Client    *geekbotclient.GeekbotClient
	Standup   *geekbotclient.Standup
	Steps     []string
	Filename  string
	Responses []*geekbotclient.QuestionAnswer
}

func (sf *SteppedFlow) IsValidStep(s string) bool {
	for _, v := range ValidSteps {
		if s == v {
			return true
		}
	}
	return false
}

func (sf *SteppedFlow) MustHaveFileName() {
	if sf.Filename == "" {
		sf.Filename = CreateStandupFile(sf.Standup)
	}
}

func (sf *SteppedFlow) MustHaveResponses() {
	if sf.Responses == nil {
		sf.MustHaveFileName()
		sf.Responses = ParseResponses(sf.Filename, sf.Standup)
	}
}

func (sf *SteppedFlow) Run(steps []string) {
	for _, s := range steps {
		if !sf.IsValidStep(s) {
			fmt.Fprintln(os.Stderr, "idk how to do '"+s+"'. typo?")
			os.Exit(1)
		}
	}
	sf.Steps = steps
	for {
		if len(sf.Steps) == 0 {
			return
		}
		var s string
		s, sf.Steps = sf.Steps[0], sf.Steps[1:]
		switch s {
		case "start":
			sf.MustHaveFileName()
			fmt.Println(sf.Filename)
		case "edit":
			sf.MustHaveFileName()
			OpenEditor(sf.Filename)
		case "check":
			sf.MustHaveResponses()
			PrintQuestionAnswers(sf.Client, sf.Responses)
			if len(sf.Steps) > 0 && !ConfirmPrompt("Confirm? [y/n]: ") {
				return
			}
		case "send":
			sf.MustHaveResponses()
			err := sf.Client.CreateReport(sf.Responses)
			Fuck(err)
			fmt.Println("Report sent.")
		}
	}
}

func NewSteppedFlow(client *geekbotclient.GeekbotClient, standup *geekbotclient.Standup) *SteppedFlow {
	return &SteppedFlow{
		Client:   client,
		Standup:  standup,
		Steps:    []string{},
		Filename: "",
	}
}

func RunSteppedFlow(client *geekbotclient.GeekbotClient, standup *geekbotclient.Standup, steps []string) {
	NewSteppedFlow(client, standup).Run(steps)
}
