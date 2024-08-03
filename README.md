# Panacea Shared Library

## Overview

This library provides shared utilities and configurations for various projects.

## Features

- Configuration management
- Logging
- Security utilities
- General utilities

## Installation

To install the library, first ensure that Go is configured to use your private GitHub repository.

### Configure Go for Private Repositories

1. Set the `GOPRIVATE` environment variable to include your GitHub repository:

    ```sh
    export GOPRIVATE=github.com/voxtmault
    ```

2. Authenticate with GitHub using a personal access token. You can create a token in your GitHub account settings:

    ```sh
    git config --global url."https://<your-token>@github.com/".insteadOf "https://github.com/"
    ```

### Install the Library

Once the above configuration is done, you can install the library by running:

```sh
go get github.com/voxtmault/panacea-shared-lib
```

### Generate .env

You can generate the env file needed to run the library by running :

```sh
make create-env
```
