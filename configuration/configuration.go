package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/kabukky/journey/filenames"
)

// Configuration settings that are neccesary for server configuration
type Configuration struct {
	HTTPHostAndPort  string
	HTTPSHostAndPort string
	HTTPSUsage       string
	URL              string
	HTTPSUrl         string
	UseLetsEncrypt   bool
	SAMLCert         string
	SAMLKey          string
	SAMLIDPUrl       string
}

// NewConfiguration loads the configuration from config.json and returns it
// It will create a new, empty configuration with defaults if it doesn't
// exist yet.  // It dies with a fatal error if the configuration file can't
// be parsed
func NewConfiguration() *Configuration {
	var config Configuration
	err := config.load()
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Error reading configuration file %s: %s",
				filenames.ConfigFilename, err)
		}
		log.Printf("%s does not exist; creating new config file", filenames.ConfigFilename)
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

// Config is thread safe and accessible from all packages
var Config = NewConfiguration()

func (c *Configuration) save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filenames.ConfigFilename, data, 0600)
}

func (c *Configuration) load() error {
	configWasChanged := false
	data, err := ioutil.ReadFile(filenames.ConfigFilename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	// Make sure the url is in the right format
	// Make sure there is no trailing slash at the end of the url
	if strings.HasSuffix(c.URL, "/") {
		c.URL = c.URL[0 : len(c.URL)-1]
		configWasChanged = true
	}
	if !strings.HasPrefix(c.URL, "http://") && !strings.HasPrefix(c.URL, "https://") {
		c.URL = "http://" + c.URL
		configWasChanged = true
	}
	// Make sure the https url is in the right format
	// Make sure there is no trailing slash at the end of the https url
	if strings.HasSuffix(c.HTTPSUrl, "/") {
		c.HTTPSUrl = c.HTTPSUrl[0 : len(c.HTTPSUrl)-1]
		configWasChanged = true
	}
	if strings.HasPrefix(c.HTTPSUrl, "http://") {
		c.HTTPSUrl = strings.Replace(c.HTTPSUrl, "http://", "https://", 1)
		configWasChanged = true
	} else if !strings.HasPrefix(c.HTTPSUrl, "https://") {
		c.HTTPSUrl = "https://" + c.HTTPSUrl
		configWasChanged = true
	}
	// Make sure there is no trailing slash at the end of the url
	if strings.HasSuffix(c.HTTPSUrl, "/") {
		c.HTTPSUrl = c.HTTPSUrl[0 : len(c.HTTPSUrl)-1]
		configWasChanged = true
	}
	// Check if all fields are filled out
	cReflected := reflect.ValueOf(*c)
	for i := 0; i < cReflected.NumField(); i++ {
		if cReflected.Field(i).Interface() == "" &&
			!strings.HasPrefix(cReflected.Type().Field(i).Name, "SAML") {
			return fmt.Errorf("Error: file %s missing required field %s", filenames.ConfigFilename, cReflected.Type().Field(i).Name)
		}
	}
	// Save the changed config
	if configWasChanged {
		err = c.save()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Configuration) create() error {
	// TODO: Change default port
	c = &Configuration{HTTPHostAndPort: ":8084", HTTPSHostAndPort: ":8085", HTTPSUsage: "None", URL: "127.0.0.1:8084", HTTPSUrl: "127.0.0.1:8085"}
	err := c.save()
	if err != nil {
		log.Println("Error: couldn't create " + filenames.ConfigFilename)
		return err
	}

	return nil
}
