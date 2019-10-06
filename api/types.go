package api

// DownloadsResponse contains information about an artifact
// and the number of times it has been downloaded
type DownloadsResponse struct {
	RepoName     string `json:"repoName"`
	ArtifactName string `json:"artifactName"`
	Downloads    int    `json:"downloads"`
}
