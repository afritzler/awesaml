package main

import (
	"os"
	"testing"

	"github.com/afritzler/awesaml/pkg/utils"
)

func TestInitVarsNoEnv(t *testing.T) {
	err := initVars()
	if err == nil {
		t.Errorf("No env should cause an error")
	}
}

var testEnvs = map[string]string{
	utils.EntityIDEnvName:       "myEntityID",
	utils.CertFileEnvName:       "myservice.cert",
	utils.KeyFileEnvName:        "myservice.key",
	utils.ServiceURLEnvName:     "http://localhost:8000",
	utils.IdpMetaDataURLEnvName: "http://myAwesomeIDP/saml2/metadata",
	utils.ContentDirEnvName:     "public",
	utils.ServingPortEnvName:    "8080",
}

func TestInitVarsWithEnvs(t *testing.T) {
	for env, value := range testEnvs {
		os.Setenv(env, value)
	}
	err := initVars()
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	if entityID != testEnvs[utils.EntityIDEnvName] {
		t.Error("myEntityID is not set")
		return
	}
	if certFile != testEnvs[utils.CertFileEnvName] {
		t.Error("certFile is not set")
		return
	}
	if certSecretName != "" {
		t.Error("certSecretName should not be set")
		return
	}
	if serviceURL != testEnvs[utils.ServiceURLEnvName] {
		t.Error("serviceURL is not set")
		return
	}
	if idpMetaDataURL != testEnvs[utils.IdpMetaDataURLEnvName] {
		t.Error("idpMetaDataURL is not set")
		return
	}
	if contentDir != testEnvs[utils.ContentDirEnvName] {
		t.Error("contentDir is not set")
		return
	}
	if servingPort != testEnvs[utils.ServingPortEnvName] {
		t.Error("servingPort is not set")
		return
	}
}
