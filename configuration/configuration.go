package configuration

import (
	"encoding/json"
	"github.com/kabukky/journey/filenames"
	"io/ioutil"
	"log"
	"strings"
)

// Configuration: settings that are neccesary for server configuration
type Configuration struct {
	HttpHostAndPort  string
	HttpsHostAndPort string
	HttpsUsage       string
	Url              string
}

func NewConfiguration() *Configuration {
	var config Configuration
	err := config.load()
	if err != nil {
		log.Println("Warning: couldn't load " + filenames.ConfigFilename + ", creating new config file.")
		err = config.create()
		if err != nil {
			log.Fatal("Fatal error: Couldn't create configuration.")
			return nil
		}
		err = config.load()
		if err != nil {
			log.Fatal("Fatal error: Couldn't load configuration.")
			return nil
		}
	}
	return &config
}

// Global config - thread safe and accessible from all packages
var Config = NewConfiguration()

func (c *Configuration) save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filenames.ConfigFilename, data, 0600)
}

func (c *Configuration) load() error {
	data, err := ioutil.ReadFile(filenames.ConfigFilename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	// Make sure the url is in the right format
	if strings.HasPrefix(c.Url, "http://") {
		c.Url = strings.Replace(c.Url, "http://", "", 1)
	} else if strings.HasPrefix(c.Url, "https://") {
		c.Url = strings.Replace(c.Url, "https://", "", 1)
	}
	// Make sure there is no trailing slash at the end of the url
	if strings.HasSuffix(c.Url, "/") {
		c.Url = c.Url[0 : len(c.Url)-1]
	}
	return nil
}

func (c *Configuration) create() error {
	// TODO: Change default port
	c = &Configuration{HttpHostAndPort: ":8084", HttpsHostAndPort: ":8085", HttpsUsage: "None", Url: "127.0.0.1:8084"}
	err := c.save()
	if err != nil {
		log.Println("Error: couldn't create " + filenames.ConfigFilename)
		return err
	}

	return nil
}
