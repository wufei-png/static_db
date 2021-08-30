package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func LoadConf(confPath string, confPtr interface{}) error {
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}

	if strings.HasSuffix(confPath, "json") {
		if err := json.Unmarshal(data, confPtr); err != nil {
			return err
		}
	} else if strings.HasSuffix(confPath, "yml") || strings.HasSuffix(confPath, "yaml") {
		if err := yaml.Unmarshal(data, confPtr); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("conf(%v) is not json or yaml", confPath)
	}

	return nil
}

func GetConfigString(conf interface{}) (configString string, prettyString string) {
	if _, ok := conf.(string); ok {
		configString = conf.(string)
	} else {
		var err error
		configString, err = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(conf)
		if err != nil {
			log.Fatal("MarshalToString ", err)
		}
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(configString), "", "\t"); err != nil {
		log.Fatal("JSON parse error: ", err)
	}

	return configString, prettyJSON.String()
}

// DeepCopy deep copy a to b using json marshaling
func DeepCopy(a, b interface{}) error {
	tmp, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return json.Unmarshal(tmp, b)
}
