<p align="center">
  <img height="150" src="https://github.com/afritzler/awesaml/blob/master/images/awesaml.png?raw=true">
</p>

# AweSAML

> "You said all I need to do, is configuring my IDP to enable SSO!!!"

![saml](images/logo.png)

## Installation

The easiest way to install AweSAML is to `go get` it into your Go bin `PATH`.

```shell script
go get -u github.com/afritzler/awesaml
```

Or you can fetch the latest binary release via

```shell script
curl -L -o awesaml "https://github.com/afritzler/awesaml/releases/download/v0.2.1/awesaml_0.2.1_linux_amd64" && chmod +x awesaml
```

All release build can be found in the release section [here](https://github.com/afritzler/awesaml/releases).

### Build locally from source

AweSAML is build using [go modules](https://github.com/golang/go/wiki/Modules). Make sure to set `GO111MODULE=on` before continuing.

```shell script
git clone https://github.com/afritzler/awesaml
cd awesaml
make
```

A Docker based build is available

```shell script
make docker-build
```

## Usage

### Running Locally

Your static web content (aka Service Provider) must have a X.509 key pair established. This typically has to be shared with your IDP provider as well. If you don't have a key pair and certificate at hand you can quickly generate it via

```shell script
openssl req -x509 -newkey rsa:2048 -keyout myservice.key -out myservice.cert -days 365 -nodes -subj "/CN=myservice.example.com"
```

Next up is the configuration of your service provider. Adapt the following configuration to your setup (`source_me.example`)

```shell script
export ENTITY_ID="myEntityID"
export SERVICE_URL="http://localhost:8000"
export SERVICE_PORT="8000" # 8000 is the default
export CONTENT_DIR="public/" # that is where you static content resides on this machine
export IDP_METADATA_URL="http://myAwesomeIDP/saml2/metadata"
export CERT_FILE="myservice.cert" # path to cert and key file
export KEY_FILE="myservice.key"
```

In a nutshell

```shell script
cp source_me.example source_me
# set the env vars in the source_me file according to your setup
source source_me
awesaml
```

You should now be able to access your SAML SSO secured web content here <http://localhost:8000>.

# Acknowledgements

AweSAML is build on the shoulders of giatns and leverages the following modules under the hood

* [crewjam/saml](https://github.com/crewjam/saml) module for the heavy lifting of the SAML flow
