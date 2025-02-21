package book

import (
	"fmt"
	"github.com/sdreger/lib-manager-go/cmd/api/errors"
	"net/url"
	"strconv"
)

const (
	queryParamLanguageFilter  = "language"
	queryParamPublisherFilter = "publisher"
	queryParamAuthorFilter    = "author"
	queryParamCategoryFilter  = "category"
	queryParamFileTypeFilter  = "file_type"
	queryParamTagFilter       = "tag"
	queryParamQueryFilter     = "query"
	queryParamSbnFilter       = "sbn"
)

type Filter struct {
	Languages  []int64
	Publishers []int64
	Authors    []int64
	Categories []int64
	FileTypes  []int64
	Tags       []int64
	Query      string
	SBN        string // Standard Book Number, one of: ISBN10 / ISBN13 / ASIN
}

func NewFilter(queryValues url.Values) (Filter, error) {
	languages := queryValues[queryParamLanguageFilter]
	publishers := queryValues[queryParamPublisherFilter]
	authors := queryValues[queryParamAuthorFilter]
	categories := queryValues[queryParamCategoryFilter]
	fileTypes := queryValues[queryParamFileTypeFilter]
	tags := queryValues[queryParamTagFilter]
	query := queryValues.Get(queryParamQueryFilter)
	sbn := queryValues.Get(queryParamSbnFilter)

	languageIDs, err := parseFilterValues(languages)
	if err != nil {
		return Filter{}, errors.ValidationError{
			Field:   "language",
			Message: fmt.Sprintf("language ID values must be a number greater than or equal to 1: %v", languages),
		}
	}

	publisherIDs, err := parseFilterValues(publishers)
	if err != nil {
		return Filter{}, errors.ValidationError{
			Field:   "publisher",
			Message: fmt.Sprintf("publisher ID values must be a number greater than or equal to 1: %v", publishers),
		}
	}

	authorIDs, err := parseFilterValues(authors)
	if err != nil {
		return Filter{}, errors.ValidationError{
			Field:   "author",
			Message: fmt.Sprintf("author ID values must be a number greater than or equal to 1: %v", authors),
		}
	}

	categoryIDs, err := parseFilterValues(categories)
	if err != nil {
		return Filter{}, errors.ValidationError{
			Field:   "category",
			Message: fmt.Sprintf("category ID values must be a number greater than or equal to 1: %v", categories),
		}
	}

	fileTypeIDs, err := parseFilterValues(fileTypes)
	if err != nil {
		return Filter{}, errors.ValidationError{
			Field:   "file_type",
			Message: fmt.Sprintf("file_type ID values must be a number greater than or equal to 1: %v", fileTypes),
		}
	}

	tagIDs, err := parseFilterValues(tags)
	if err != nil {
		return Filter{}, errors.ValidationError{
			Field:   "tag",
			Message: fmt.Sprintf("tag ID values must be a number greater than or equal to 1: %v", tags),
		}
	}

	return Filter{
		Languages:  languageIDs,
		Publishers: publisherIDs,
		Authors:    authorIDs,
		Categories: categoryIDs,
		FileTypes:  fileTypeIDs,
		Tags:       tagIDs,
		Query:      query,
		SBN:        sbn,
	}, nil
}

func parseFilterValues(input []string) ([]int64, error) {
	result := make([]int64, len(input))
	for i, stringValue := range input {
		filterValue, err := strconv.Atoi(stringValue)
		if err != nil {
			return nil, err
		}
		if filterValue < 1 {
			return nil, fmt.Errorf("invalid filter value: %s", stringValue)
		}
		result[i] = int64(filterValue)
	}

	return result, nil
}
