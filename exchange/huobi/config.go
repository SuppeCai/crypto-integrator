package huobi

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strings"
	"integrator/exchange"
)

var DefaultConfigFile = "conf/huobi.yaml"
var DefaultSymbolPlaceHolder = "$symbol"
var DefaultPeriodPlaceHolder = "$period"
var DefaultPeriodSeparator = "->"

type Config struct {
	S Sub `yaml:"sub"`
}

type Sub struct {
	Name        string              `yaml:"name"`
	Symbol      []string            `yaml:"symbol"`
	Period      []string            `yaml:"period"`
	ExtraPeriod map[string][]string `yaml:"extra_period"`
}

func (huobi *Huobi) initConfig() error {

	if huobi.configFile == "" {
		huobi.configFile = DefaultConfigFile
	}

	configFile, err := ioutil.ReadFile(huobi.configFile)
	if err != nil {
		exchange.LogErr.Error("config file error:" + err.Error())
	}

	err = yaml.Unmarshal(configFile, &huobi.config)
	if err != nil {
		exchange.LogErr.Error("yaml unmarshal error:" + err.Error())
	}

	if huobi.config.S.Symbol == nil || huobi.config.S.Period == nil {
		exchange.LogErr.Error("no subscription in config:" + err.Error())
		return err
	}

	huobi.subList = huobi.subList[:0:0]
	for _, symbol := range huobi.config.S.Symbol {
		for _, period := range huobi.config.S.Period {
			name := strings.Replace(huobi.config.S.Name, DefaultSymbolPlaceHolder, symbol, 1)
			name = strings.Replace(name, DefaultPeriodPlaceHolder, period, 1)
			huobi.subList = append(huobi.subList, name)
		}
	}

	return err
}

func (huobi *Huobi) subscribe() error {

	if huobi.config.S.Symbol == nil || huobi.config.S.Period == nil {
		exchange.LogErr.Error("no subscription in config")
		return nil
	}

	for _, name := range huobi.subList {
		sub := SubKline{Sub: name, Id: exchange.AppName}
		huobi.Send(sub)
	}

	return nil
}
