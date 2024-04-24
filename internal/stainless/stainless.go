package stainless

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

const BASE_URL string = "https://api.stainlessapi.com/api"

type Stainless struct {
	AuthCookies string
	OrgName     string
	ProjectName string
}

func New(options ...func(*Stainless) error) (*Stainless, error) {
	stl := &Stainless{}
	for _, option := range options {
		err := option(stl)
		if err != nil {
			return nil, err
		}
	}

	return stl, nil
}

func WithAuthCookies(cookieString string) func(*Stainless) error {
	return func(stl *Stainless) error {
		if cookieString == "" {
			return errors.New("cookie string cannot be empty")
		}
		stl.AuthCookies = cookieString
		return nil
	}
}

func WithDefaultOrgName() func(*Stainless) error {
	return func(stl *Stainless) error {
		orgs, err := stl.ListOrgs()
		if err != nil {
			return err
		}

		if len(orgs.Items) < 1 {
			return errors.New("no orgs found for account")
		}

		stl.OrgName = orgs.Items[0].Name

		return nil
	}
}

func WithDefaultProjectName() func(*Stainless) error {
	return func(stl *Stainless) error {
		if stl.OrgName == "" {
			return errors.New("cannot set default project without default org set")
		}

		projects, err := stl.ListProjects(stl.OrgName)
		if err != nil {
			return err
		}

		if len(projects.Items) < 1 {
			return errors.New("no orgs found for account")
		}

		stl.ProjectName = projects.Items[0].Name

		return nil
	}
}

func (stl *Stainless) ListOrgs() (*ListOrgsResponse, error) {
	url := stl.formatURL(&StainlessURLParams{
		Path: "/orgs",
	})

	req, err := stl.newAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
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

func (stl *Stainless) ListProjects(orgName string) (*ListProjectsResponse, error) {
	url := stl.formatURL(&StainlessURLParams{
		OrgName: orgName,
		Path:    "/projects",
	})

	req, err := stl.newAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
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

func (stl *Stainless) ListBuilds(orgName, projectName string) (*ListBuildsResponse, error) {
	url := stl.formatURL(&StainlessURLParams{
		OrgName:     orgName,
		ProjectName: projectName,
		Path:        "/builds",
	})
	req, err := stl.newAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
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

func (stl *Stainless) Generate(orgName, projectName, openApiSpec, stainlessConfig, outDir, language string) (*GenerateResponse, error) {
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

	url := stl.formatURL(&StainlessURLParams{
		OrgName:     orgName,
		ProjectName: projectName,
		Path:        "/generate",
	})
	req, err := stl.newAuthenticatedRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := stl.doAuthenticatedRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	response := &GenerateResponse{}
	json.Unmarshal(bodyBytes, response)

	return response, nil
}

func (stl *Stainless) ListSdks(orgName, projectName string) (*ListSdkResponse, error) {
	url := stl.formatURL(&StainlessURLParams{
		OrgName:     orgName,
		ProjectName: projectName,
		Path:        "/sdks",
	})

	req, err := stl.newAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := stl.doAuthenticatedRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	response := &ListSdkResponse{}
	json.Unmarshal(bodyBytes, response)

	return response, nil
}

func (stl *Stainless) RetrieveSdkStatus(orgName, projectName, language, branch string) (*RetrieveSdkBuildResponse, error) {
	path := fmt.Sprintf("/sdks/%s/builds/%s?fallbackBranch=main", language, branch)
	url := stl.formatURL(&StainlessURLParams{
		OrgName:     orgName,
		ProjectName: projectName,
		Path:        path,
	})

	req, err := stl.newAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := stl.doAuthenticatedRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	response := &RetrieveSdkBuildResponse{}
	json.Unmarshal(bodyBytes, response)

	return response, nil
}

func (stainless *Stainless) newAuthenticatedRequest(method, url string, body io.Reader) (*http.Request, error) {
	if stainless.AuthCookies == "" {
		log.Fatal("missing auth session. did you login?")
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("cookie", stainless.AuthCookies)

	return req, nil
}

func (stl *Stainless) doAuthenticatedRequest(r *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 401 {
			return nil, fmt.Errorf("check credentials, server responded with %s", resp.Status)
		}

		return nil, fmt.Errorf("server did not respond with 200 OK got %s", resp.Status)
	}

	return resp, nil
}

type StainlessURLParams struct {
	OrgName     string
	ProjectName string
	Path        string
}

func (stl *Stainless) formatURL(params *StainlessURLParams) string {
	if params.OrgName != "" && params.ProjectName != "" && params.Path != "" {
		return fmt.Sprintf("%s/orgs/%s/projects/%s%s", BASE_URL, params.OrgName, params.ProjectName, params.Path)
	} else if params.OrgName != "" {
		return fmt.Sprintf("%s/orgs/%s%s", BASE_URL, params.OrgName, params.Path)
	}

	return fmt.Sprintf("%s%s", BASE_URL, params.Path)
}
