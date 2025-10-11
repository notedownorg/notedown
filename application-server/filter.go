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
	"fmt"
	"reflect"
	"strings"

	pb "github.com/notedownorg/notedown/notedown/application_server/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// filterEngine provides document filtering capabilities
type filterEngine struct {
}

// newFilterEngine creates a new filter engine
func newFilterEngine() *filterEngine {
	return &filterEngine{}
}

// FilterDocuments filters documents from a channel based on the given filter expression
// Returns a channel of filtered documents and an error channel
func (fe *filterEngine) FilterDocuments(docChan <-chan *pb.Document, filter *pb.FilterExpression) (<-chan *pb.Document, <-chan error) {
	filteredChan := make(chan *pb.Document)
	errChan := make(chan error)

	go func() {
		defer close(filteredChan)
		defer close(errChan)

		for doc := range docChan {
			if filter == nil {
				filteredChan <- doc
				continue
			}

			matches, err := fe.evaluateFilter(doc, filter)
			if err != nil {
				errChan <- err
				return
			}
			if matches {
				filteredChan <- doc
			}
		}
	}()

	return filteredChan, errChan
}

// evaluateFilter evaluates a filter expression against a document
func (fe *filterEngine) evaluateFilter(doc *pb.Document, filter *pb.FilterExpression) (bool, error) {
	if filter == nil {
		return true, nil
	}

	switch expr := filter.Expression.(type) {
	case *pb.FilterExpression_MetadataFilter:
		return fe.evaluateMetadataFilter(doc, expr.MetadataFilter)
	case *pb.FilterExpression_AndFilter:
		return fe.evaluateAndFilter(doc, expr.AndFilter)
	case *pb.FilterExpression_OrFilter:
		return fe.evaluateOrFilter(doc, expr.OrFilter)
	case *pb.FilterExpression_NotFilter:
		return fe.evaluateNotFilter(doc, expr.NotFilter)
	default:
		return false, fmt.Errorf("unknown filter expression type: %T", expr)
	}
}

// evaluateMetadataFilter evaluates a metadata filter
func (fe *filterEngine) evaluateMetadataFilter(doc *pb.Document, filter *pb.MetadataFilter) (bool, error) {
	if filter == nil {
		return true, nil
	}

	// Get the metadata field value from protobuf Struct
	var value any
	var exists bool

	if doc.Metadata != nil && doc.Metadata.Fields != nil {
		if pbValue, found := doc.Metadata.Fields[filter.Field]; found {
			value = pbValue.AsInterface()
			exists = true
		}
	}

	switch filter.Operator {
	case pb.MetadataOperator_METADATA_OPERATOR_EXISTS:
		return exists, nil
	case pb.MetadataOperator_METADATA_OPERATOR_NOT_EXISTS:
		return !exists, nil
	}

	// For other operators, the field must exist
	if !exists {
		return false, nil
	}

	// Convert protobuf Value to Go value
	filterValue := protoValueToGo(filter.Value)

	switch filter.Operator {
	case pb.MetadataOperator_METADATA_OPERATOR_EQUALS:
		// Handle numeric equality specially to account for type differences
		if v1, ok1 := toFloat64(value); ok1 {
			if v2, ok2 := toFloat64(filterValue); ok2 {
				return v1 == v2, nil
			}
		}
		return reflect.DeepEqual(value, filterValue), nil
	case pb.MetadataOperator_METADATA_OPERATOR_NOT_EQUALS:
		return !reflect.DeepEqual(value, filterValue), nil
	case pb.MetadataOperator_METADATA_OPERATOR_CONTAINS:
		return fe.containsValue(value, filterValue)
	case pb.MetadataOperator_METADATA_OPERATOR_STARTS_WITH:
		return fe.startsWithValue(value, filterValue)
	case pb.MetadataOperator_METADATA_OPERATOR_ENDS_WITH:
		return fe.endsWithValue(value, filterValue)
	case pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN:
		return fe.compareNumbers(value, filterValue, ">")
	case pb.MetadataOperator_METADATA_OPERATOR_GREATER_THAN_OR_EQUAL:
		return fe.compareNumbers(value, filterValue, ">=")
	case pb.MetadataOperator_METADATA_OPERATOR_LESS_THAN:
		return fe.compareNumbers(value, filterValue, "<")
	case pb.MetadataOperator_METADATA_OPERATOR_LESS_THAN_OR_EQUAL:
		return fe.compareNumbers(value, filterValue, "<=")
	case pb.MetadataOperator_METADATA_OPERATOR_IN:
		return fe.inArray(value, filterValue)
	case pb.MetadataOperator_METADATA_OPERATOR_NOT_IN:
		result, err := fe.inArray(value, filterValue)
		return !result, err
	default:
		return false, fmt.Errorf("unknown metadata operator: %v", filter.Operator)
	}
}

// evaluateAndFilter evaluates an AND filter
func (fe *filterEngine) evaluateAndFilter(doc *pb.Document, filter *pb.AndFilter) (bool, error) {
	if filter == nil || len(filter.Filters) == 0 {
		return true, nil
	}

	for _, subFilter := range filter.Filters {
		matches, err := fe.evaluateFilter(doc, subFilter)
		if err != nil {
			return false, err
		}
		if !matches {
			return false, nil
		}
	}
	return true, nil
}

// evaluateOrFilter evaluates an OR filter
func (fe *filterEngine) evaluateOrFilter(doc *pb.Document, filter *pb.OrFilter) (bool, error) {
	if filter == nil || len(filter.Filters) == 0 {
		return true, nil
	}

	for _, subFilter := range filter.Filters {
		matches, err := fe.evaluateFilter(doc, subFilter)
		if err != nil {
			return false, err
		}
		if matches {
			return true, nil
		}
	}
	return false, nil
}

// evaluateNotFilter evaluates a NOT filter
func (fe *filterEngine) evaluateNotFilter(doc *pb.Document, filter *pb.NotFilter) (bool, error) {
	if filter == nil {
		return true, nil
	}

	matches, err := fe.evaluateFilter(doc, filter.Filter)
	if err != nil {
		return false, err
	}
	return !matches, nil
}

// Helper functions for metadata comparison

func (fe *filterEngine) containsValue(value, searchValue any) (bool, error) {
	// If value is an array/slice, check if it contains the search value
	if valueReflect := reflect.ValueOf(value); valueReflect.Kind() == reflect.Slice || valueReflect.Kind() == reflect.Array {
		for i := 0; i < valueReflect.Len(); i++ {
			item := valueReflect.Index(i).Interface()
			// Try direct equality first
			if reflect.DeepEqual(item, searchValue) {
				return true, nil
			}
			// Try numeric equality for numbers
			if v1, ok1 := toFloat64(item); ok1 {
				if v2, ok2 := toFloat64(searchValue); ok2 {
					if v1 == v2 {
						return true, nil
					}
				}
			}
		}
		return false, nil
	}

	// For non-array types, only do string contains on strings
	if _, ok := value.(string); ok {
		valueStr := fmt.Sprintf("%v", value)
		searchStr := fmt.Sprintf("%v", searchValue)
		return strings.Contains(valueStr, searchStr), nil
	}

	// For non-string, non-array types, contains doesn't make sense
	return false, nil
}

func (fe *filterEngine) startsWithValue(value, searchValue any) (bool, error) {
	valueStr := fmt.Sprintf("%v", value)
	searchStr := fmt.Sprintf("%v", searchValue)
	return strings.HasPrefix(valueStr, searchStr), nil
}

func (fe *filterEngine) endsWithValue(value, searchValue any) (bool, error) {
	valueStr := fmt.Sprintf("%v", value)
	searchStr := fmt.Sprintf("%v", searchValue)
	return strings.HasSuffix(valueStr, searchStr), nil
}

func (fe *filterEngine) compareNumbers(value, compareValue any, operator string) (bool, error) {
	// Try numeric comparison first
	v1, ok1 := toFloat64(value)
	v2, ok2 := toFloat64(compareValue)

	if ok1 && ok2 {
		switch operator {
		case ">":
			return v1 > v2, nil
		case ">=":
			return v1 >= v2, nil
		case "<":
			return v1 < v2, nil
		case "<=":
			return v1 <= v2, nil
		default:
			return false, fmt.Errorf("unknown operator: %s", operator)
		}
	}

	// Try string comparison
	s1, ok1 := value.(string)
	s2, ok2 := compareValue.(string)

	if ok1 && ok2 {
		switch operator {
		case ">":
			return s1 > s2, nil
		case ">=":
			return s1 >= s2, nil
		case "<":
			return s1 < s2, nil
		case "<=":
			return s1 <= s2, nil
		default:
			return false, fmt.Errorf("unknown operator: %s", operator)
		}
	}

	// Cannot compare these types - return false instead of error
	return false, nil
}

func (fe *filterEngine) inArray(value, arrayValue any) (bool, error) {
	// Convert arrayValue to slice
	arrayReflect := reflect.ValueOf(arrayValue)
	if arrayReflect.Kind() != reflect.Slice && arrayReflect.Kind() != reflect.Array {
		// Return false instead of error for non-array values
		return false, nil
	}

	for i := 0; i < arrayReflect.Len(); i++ {
		item := arrayReflect.Index(i).Interface()
		// Try direct equality first
		if reflect.DeepEqual(value, item) {
			return true, nil
		}
		// Try numeric equality for numbers
		if v1, ok1 := toFloat64(value); ok1 {
			if v2, ok2 := toFloat64(item); ok2 {
				if v1 == v2 {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// Utility functions

func protoValueToGo(value *structpb.Value) any {
	if value == nil {
		return nil
	}

	switch v := value.GetKind().(type) {
	case *structpb.Value_NullValue:
		return nil
	case *structpb.Value_NumberValue:
		return v.NumberValue
	case *structpb.Value_StringValue:
		return v.StringValue
	case *structpb.Value_BoolValue:
		return v.BoolValue
	case *structpb.Value_StructValue:
		// Convert struct to map
		result := make(map[string]any)
		for key, val := range v.StructValue.GetFields() {
			result[key] = protoValueToGo(val)
		}
		return result
	case *structpb.Value_ListValue:
		// Convert list to slice
		result := make([]any, len(v.ListValue.GetValues()))
		for i, val := range v.ListValue.GetValues() {
			result[i] = protoValueToGo(val)
		}
		return result
	default:
		return nil
	}
}

func toFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}
