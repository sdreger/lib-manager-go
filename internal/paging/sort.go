package paging

import (
	"fmt"
	"github.com/sdreger/lib-manager-go/cmd/api/errors"
	"net/url"
	"strings"
)

const (
	sortAscending        = "ASC"
	sortDescending       = "DESC"
	defaultSortBy        = "updated_at"
	defaultSortDirection = sortDescending
	queryParamSort       = "sort"
)

var (
	allowedSortDirections = []string{sortAscending, sortDescending}
)

type Sort struct {
	field     string
	direction string
}

func NewSort(queryValues url.Values, allowedSortFields []string) (Sort, error) {
	sortString := queryValues.Get(queryParamSort)
	if sortString == "" {
		return Sort{
			field:     defaultSortBy,
			direction: defaultSortDirection,
		}, nil
	}

	sortStringParts := strings.Split(sortString, ",")
	if len(sortStringParts) != 2 {
		return Sort{}, errors.ValidationError{Field: queryParamSort, Message: "wrong sort request: " + sortString}
	}

	sortField := sortStringParts[0]
	sortDirection := sortStringParts[1]
	if !isFieldAllowed(sortField, allowedSortFields) {
		return Sort{}, errors.ValidationError{
			Field:   queryParamSort,
			Message: fmt.Sprintf("sort field %q is not allowed", sortField),
		}

	}
	if !isDirectionAllowed(sortDirection) {
		return Sort{}, errors.ValidationError{
			Field:   queryParamSort,
			Message: fmt.Sprintf("sort direction %q is not allowed", sortDirection),
		}
	}

	return Sort{
		field:     strings.ToLower(sortField),
		direction: strings.ToUpper(sortDirection),
	}, nil
}

func (s Sort) GetOrderBy() string {
	return fmt.Sprintf("%s %s", s.field, s.direction)
}

func isFieldAllowed(field string, allowedSortFields []string) bool {
	lowerField := strings.ToLower(field)
	for _, allowedField := range allowedSortFields {
		if allowedField == lowerField {
			return true
		}
	}

	return false
}

func isDirectionAllowed(direction string) bool {
	upperDirection := strings.ToUpper(direction)
	for _, allowedDirection := range allowedSortDirections {
		if allowedDirection == upperDirection {
			return true
		}
	}

	return false
}
