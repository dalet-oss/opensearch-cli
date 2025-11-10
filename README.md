[![Build](https://github.com/dalet-oss/opensearch-cli/actions/workflows/build.yml/badge.svg)](https://github.com/dalet-oss/opensearch-cli/actions/workflows/build.yml)
# OpenSearch CLI

A command-line interface tool for interacting with OpenSearch.


## Development

New developers must have [`gobrew`](https://github.com/kevincobain2000/gobrew) installed for their convenience + [Task to build cli](https://taskfile.dev/docs/installation)

## Overview

OpenSearch CLI provides a convenient way to manage and interact with OpenSearch clusters directly from your terminal. The tool supports various operations including querying, index management, and cluster administration tasks.

## Installation

### Prerequisites

- Go 1.24.6 or later

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/dalet-oss/opensearch-cli.git
   cd opensearch-cli
   ```

2. Build the project:
   ```
   go build -o opensearch-cli
   ```

## Usage

Run `opensearch-cli --help` to see available commands and options.

## Features

- Connect to OpenSearch clusters
- Execute queries and view results
- Manage indexes and mappings
- Monitor cluster health and stats
- Secure credential management

## Dependencies

This project utilizes the following libraries:
- github.com/opensearch-project/opensearch-go - OpenSearch Go client
- github.com/spf13/cobra - Command line interface framework
- github.com/manifoldco/promptui - Interactive prompt for command-line applications
- github.com/zalando/go-keyring - Secure credential storage

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms of the LICENSE file included in the repository.

## Acknowledgments

- OpenSearch Project for providing the Go client library
- All contributors who have helped with the development of this tool