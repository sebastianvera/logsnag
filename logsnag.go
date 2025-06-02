package logsnag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type LogSnag struct {
	Token   string
	Project string
}

func (logsnag *LogSnag) GetProject() string {
	return logsnag.Project
}

type PublishRequest struct {
	Project     string         `json:"project,omitempty"`
	Channel     string         `json:"channel,omitempty"`
	Event       string         `json:"event,omitempty"`
	Description string         `json:"description,omitempty"`
	Icon        string         `json:"icon,omitempty"`
	Notify      bool           `json:"notify,omitempty"`
	Tags        map[string]any `json:"tags,omitempty"`
}

func (logsnag *LogSnag) Publish(input *PublishRequest) error {
	if input.Channel == "" || input.Event == "" {
		return errors.New("logsnag: LogSnag.Publish missing one of required fields channel, or event")
	}
	input.Project = logsnag.GetProject()

	baseURL := "https://api.logsnag.com/v1/log"

	body, err := json.Marshal(input)
	if err != nil {
		return errors.Wrap(err, "logsnag: LogSnag.Publish json.Marshal error")
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(body))

	if err != nil {
		return errors.Wrap(err, "logsnag: LogSnag.Publish http.NewRequest error")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", logsnag.Token))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "logsnag: LogSnag.Publish http.DefaultClient.Do error")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(res.Body)
		fmt.Println("body response error: ", string(b))
		return fmt.Errorf("logsnag: LogSnag.Publish unexpected http response status <%d> error", res.StatusCode)
	}

	return nil
}

func NewLogSnag(token string, project string) LogSnag {
	return LogSnag{
		Token:   token,
		Project: project,
	}
}
