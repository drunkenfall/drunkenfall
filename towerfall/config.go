package towerfall

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type Config struct {
	Production bool `default:"true"`

	DbPath    string `default:"data/dev.db"`
	DbUser    string `default:"postgres"`
	DbName    string `default:"drunkenfall"`
	DbVerbose bool   `default:"false"`

	// Pointing to the test app
	FacebookID          string `default:"668534419991204"`
	FacebookSecret      string `default:"e74696c890216108c69d55d0e1b7866f"`
	FacebookCallback    string `default:"http://localhost/api/facebook/oauth2callback"`
	Port                int    `default:"42001"`
	RabbitURL           string `default:"amqp://rabbitmq:thiderman@drunkenfall.com:5672/"`
	RabbitIncomingQueue string `default:"drunkenfall-app-dev"`
	RabbitOutgoingQueue string `default:"drunkenfall-game-dev"`
	oauthConf           *oauth2.Config
	log                 *zap.Logger
}

func ParseConfig() *Config {
	c := Config{}

	envconfig.MustProcess("drunkenfall", &c)

	if c.Production {
		c.log, _ = zap.NewProduction()
	} else {
		c.log, _ = zap.NewDevelopment()
	}

	c.parseOauth()
	return &c
}

// Print prints a visualization of what's going on
func (c *Config) Print() {
	c.log.Info("Configuration loaded", zap.Any("config", c))
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
