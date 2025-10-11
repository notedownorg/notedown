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

package appserver

import (
	"testing"

	pb "github.com/notedownorg/notedown/notedown/application_server/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

// Helper function to convert slice-based testing to channel-based API
func filterDocumentsSlice(engine *filterEngine, docs []*pb.Document, filter *pb.FilterExpression) ([]*pb.Document, error) {
	// Create an unbuffered channel and send documents in a goroutine
	docChan := make(chan *pb.Document)
	go func() {
		defer close(docChan)
		for _, doc := range docs {
			docChan <- doc
		}
	}()

	// Filter documents using the channel-based API
	filteredChan, errChan := engine.FilterDocuments(docChan, filter)

	// Collect results
	var result []*pb.Document
	for doc := range filteredChan {
		result = append(result, doc)
	}

	// Check for errors after collecting all results
	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	default:
		// No error
	}

	return result, nil
}

func TestFilterEngine_New(t *testing.T) {
	engine := newFilterEngine()
	assert.NotNil(t, engine)
}

func TestFilterEngine_MetadataFilters(t *testing.T) {
	engine := newFilterEngine()

	// Create metadata as protobuf Struct
	metadata, err := structpb.NewStruct(map[string]any{
		"title":       "Test pb.Document",
		"author":      "Test Author",
		"tags":        []any{"test", "example", "docs"},
		"priority":    5,
		"score":       3.14,
		"published":   true,
		"draft":       false,
		"count":       0,
		"empty_array": []any{},
		"mixed_array": []any{"string", 42, true},
	})
	require.NoError(t, err)

	testDoc := &pb.Document{
		Path:     "/path/to/test.md",
		Metadata: metadata,
	}

	tests := []struct {
		name        string
		field       string
		operator    pb.MetadataOperator
		value       any
		expectMatch bool
		expectError bool
	}{
		// EXISTS tests
		{
			name:        "exists - field present",
			field:       "title",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EXISTS,
			value:       nil,
			expectMatch: true,
		},
		{
			name:        "exists - field absent",
			field:       "nonexistent",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EXISTS,
			value:       nil,
			expectMatch: false,
		},

		// NOT_EXISTS tests
		{
			name:        "not exists - field absent",
			field:       "nonexistent",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_NOT_EXISTS,
			value:       nil,
			expectMatch: true,
		},
		{
			name:        "not exists - field present",
			field:       "title",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_NOT_EXISTS,
			value:       nil,
			expectMatch: false,
		},

		// EQUALS tests
		{
			name:        "equals - string match",
			field:       "title",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
			value:       "Test pb.Document",
			expectMatch: true,
		},
		{
			name:        "equals - string no match",
			field:       "title",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
			value:       "Other pb.Document",
			expectMatch: false,
		},
		{
			name:        "equals - int match",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
			value:       5,
			expectMatch: true,
		},
		{
			name:        "equals - int no match",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
			value:       3,
			expectMatch: false,
		},
		{
			name:        "equals - bool match true",
			field:       "published",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
			value:       true,
			expectMatch: true,
		},
		{
			name:        "equals - bool match false",
			field:       "draft",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
			value:       false,
			expectMatch: true,
		},
		{
			name:        "equals - bool no match",
			field:       "published",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
			value:       false,
			expectMatch: false,
		},
		{
			name:        "equals - float match",
			field:       "score",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
			value:       3.14,
			expectMatch: true,
		},
		{
			name:        "equals - zero value",
			field:       "count",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
			value:       0,
			expectMatch: true,
		},
		{
			name:        "equals - numeric type conversion (int to float)",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
			value:       5.0,
			expectMatch: true,
		},
		{
			name:        "equals - numeric type conversion (float to int)",
			field:       "score",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_EQUALS,
			value:       3,
			expectMatch: false, // 3.14 != 3
		},

		// NOT_EQUALS tests
		{
			name:        "not equals - string no match",
			field:       "title",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_NOT_EQUALS,
			value:       "Other pb.Document",
			expectMatch: true,
		},
		{
			name:        "not equals - string match",
			field:       "title",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_NOT_EQUALS,
			value:       "Test pb.Document",
			expectMatch: false,
		},

		// CONTAINS tests
		{
			name:        "contains - string substring match",
			field:       "title",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_CONTAINS,
			value:       "Test",
			expectMatch: true,
		},
		{
			name:        "contains - string substring no match",
			field:       "title",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_CONTAINS,
			value:       "xyz",
			expectMatch: false,
		},
		{
			name:        "contains - array element match",
			field:       "tags",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_CONTAINS,
			value:       "test",
			expectMatch: true,
		},
		{
			name:        "contains - array element no match",
			field:       "tags",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_CONTAINS,
			value:       "notfound",
			expectMatch: false,
		},
		{
			name:        "contains - empty array",
			field:       "empty_array",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_CONTAINS,
			value:       "anything",
			expectMatch: false,
		},
		{
			name:        "contains - mixed array types",
			field:       "mixed_array",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_CONTAINS,
			value:       42,
			expectMatch: true,
		},
		{
			name:        "contains - non-string non-array field",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_CONTAINS,
			value:       "5",
			expectMatch: false,
		},

		// More CONTAINS tests
		{
			name:        "contains - case sensitivity",
			field:       "title",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_CONTAINS,
			value:       "test", // should not match "Test" (case sensitive)
			expectMatch: false,
		},

		// GREATER_THAN tests
		{
			name:        "greater than - int true",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN,
			value:       3,
			expectMatch: true,
		},
		{
			name:        "greater than - int false",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN,
			value:       10,
			expectMatch: false,
		},
		{
			name:        "greater than - float true",
			field:       "score",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN,
			value:       3.0,
			expectMatch: true,
		},
		{
			name:        "greater than - string comparison",
			field:       "title",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN,
			value:       "A",
			expectMatch: true, // "Test pb.Document" > "A"
		},
		{
			name:        "greater than - non-comparable type",
			field:       "published",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN,
			value:       true,
			expectMatch: false,
		},

		// LESS_THAN tests
		{
			name:        "less than - int true",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_LESS_THAN,
			value:       10,
			expectMatch: true,
		},
		{
			name:        "less than - int false",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_LESS_THAN,
			value:       3,
			expectMatch: false,
		},

		// GREATER_THAN_OR_EQUAL tests
		{
			name:        "greater than or equal - equal",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN_OR_EQUAL,
			value:       5,
			expectMatch: true,
		},
		{
			name:        "greater than or equal - greater",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN_OR_EQUAL,
			value:       3,
			expectMatch: true,
		},
		{
			name:        "greater than or equal - less",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN_OR_EQUAL,
			value:       10,
			expectMatch: false,
		},

		// LESS_THAN_OR_EQUAL tests
		{
			name:        "less than or equal - equal",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_LESS_THAN_OR_EQUAL,
			value:       5,
			expectMatch: true,
		},
		{
			name:        "less than or equal - less",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_LESS_THAN_OR_EQUAL,
			value:       10,
			expectMatch: true,
		},
		{
			name:        "less than or equal - greater",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_LESS_THAN_OR_EQUAL,
			value:       3,
			expectMatch: false,
		},

		// IN tests
		{
			name:        "in - array contains value",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_IN,
			value:       []any{3, 5, 7},
			expectMatch: true,
		},
		{
			name:        "in - array does not contain value",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_IN,
			value:       []any{1, 2, 3},
			expectMatch: false,
		},
		{
			name:        "in - string in array",
			field:       "author",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_IN,
			value:       []any{"Test Author", "Other Author"},
			expectMatch: true,
		},
		{
			name:        "in - non-array value",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_IN,
			value:       5,
			expectMatch: false,
		},

		// NOT_IN tests
		{
			name:        "not in - array does not contain value",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_NOT_IN,
			value:       []any{1, 2, 3},
			expectMatch: true,
		},
		{
			name:        "not in - array contains value",
			field:       "priority",
			operator:    pb.MetadataOperator_METADATA_OPERATOR_NOT_IN,
			value:       []any{3, 5, 7},
			expectMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := NewMetadataFilter(tt.field, tt.operator, tt.value)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			docs, err := filterDocumentsSlice(engine, []*pb.Document{testDoc}, filter)
			require.NoError(t, err)

			if tt.expectMatch {
				assert.Len(t, docs, 1, "expected document to match filter")
			} else {
				assert.Len(t, docs, 0, "expected document to not match filter")
			}
		})
	}
}

func TestFilterEngine_ComplexFilters(t *testing.T) {
	engine := newFilterEngine()

	// Create test documents with protobuf Struct metadata
	createDoc := func(path string, metadata map[string]any) *pb.Document {
		pbMetadata, err := structpb.NewStruct(metadata)
		require.NoError(t, err)
		return &pb.Document{
			Path:     path,
			Metadata: pbMetadata,
		}
	}

	testDocs := []*pb.Document{
		createDoc("/project/doc1.md", map[string]any{"author": "Alice", "published": true, "priority": 1}),
		createDoc("/project/doc2.md", map[string]any{"author": "Bob", "published": true, "priority": 2}),
		createDoc("/project/doc3.md", map[string]any{"author": "Alice", "published": false, "priority": 3}),
		createDoc("/project/doc4.md", map[string]any{"author": "Charlie", "published": true, "priority": 1}),
	}

	t.Run("AND filter", func(t *testing.T) {
		authorFilter, err := NewMetadataFilter("author", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, "Alice")
		require.NoError(t, err)

		publishedFilter, err := NewMetadataFilter("published", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, true)
		require.NoError(t, err)

		andFilter := NewAndFilter(authorFilter, publishedFilter)

		docs, err := filterDocumentsSlice(engine, testDocs, andFilter)
		require.NoError(t, err)

		assert.Len(t, docs, 1) // Only doc1.md matches both conditions
		assert.Equal(t, "/project/doc1.md", docs[0].Path)
	})

	t.Run("OR filter", func(t *testing.T) {
		authorFilter, err := NewMetadataFilter("author", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, "Alice")
		require.NoError(t, err)

		priorityFilter, err := NewMetadataFilter("priority", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, 2)
		require.NoError(t, err)

		orFilter := NewOrFilter(authorFilter, priorityFilter)

		docs, err := filterDocumentsSlice(engine, testDocs, orFilter)
		require.NoError(t, err)

		assert.Len(t, docs, 3) // doc1.md, doc2.md, doc3.md
		paths := make([]string, len(docs))
		for i, doc := range docs {
			paths[i] = doc.Path
		}
		assert.Contains(t, paths, "/project/doc1.md") // Alice
		assert.Contains(t, paths, "/project/doc2.md") // Priority 2
		assert.Contains(t, paths, "/project/doc3.md") // Alice
	})

	t.Run("NOT filter", func(t *testing.T) {
		publishedFilter, err := NewMetadataFilter("published", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, true)
		require.NoError(t, err)

		notFilter := NewNotFilter(publishedFilter)

		docs, err := filterDocumentsSlice(engine, testDocs, notFilter)
		require.NoError(t, err)

		assert.Len(t, docs, 1) // Only doc3.md is not published
		assert.Equal(t, "/project/doc3.md", docs[0].Path)
	})

	t.Run("nested AND/OR filter", func(t *testing.T) {
		// (author = "Alice" OR priority = 1) AND published = true
		authorFilter, err := NewMetadataFilter("author", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, "Alice")
		require.NoError(t, err)

		priorityFilter, err := NewMetadataFilter("priority", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, 1)
		require.NoError(t, err)

		publishedFilter, err := NewMetadataFilter("published", pb.MetadataOperator_METADATA_OPERATOR_EQUALS, true)
		require.NoError(t, err)

		orFilter := NewOrFilter(authorFilter, priorityFilter)
		andFilter := NewAndFilter(orFilter, publishedFilter)

		docs, err := filterDocumentsSlice(engine, testDocs, andFilter)
		require.NoError(t, err)

		assert.Len(t, docs, 2) // doc1.md (Alice + published) and doc4.md (priority 1 + published)
		paths := make([]string, len(docs))
		for i, doc := range docs {
			paths[i] = doc.Path
		}
		assert.Contains(t, paths, "/project/doc1.md")
		assert.Contains(t, paths, "/project/doc4.md")
	})
}

func TestFilterEngine_EdgeCases(t *testing.T) {
	engine := newFilterEngine()

	t.Run("empty document list", func(t *testing.T) {
		filter, err := NewMetadataFilter("any", pb.MetadataOperator_METADATA_OPERATOR_EXISTS, nil)
		require.NoError(t, err)

		docs, err := filterDocumentsSlice(engine, []*pb.Document{}, filter)
		require.NoError(t, err)
		assert.Empty(t, docs)
	})

	t.Run("nil filter", func(t *testing.T) {
		testDoc := &pb.Document{Path: "/test.md"}

		docs, err := filterDocumentsSlice(engine, []*pb.Document{testDoc}, nil)
		require.NoError(t, err)
		assert.Len(t, docs, 1) // No filter means all documents pass
	})

	t.Run("document with nil metadata", func(t *testing.T) {
		testDoc := &pb.Document{
			Path:     "/test.md",
			Metadata: nil,
		}

		filter, err := NewMetadataFilter("title", pb.MetadataOperator_METADATA_OPERATOR_EXISTS, nil)
		require.NoError(t, err)

		docs, err := filterDocumentsSlice(engine, []*pb.Document{testDoc}, filter)
		require.NoError(t, err)
		assert.Empty(t, docs) // No metadata means field doesn't exist
	})

}

func TestNewMetadataFilter_InvalidOperator(t *testing.T) {
	_, err := NewMetadataFilter("field", pb.MetadataOperator(999), "value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported metadata operator")
}

func TestFilterHelpers(t *testing.T) {
	t.Run("NewAndFilter", func(t *testing.T) {
		filter1, err := NewMetadataFilter("field1", pb.MetadataOperator_METADATA_OPERATOR_EXISTS, nil)
		require.NoError(t, err)

		filter2, err := NewMetadataFilter("field2", pb.MetadataOperator_METADATA_OPERATOR_EXISTS, nil)
		require.NoError(t, err)

		andFilter := NewAndFilter(filter1, filter2)
		require.NotNil(t, andFilter)
		require.NotNil(t, andFilter.GetAndFilter())
		assert.Len(t, andFilter.GetAndFilter().Filters, 2)
	})

	t.Run("NewOrFilter", func(t *testing.T) {
		filter1, err := NewMetadataFilter("field1", pb.MetadataOperator_METADATA_OPERATOR_EXISTS, nil)
		require.NoError(t, err)

		filter2, err := NewMetadataFilter("field2", pb.MetadataOperator_METADATA_OPERATOR_EXISTS, nil)
		require.NoError(t, err)

		orFilter := NewOrFilter(filter1, filter2)
		require.NotNil(t, orFilter)
		require.NotNil(t, orFilter.GetOrFilter())
		assert.Len(t, orFilter.GetOrFilter().Filters, 2)
	})

	t.Run("NewNotFilter", func(t *testing.T) {
		filter, err := NewMetadataFilter("field", pb.MetadataOperator_METADATA_OPERATOR_EXISTS, nil)
		require.NoError(t, err)

		notFilter := NewNotFilter(filter)
		require.NotNil(t, notFilter)
		require.NotNil(t, notFilter.GetNotFilter())
		assert.Equal(t, filter, notFilter.GetNotFilter().Filter)
	})

}
