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
	"testing"

	"github.com/notedownorg/notedown/pkg/configuration"
	"github.com/stretchr/testify/assert"
)

func TestSanitizeTags(t *testing.T) {
	tests := []struct {
		name     string
		metadata Metadata
		format   configuration.TagFormat
		expected []string
	}{
		{
			name: "No tags",
			metadata: Metadata{
				MetadataTagsKey: []string{},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{},
		},
		{
			name: "Single tag, no spaces",
			metadata: Metadata{
				MetadataTagsKey: []string{"validTag"},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{"valid-tag"},
		},
		{
			name: "Multiple tags, no spaces",
			metadata: Metadata{
				MetadataTagsKey: []string{"validTag1", "validTag2"},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{"valid-tag1", "valid-tag2"},
		},
		{
			name: "Single tag with spaces, kebab case",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid tag"},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{"valid-tag"},
		},
		{
			name: "Single tag with spaces, snake case",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid tag"},
			},
			format:   configuration.TagFormatSnakeCase,
			expected: []string{"valid_tag"},
		},
		{
			name: "Single tag with spaces, camel case",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid tag"},
			},
			format:   configuration.TagFormatCamelCase,
			expected: []string{"validTag"},
		},
		{
			name: "Single tag with spaces, pascal case",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid tag"},
			},
			format:   configuration.TagFormatPascalCase,
			expected: []string{"ValidTag"},
		},
		{
			name: "Kebab to pascal case",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid-tag"},
			},
			format:   configuration.TagFormatPascalCase,
			expected: []string{"ValidTag"},
		},
		{
			name: "Kebab to snake case",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid-tag"},
			},
			format:   configuration.TagFormatSnakeCase,
			expected: []string{"valid_tag"},
		},
		{
			name: "Kebab to camel case",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid-tag"},
			},
			format:   configuration.TagFormatCamelCase,
			expected: []string{"validTag"},
		},
		{
			name: "Snake to kebab case",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid_tag"},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{"valid-tag"},
		},
		{
			name: "Snake to pascal case",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid_tag"},
			},
			format:   configuration.TagFormatPascalCase,
			expected: []string{"ValidTag"},
		},
		{
			name: "Snake to camel case",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid_tag"},
			},
			format:   configuration.TagFormatCamelCase,
			expected: []string{"validTag"},
		},
		{
			name: "Camel to kebab case",
			metadata: Metadata{
				MetadataTagsKey: []string{"validTag"},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{"valid-tag"},
		},
		{
			name: "Camel to snake case",
			metadata: Metadata{
				MetadataTagsKey: []string{"validTag"},
			},
			format:   configuration.TagFormatSnakeCase,
			expected: []string{"valid_tag"},
		},
		{
			name: "Camel to pascal case",
			metadata: Metadata{
				MetadataTagsKey: []string{"validTag"},
			},
			format:   configuration.TagFormatPascalCase,
			expected: []string{"ValidTag"},
		},
		{
			name: "Pascal to kebab case",
			metadata: Metadata{
				MetadataTagsKey: []string{"ValidTag"},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{"valid-tag"},
		},
		{
			name: "Pascal to snake case",
			metadata: Metadata{
				MetadataTagsKey: []string{"ValidTag"},
			},
			format:   configuration.TagFormatSnakeCase,
			expected: []string{"valid_tag"},
		},
		{
			name: "Pascal to camel case",
			metadata: Metadata{
				MetadataTagsKey: []string{"ValidTag"},
			},
			format:   configuration.TagFormatCamelCase,
			expected: []string{"validTag"},
		},
		{
			name: "Invalid characters removed",
			metadata: Metadata{
				MetadataTagsKey: []string{"invalid!@#tag"},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{"invalidtag"},
		},
		{
			name: "Single tag with only whitespace",
			metadata: Metadata{
				MetadataTagsKey: []string{"   "},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{},
		},
		{
			name: "Single tag with special characters",
			metadata: Metadata{
				MetadataTagsKey: []string{"!@#$%^&*()_+{}|:\"<>?~"},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{},
		},
		{
			name: "Single tag with empty string",
			metadata: Metadata{
				MetadataTagsKey: []string{""},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{},
		},
		{
			name: "Single very long tag",
			metadata: Metadata{
				MetadataTagsKey: []string{"this-is-a-very-long-tag-that-might-exceed-normal-length-limits"},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{"this-is-a-very-long-tag-that-might-exceed-normal-length-limits"},
		},
		{
			name: "Single tag with spaces, camel case (multiple uppercase)",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid TAG Name"},
			},
			format:   configuration.TagFormatCamelCase,
			expected: []string{"validTagName"},
		},
		{
			name: "Single tag with spaces, pascal case (multiple uppercase)",
			metadata: Metadata{
				MetadataTagsKey: []string{"VALID TAG Name"},
			},
			format:   configuration.TagFormatPascalCase,
			expected: []string{"ValidTagName"},
		},
		{
			name: "Single tag with multiple spaces",
			metadata: Metadata{
				MetadataTagsKey: []string{"valid   tag   name"},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{"valid-tag-name"},
		},
		{
			name: "Single tag starting with number",
			metadata: Metadata{
				MetadataTagsKey: []string{"123validTag"},
			},
			format:   configuration.TagFormatCamelCase,
			expected: []string{"123validTag"},
		},
		{
			name: "Number only tag",
			metadata: Metadata{
				MetadataTagsKey: []string{"123"},
			},
			format:   configuration.TagFormatKebabCase,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sanitizeTags(tt.metadata, tt.format)
			actual := tt.metadata[MetadataTagsKey].([]string)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
