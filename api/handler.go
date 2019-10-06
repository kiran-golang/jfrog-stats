package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/kiran-golang/jfrog-stats/config"

	"github.com/gorilla/mux"
	"github.com/jfrog/jfrog-client-go/artifactory"
	jfrogAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	jfroglog "github.com/jfrog/jfrog-client-go/utils/log"
	pkgerrors "github.com/pkg/errors"
)

type handler struct {
	rtManager *artifactory.ArtifactoryServicesManager
}

// NewHandler creates a handler that supports the typical CRUD operations
func newHandler() *handler {
	return &handler{}
}

func (h *handler) initArtifactoryConnection() error {

	if h.rtManager != nil {
		return nil
	}

	jfroglog.SetLogger(jfroglog.NewLogger(jfroglog.INFO, os.Stdout))

	var err error
	artifactoryDetails := jfrogAuth.NewArtifactoryDetails()
	artifactoryDetails.SetUrl(config.GetConfiguration().ArtifactoryURL)
	artifactoryDetails.SetUser(config.GetConfiguration().User)
	artifactoryDetails.SetPassword(config.GetConfiguration().Password)

	serviceConfig, err := artifactory.NewConfigBuilder().SetArtDetails(artifactoryDetails).Build()
	if err != nil {
		return pkgerrors.Wrap(err, "Initializing Service Configuration")
	}

	h.rtManager, err = artifactory.New(&artifactoryDetails, serviceConfig)
	if err != nil {
		return pkgerrors.Wrap(err, "Initializing Artifactory Service Manager")
	}

	return nil
}

func (h *handler) getDownloadsHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	repoName := vars["repo-name"]

	limitStr := r.FormValue("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ret, err := h.processGetDownloads(repoName, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// processGetDownloads builds an AQL query to get the desired information
// from the artifacts' server and returns the relevant information
func (h *handler) processGetDownloads(repoName string, limit int) ([]DownloadsResponse, error) {

	err := h.initArtifactoryConnection()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Connecting to Artifactory")
	}

	aqlString := fmt.Sprintf(`items.find(
		{"repo":"%s"},
		{"stat.downloads":{"$gte":"1"}}
	)
	.include("stat")
	.sort({"$desc":["stat.downloads"]})
	.limit(%d)`, repoName, limit)

	fmt.Println(aqlString)

	data, err := h.rtManager.Aql(aqlString)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Executing AQL query")
	}

	results := struct {
		Results []struct {
			Name  string `json:"name"`
			Stats []struct {
				Downloads int `json:"downloads"`
			} `json:"stats"`
		} `json:"results"`
	}{}

	err = json.Unmarshal(data, &results)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Unmarshaling AQL Results")
	}

	ret := []DownloadsResponse{}

	for _, artifact := range results.Results {
		ret = append(ret, DownloadsResponse{
			RepoName:     repoName,
			ArtifactName: artifact.Name,
			Downloads:    artifact.Stats[0].Downloads,
		})
	}

	return ret, nil
}
