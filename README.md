# jfrog-stats
Project that queries Jfrog REST API to query statistics on artifacts.

The following concepts are implemented in this demo project.

* A go microservice with a REST API
* Integration of the microservice with a Jenkins build
* Creation and Distribution of Docker Image artifacts to an Artifactory Docker Registry

## Installation

### Run as a Docker Container
A Dockerfile is included in the repo. It uses a multi-stage docker pipeline to build the docker image.

```bash
# BUILD IMAGE
docker build -t jfrog-stats:latest -f Dockerfile .

# RUN the Service with a config.json
cat << EOF > config.json
{
    "artifactoryURL": "http://34.68.140.218/artifactory/",
    "user": "<USERNAME>",
    "password": "<PASSWORD>"
}
EOF
docker create --name jfrog-stats jfrog-stats:latest
docker cp config.json jfrog-stats:/opt/jfrog-stats/
docker start jfrog-stats

# View Logs
docker logs -f jfrog-stats
```

### Run in a Kubernetes Environment
Helm charts have been provided and they can be used to run the service in Kubernetes.

```bash
cd deployment/helm
mkdir manifests

# Create the Kubernetes manifest files. Assume Tiller is not available.
helm template --namespace testns -n test jfrog-stats --set config.user="admin" --set config.password="eUhCbPG3mK" --output-dir manifests

# Apply the created manifests
kubectl -n testns apply --recursive -f manifests
```

## Configuration
The application supports the following configuration parameters provided in a config.json file in its WORKDIR

```bash
// Certificates for starting the service with https
"caFile"
"serverCert"
"serverKey"

// Port on which to start the service. Default is 9000
"servicePort"

// URL to artifactory
"artifactoryURL"

// Username and Password to access the artifactory
"user"
"password"
```

### Enable HTTPS endpoint

The service can be started as a https service if the caFile, serverCert, serverKey files are provided. If the serverKey is encrypted, the module looks for a file with the name <serverkey>.pass in the WORKDIR and then loads that to decrypt the .key

## APIs

The API has the following endpoints

```go
/v1/stats/downloads/{repo-name}
/v1/healthcheck
```

### Example Calls

```bash
curl -s -X GET http://localhost:9000/v1/stats/download/jcenter-cache
```

The API response will look something like this:

```json
[
  {
    "repoName": "jcenter-cache",
    "artifactName": "struts2-core-2.3.14.pom",
    "downloads": 27
  },
  {
    "repoName": "jcenter-cache",
    "artifactName": "struts-master-9.pom",
    "downloads": 27
  }
]
```

## Third Party Packages used

A few third party packages were used in the development of this micro service.
They are detailed below:
- [Gorilla Mux and Handler](github.com/gorilla)
- [Errors Formatter](github.com/pkg/errors)
- [JFrog Go Client](github.com/jfrog/jfrog-client-go)
