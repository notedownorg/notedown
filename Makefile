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

test: test-pkg test-lsp test-nvim

test-pkg:
	go test ./pkg/...

test-lsp:
	go test ./language-server/...

test-nvim:
	cd neovim && nvim --headless --noplugin -u tests/helpers/minimal_init.lua -c "lua MiniTest.run()" -c "qall!"

test-vhs:
	@echo "=== Starting VHS test execution ==="
	@echo "Current time: $$(date)"
	@echo "Memory usage: $$(free -h || echo 'free command not available')"
	@echo "Disk usage: $$(df -h . || echo 'df command not available')"
	@echo "Display: $${DISPLAY:-'No DISPLAY set'}"
	@echo "VHS availability: $$(which vhs || echo 'vhs not found in PATH')"
	@echo "VHS version: $$(timeout 10s vhs --version 2>&1 || echo 'vhs version check failed/timeout')"
	@echo "Process limits: $$(ulimit -a || echo 'ulimit check failed')"
	@echo "Environment check: $$(timeout 5s vhs --help >/dev/null 2>&1 && echo 'VHS accessible' || echo 'VHS not accessible')"
	GOMAXPROCS=2 go test -parallel 2 -v -x -count=1 ./vhs-tests/... 2>&1 | tee /tmp/vhs-test.log || (echo "Go test failed with exit code: $$?"; echo "=== Last 50 lines of output ==="; tail -50 /tmp/vhs-test.log; exit 1)
	@echo "=== VHS test execution completed at $$(date) ==="

test-vhs-golden:
	rm -f vhs-tests/golden/*.ascii
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
