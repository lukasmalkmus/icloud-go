# iCloud Go

[![Go Reference][gopkg_badge]][gopkg]
[![Go Workflow][go_workflow_badge]][go_workflow]
[![Coverage Status][coverage_badge]][coverage]
[![Go Report][report_badge]][report]
[![Latest Release][release_badge]][release]
[![License][license_badge]][license]

---

## Table of Contents

1. [Introduction](#introduction)
1. [Installation](#Installation)
1. [Authentication](#authentication)
1. [Usage](#usage)
1. [Contributing](#contributing)
1. [License](#license)

## Introduction

_iCloud Go_ is a Go client library for using the [CloudKit Web Services][1] API.
While the foundation of this package is powerfull enough to support the whole
API, its main purpose is to provide enough features to create records in
CloudKit. This serves my personal purpose of using this package for building
various kinds of importers.

  [1]: https://developer.apple.com/library/archive/documentation/DataManagement/Conceptual/CloudKitWebServicesReference/index.html

## Installation

### Install using `go get`

```shell
$ go get github.com/lukasmalkmus/icloud-go/icloud
```

### Install from source

```shell
$ git clone https://github.com/lukasmalkmus/icloud-go.git
$ cd icloud-go
$ make
```

## Authentication

This package only supports Server-to-Server authentication.

You can create a keypair using `openssl`:

```shell
$ openssl ecparam -name prime256v1 -genkey -noout -out eckey.pem
```

To get the public key (to enter it on iCloud Dashboard):

```shell
$ openssl ec -in eckey.pem -pubout
```

To get the private key (used by the client to sign requests):

```shell
$ openssl ec -in eckey.pem
```

## Usage

```go
import "github.com/lukasmalkmus/icloud-go/icloud"

var (
	keyID         = os.Getenv("ICLOUD_KEY_ID")
	container     = os.Getenv("ICLOUD_CONTAINER")
	rawPrivateKey = os.Getenv("ICLOUD_PRIVATE_KEY")
)

// 1. Parse the private key.
privateKey, err := x509.ParseECPrivateKey([]byte(rawPrivateKey))
if err != nil {
	log.Fatal(err)
}

// 2. Create the iCloud client.
client, err := icloud.NewClient(container, keyID, privateKey, icloud.Development)
if err != nil {
	log.Fatal(err)
}

// 3. Create a record.
if _, err = client.Records.Modify(context.Background(), icloud.Public, icloud.RecordsRequest{
	Operations: []icloud.RecordOperation{
		{
			Type: icloud.Create,
			Record: icloud.Record{
				Type: "MyRecord",
				Fields: icloud.Fields{
					{
						Name:  "MyField",
						Value: "Hello, World!",
					},
					{
						Name:  "MyOtherField",
						Value: 1000,
					},
				},
			},
		},
	},
}); err != nil {
	log.Fatal(err)
}
```

## Contributing

Feel free to submit PRs or to fill issues. Every kind of help is appreciated.

Before committing, `make` should run without any issues.

## License

&copy; Lukas Malkmus, 2021

Distributed under MIT License (`The MIT License`).

See [LICENSE](LICENSE) for more information.

[![License Status][license_status_badge]][license_status]

<!-- Badges -->

[gopkg]: https://pkg.go.dev/github.com/lukasmalkmus/icloud-go
[gopkg_badge]: https://img.shields.io/badge/doc-reference-007d9c?logo=go&logoColor=white&style=flat-square
[go_workflow]: https://github.com/lukasmalkmus/icloud-go/actions?query=workflow%3Ago
[go_workflow_badge]: https://img.shields.io/github/workflow/status/lukasmalkmus/icloud-go/go?style=flat-square&ghcache=unused
[coverage]: https://codecov.io/gh/lukasmalkmus/icloud-go
[coverage_badge]: https://img.shields.io/codecov/c/github/lukasmalkmus/icloud-go.svg?style=flat-square&ghcache=unused
[report]: https://goreportcard.com/report/github.com/lukasmalkmus/icloud-go
[report_badge]: https://goreportcard.com/badge/github.com/lukasmalkmus/icloud-go?style=flat-square&ghcache=unused
[release]: https://github.com/lukasmalkmus/icloud-go/releases/latest
[release_badge]: https://img.shields.io/github/release/lukasmalkmus/icloud-go.svg?style=flat-square&ghcache=unused
[license]: https://opensource.org/licenses/MIT
[license_badge]: https://img.shields.io/github/license/lukasmalkmus/icloud-go.svg?color=blue&style=flat-square&ghcache=unused
[license_status]: https://app.fossa.com/projects/git%2Bgithub.com%2Flukasmalkmus%2Ficloud-go
[license_status_badge]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Flukasmalkmus%2Ficloud-go.svg?type=large&ghcache=unused
