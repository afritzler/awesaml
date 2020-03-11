package main

import (
	"os"
	"testing"

	"github.com/afritzler/awesaml/pkg/types"
)

func TestInitVarsNoEnv(t *testing.T) {
	err := initVars()
	if err == nil {
		t.Errorf("No env should cause an error")
	}
}

var testEnvs = map[string]string{
	types.EntityIDEnvName: "myEntityID",
	types.CertFileEnvName: "myservice.cert",
	types.KeyFileEnvName: "myservice.key",
	types.ServiceURLEnvName: "http://localhost:8000",
	types.IdpMetaDataURLEnvName: "http://myAwesomeIDP/saml2/metadata",
	types.ContentDirEnvName: "public",
	types.ServingPortEnvName: "8080",
}

func TestInitVarsWithEnvs(t *testing.T) {
	for env, value := range testEnvs{
		os.Setenv(env, value)
	}
	err := initVars()
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	if entityID != testEnvs[types.EntityIDEnvName] {
		t.Error("myEntityID is not set")
		return
	}
	if certFile != testEnvs[types.CertFileEnvName] {
		t.Error("certFile is not set")
		return
	}
	if certSecretName != "" {
		t.Error("certSecretName should not be set")
		return
	}
	if serviceURL != testEnvs[types.ServiceURLEnvName]{
		t.Error("serviceURL is not set")
		return
	}
	if idpMetaDataURL != testEnvs[types.IdpMetaDataURLEnvName]{
		t.Error("idpMetaDataURL is not set")
		return
	}
	if contentDir != testEnvs[types.ContentDirEnvName]{
		t.Error("contentDir is not set")
		return
	}
	if servingPort != testEnvs[types.ServingPortEnvName]{
		t.Error("servingPort is not set")
		return
	}
}