package types

type RepoResponse struct {
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Author       string `json:"author"`
	CreatedAt    string `json:"created_at"`
	LastCommitAt string `json:"last_commit_at"`
}
