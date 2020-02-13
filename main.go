package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/afritzler/awesaml/pkg/types"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/crewjam/saml/samlsp"
)

var (
	entityID, certFile, certSecretName, keyFile, keySecretName, serviceURL, idpMetaDataURL, contentDir, servingPort string
)

func main() {
	// init vars or panic
	initVars()

	keyPair, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err) // TODO handle error
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err) // TODO handle error
	}

	rootURL, _ := url.Parse(serviceURL)
	idpMetadataURL, _ := url.Parse(idpMetaDataURL)

	idpMetadata, err := samlsp.FetchMetadata(
		context.Background(),
		http.DefaultClient,
		*idpMetadataURL)
	if err != nil {
		panic(err) // TODO handle error
	}

	samlSP, err := samlsp.New(samlsp.Options{
		EntityID:    entityID,
		ForceAuthn:  true,
		URL:         *rootURL,
		IDPMetadata: idpMetadata,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
	})
	if err != nil {
		panic(err) // TODO handle error
	}

	log.Println("Starting Service Provider ...")
	log.Printf("EntityID: %s\n", entityID)
	log.Printf("IDPMetaDataURL: %s\n", idpMetaDataURL)
	log.Printf("Listening on %s\n", serviceURL)
	log.Printf("Serving content from %s\n", contentDir)
	fs := http.FileServer(http.Dir(contentDir))

	http.Handle("/", samlSP.RequireAccount(fs))
	http.Handle("/saml/", samlSP)
	http.ListenAndServe(fmt.Sprintf(":%s", servingPort), nil)
}

func initVars() {
	entityID = os.Getenv(types.EntityIDEnvName)
	if len(entityID) == 0 {
		panic(fmt.Sprintf("you need to provide the entityID by exporting it via the following env var %s", types.EntityIDEnvName))
	}
	certFile = os.Getenv(types.CertFileEnvName)
	if len(certFile) == 0 {
		certSecretName = os.Getenv(types.CertSecretNameEnvName)
		if len(certSecretName) == 0 {
			panic(fmt.Sprintf("you need to provide either a location of the cert file via %s or cert secret name via %s", types.CertFileEnvName, types.CertSecretNameEnvName))
		}
	}
	keyFile = os.Getenv(types.KeyFileEnvName)
	if len(keyFile) == 0 {
		keySecretName = os.Getenv(types.KeySecretNameEnvName)
		if len(keySecretName) == 0 {
			panic(fmt.Sprintf("you need to provide either a location of the key file via %s or key secret name via %s", types.KeyFileEnvName, types.KeySecretNameEnvName))
		}
	}
	serviceURL = os.Getenv(types.ServiceURLEnvName)
	if len(serviceURL) == 0 {
		panic(fmt.Sprintf("you need to provide the serviceURL by exporting it via the following env var %s", types.ServiceURLEnvName))
	}
	idpMetaDataURL = os.Getenv(types.IdpMetaDataURLEnvName)
	if len(idpMetaDataURL) == 0 {
		panic(fmt.Sprintf("you need to provide the idpMetaDataURL by exporting it via the following env var %s", types.IdpMetaDataURLEnvName))
	}
	contentDir = os.Getenv(types.ContentDirEnvName)
	if len(contentDir) == 0 {
		panic(fmt.Sprintf("you need to provide the contentDir by exporting it via the following env var %s", types.ContentDirEnvName))
	}
	servingPort = os.Getenv(types.ServingPortEnvName)
	if len(servingPort) == 0 {
		servingPort = types.DefaultServingPort
	}
}
