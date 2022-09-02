package app_test

import (
	"github.com/keskad/webhook-conversion-service/pkg/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigurationFromFile_LoadsExample(t *testing.T) {
	config, err := app.LoadConfigurationFromFile("../../example-config.yaml")

	assert.Nil(t, err)
	assert.Len(t, config.Endpoints, 2)
	assert.Len(t, config.Endpoints[0].Replacements, 1)
	assert.Equal(t, config.Endpoints[0].Path, "/api/webhook")
	assert.Equal(t, config.Endpoints[0].TargetUrl, "http://argocd-server.argocd.svc.cluster.local")
}

func TestLoadConfigurationFromFile_MissingFile(t *testing.T) {
	_, err := app.LoadConfigurationFromFile("invalid-file-path.yaml")
	assert.Error(t, err)
}

func TestLoadConfigurationFromFile_IncorrectSyntax(t *testing.T) {
	_, err := app.LoadConfigurationFromFile("config_test.go")
	assert.Error(t, err)
}
