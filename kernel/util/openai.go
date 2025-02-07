// SiYuan - Build Your Eternal Digital Garden
// Copyright (c) 2020-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package util

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/siyuan-note/logging"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func ChatGPT(msg string, contextMsgs []string, c *gogpt.Client, maxTokens, timeout int) (ret string, stop bool, err error) {
	var reqMsgs []gogpt.ChatCompletionMessage

	for _, ctxMsg := range contextMsgs {
		reqMsgs = append(reqMsgs, gogpt.ChatCompletionMessage{
			Role:    "user",
			Content: ctxMsg,
		})
	}
	reqMsgs = append(reqMsgs, gogpt.ChatCompletionMessage{
		Role:    "user",
		Content: msg,
	})

	req := gogpt.ChatCompletionRequest{
		Model:     gogpt.GPT3Dot5Turbo,
		MaxTokens: maxTokens,
		Messages:  reqMsgs,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	resp, err := c.CreateChatCompletion(ctx, req)
	if nil != err {
		PushErrMsg("Requesting failed, please check kernel log for more details", 3000)
		logging.LogErrorf("create chat completion failed: %s", err)
		stop = true
		return
	}

	if 1 > len(resp.Choices) {
		stop = true
		return
	}

	buf := &strings.Builder{}
	choice := resp.Choices[0]
	buf.WriteString(choice.Message.Content)
	if "length" == choice.FinishReason {
		stop = false
	} else {
		stop = true
	}

	ret = buf.String()
	ret = strings.TrimSpace(ret)
	return
}

func NewOpenAIClient(apiKey, apiProxy, apiBaseURL string) *gogpt.Client {
	config := gogpt.DefaultConfig(apiKey)
	if "" != apiProxy {
		proxyUrl, err := url.Parse(apiProxy)
		if nil != err {
			logging.LogErrorf("OpenAI API proxy failed: %v", err)
		} else {
			config.HTTPClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		}
	}

	config.BaseURL = apiBaseURL
	return gogpt.NewClientWithConfig(config)
}
