package binance

import (
	"gopkg.in/yaml.v2"
	"integrator/exchange"
	"io/ioutil"
	"strings"
)

var DefaultConfigFile = "conf/binance.yaml"
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

func (binance *Binance) initConfig() error {

	if binance.configFile == "" {
		binance.configFile = DefaultConfigFile
	}

	configFile, err := ioutil.ReadFile(binance.configFile)
	if err != nil {
		exchange.LogErr.Error("config file error:" + err.Error())
	}

	err = yaml.Unmarshal(configFile, &binance.config)
	if err != nil {
		exchange.LogErr.Error("yaml unmarshal error:" + err.Error())
	}

	if binance.config.S.Symbol == nil || binance.config.S.Period == nil {
		exchange.LogErr.Error("no subscription in config:" + err.Error())
		return err
	}

	binance.subList = binance.subList[:0:0]
	for _, symbol := range binance.config.S.Symbol {
		for _, period := range binance.config.S.Period {
			name := strings.Replace(binance.config.S.Name, DefaultSymbolPlaceHolder, symbol, 1)
			name = strings.Replace(name, DefaultPeriodPlaceHolder, period, 1)
			binance.subList = append(binance.subList, name)
		}
	}

	return err
}
