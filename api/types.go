package api

// StatDownloads contains information about an artifact
// and the number of times it has been downloaded
type StatDownloads struct {
	RepoName     string `json:"repoName"`
	ArtifactName string `json:"artifactName"`
	Downloads    int    `json:"downloads"`
}
