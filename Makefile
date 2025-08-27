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
define NIX_SETTINGS
warn-dirty = false
download-buffer-size = 134217728
endef
export NIX_CONFIG := $(NIX_SETTINGS)
ifneq ($(shell command -v nix 2> /dev/null),)
SHELL := nix develop --command bash
endif

# Version information
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

check: clean format mod lint test

all: hygiene test dirty

hygiene: format mod

dirty:
	git diff --exit-code

mod:
	go mod tidy

format: licenser
	gofmt -w .
	stylua neovim/

lint:
	golangci-lint run

test: deps test-pkg test-lsp test-nvim

deps:
	go mod download

test-pkg:
	go test ./pkg/...

test-lsp:
	go test ./language-server/...

test-nvim:
	cd neovim && nvim --headless --noplugin -u tests/helpers/minimal_init.lua -c "lua MiniTest.run()" -c "qall!"

# VHS tests require large dependencies (Chromium, FFmpeg) via nix
# Pre-download Go modules to avoid CI timeout during test execution
test-vhs:
	go mod download
	go test -parallel 4 -v ./vhs-tests/...

test-vhs-golden:
	rm -f vhs-tests/golden/*.ascii
	go mod download
	go test -parallel 4 -v ./vhs-tests/...

install: clean
	go build -ldflags "\
		-w -s \
		-X github.com/notedownorg/notedown/pkg/version.version=$(VERSION) \
		-X github.com/notedownorg/notedown/pkg/version.commit=$(COMMIT) \
		-X github.com/notedownorg/notedown/pkg/version.date=$(DATE)" \
		-o $(shell go env GOPATH)/bin/notedown-language-server \
		./language-server/
	mkdir -p ~/.config/notedown/nvim
	cp -r neovim/* ~/.config/notedown/nvim/

clean:
	rm -f $(shell go env GOPATH)/bin/notedown-language-server
	rm -rf ~/.config/notedown/nvim

licenser:
	licenser apply -r "Notedown Authors"

dev: install
	rm -rf /tmp/notedown_demo_workspace
	cp -r demo_workspace /tmp/notedown_demo_workspace
	cd /tmp/notedown_demo_workspace && nvim .
