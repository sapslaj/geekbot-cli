package geekbotclient

type Question struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type QuestionAnswer struct {
	Standup  *Standup
	Question *Question
	Answer   string
}

type Standup struct {
	Id        int         `json:"id"`
	Name      string      `json:"name"`
	Questions []*Question `json:"questions"`
}
