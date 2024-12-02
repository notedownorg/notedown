// Copyright 2024 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package source

import (
	"strings"
)

const MetadataKey = "source"

type identifier struct {
	path    string
	version string
}

// By default we will set line to -1 to default to end of file
func NewIdentifier(path string, version string) identifier {
	return identifier{path: path, version: version}
}

func (i identifier) String() string {
	// Pipe separators are good enough for now but may need to be changed as pipes
	// are technically valid (although unlikely to actually be used) in unix file paths
	// We may want to consider an actual encoding scheme for this in the future.
	var builder strings.Builder
	builder.WriteString(i.path)
	builder.WriteString("|")
	builder.WriteString(i.version)
	return builder.String()
}

const (
	TitleKey  = "title"
	FormatKey = "format"
	UrlKey    = "url"
)

type Format string

const (
	Article Format = "article"
	Video   Format = "video"

	Unknown Format = ""
)

var formatMap = map[string]Format{
	"article": Article,
	"video":   Video,
}

type Source struct {
	title      string
	identifier identifier
	format     Format
	url        string
}

type SourceOption func(*Source)

func NewArticle(identifier identifier, title string, url string, opts ...SourceOption) Source {
	return NewSource(identifier, title, Article, append(opts, WithUrl(url))...)
}

func NewVideo(identifier identifier, title string, url string, opts ...SourceOption) Source {
	return NewSource(identifier, title, Video, append(opts, WithUrl(url))...)
}

func NewSource(identifier identifier, title string, format Format, opts ...SourceOption) Source {
	p := Source{
		identifier: identifier,
		title:      title,
		format:     format,
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// If no name is provided, we will attempt to infer it from the file's basename
func WithName(name string) SourceOption {
	return func(p *Source) {
		p.title = name
	}
}

func WithFormat(format Format) SourceOption {
	return func(p *Source) {
		p.format = format
	}
}

func WithUrl(url string) SourceOption {
	return func(p *Source) {
		p.url = url
	}
}

func (p Source) Identifier() identifier {
	return p.identifier
}

func (p Source) Name() string {
	return p.title
}

func (p Source) Path() string {
	return p.identifier.path
}

func (p Source) Format() Format {
	return p.format
}

func (p Source) String() string {
	return p.Name()
}
