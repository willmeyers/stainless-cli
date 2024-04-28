package stainless

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
)

func (client *Client) ListOrgs() (*ListOrgsResponse, error) {
	endppint := client.NewStainlessApiURL("", "", "/orgs")
	req, err := client.Request("GET", endppint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	response := &ListOrgsResponse{}
	json.Unmarshal(bodyBytes, response)

	return response, nil
}

func (client *Client) ListProjects(orgName string) (*ListProjectsResponse, error) {
	endppint := client.NewStainlessApiURL(orgName, "", "/projects")
	req, err := client.Request("GET", endppint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	response := &ListProjectsResponse{}
	json.Unmarshal(bodyBytes, response)

	return response, nil
}

func (client *Client) ListMembers(orgName string) (*ListMembersResponse, error) {
	url := client.NewStainlessApiURL(
		orgName,
		"",
		"/members",
	)

	req, err := client.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	response := &ListMembersResponse{}
	json.Unmarshal(bodyBytes, response)

	return response, nil
}

func (client *Client) ListApiKeys(orgName string) (*ListApiKeysResponse, error) {
	url := client.NewStainlessApiURL(
		orgName,
		"",
		"/api-keys",
	)

	req, err := client.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	response := &ListApiKeysResponse{}
	json.Unmarshal(bodyBytes, response)

	return response, nil
}

func (client *Client) ListBuilds(orgName, projectName string) (*ListBuildsResponse, error) {
	// TODO (willmeyers) add support for pagination
	url := client.NewStainlessApiURL(
		orgName,
		projectName,
		"/builds",
	)

	req, err := client.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	response := &ListBuildsResponse{}
	json.Unmarshal(bodyBytes, response)

	return response, nil
}

func (client *Client) GetBuildStatus(orgName, projectName, language, branch string) (*GetBuildStatusResponse, error) {
	url := client.NewStainlessApiURL(
		orgName,
		projectName,
		fmt.Sprintf("/sdks/%s/builds/%s?fallbackBranch=main", language, branch),
	)

	req, err := client.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	response := &GetBuildStatusResponse{}
	json.Unmarshal(bodyBytes, response)

	return response, nil
}

func (client *Client) ListSdks(orgName, projectName string) (*ListSdksResponse, error) {
	url := client.NewStainlessApiURL(
		orgName,
		projectName,
		"/sdks",
	)

	req, err := client.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	response := &ListSdksResponse{}
	json.Unmarshal(bodyBytes, response)

	return response, nil
}

func (client *Client) GenerateSdk(orgName, projectName, openApiSpec, stainlessConfig, outDir, language string) (*GenerateResponse, error) {
	var b bytes.Buffer

	w := multipart.NewWriter(&b)

	w.WriteField("projectName", projectName)
	w.WriteField("branch", "main")

	f, err := os.Open(openApiSpec)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "openApiSpec", filepath.Base(openApiSpec)))
	header.Set("Content-Type", "application/x-yaml")
	fw, err := w.CreatePart(header)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return nil, err
	}

	f, err = os.Open(stainlessConfig)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfgHeader := make(textproto.MIMEHeader)
	cfgHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "stainlessConfig", filepath.Base(stainlessConfig)))
	cfgHeader.Set("Content-Type", "application/x-yaml")
	fw, err = w.CreatePart(cfgHeader)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return nil, err
	}

	w.Close()

	url := client.NewStainlessApiURL(
		orgName,
		projectName,
		"/generate",
	)

	req, err := client.Request("POST", url, &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	response := &GenerateResponse{}
	json.Unmarshal(bodyBytes, response)

	return response, nil
}

func (client *Client) GetDefaultOrg() (string, error) {
	orgs, err := client.ListOrgs()
	if err != nil {
		return "", err
	}

	if len(orgs.Items) < 1 {
		return "", errors.New("user has no orgs")
	}

	defaultOrg := orgs.Items[0]
	return defaultOrg.Name, nil
}

func (client *Client) GetDefaultProject(orgName string) (string, error) {
	projects, err := client.ListProjects(orgName)
	if err != nil {
		return "", err
	}

	if len(projects.Items) < 1 {
		return "", errors.New("user has no projects")
	}

	defaultProject := projects.Items[0]
	return defaultProject.Name, nil
}
