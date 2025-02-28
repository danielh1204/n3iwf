/*
 * N3IWF Configuration Factory
 */

package factory

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/asaskevich/govalidator"
	yaml "gopkg.in/yaml.v2"

	"github.com/free5gc/n3iwf/internal/logger"
)

var N3iwfConfig *Config

// TODO: Support configuration update from REST api
func InitConfigFactory(f string, cfg *Config) error {
	if f == "" {
		// Use default config path
		f = N3iwfDefaultConfigPath
	}

	if content, err := ioutil.ReadFile(f); err != nil {
		return fmt.Errorf("[Factory] %+v", err)
	} else {
		logger.CfgLog.Infof("Read config from [%s]", f)
		if yamlErr := yaml.Unmarshal(content, cfg); yamlErr != nil {
			return fmt.Errorf("[Factory] %+v", yamlErr)
		}
	}

	//change sd to lowercase
	SupportedTAList := cfg.Configuration.N3IWFInfo.SupportedTAList
	for i := range SupportedTAList {
		BroadcastPLMNList := SupportedTAList[i].BroadcastPLMNList
		for j := range BroadcastPLMNList {
			TAISliceSupportList := BroadcastPLMNList[j].TAISliceSupportList
			for k := range TAISliceSupportList {
				TAISliceSupportList[k].SNSSAI.SD = strings.ToLower(TAISliceSupportList[k].SNSSAI.SD)
			}
		}
	}

	return nil
}

func ReadConfig(cfgPath string) (*Config, error) {
	cfg := &Config{}
	if err := InitConfigFactory(cfgPath, cfg); err != nil {
		return nil, fmt.Errorf("ReadConfig [%s] Error: %+v", cfgPath, err)
	}
	if _, err := cfg.Validate(); err != nil {
		validErrs := err.(govalidator.Errors).Errors()
		for _, validErr := range validErrs {
			logger.CfgLog.Errorf("%+v", validErr)
		}
		logger.CfgLog.Errorf("[-- PLEASE REFER TO SAMPLE CONFIG FILE COMMENTS --]")
		return nil, fmt.Errorf("Config validate Error")
	}

	return cfg, nil
}
