package types

type RepoDetails struct {
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	DefaultBranch string `json:"defaultBranch"`
	Author        string `json:"author"`
	CreatedAt     string `json:"created_at"`
	LastCommitAt  string `json:"last_commit_at"`
}

type JsonResponse[T any] struct {
	RequestStatus int    `json:"requestStatus"`
	StatusCode    int    `json:"statusCode"`
	Msg           string `json:"message"`
	Data          T      `json:"data"`
}

type RepoObjects struct {
	Name     string    `json:"name"`
	Desc     string    `json:"desc"`
	Branches []string  `json:"branches"`
	Data     []Objects `json:"objects"`
	Blob     []string  `json:"blob"`
}

type Objects struct {
	Type     string `json:"type"`
	FullPath string `json:"fullPath"`
	Path     string `json:"path"`
	Size     string `json:"size"`
}

type Log struct {
	Hash      string `json:"commit_hash"`
	Author    string `json:"commit_author"`
	Message   string `json:"commit_message"`
	Timestamp string `json:"commit_timestamp"`
	LogStat
}

type CommitHistory struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Logs []Log  `json:"logs"`
}

type LogStat struct {
	FilesChanged string `json:"files_changed"`
	Deletions    string `json:"deletions"`
	Insertions   string `json:"insertions"`
}

type DiffResponse struct {
    FilePath  string `json:"filePath"`
    CodeLines []string `json:"Code_lines"`
}
