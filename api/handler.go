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

// artifactoryServicesManagerInterface exists to allow mocking the manager calls
// by providing an interface that can be implemented to return static data in tests
type artifactoryServicesManagerInterface interface {
	Aql(aql string) ([]byte, error)
}

type handler struct {
	rtManager artifactoryServicesManagerInterface
}

// initArtifactoryConnection initializes the manager that connects to jfrog
// artifactory
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

// getDownloadsHandler handles both the queries based request as well as
// requests where a limit filter is not provided
func (h *handler) getDownloadsHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	repoName := vars["repo-name"]

	// init limit to default of 2
	limit := 2
	limitStr := r.FormValue("limit")
	if limitStr != "" {
		// Override limit if provided
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
func (h *handler) processGetDownloads(repoName string, limit int) ([]StatDownloads, error) {

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

	data, err := h.rtManager.Aql(aqlString)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Executing AQL query")
	}

	// We are only interested in the Name and Stats.Downloads fields
	// from the returned data
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

	ret := []StatDownloads{}

	for _, artifact := range results.Results {
		ret = append(ret, StatDownloads{
			RepoName:     repoName,
			ArtifactName: artifact.Name,
			Downloads:    artifact.Stats[0].Downloads,
		})
	}

	return ret, nil
}
