package community_searcher

import (
	"reflect"
	"testing"
	"time"
)

func mustParseDate(value string) time.Time {
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		panic(err)
	}
	return date
}

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expected    ParsedSearchQuery
		expectError bool
	}{
		{
			name:  "Basic integer query with two filters",
			query: "TwirApp messages:>=50 usedChannelPoints:<=100",
			expected: ParsedSearchQuery{
				Username: "TwirApp",
				Conditions: []Condition{
					{Field: "messages", Operator: ">=", Value: 50, Type: IntegerType},
					{Field: "usedChannelPoints", Operator: "<=", Value: 100, Type: IntegerType},
				},
			},
			expectError: false,
		},
		{
			name:  "Two fields without username",
			query: "messages:>2 emotes:>30",
			expected: ParsedSearchQuery{
				Username: "",
				Conditions: []Condition{
					{Field: "messages", Operator: ">", Value: 2, Type: IntegerType},
					{Field: "emotes", Operator: ">", Value: 30, Type: IntegerType},
				},
			},
			expectError: false,
		},
		{
			name:  "Simple query with only one filter",
			query: "TwirApp messages:>=50",
			expected: ParsedSearchQuery{
				Username: "TwirApp",
				Conditions: []Condition{
					{Field: "messages", Operator: ">=", Value: 50, Type: IntegerType},
				},
			},
			expectError: false,
		},
		{
			name:  "Multiple usernames and one filter, so, just take first one",
			query: "TwirApp messages:>=50 MellKam",
			expected: ParsedSearchQuery{
				Username: "TwirApp",
				Conditions: []Condition{
					{Field: "messages", Operator: ">=", Value: 50, Type: IntegerType},
				},
			},
			expectError: false,
		},
		{
			name:        "Unsupported field",
			query:       "TwirApp unknown_field:>=50",
			expected:    ParsedSearchQuery{},
			expectError: true,
		},
		{
			name:        "Invalid integer value",
			query:       "messages:>=abc",
			expected:    ParsedSearchQuery{},
			expectError: true,
		},
		{
			name:  "Username without explicit field",
			query: "TwirApp",
			expected: ParsedSearchQuery{
				Username:   "TwirApp",
				Conditions: []Condition{},
			},
			expectError: false,
		},
		{
			name:  "Empty query",
			query: "",
			expected: ParsedSearchQuery{
				Username:   "",
				Conditions: []Condition{},
			},
			expectError: false,
		},
		{
			name:        "Invalid syntax must behave as empty search",
			query:       "TwirAppunknown_field:>=50 gsgs",
			expected:    ParsedSearchQuery{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := NewCommunitySearcher()
			result, err := cs.ParseSearchQuery(tt.query)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Username != tt.expected.Username {
				t.Errorf("Expected username %q, got %q", tt.expected.Username, result.Username)
			}

			if len(result.Conditions) != len(tt.expected.Conditions) {
				t.Fatalf(
					"Expected %d conditions, got %d",
					len(tt.expected.Conditions),
					len(result.Conditions),
				)
			}
			for i, expectedCondition := range tt.expected.Conditions {
				actualCondition := result.Conditions[i]
				if actualCondition.Field != expectedCondition.Field ||
					actualCondition.Operator != expectedCondition.Operator ||
					actualCondition.Type != expectedCondition.Type ||
					!reflect.DeepEqual(actualCondition.Value, expectedCondition.Value) {
					t.Errorf(
						"Condition mismatch at index %d: expected %+v, got %+v",
						i,
						expectedCondition,
						actualCondition,
					)
				}
			}
		})
	}
}
