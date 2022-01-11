package cmd

import (
	"forklift/pkg/config"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

var (
	cfgFile string
	RootCmd = &cobra.Command{}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra-example.yaml)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	var (
		err           error
		configYaml    []byte
		configuration config.Configuration
	)

	if cfgFile == "" {
		log.Fatalf("Config file name is empty!")
	}

	if configYaml, err = ioutil.ReadFile(cfgFile); err != nil {
		log.Fatalf("Configuration could not be read!")
	}

	if err = yaml.Unmarshal(configYaml, &configuration); err != nil {
		panic(err)
	}

	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err == nil {
		log.Warningf("Using config file: %s", viper.ConfigFileUsed())
	}
}
