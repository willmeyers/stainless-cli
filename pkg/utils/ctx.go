package utils

type key int

const (
	CredentialsCtxKey key = iota
	OrgCtxKey         key = iota
	ProjectCtxKey     key = iota
)
