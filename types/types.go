package types

type Config struct {
	Username string
	Password string
	Shell    string
	Apps     map[string]App
}

type App struct {
	Name     string
	Path     string
	Commands []string
	Secret   string
}

type GithubRelease struct {
	TagName string `json:"tag_name"`
	Url     string `json:"url"`
}

type GithubReleasePayload struct {
	Action  string `json:"action"`
	Release GithubRelease
}
