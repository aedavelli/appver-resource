package models

type AppVersion struct {
	Version string `json:"version"`
}

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Metadata []MetadataField

type InRequest struct {
	Source  Source     `json:"source"`
	Version AppVersion `json:"version"`
}

type InResponse struct {
	Version  AppVersion `json:"version"`
	Metadata Metadata   `json:"metadata"`
}

type CheckRequest struct {
	Source  Source     `json:"source"`
	Version AppVersion `json:"version"`
}

type CheckResponse []AppVersion

type Source struct {
	Url          string `json:"url"`
	VersionField string `json:"version_field"`
	UserName     string `json:"username"`
	Password     string `json:"password"`
	Accept       string `json:"accept"`
	Insecure     bool   `json:"insecure"`
}
