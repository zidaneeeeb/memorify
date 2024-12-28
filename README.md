# hbdtoyou

hbdtoyou provides functionalities for clients (school and user) in Memorify. It is also responsible for storing school-level config that can be used to supports single-tenancy in data stores.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

#### Golang

You need to have [Go v1.19.4](https://golang.org/dl/) installed on your machine. Follow the [official installation guide](https://golang.org/doc/install) to install Go. Or, follow [managing installations guide](https://go.dev/doc/manage-install) to have multiple Go versions on your machine.

#### PostgreSQL

This service has dependency with PostgreSQL. For development environment, you need to have a PostgreSQL server running on your machine.

#### Secret File

Currently, this service stores credentials on a secret file. For development environment, the files are stored on `files/ets/<service-name>/secret.development.yaml`. You need to modify the file accordingly.

### Building

1. Once you have all the prerequisites, you can start by cloning this repository into your machine.

```sh
$ mkdir -p $GOPATH/src/github.com/temui-sisva/
$ cd $GOPATH/src/github.com/temui-sisva
$ git clone https://github.com/temui-sisva/hbdtoyou.git
$ cd hbdtoyou
```

> The rest of this instructions assumes that your current working directory is on `$GOPATH/src/github.com/temui-sisva/hbdtoyou`

2. Build binaries using the `go build` command.

```sh
$ go build ./cmd/hbdtoyou-api-grpc    # For gRPC server binary
$ go build ./cmd/hbdtoyou-api-http    # For HTTP server binary
```

### Running

1. If needed, you can modify the app config for development environment through this file `files/ets/<service-name>/config.development.yaml`.

2. Execute the binary to start the service

```sh
$ ./hbdtoyou-api-grpc -secret-path files/etc/hbdtoyou-api-grpc/secret.development.yaml  # For gRPC server binary
$ ./hbdtoyou-api-http -secret-path files/etc/hbdtoyou-api-http/secret.development.yaml  # For HTTP server binary
```

## Directory Structure

This repository is organized with the following structure

```
hbdtoyou
|-- cmd                         # Contains executables codes
|   |-- hbdtoyou-api-grpc         # gRPC server
|   |-- hbdtoyou-api-http         # HTTP server
|-- files
|   |-- etc                     # Contains config files
|   |   |-- hbdtoyou-api-grpc         
|   |   |-- hbdtoyou-api-http   
|-- internal                    # Application service packages      
```

## Contributing

### Code

Application service packages should be developed in the `internal` directory, as those logic should not be used/imported by external repositories.

Application service packages are made using the domain-driven design concept. Some articles to read:
* https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
* https://medium.easyread.co/golang-clean-archithecture-efd6d7c43047

Application service package's naming should be self-explanatory about its purpose, so that other developers would not misinterpret the package.

### Sending Changes

Commits must follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/), and multiple commits must be squashed into one. Commit description is mandatory, but body is optional; you can omit if you think the title is clear enough.

Commit description must start with lowercase letter. Use [imperative mood](https://www.freecodecamp.org/news/how-to-write-better-git-commit-messages/) in the description.

Once you're done, push your branch and open a Pull Request to `main`.
