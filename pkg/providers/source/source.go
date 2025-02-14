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

const MetadataTagsKey = "source"

// Ordered tiers with name last. Source is added automatically
func SourceTag(tiers ...string) string {
	return strings.Join(append([]string{MetadataTagsKey}, tiers...), "/")
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
	// Required to lookup the document in the workspace when we need to write it
	path string

	Format Format
	Title  string
	Url    string
}

type SourceOption func(*Source)

func NewArticle(title string, url string, opts ...SourceOption) Source {
	return NewSource(title, Article, append(opts, WithUrl(url))...)
}

func NewVideo(title string, url string, opts ...SourceOption) Source {
	return NewSource(title, Video, append(opts, WithUrl(url))...)
}

func NewSource(title string, format Format, opts ...SourceOption) Source {
	src := Source{Title: title, Format: format}
	for _, opt := range opts {
		opt(&src)
	}
	return src
}

func WithFormat(format Format) SourceOption {
	return func(p *Source) {
		p.Format = format
	}
}

func WithUrl(url string) SourceOption {
	return func(p *Source) {
		p.Url = url
	}
}

func (p Source) Name() string {
	return p.Title
}

func (p Source) String() string {
	return p.Title
}
