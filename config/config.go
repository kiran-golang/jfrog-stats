package config

import (
	"encoding/json"
	"log"
	"os"
)

// Configuration loads up all the values that are used to configure
// service
type Configuration struct {
	CertAuthorityFile string `json:"caFile"`
	ServerCertificate string `json:"serverCert"`
	ServerPrivateKey  string `json:"serverKey"`
	ServicePort       string `json:"servicePort"`
	ArtifactoryURL    string `json:"artifactoryURL"`
	User              string `json:"user"`
	Password          string `json:"password"`
}

// Config is the structure that stores the configuration
var gConfig *Configuration

// readConfigFile reads the specified file to setup some env variables
func readConfigFile(file string) (*Configuration, error) {
	f, err := os.Open(file)
	if err != nil {
		return defaultConfiguration(), err
	}
	defer f.Close()

	// Setup some defaults here
	// If the json file has values in it, the defaults will be overwritten
	conf := defaultConfiguration()

	// Read the configuration from json file
	decoder := json.NewDecoder(f)
	err = decoder.Decode(conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

// defaultConfiguration sets the config to some default values
func defaultConfiguration() *Configuration {

	return &Configuration{
		CertAuthorityFile: "ca.cert",
		ServerCertificate: "server.cert",
		ServerPrivateKey:  "server.key",
		ServicePort:       "9000",
		ArtifactoryURL:    "http://localhost:8081/artifactory",
		User:              "",
		Password:          "",
	}
}

// GetConfiguration returns the configuration for the app.
// It will try to load it if it is not already loaded.
func GetConfiguration() *Configuration {
	if gConfig == nil {
		conf, err := readConfigFile("config.json")
		if err != nil {
			log.Println("Did not find a config.json file. Using defaults.")
		}
		gConfig = conf
	}

	return gConfig
}
