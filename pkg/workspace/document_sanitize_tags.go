// Copyright 2025 Notedown Authors
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

package workspace

import (
	"strconv"
	"strings"

	"github.com/notedownorg/notedown/pkg/configuration"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func sanitizeTags(m Metadata, format configuration.TagFormat) {
	tags, ok := m[MetadataTagsKey]
	if !ok {
		return
	}

	res, ok := tags.([]string)
	if !ok {
		// Check if tags is set to a single string instead of a list
		// Convert it to a list if it is
		tag, ok := tags.(string)
		if !ok {
			return
		}
		res = []string{tag}
	}

	validTags := make([]string, 0, len(res))
	for _, tag := range res {
		if tag == "" {
			continue
		}
		// Tags can only contain alphanumeric characters, hyphens, underscores and forward slashes
		// Remove anything thats not in that set (except spaces, which we will get to)
		tag = strings.Map(func(r rune) rune {
			if r >= 'a' && r <= 'z' {
				return r
			}
			if r >= 'A' && r <= 'Z' {
				return r
			}
			if r >= '0' && r <= '9' {
				return r
			}
			if r == '-' || r == '_' || r == '/' || r == ' ' {
				return r
			}
			return -1
		}, tag)

		// Now remove spaces based on the format
		switch format {
		case configuration.TagFormatKebabCase:
			tag = toKebabCase(tag)
		case configuration.TagFormatSnakeCase:
			tag = toSnakeCase(tag)
		case configuration.TagFormatCamelCase:
			tag = toCamelCase(tag)
		case configuration.TagFormatPascalCase:
			tag = toPascalCase(tag)
		}

		// Filter out empty tags and number only tags
		if tag == "" {
			continue
		}
		if _, err := strconv.Atoi(tag); err == nil {
			continue
		}

		validTags = append(validTags, tag)
	}

	m[MetadataTagsKey] = validTags
}

var title = cases.Title(language.Und).String

func toKebabCase(s string) string {
	words := split(s)
	if len(words) == 0 {
		return ""
	}
	if len(words) == 1 {
		return s
	}

	kebab := strings.ToLower(words[0])

	for _, word := range words[1:] {
		kebab += "-" + strings.ToLower(word)
	}

	return kebab
}

func toSnakeCase(s string) string {
	words := split(s)
	if len(words) == 0 {
		return ""
	}
	if len(words) == 1 {
		return s
	}

	snake := strings.ToLower(words[0])

	for _, word := range words[1:] {
		snake += "_" + strings.ToLower(word)
	}

	return snake
}

func toCamelCase(s string) string {
	words := split(s)
	if len(words) == 0 {
		return ""
	}
	if len(words) == 1 {
		return s
	}

	// Convert first word to lowercase
	camel := strings.ToLower(words[0])

	// Capitalize the rest
	for _, word := range words[1:] {
		camel += title(strings.ToLower(word)) // Capitalize first letter
	}

	return camel
}

func toPascalCase(s string) string {
	words := split(s)
	if len(words) == 0 {
		return ""
	}
	if len(words) == 1 {
		return s
	}

	pascal := title(strings.ToLower(words[0])) // Capitalize first letter

	for _, word := range words[1:] {
		pascal += title(strings.ToLower(word)) // Capitalize first letter
	}

	return pascal
}

// Split on spaces, underscores and hyphens or capital letters
func split(s string) []string {
	words := make([]string, 0)

	// Split on spaces, underscores and hyphens
	split := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '-' || r == '_'
	})

	for _, word := range split {
		words = append(words, splitCamel(word)...)
	}
	return words
}

func splitCamel(s string) []string {
	words := make([]string, 0)

	// Split on capital letters only if the previous character is a lowercase letter
	start := 0
	for i := 1; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' && s[i-1] >= 'a' && s[i-1] <= 'z' {
			words = append(words, s[start:i])
			start = i
		}
	}
	words = append(words, s[start:])

	return words
}

func filterEmptyStrings(s []string) []string {
	filtered := make([]string, 0, len(s))
	for _, str := range s {
		if str != "" {
			filtered = append(filtered, str)
		}
	}
	return filtered
}
