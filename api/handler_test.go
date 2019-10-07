package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRTManager struct {
	items []byte
	err   error
}

func (m *mockRTManager) Aql(aqlstring string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.items, nil
}

func TestGetDownloadsHandler(t *testing.T) {
	testCases := []struct {
		label        string
		reponame     string
		mockItems    []byte
		mockError    error
		expectedCode int
		expectedJSON []StatDownloads
	}{
		{
			label:    "Positive test with successful GET",
			reponame: "jcenter-cache",
			mockItems: []byte(`
			{
				"results" : [ {
				  "repo" : "jcenter-cache",
				  "path" : "org/apache/struts/struts2-core/2.3.14",
				  "name" : "struts2-core-2.3.14.pom",
				  "type" : "file",
				  "size" : 12241,
				  "created" : "2019-04-22T22:25:35.107Z",
				  "created_by" : "anonymous",
				  "modified" : "2013-03-28T21:55:51.000Z",
				  "modified_by" : "anonymous",
				  "updated" : "2019-04-22T22:25:35.108Z",
				  "stats" : [ {
					"downloaded" : "2019-10-04T16:40:43.191Z",
					"downloaded_by" : "anonymous",
					"downloads" : 27,
					"remote_downloads" : 0
				  } ]
				},{
				  "repo" : "jcenter-cache",
				  "path" : "org/apache/struts/struts-master/9",
				  "name" : "struts-master-9.pom",
				  "type" : "file",
				  "size" : 10260,
				  "created" : "2019-04-22T22:25:35.774Z",
				  "created_by" : "anonymous",
				  "modified" : "2012-02-28T11:11:07.000Z",
				  "modified_by" : "anonymous",
				  "updated" : "2019-04-22T22:25:35.775Z",
				  "stats" : [ {
					"downloaded" : "2019-10-04T16:40:43.245Z",
					"downloaded_by" : "anonymous",
					"downloads" : 27,
					"remote_downloads" : 0
				  } ]
				} ],
				"range" : {
				  "start_pos" : 0,
				  "end_pos" : 2,
				  "total" : 2,
				  "limit" : 2
				}
			}`),
			expectedCode: http.StatusOK,
			mockError:    nil,
		},
		{
			label:        "Failure case where the repo does not exist",
			reponame:     "repo",
			expectedCode: http.StatusInternalServerError,
			mockError:    errors.New("Error processing request"),
		},
	}

	h := &handler{}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/v1/stats/downloads/"+tc.reponame, nil)
			recorder := httptest.NewRecorder()
			h.rtManager = &mockRTManager{
				items: tc.mockItems,
			}
			NewRouter(h).ServeHTTP(recorder, request)
			resp := recorder.Result()

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("Expected %d; Got %d status code", tc.expectedCode, resp.StatusCode)
			}
		})
	}
}
