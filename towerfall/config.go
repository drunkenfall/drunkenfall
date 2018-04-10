package towerfall

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type Config struct {
	DbPath string `default:"data/dev.db"`
	// Pointing to the test app
	FacebookID       string `default:"668534419991204"`
	FacebookSecret   string `default:"e74696c890216108c69d55d0e1b7866f"`
	FacebookCallback string `default:"http://localhost/api/facebook/oauth2callback"`
	Port             int    `default:"42001"`
	oauthConf        *oauth2.Config
}

func ParseConfig() *Config {
	ret := Config{}

	envconfig.MustProcess("drunkenfall", &ret)
	ret.parseOauth()
	log.Printf("Configuration loaded: %+v", ret)
	return &ret
}

func (c *Config) parseOauth() {
	c.oauthConf = &oauth2.Config{
		ClientID:     c.FacebookID,
		ClientSecret: c.FacebookSecret,
		RedirectURL:  c.FacebookCallback,
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
}
