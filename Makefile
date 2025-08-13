# Copyright 2024 Notedown Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Version information
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

all: format mod test dirty

hygiene: format mod

dirty:
	git diff --exit-code

mod:
	go mod tidy

format:
	gofmt -w .

test:
	go test ./...

install:
	go build -ldflags "\
		-w -s \
		-X github.com/notedownorg/notedown/pkg/version.version=$(VERSION) \
		-X github.com/notedownorg/notedown/pkg/version.commit=$(COMMIT) \
		-X github.com/notedownorg/notedown/pkg/version.date=$(DATE)" \
		-o $(shell go env GOPATH)/bin/notedown-language-server \
		./lsp/

licenser:
	licenser apply -r "Notedown Authors"
