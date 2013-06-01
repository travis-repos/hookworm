package hookworm

type Payload struct {
	Ref        string      `json:"ref"`
	After      string      `json:"after"`
	Before     string      `json:"before"`
	Created    bool        `json:"created"`
	Deleted    bool        `json:"deleted"`
	Forced     bool        `json:"forced"`
	Compare    string      `json:"compare"`
	Commits    []string    `json:"commits"`
	HeadCommit *HeadCommit `json:"head_commit"`
	Repository *Repository `json:"repository"`
	Pusher     *Pusher     `json:"pusher"`
}

type HeadCommit struct {
	Id        string   `json:"id"`
	Distinct  bool     `json:"distinct"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	Url       string   `json:"url"`
	Author    *Author  `json:"author"`
	Committer *Author  `json:"committer"`
	Added     []string `json:"added"`
	Removed   []string `json:"removed"`
	Modified  []string `json:"modified"`
}

type Repository struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Url          string `json:"url"`
	Description  string `json:"description"`
	Watchers     int64  `json:"watchers"`
	Stargazers   int64  `json:"stargazers"`
	Forks        int64  `json:"forks"`
	Fork         bool   `json:"fork"`
	Size         int64  `json:"size"`
	Owner        *Owner `json:"owner"`
	Private      bool   `json:"private"`
	OpenIssues   int64  `json:"open_issues"`
	HasIssues    bool   `json:"has_issues"`
	HasDownloads bool   `json:"has_downloads"`
	HasWiki      bool   `json:"has_wiki"`
	Language     string `json:"language"`
	CreatedAt    int64  `json:"created_at"`
	PushedAt     int64  `json:"pushed_at"`
	MasterBranch string `json:"master_branch"`
	Organization string `json:"organization"`
}

type Pusher struct {
	Name string `json:"name"`
}

type Author struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Owner struct {
	Name  string          `json:"name"`
	Email *NullableString `json:"email"`
}

type NullableString struct {
	string
}
