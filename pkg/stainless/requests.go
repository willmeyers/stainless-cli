package stainless

import (
	"fmt"
	"io"
	"net/http"
)

const BASE_URL string = "https://api.stainlessapi.com/api"

func (client *Client) Request(method, url string, body io.Reader) (*http.Request, error) {
	if client.credentials == "" {
		err := fmt.Errorf("missing authentication credentials: did you run \n\n $ stainless-cli login\n\n")
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("cookie", client.credentials)

	return req, nil
}

func (client *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 401 {
			err := fmt.Errorf("check credentials, server responded with %s", resp.Status)
			return resp, err
		}

		if resp.StatusCode >= 400 {
			err := fmt.Errorf("server responded with %s", resp.Status)
			return resp, err
		}
	}

	return resp, nil
}

func (client *Client) NewStainlessApiURL(orgName, projectName, path string) string {
	if orgName != "" && projectName != "" && path != "" {
		return fmt.Sprintf("%s/orgs/%s/projects/%s%s", BASE_URL, orgName, projectName, path)
	} else if orgName != "" {
		return fmt.Sprintf("%s/orgs/%s%s", BASE_URL, orgName, path)
	}

	return fmt.Sprintf("%s%s", BASE_URL, path)
}
