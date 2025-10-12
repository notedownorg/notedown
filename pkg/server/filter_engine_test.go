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

package server

import (
	"testing"

	"github.com/notedownorg/notedown/apis/go/application_server/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestEvaluateFilter(t *testing.T) {

	// Sample metadata for testing
	metadata := map[string]any{
		"title":    "Test Document",
		"status":   "active",
		"priority": "high",
		"tags":     []any{"project", "important"},
		"version":  1.5,
		"count":    42,
		"draft":    false,
	}

	t.Run("nil filter returns true", func(t *testing.T) {
		result, err := EvaluateFilter(nil, metadata)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("equals operator", func(t *testing.T) {
		value, err := structpb.NewValue("active")
		require.NoError(t, err)

		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_MetadataFilter{
				MetadataFilter: &v1alpha1.MetadataFilter{
					Field:    "status",
					Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_EQUALS,
					Value:    value,
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result)

		// Test non-matching value
		wrongValue, err := structpb.NewValue("inactive")
		require.NoError(t, err)
		filter.GetMetadataFilter().Value = wrongValue

		result, err = EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("not equals operator", func(t *testing.T) {
		value, err := structpb.NewValue("inactive")
		require.NoError(t, err)

		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_MetadataFilter{
				MetadataFilter: &v1alpha1.MetadataFilter{
					Field:    "status",
					Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_NOT_EQUALS,
					Value:    value,
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("contains operator - string", func(t *testing.T) {
		value, err := structpb.NewValue("Test")
		require.NoError(t, err)

		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_MetadataFilter{
				MetadataFilter: &v1alpha1.MetadataFilter{
					Field:    "title",
					Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_CONTAINS,
					Value:    value,
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("contains operator - array", func(t *testing.T) {
		value, err := structpb.NewValue("project")
		require.NoError(t, err)

		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_MetadataFilter{
				MetadataFilter: &v1alpha1.MetadataFilter{
					Field:    "tags",
					Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_CONTAINS,
					Value:    value,
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("starts with operator", func(t *testing.T) {
		value, err := structpb.NewValue("Test")
		require.NoError(t, err)

		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_MetadataFilter{
				MetadataFilter: &v1alpha1.MetadataFilter{
					Field:    "title",
					Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_STARTS_WITH,
					Value:    value,
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("ends with operator", func(t *testing.T) {
		value, err := structpb.NewValue("Document")
		require.NoError(t, err)

		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_MetadataFilter{
				MetadataFilter: &v1alpha1.MetadataFilter{
					Field:    "title",
					Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_ENDS_WITH,
					Value:    value,
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("numeric comparisons", func(t *testing.T) {
		tests := []struct {
			operator v1alpha1.MetadataOperator
			value    float64
			expected bool
		}{
			{v1alpha1.MetadataOperator_METADATA_OPERATOR_GREATER_THAN, 1.0, true},
			{v1alpha1.MetadataOperator_METADATA_OPERATOR_GREATER_THAN, 2.0, false},
			{v1alpha1.MetadataOperator_METADATA_OPERATOR_GREATER_THAN_OR_EQUAL, 1.5, true},
			{v1alpha1.MetadataOperator_METADATA_OPERATOR_LESS_THAN, 2.0, true},
			{v1alpha1.MetadataOperator_METADATA_OPERATOR_LESS_THAN, 1.0, false},
			{v1alpha1.MetadataOperator_METADATA_OPERATOR_LESS_THAN_OR_EQUAL, 1.5, true},
		}

		for _, test := range tests {
			value, err := structpb.NewValue(test.value)
			require.NoError(t, err)

			filter := &v1alpha1.FilterExpression{
				Expression: &v1alpha1.FilterExpression_MetadataFilter{
					MetadataFilter: &v1alpha1.MetadataFilter{
						Field:    "version",
						Operator: test.operator,
						Value:    value,
					},
				},
			}

			result, err := EvaluateFilter(filter, metadata)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result, "operator %v with value %v", test.operator, test.value)
		}
	})

	t.Run("in array operator", func(t *testing.T) {
		arrayValue := &structpb.ListValue{
			Values: []*structpb.Value{
				structpb.NewStringValue("active"),
				structpb.NewStringValue("pending"),
			},
		}

		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_MetadataFilter{
				MetadataFilter: &v1alpha1.MetadataFilter{
					Field:    "status",
					Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_IN,
					Value:    structpb.NewListValue(arrayValue),
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("exists operator", func(t *testing.T) {
		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_MetadataFilter{
				MetadataFilter: &v1alpha1.MetadataFilter{
					Field:    "title",
					Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_EXISTS,
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result)

		// Test non-existent field
		filter.GetMetadataFilter().Field = "nonexistent"
		result, err = EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("not exists operator", func(t *testing.T) {
		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_MetadataFilter{
				MetadataFilter: &v1alpha1.MetadataFilter{
					Field:    "nonexistent",
					Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_NOT_EXISTS,
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("AND filter", func(t *testing.T) {
		statusValue, err := structpb.NewValue("active")
		require.NoError(t, err)

		priorityValue, err := structpb.NewValue("high")
		require.NoError(t, err)

		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_AndFilter{
				AndFilter: &v1alpha1.AndFilter{
					Filters: []*v1alpha1.FilterExpression{
						{
							Expression: &v1alpha1.FilterExpression_MetadataFilter{
								MetadataFilter: &v1alpha1.MetadataFilter{
									Field:    "status",
									Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_EQUALS,
									Value:    statusValue,
								},
							},
						},
						{
							Expression: &v1alpha1.FilterExpression_MetadataFilter{
								MetadataFilter: &v1alpha1.MetadataFilter{
									Field:    "priority",
									Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_EQUALS,
									Value:    priorityValue,
								},
							},
						},
					},
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("OR filter", func(t *testing.T) {
		statusValue, err := structpb.NewValue("inactive")
		require.NoError(t, err)

		priorityValue, err := structpb.NewValue("high")
		require.NoError(t, err)

		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_OrFilter{
				OrFilter: &v1alpha1.OrFilter{
					Filters: []*v1alpha1.FilterExpression{
						{
							Expression: &v1alpha1.FilterExpression_MetadataFilter{
								MetadataFilter: &v1alpha1.MetadataFilter{
									Field:    "status",
									Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_EQUALS,
									Value:    statusValue,
								},
							},
						},
						{
							Expression: &v1alpha1.FilterExpression_MetadataFilter{
								MetadataFilter: &v1alpha1.MetadataFilter{
									Field:    "priority",
									Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_EQUALS,
									Value:    priorityValue,
								},
							},
						},
					},
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result) // Should match because priority is high
	})

	t.Run("NOT filter", func(t *testing.T) {
		statusValue, err := structpb.NewValue("inactive")
		require.NoError(t, err)

		filter := &v1alpha1.FilterExpression{
			Expression: &v1alpha1.FilterExpression_NotFilter{
				NotFilter: &v1alpha1.NotFilter{
					Filter: &v1alpha1.FilterExpression{
						Expression: &v1alpha1.FilterExpression_MetadataFilter{
							MetadataFilter: &v1alpha1.MetadataFilter{
								Field:    "status",
								Operator: v1alpha1.MetadataOperator_METADATA_OPERATOR_EQUALS,
								Value:    statusValue,
							},
						},
					},
				},
			},
		}

		result, err := EvaluateFilter(filter, metadata)
		require.NoError(t, err)
		assert.True(t, result) // Should be true because status is NOT "inactive"
	})
}
