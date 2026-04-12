#!/bin/bash

set -e

golangci-lint run

go vet ./...