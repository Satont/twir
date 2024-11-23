package community_searcher

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FieldType int

const (
	IntegerType FieldType = iota
	DateType
	StringType
)

type CommunitySearcher struct {
	SearchFields  map[string]FieldType
	searchPattern string
}

type Condition struct {
	Field    string
	Operator string
	Value    interface{}
	Type     FieldType
}
type ParsedSearchQuery struct {
	Conditions []Condition
	Username   string
}

var defaultFields = map[string]FieldType{
	"messages":          IntegerType,
	"usedChannelPoints": IntegerType,
	"emotes":            IntegerType,
	"watched":           IntegerType,
}

const defaultSearchPattern = `(\w+):([><]=?|=)(\S+)`

func NewCommunitySearcher() *CommunitySearcher {
	return &CommunitySearcher{
		SearchFields:  defaultFields,
		searchPattern: defaultSearchPattern,
	}
}

func (cs *CommunitySearcher) ParseSearchQuery(query string) (*ParsedSearchQuery, error) {
	re := regexp.MustCompile(cs.searchPattern)
	matches := re.FindAllStringSubmatch(query, -1)

	var conditions []Condition
	var username string

	remainingTokens := query
	for _, match := range matches {
		field := match[1]
		operator := match[2]
		rawValue := match[3]

		fieldType, exists := cs.SearchFields[field]
		if !exists {
			return nil, fmt.Errorf("unknown field: %s", field)
		}

		var value interface{}
		var err error
		switch fieldType {
		case IntegerType:
			value, err = strconv.Atoi(rawValue)
		case DateType:
			value, err = time.Parse("2006-01-02", rawValue) // "YYYY-MM-DD"
		case StringType:
			value = rawValue
		}
		if err != nil {
			return nil, err
		}

		conditions = append(conditions, Condition{
			Field:    field,
			Operator: operator,
			Value:    value,
			Type:     fieldType,
		})

		remainingTokens = strings.Replace(remainingTokens, match[0], "", 1)
		remainingTokens = strings.TrimSpace(remainingTokens)
	}

	if remainingTokens != "" && username == "" {
		username = strings.Fields(remainingTokens)[0]
	}

	return &ParsedSearchQuery{
		Conditions: conditions,
		Username:   username,
	}, nil
}
