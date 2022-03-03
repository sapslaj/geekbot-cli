package flows

import (
	"fmt"

	"github.com/sapslaj/geekbot-cli/internal/geekbotclient"
)

type QuickFlow struct {
	Client  *geekbotclient.GeekbotClient
	Standup *geekbotclient.Standup
}

func (qf *QuickFlow) Run(_ []string) {
	filename := CreateStandupFile(qf.Standup)
	OpenEditor(filename)

	responses := ParseResponses(filename, qf.Standup)
	sent, err := ConfirmAndSend(qf.Client, responses)
	Fuck(err)
	if sent {
		fmt.Println("Report sent.")
	} else {
		fmt.Println("aborting.")
	}
}

func NewQuickFlow(client *geekbotclient.GeekbotClient, standup *geekbotclient.Standup) *QuickFlow {
	return &QuickFlow{
		Client:  client,
		Standup: standup,
	}
}

func RunQuickFlow(client *geekbotclient.GeekbotClient, standup *geekbotclient.Standup) {
	NewQuickFlow(client, standup).Run([]string{})
}
