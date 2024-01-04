# Duplicate File Scanner

A command-line tool to efficiently locate and display duplicate files in a given directory.

## Overview

This Go program uses concurrent processing to scan a directory, compute the MD5 hash of each file, and identify duplicates based on their hash values.

## Features

- Efficiently identifies duplicate files in a specified directory.
- Utilizes Go's concurrency for faster processing.
- Outputs hash values and paths of duplicate files.

## Getting Started

### Prerequisites

- Go installed on your machine
- Docker (if you prefer running the program in a Docker container)

### Installation

1. Clone the repository.
2. Navigate to the project directory.
3. Build the Docker image (optional).

#### Running with Go:

go run main.go -path /your/directory/to/scan

#### Running with Docker:

docker run duplicate-file-scanner -path /your/directory/to/scan

## Contributing

Contributions are welcome!

## License

This project is licensed under the MIT License.
