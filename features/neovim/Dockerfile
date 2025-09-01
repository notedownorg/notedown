# Copyright 2025 Notedown Authors
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

# Multi-stage Dockerfile for VHS testing with Notedown Language Server
FROM golang:1.24.4-bullseye AS builder

# Build the language server
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o notedown-language-server ./language-server/

# Final stage: Debian with VHS, Neovim, and our language server
FROM debian:bookworm-slim

# Install system dependencies (without neovim - we'll install 0.11 manually)
RUN apt-get update && \
    apt-get install -y \
    curl \
    ca-certificates \
    ffmpeg \
    chromium \
    xvfb \
    fuse \
    && rm -rf /var/lib/apt/lists/*

# Install Neovim 0.10.2 from binary tar.gz
RUN curl -L https://github.com/neovim/neovim/releases/download/v0.10.2/nvim-linux64.tar.gz -o nvim.tar.gz && \
    tar -xzf nvim.tar.gz -C /opt/ && \
    ln -sf /opt/nvim-linux64/bin/nvim /usr/local/bin/nvim && \
    rm nvim.tar.gz

# Install ttyd from GitHub releases
RUN curl -L https://github.com/tsl0922/ttyd/releases/download/1.7.4/ttyd.x86_64 -o /usr/local/bin/ttyd && \
    chmod +x /usr/local/bin/ttyd

# Install VHS from GitHub releases
RUN VHS_VERSION=0.10.0 && \
    curl -L "https://github.com/charmbracelet/vhs/releases/download/v${VHS_VERSION}/vhs_${VHS_VERSION}_linux_x86_64.tar.gz" | \
    tar -xz -C /tmp && \
    mv "/tmp/vhs_${VHS_VERSION}_Linux_x86_64/vhs" /usr/local/bin/vhs && \
    chmod +x /usr/local/bin/vhs && \
    rm -rf /tmp/vhs_${VHS_VERSION}_Linux_x86_64

# Copy the built language server
COPY --from=builder /build/notedown-language-server /usr/local/bin/

# Copy Neovim plugin code
COPY neovim /opt/notedown/nvim

# Install additional Neovim plugins manually (no plugin manager needed)
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /root/.local/share/nvim/site/pack/plugins/start

# Clone Catppuccin theme
RUN git clone https://github.com/catppuccin/nvim.git \
    /root/.local/share/nvim/site/pack/plugins/start/catppuccin

# Clone nvim-web-devicons (useful for completion menu icons)
RUN git clone https://github.com/nvim-tree/nvim-web-devicons.git \
    /root/.local/share/nvim/site/pack/plugins/start/nvim-web-devicons

# Clone nvim-cmp ecosystem for automatic completion
RUN git clone https://github.com/hrsh7th/nvim-cmp.git \
    /root/.local/share/nvim/site/pack/plugins/start/nvim-cmp

RUN git clone https://github.com/hrsh7th/cmp-nvim-lsp.git \
    /root/.local/share/nvim/site/pack/plugins/start/cmp-nvim-lsp

RUN git clone https://github.com/hrsh7th/cmp-buffer.git \
    /root/.local/share/nvim/site/pack/plugins/start/cmp-buffer

RUN git clone https://github.com/hrsh7th/cmp-path.git \
    /root/.local/share/nvim/site/pack/plugins/start/cmp-path

# Clone LuaSnip for snippet support
RUN git clone https://github.com/L3MON4D3/LuaSnip.git \
    /root/.local/share/nvim/site/pack/plugins/start/LuaSnip

RUN git clone https://github.com/saadparwaiz1/cmp_luasnip.git \
    /root/.local/share/nvim/site/pack/plugins/start/cmp_luasnip

# Create Neovim config directory and copy init.lua
RUN mkdir -p /root/.config/nvim
COPY features/neovim/init.lua /root/.config/nvim/init.lua

# Set up environment for containerized VHS (matching official VHS Docker image)
ENV DISPLAY=:99
ENV PATH="/usr/local/bin:${PATH}"
ENV VHS_NO_SANDBOX="true"

# Create workspace directory
WORKDIR /vhs

# Create entrypoint script that starts Xvfb and runs VHS
RUN echo '#!/bin/bash\n\
# Start Xvfb in background\n\
Xvfb :99 -screen 0 1024x768x24 &\n\
\n\
# Wait a moment for Xvfb to start\n\
sleep 2\n\
\n\
# Run VHS with provided arguments\n\
exec vhs "$@"\n\
' > /usr/local/bin/entrypoint.sh && chmod +x /usr/local/bin/entrypoint.sh

# Default command runs VHS with the provided tape file
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
