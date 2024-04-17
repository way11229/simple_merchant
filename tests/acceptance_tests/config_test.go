package acceptance_tests

import "github.com/way11229/simple_merchant/config"

var (
	accessToken = ""
)

type testConfig struct {
	config.Config

	AccessToken string `mapstructure:"ACCESS_TOKEN"`
}
