package smtp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	gosmtp "net/smtp"
)

type MailHogAPI struct {
	url *url.URL
}

func statusCodeToError(sc int) error {
	if sc/100 != 2 {
		txt := fmt.Sprintf(
			"Server returned unexepected status code: %d",
			sc,
		)
		return errors.New(txt)
	} else {
		return nil
	}
}

func NewMailHogAPI(api string) (*MailHogAPI, error) {
	if len(api) > 0 && api[len(api)-1] != '/' {
		api += "/"
	}
	parsed, err := url.Parse(api)
	if err != nil {
		return nil, err
	}

	res := &MailHogAPI{url: parsed}
	return res, nil
}

func (mh *MailHogAPI) DeleteAllEmails() error {

	deleteRequest := func() (*http.Request, error) {
		ctxt, _ := url.Parse("v1/messages")
		endPoint := mh.url.ResolveReference(ctxt)
		req, err := http.NewRequest(http.MethodDelete, endPoint.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Accept", "application/json")
		return req, nil
	}

	req, err := deleteRequest()
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	err = statusCodeToError(resp.StatusCode)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func (mh *MailHogAPI) MessageCount() (int, error) {
	ctxt, _ := url.Parse("v1/messages")
	endPoint := mh.url.ResolveReference(ctxt)
	client := &http.Client{}
	req, err := http.NewRequest("GET", endPoint.String(), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var msgs []interface{}
	err = json.Unmarshal(body, &msgs)
	if err != nil {
		return 0, err
	}
	return len(msgs), nil
}

type mhSearchResp struct {
	Total int `json:"total"`
	Count int `json:"count"`
	Start int `json:"start"`
	Items []struct {
		From struct {
			Mailbox string `json:"Mailbox"`
			Domain  string `json:"Domain"`
		} `json:"From"`
		To struct {
			Mailbox string `json:"Mailbox"`
			Domain  string `json:"Domain"`
		} `json:"To"`
		Content struct {
			Headers map[string][]string `json:"Headers"`
			Body    string              `json:"Body"`
		} `json:"Content"`
	} `json:"items"`
}

func (mh *MailHogAPI) CheckEmail(dest, from, subject, body string) (bool, error) {
	ctxt, _ := url.Parse("v2/search")
	done := false
	found := false
	start := 0
	for !done {
		reqUrl := mh.url.ResolveReference(ctxt)
		query := reqUrl.Query()
		query.Add("kind", "to")
		query.Add("query", dest)
		query.Add("start", fmt.Sprintf("%d", start))
		reqUrl.RawQuery = query.Encode()
		req, err := http.NewRequest(http.MethodGet, reqUrl.String(), nil)
		if err != nil {
			return false, err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		var msgs mhSearchResp
		respBody, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal(respBody, &msgs)
		if err != nil {
			return false, err
		}
		for _, m := range msgs.Items {
			msgFrom := fmt.Sprintf("%s@%s", m.From.Mailbox, m.From.Domain)
			msgSubject, ok := m.Content.Headers["Subject"]
			if !ok {
				continue
			}
			if len(msgSubject) != 1 {
				continue
			}
			msgBody := m.Content.Body
			if msgFrom == from && msgSubject[0] == subject && msgBody == body {
				found = true
				done = true
				break
			}
		}
		if !done {
			start += len(msgs.Items)
			if start >= msgs.Total {
				done = true
			}
		}
	}
	return found, nil
}

func SendEmail(host string, port int, from, to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	msg := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))
	return gosmtp.SendMail(addr, nil, from, []string{to}, msg)
}
