package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type SlackManager struct {
	apiUrl   string
	apiToken string
}

func NewSlackManager(apiUrl string, apiToken string) *SlackManager {
	return &SlackManager{
		apiUrl:   apiUrl,
		apiToken: apiToken,
	}
}

type slackChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type slackChannelListResponse struct {
	OK       bool           `json:"ok"`
	Channels []slackChannel `json:"channels"`
}

func (s SlackManager) prepareHeader(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiToken))
}

func (s SlackManager) GetChannelId(channelName string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/conversations.list", s.apiUrl), nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create http request")
	}
	s.prepareHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to do request")
	}
	defer resp.Body.Close()

	var channelListResponse slackChannelListResponse
	if err := json.NewDecoder(resp.Body).Decode(&channelListResponse); err != nil || !channelListResponse.OK {
		return "", errors.Wrap(err, "failed to decode http response")
	}

	for _, channel := range channelListResponse.Channels {
		if channel.Name == channelName {
			return channel.ID, nil
		}
	}

	return "", errors.New("channel not found")
}

type sendMessageResponse struct {
	OK    bool   `json:"ok"`
	Ts    string `json:"ts"`
	Error string `json:"error"`
}

func (s SlackManager) SendMessage(channelName string, text string, threadId string) (string, error) {
	channelId, err := s.GetChannelId(channelName)
	if err != nil {
		return "", errors.Wrap(err, "failed to get channel id")
	}

	client := &http.Client{}

	reqBody, err := json.Marshal(map[string]string{
		"channel":   channelId,
		"text":      text,
		"thread_ts": threadId,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal chat.postMessage body")
	}

	reqBodyBuffer := bytes.NewBuffer(reqBody)
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/chat.postMessage", s.apiUrl),
		reqBodyBuffer,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to send chat.postMessage request")
	}
	s.prepareHeader(req)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to do request")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read chat.postMessage response")
	}

	var sendMessageResponse sendMessageResponse
	if err := json.Unmarshal(respBody, &sendMessageResponse); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal sendMessageResponse")
	}

	if !sendMessageResponse.OK {
		return "", errors.New(sendMessageResponse.Error)
	}

	return sendMessageResponse.Ts, nil
}
