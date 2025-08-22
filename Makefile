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

# Use nix develop shell if nix is available
export NIX_CONFIG := warn-dirty = false
ifneq ($(shell command -v nix 2> /dev/null),)
SHELL := nix develop --command bash
endif

# Version information
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

check: clean format mod lint test

hygiene: format mod

dirty:
	git diff --exit-code

mod:
	go mod tidy

format: licenser
	gofmt -w .
	stylua nvim/

lint:
	golangci-lint run

test: test-lsp test-nvim

test-lsp:
	go test ./...

test-nvim:
	cd nvim && nvim --headless --noplugin -u tests/helpers/minimal_init.lua -c "lua MiniTest.run()" -c "qall!"

install: clean
	go build -ldflags "\
		-w -s \
		-X github.com/notedownorg/notedown/pkg/version.version=$(VERSION) \
		-X github.com/notedownorg/notedown/pkg/version.commit=$(COMMIT) \
		-X github.com/notedownorg/notedown/pkg/version.date=$(DATE)" \
		-o $(shell go env GOPATH)/bin/notedown-language-server \
		./language-server/
	mkdir -p ~/.config/notedown/nvim
	cp -r nvim/* ~/.config/notedown/nvim/

clean:
	rm -f $(shell go env GOPATH)/bin/notedown-language-server
	rm -rf ~/.config/notedown/nvim

licenser:
	licenser apply -r "Notedown Authors"

dev: install
	rm -rf /tmp/notedown_demo_workspace
	cp -r demo_workspace /tmp/notedown_demo_workspace
	cd /tmp/notedown_demo_workspace && nvim .
