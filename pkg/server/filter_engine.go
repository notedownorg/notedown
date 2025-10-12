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
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/notedownorg/notedown/apis/go/application_server/v1alpha1"
	"google.golang.org/protobuf/types/known/structpb"
)

// EvaluateFilter evaluates a filter expression against document metadata
func EvaluateFilter(filter *v1alpha1.FilterExpression, metadata map[string]any) (bool, error) {
	if filter == nil {
		return true, nil // No filter means all documents match
	}

	switch expr := filter.Expression.(type) {
	case *v1alpha1.FilterExpression_MetadataFilter:
		return evaluateMetadataFilter(expr.MetadataFilter, metadata)
	case *v1alpha1.FilterExpression_AndFilter:
		return evaluateAndFilter(expr.AndFilter, metadata)
	case *v1alpha1.FilterExpression_OrFilter:
		return evaluateOrFilter(expr.OrFilter, metadata)
	case *v1alpha1.FilterExpression_NotFilter:
		return evaluateNotFilter(expr.NotFilter, metadata)
	default:
		return false, fmt.Errorf("unknown filter expression type: %T", expr)
	}
}

// evaluateMetadataFilter evaluates a metadata filter
func evaluateMetadataFilter(filter *v1alpha1.MetadataFilter, metadata map[string]any) (bool, error) {
	if filter == nil {
		return true, nil
	}

	// Get the field value from metadata
	fieldValue, exists := metadata[filter.Field]

	switch filter.Operator {
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_EXISTS:
		return exists, nil
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_NOT_EXISTS:
		return !exists, nil
	}

	// For all other operators, the field must exist
	if !exists {
		return false, nil
	}

	// Convert protobuf Value to Go value
	filterValue, err := protoValueToGoValue(filter.Value)
	if err != nil {
		return false, fmt.Errorf("failed to convert filter value: %w", err)
	}

	return compareValues(fieldValue, filterValue, filter.Operator)
}

// evaluateAndFilter evaluates an AND filter (all must be true)
func evaluateAndFilter(filter *v1alpha1.AndFilter, metadata map[string]any) (bool, error) {
	if filter == nil || len(filter.Filters) == 0 {
		return true, nil
	}

	for _, subFilter := range filter.Filters {
		result, err := EvaluateFilter(subFilter, metadata)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil // Early exit on first false
		}
	}
	return true, nil
}

// evaluateOrFilter evaluates an OR filter (any must be true)
func evaluateOrFilter(filter *v1alpha1.OrFilter, metadata map[string]any) (bool, error) {
	if filter == nil || len(filter.Filters) == 0 {
		return true, nil
	}

	for _, subFilter := range filter.Filters {
		result, err := EvaluateFilter(subFilter, metadata)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil // Early exit on first true
		}
	}
	return false, nil
}

// evaluateNotFilter evaluates a NOT filter
func evaluateNotFilter(filter *v1alpha1.NotFilter, metadata map[string]any) (bool, error) {
	if filter == nil || filter.Filter == nil {
		return true, nil
	}

	result, err := EvaluateFilter(filter.Filter, metadata)
	if err != nil {
		return false, err
	}
	return !result, nil
}

// compareValues compares two values using the specified operator
func compareValues(fieldValue, filterValue any, operator v1alpha1.MetadataOperator) (bool, error) {
	switch operator {
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_EQUALS:
		return equalValues(fieldValue, filterValue), nil
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_NOT_EQUALS:
		return !equalValues(fieldValue, filterValue), nil
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_CONTAINS:
		return containsValue(fieldValue, filterValue)
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_STARTS_WITH:
		return startsWithValue(fieldValue, filterValue)
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_ENDS_WITH:
		return endsWithValue(fieldValue, filterValue)
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_GREATER_THAN:
		return compareNumeric(fieldValue, filterValue, ">")
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_GREATER_THAN_OR_EQUAL:
		return compareNumeric(fieldValue, filterValue, ">=")
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_LESS_THAN:
		return compareNumeric(fieldValue, filterValue, "<")
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_LESS_THAN_OR_EQUAL:
		return compareNumeric(fieldValue, filterValue, "<=")
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_IN:
		return inArray(fieldValue, filterValue)
	case v1alpha1.MetadataOperator_METADATA_OPERATOR_NOT_IN:
		result, err := inArray(fieldValue, filterValue)
		return !result, err
	default:
		return false, fmt.Errorf("unsupported operator: %v", operator)
	}
}

// equalValues checks if two values are equal
func equalValues(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

// containsValue checks if field contains filter value (for strings and arrays)
func containsValue(fieldValue, filterValue any) (bool, error) {
	fieldStr, ok1 := fieldValue.(string)
	filterStr, ok2 := filterValue.(string)

	if ok1 && ok2 {
		return strings.Contains(fieldStr, filterStr), nil
	}

	// Check if fieldValue is an array containing filterValue
	fieldSlice := reflect.ValueOf(fieldValue)
	if fieldSlice.Kind() == reflect.Slice || fieldSlice.Kind() == reflect.Array {
		for i := 0; i < fieldSlice.Len(); i++ {
			if equalValues(fieldSlice.Index(i).Interface(), filterValue) {
				return true, nil
			}
		}
		return false, nil
	}

	return false, fmt.Errorf("contains operator requires string or array field")
}

// startsWithValue checks if field starts with filter value
func startsWithValue(fieldValue, filterValue any) (bool, error) {
	fieldStr, ok1 := fieldValue.(string)
	filterStr, ok2 := filterValue.(string)

	if !ok1 || !ok2 {
		return false, fmt.Errorf("starts_with operator requires string values")
	}

	return strings.HasPrefix(fieldStr, filterStr), nil
}

// endsWithValue checks if field ends with filter value
func endsWithValue(fieldValue, filterValue any) (bool, error) {
	fieldStr, ok1 := fieldValue.(string)
	filterStr, ok2 := filterValue.(string)

	if !ok1 || !ok2 {
		return false, fmt.Errorf("ends_with operator requires string values")
	}

	return strings.HasSuffix(fieldStr, filterStr), nil
}

// compareNumeric compares numeric values
func compareNumeric(fieldValue, filterValue any, operator string) (bool, error) {
	fieldNum, err1 := toFloat64(fieldValue)
	filterNum, err2 := toFloat64(filterValue)

	if err1 != nil || err2 != nil {
		return false, fmt.Errorf("numeric comparison requires numeric values")
	}

	switch operator {
	case ">":
		return fieldNum > filterNum, nil
	case ">=":
		return fieldNum >= filterNum, nil
	case "<":
		return fieldNum < filterNum, nil
	case "<=":
		return fieldNum <= filterNum, nil
	default:
		return false, fmt.Errorf("unknown numeric operator: %s", operator)
	}
}

// inArray checks if fieldValue is in the filterValue array
func inArray(fieldValue, filterValue any) (bool, error) {
	filterSlice := reflect.ValueOf(filterValue)
	if filterSlice.Kind() != reflect.Slice && filterSlice.Kind() != reflect.Array {
		return false, fmt.Errorf("in operator requires array filter value")
	}

	for i := 0; i < filterSlice.Len(); i++ {
		if equalValues(fieldValue, filterSlice.Index(i).Interface()) {
			return true, nil
		}
	}
	return false, nil
}

// toFloat64 converts various numeric types to float64
func toFloat64(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// protoValueToGoValue converts protobuf Value to Go value
func protoValueToGoValue(value *structpb.Value) (any, error) {
	if value == nil {
		return nil, nil
	}

	switch v := value.Kind.(type) {
	case *structpb.Value_NullValue:
		return nil, nil
	case *structpb.Value_NumberValue:
		return v.NumberValue, nil
	case *structpb.Value_StringValue:
		return v.StringValue, nil
	case *structpb.Value_BoolValue:
		return v.BoolValue, nil
	case *structpb.Value_ListValue:
		var result []any
		for _, item := range v.ListValue.Values {
			goValue, err := protoValueToGoValue(item)
			if err != nil {
				return nil, err
			}
			result = append(result, goValue)
		}
		return result, nil
	case *structpb.Value_StructValue:
		result := make(map[string]any)
		for key, val := range v.StructValue.Fields {
			goValue, err := protoValueToGoValue(val)
			if err != nil {
				return nil, err
			}
			result[key] = goValue
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unknown protobuf value type: %T", v)
	}
}
