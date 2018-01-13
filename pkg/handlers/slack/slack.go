/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package slack

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"

	"github.com/skippbox/kubewatch/config"
	"github.com/skippbox/kubewatch/pkg/event"
	kbEvent "github.com/skippbox/kubewatch/pkg/event"
	"net/http"
	"net/url"
	"io/ioutil"
	"time"
)

var slackColors = map[string]string{
	"Normal":  "good",
	"Warning": "warning",
	"Danger":  "danger",
}

var slackErrMsg = `
%s

You need to set both slack token and channel for slack notify,
using "--token/-t" and "--channel/-c", or using environment variables:

export KW_SLACK_TOKEN=slack_token
export KW_SLACK_CHANNEL=slack_channel

Command line flags will override environment variables

`

// Slack handler implements handler.Handler interface,
// Notify event to slack channel
type Slack struct {
	Token   string
	Channel string
	Url 	string
	Connected	bool
}

// Init prepares slack configuration
func (s *Slack) Init(c *config.Config) error {
	token := c.Handler.Slack.Token
	channel := c.Handler.Slack.Channel
	url := c.Handler.Slack.Url

	if token == "" {
		token = os.Getenv("KW_SLACK_TOKEN")
	}

	if channel == "" {
		channel = os.Getenv("KW_SLACK_CHANNEL")
	}

	if url == "" {
		url = os.Getenv("KW_URL")
	}

	s.Token = token
	s.Channel = channel
	s.Url = url
	s.Connected = false

	go s.getServer()

	return nil
}

func (s *Slack) getServer() {
	fmt.Println("Waiting for receiving server pong response...")

	for {
		time.Sleep(time.Second * 3)
		r, err := http.Get(s.Url + "/ping")
		if err != nil {
			fmt.Println(err)
			continue
		}
		responseData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if string(responseData) == "pong" {
			s.Connected = true
			fmt.Println("Response received")
			break
		}
	}
}

func (s *Slack) ObjectCreated(obj interface{}) {
	notifySlack(s, obj, "created")
}

func (s *Slack) ObjectDeleted(obj interface{}) {
	notifySlack(s, obj, "deleted")
}

func (s *Slack) ObjectUpdated(oldObj, newObj interface{}) {
	notifySlack(s, newObj, "updated")
}

func notifySlack(s *Slack, obj interface{}, action string) {
	e := kbEvent.New(obj, action)
	msg := prepareHTTPMessage(e)

	for {
		if s.Connected {
			break
		}
		time.Sleep(time.Second * 3)
	}
	http.PostForm(s.Url, url.Values{
		"event": {msg},
		"kind": {e.Kind},
		"namespace": {e.Namespace},
		"reason": {e.Reason},
		"name": {e.Name},
		"status": {e.Status},
		"host": {e.Host},
		"component": {e.Component},
	})
}

func prepareHTTPMessage(e event.Event) string {
	msg := fmt.Sprintf(
		"A %s in namespace %s has been %s: %s",
		e.Kind,
		e.Namespace,
		e.Reason,
		e.Name,
	)
	return msg
}

func checkMissingSlackVars(s *Slack) error {
	if s.Token == "" || s.Channel == "" {
		return fmt.Errorf(slackErrMsg, "Missing slack token or channel")
	}

	return nil
}

func prepareSlackAttachment(e event.Event) slack.Attachment {
	msg := fmt.Sprintf(
		"A %s in namespace %s has been %s: %s",
		e.Kind,
		e.Namespace,
		e.Reason,
		e.Name,
	)

	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "kubewatch",
				Value: msg,
			},
		},
	}

	if color, ok := slackColors[e.Status]; ok {
		attachment.Color = color
	}

	attachment.MarkdownIn = []string{"fields"}

	return attachment
}
