package geekbotclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	apiHostname                 = "https://api.geekbot.com"
	apiTokenEnvironmentVariable = "GEEKBOT_API_KEY"
)

type ReportCreateAnswer struct {
	Text string `json:"text"`
}

type ReportCreate struct {
	StandupID int                           `json:"standup_id"`
	Answers   map[string]ReportCreateAnswer `json:"answers"`
}

type GeekbotClient struct {
	Ctx      context.Context
	ApiToken *string
	Client   *http.Client
}

func NewGeekbotClient(input *GeekbotClient) *GeekbotClient {
	var ctx context.Context
	var apiToken *string
	var client *http.Client

	if input.Ctx == nil {
		ctx = context.TODO()
	} else {
		ctx = input.Ctx
	}

	if input.ApiToken == nil {
		token, found := os.LookupEnv(apiTokenEnvironmentVariable)
		if !found {
			panic(apiTokenEnvironmentVariable + " not found in ENV")
		}
		apiToken = &token
	} else {
		apiToken = input.ApiToken
	}

	if input.Client == nil {
		client = &http.Client{}
	} else {
		client = input.Client
	}

	return &GeekbotClient{
		Ctx:      ctx,
		ApiToken: apiToken,
		Client:   client,
	}
}

func (c *GeekbotClient) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	if path[0] != '/' {
		path = "/" + path
	}
	req, err := http.NewRequest(method, apiHostname+path, body)
	if err != nil {
		return req, err
	}
	req.Header.Add("Authorization", *c.ApiToken)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c *GeekbotClient) newV1Request(method, path string, body io.Reader) (*http.Request, error) {
	if path[0] != '/' {
		path = "/" + path
	}
	return c.newRequest(method, "/v1"+path, body)
}

func (c *GeekbotClient) rawResponse(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("GeekbotClient.jsonResponse: request is nil")
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GeekbotClient) jsonResponse(req *http.Request, v interface{}) (*http.Response, []byte, error) {
	resp, err := c.rawResponse(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(body, v)
	return resp, body, err
}

func (c *GeekbotClient) StandupList() ([]*Standup, error) {
	req, err := c.newV1Request("GET", "/standups", nil)
	if err != nil {
		return nil, err
	}
	standups := make([]*Standup, 2)
	_, _, err = c.jsonResponse(req, &standups)
	return standups, err
}

func (c *GeekbotClient) QuestionAnswersToJson(qas []*QuestionAnswer) ([]byte, error) {
	// Geekbots API is kinda terrible. Mismatched types, no consistent structure
	// or schema, bad documentation, etc. Geekbot devs, if you see this, please
	// contact me. I don't want to shit talk your work.
	data := ReportCreate{
		// StandupID: fmt.Sprintf("%d", qas[0].Standup.Id),
		StandupID: qas[0].Standup.Id,
		Answers:   make(map[string]ReportCreateAnswer),
	}
	for _, qa := range qas {
		data.Answers[fmt.Sprintf("%d", qa.Question.Id)] = ReportCreateAnswer{
			Text: qa.Answer,
		}
	}
	return json.Marshal(data)
}

func (c *GeekbotClient) CreateReport(qas []*QuestionAnswer) error {
	body, err := c.QuestionAnswersToJson(qas)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	req, err := c.newV1Request("POST", "/reports", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	resp, err := c.rawResponse(req)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	rbody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println(rbody)
		return errors.New("report failed to send")
	}
	return nil
}
