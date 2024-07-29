package types

type RepoDetails struct {
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Author       string `json:"author"`
	CreatedAt    string `json:"created_at"`
	LastCommitAt string `json:"last_commit_at"`
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
}

type Objects struct {
	Type string `json:"type"`
	Path string `json:"path"`
}
