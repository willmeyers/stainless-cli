package stainless

import "time"

// Builds

type Build struct {
	ID          int        `json:"id"`
	Status      string     `json:"status"`
	StartedAt   string     `json:"startedAt"`
	EndedAt     string     `json:"endedAt"`
	Org         string     `json:"org"`
	Project     string     `json:"project"`
	Sdks        []BuildSDK `json:"sdks"`
	TriggeredBy string     `json:"triggeredBy"`
}

type BuildSDK struct {
	Language string `json:"language"`
	Status   string `json:"status"`
}

type ListBuildsResponse struct {
	Builds          []Build `json:"builds"`
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
}

type GetBuildStatusResponse struct {
	ID                  uint64    `json:"id"`
	Status              string    `json:"status"`
	DiagnosticsFileURL  string    `json:"diagnosticsFileURL"`
	StartedGeneratingAt time.Time `json:"startedGeneratingAt"`
	EndedAt             time.Time `json:"endedAt"`
	HasGenerated        bool      `json:"hasGenerated"`
}

// Orgs

type OrgItem struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type ListOrgsResponse struct {
	StartCursor     string    `json:"startCursor"`
	EndCursor       string    `json:"EndCursor"`
	HasNextPage     bool      `json:"hasNextPage"`
	HasPreviousPage bool      `json:"hasPreviousPage"`
	Items           []OrgItem `json:"items"`
}

// Projects

type ProjectItem struct {
	Name string `json:"name"`
	Org  string `json:"org"`
}

type ListProjectsResponse struct {
	StartCursor     string        `json:"startCursor"`
	EndCursor       string        `json:"EndCursor"`
	HasNextPage     bool          `json:"hasNextPage"`
	HasPreviousPage bool          `json:"hasPreviousPage"`
	Items           []ProjectItem `json:"items"`
}

// API Keys

type ListApiKeysResponse []struct {
	ID            string    `json:"id"`
	Key           string    `json:"key"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	DeletedAt     time.Time `json:"deletedAt"`
	CreatedByUser struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		ImageURL string `json:"imageURL"`
	} `json:"createdByUser"`
}

// Members

type ListMembersResponse []struct {
	UserID   string `json:"userID"`
	InviteID string `json:"inviteID"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageURL string `json:"imageURL"`
	Role     string `json:"role"`
}

// SDKs

type SdkListItem struct {
	ID                    int    `json:"id"`
	Org                   string `json:"org"`
	Project               string `json:"project"`
	Language              string `json:"language"`
	InternalRepositoryURL string `json:"internalRepositoryUrl"`
	AutoMergeEnabled      bool   `json:"autoMergeEnabled"`
	IsLive                bool   `json:"isLive"`
	ReleaseInformation    struct {
		PackageName string `json:"packageName"`
		Repo        string `json:"repo"`
	} `json:"releaseInformation"`
	CustomCodeInformation struct {
		GistURL      string `json:"gistUrl"`
		LinesAdded   int    `json:"linesAdded"`
		LinesRemoved int    `json:"linesRemoved"`
	} `json:"customCodeInformation"`
}

type ListSdksResponse struct {
	Items []SdkListItem `json:"items"`
}

// Generate

type GenerateResponse struct {
	ProjectName string `json:"projectName"`
	Generate    bool   `json:"generated"`
	NewSha      string `json:"newSha"`
	HasMerged   bool   `json:"hasMerged"`
}
