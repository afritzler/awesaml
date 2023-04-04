package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	secretmanagerapiv1beta1 "cloud.google.com/go/secretmanager/apiv1beta1"
	"github.com/afritzler/awesaml/pkg/utils"
	"github.com/crewjam/saml/samlsp"
	secretmanagerv1beta1 "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
)

var (
	entityID,
	certFile,
	certSecretName,
	keyFile,
	keySecretName,
	serviceURL,
	idpMetaDataURL,
	contentDir,
	servingPort string
	keyPair tls.Certificate
)

func main() {
	// init vars or panic
	if err := initVars(); err != nil {
		log.Fatalf("failed to initialize vars %+v", err)
		os.Exit(1)
	}
	setVarsOrFail()

	rootURL, _ := url.Parse(serviceURL)
	idpMetadataURL, _ := url.Parse(idpMetaDataURL)

	idpMetadata, err := samlsp.FetchMetadata(
		context.Background(),
		http.DefaultClient,
		*idpMetadataURL)
	if err != nil {
		log.Fatalf("failed to fetch metadata from idp: %+v", err)
		os.Exit(1)
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
		log.Fatalf("failed to create service provider instance: %+v", err)
		os.Exit(1)
	}

	log.Println("starting Service Provider ...")
	log.Printf("entityID: %s\n", entityID)
	log.Printf("IDPMetaDataURL: %s\n", idpMetaDataURL)
	log.Printf("listening on %s\n", serviceURL)
	log.Printf("serving content from %s\n", contentDir)
	fs := http.FileServer(http.Dir(contentDir))

	log.Println("using net/http ...")
	http.Handle("/", samlSP.RequireAccount(fs))
	http.Handle("/saml/acs", samlSP)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", servingPort), nil); err != nil {
		log.Fatalf("failed to listen and serve: %+v", err)
		os.Exit(1)
	}
}

func initVars() error {
	entityID = os.Getenv(utils.EntityIDEnvName)
	if len(entityID) == 0 {
		return fmt.Errorf("you need to provide the entityID by exporting it via the following env var %s", utils.EntityIDEnvName)
	}
	certFile = os.Getenv(utils.CertFileEnvName)
	if len(certFile) == 0 {
		certSecretName = os.Getenv(utils.CertSecretNameEnvName)
		if len(certSecretName) == 0 {
			return fmt.Errorf(fmt.Sprintf("you need to provide either a location of the cert file via %s or cert secret name via %s", utils.CertFileEnvName, utils.CertSecretNameEnvName))
		}
	}
	keyFile = os.Getenv(utils.KeyFileEnvName)
	if len(keyFile) == 0 {
		keySecretName = os.Getenv(utils.KeySecretNameEnvName)
		if len(keySecretName) == 0 {
			return fmt.Errorf(fmt.Sprintf("you need to provide either a location of the key file via %s or key secret name via %s", utils.KeyFileEnvName, utils.KeySecretNameEnvName))
		}
	}
	serviceURL = os.Getenv(utils.ServiceURLEnvName)
	if len(serviceURL) == 0 {
		return fmt.Errorf(fmt.Sprintf("you need to provide the serviceURL by exporting it via the following env var %s", utils.ServiceURLEnvName))
	}
	idpMetaDataURL = os.Getenv(utils.IdpMetaDataURLEnvName)
	if len(idpMetaDataURL) == 0 {
		return fmt.Errorf(fmt.Sprintf("you need to provide the idpMetaDataURL by exporting it via the following env var %s", utils.IdpMetaDataURLEnvName))
	}
	contentDir = os.Getenv(utils.ContentDirEnvName)
	if len(contentDir) == 0 {
		return fmt.Errorf(fmt.Sprintf("you need to provide the contentDir by exporting it via the following env var %s", utils.ContentDirEnvName))
	}
	servingPort = os.Getenv(utils.ServingPortEnvName)
	if len(servingPort) == 0 {
		servingPort = utils.DefaultServingPort
	}
	return nil
}

func setVarsOrFail() {
	if len(certSecretName) > 0 && len(keySecretName) > 0 {
		log.Println("using cert/key from secret manager")
		certBytes, err := getSecret(certSecretName)
		if err != nil {
			log.Fatalf("failed to get secret from secret manager: %+v", err)
			os.Exit(1)
		}
		keyBytes, err := getSecret(keySecretName)
		if err != nil {
			log.Fatalf("failed to get key secret from secret manager: %+v", err)
			os.Exit(1)
		}
		keyPair, err = tls.X509KeyPair(certBytes, keyBytes)
		if err != nil {
			log.Fatalf("failed to load key pair: %+v", err)
			os.Exit(1)
		}
		keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
		if err != nil {
			log.Fatalf("failed to parse certificate: %+v", err)
			os.Exit(1)
		}
	} else {
		log.Println("using cert/key from disk")
		var err error
		keyPair, err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("failed to load key pair: %+v", err)
			os.Exit(1)
		}
		keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
		if err != nil {
			log.Fatalf("failed to parse certificate: %+v", err)
			os.Exit(1)
		}
	}
}

func getSecret(name string) ([]byte, error) {
	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanagerapiv1beta1.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	// Build the request.
	req := &secretmanagerv1beta1.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %v", err)
	}

	return result.Payload.Data, nil
}
