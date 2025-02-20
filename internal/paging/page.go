package paging

import (
	"github.com/sdreger/lib-manager-go/cmd/api/errors"
	"net/url"
	"strconv"
)

const (
	defaultPage    = 1
	defaultSize    = 10
	queryParamPage = "page"
	queryParamSize = "size"
)

type Page[T any] struct {
	Page       int64 `json:"page"`
	Size       int64 `json:"size"`
	TotalPages int64 `json:"total_pages"`
	TotalItems int64 `json:"total_elements"`
	Content    []T   `json:"content"`
}

func NewPage[T any](pageRequest PageRequest, totalItems int64, content []T) Page[T] {
	totalPages := totalItems / pageRequest.size
	if totalItems%pageRequest.size != 0 {
		totalPages = totalItems/pageRequest.size + 1
	}

	return Page[T]{
		Page:       pageRequest.page,
		Size:       int64(len(content)),
		TotalPages: totalPages,
		TotalItems: totalItems,
		Content:    content,
	}
}

type PageRequest struct {
	page int64
	size int64
}

func NewPageRequest(queryValues url.Values) (PageRequest, error) {
	pageNumber := defaultPage
	pageSize := defaultSize
	pageNumberString := queryValues.Get(queryParamPage)
	pageSizeString := queryValues.Get(queryParamSize)

	if pageNumberString != "" {
		pageInt, err := strconv.Atoi(pageNumberString)
		if err != nil {
			return PageRequest{}, errors.ValidationError{
				Field:   "page",
				Message: "wrong page value: " + pageNumberString,
			}
		}
		if pageInt < 1 {
			return PageRequest{}, errors.ValidationError{
				Field:   "page",
				Message: "page must be greater than or equal to 1",
			}
		}
		pageNumber = pageInt
	}

	if pageSizeString != "" {
		sizeInt, err := strconv.Atoi(pageSizeString)
		if err != nil || sizeInt < 1 {
			return PageRequest{}, errors.ValidationError{
				Field:   "size",
				Message: "wrong page size value: " + pageSizeString,
			}
		}
		pageSize = sizeInt
	}

	return PageRequest{
		page: int64(pageNumber),
		size: int64(pageSize),
	}, nil
}

func (req PageRequest) Offset() uint64 {
	return uint64((req.page - 1) * req.size)
}

func (req PageRequest) Limit() uint64 {
	return uint64(req.size)
}
