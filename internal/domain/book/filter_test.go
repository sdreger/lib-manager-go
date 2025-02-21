package book

import (
	"github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewFilter(t *testing.T) {

	tt := []struct {
		language           []string
		publisher          []string
		author             []string
		category           []string
		fileType           []string
		tag                []string
		query              []string
		sbn                []string
		expectedLanguages  []int64
		expectedPublishers []int64
		expectedAuthors    []int64
		expectedCategories []int64
		expectedFileTypes  []int64
		expectedTags       []int64
		expectedQuery      string
		expectedSbn        string
		err                bool
	}{
		{
			language: []string{"1", "2", "3"}, expectedLanguages: []int64{1, 2, 3},
			expectedPublishers: []int64{}, expectedAuthors: []int64{}, expectedCategories: []int64{},
			expectedFileTypes: []int64{}, expectedTags: []int64{}, expectedQuery: "", expectedSbn: "",
			err: false,
		},
		{
			publisher: []string{"2", "3", "4"}, expectedPublishers: []int64{2, 3, 4},
			expectedLanguages: []int64{}, expectedAuthors: []int64{}, expectedCategories: []int64{},
			expectedFileTypes: []int64{}, expectedTags: []int64{}, expectedQuery: "", expectedSbn: "",
			err: false,
		},
		{
			author: []string{"3", "4", "5"}, expectedAuthors: []int64{3, 4, 5},
			expectedLanguages: []int64{}, expectedPublishers: []int64{}, expectedCategories: []int64{},
			expectedFileTypes: []int64{}, expectedTags: []int64{}, expectedQuery: "", expectedSbn: "",
			err: false,
		},
		{
			category: []string{"4", "5", "6"}, expectedCategories: []int64{4, 5, 6},
			expectedLanguages: []int64{}, expectedPublishers: []int64{}, expectedAuthors: []int64{},
			expectedFileTypes: []int64{}, expectedTags: []int64{}, expectedQuery: "", expectedSbn: "",
			err: false,
		},
		{
			fileType: []string{"7", "8", "9"}, expectedFileTypes: []int64{7, 8, 9},
			expectedLanguages: []int64{}, expectedPublishers: []int64{}, expectedAuthors: []int64{},
			expectedCategories: []int64{}, expectedTags: []int64{}, expectedQuery: "", expectedSbn: "",
			err: false,
		},
		{
			tag: []string{"8", "9", "10"}, expectedTags: []int64{8, 9, 10},
			expectedLanguages: []int64{}, expectedPublishers: []int64{}, expectedAuthors: []int64{},
			expectedCategories: []int64{}, expectedFileTypes: []int64{}, expectedQuery: "", expectedSbn: "",
			err: false,
		},
		{
			query: []string{"computers"}, expectedQuery: "computers",
			expectedLanguages: []int64{}, expectedPublishers: []int64{}, expectedAuthors: []int64{},
			expectedCategories: []int64{}, expectedFileTypes: []int64{}, expectedTags: []int64{}, expectedSbn: "",
			err: false,
		},
		{
			query: []string{"computers", "rockets", "cars"}, expectedQuery: "computers",
			expectedLanguages: []int64{}, expectedPublishers: []int64{}, expectedAuthors: []int64{},
			expectedCategories: []int64{}, expectedFileTypes: []int64{}, expectedTags: []int64{}, expectedSbn: "",
			err: false,
		},
		{
			sbn: []string{"1111111111"}, expectedSbn: "1111111111",
			expectedLanguages: []int64{}, expectedPublishers: []int64{}, expectedAuthors: []int64{},
			expectedCategories: []int64{}, expectedFileTypes: []int64{}, expectedTags: []int64{}, expectedQuery: "",
			err: false,
		},
		{
			sbn: []string{"1111111111", "2222222222", "3333333333"}, expectedSbn: "1111111111",
			expectedLanguages: []int64{}, expectedPublishers: []int64{}, expectedAuthors: []int64{},
			expectedCategories: []int64{}, expectedFileTypes: []int64{}, expectedTags: []int64{}, expectedQuery: "",
			err: false,
		},
		{
			// no filters
			expectedLanguages: []int64{}, expectedPublishers: []int64{}, expectedAuthors: []int64{},
			expectedCategories: []int64{}, expectedFileTypes: []int64{}, expectedTags: []int64{}, expectedQuery: "",
			expectedSbn: "", err: false,
		},
		{
			// all filters
			language: []string{"1"}, publisher: []string{"2"}, author: []string{"3"}, category: []string{"4"},
			fileType: []string{"5"}, tag: []string{"6"}, query: []string{"computers"}, sbn: []string{"1111111111"},
			expectedLanguages: []int64{1}, expectedPublishers: []int64{2}, expectedAuthors: []int64{3},
			expectedCategories: []int64{4}, expectedFileTypes: []int64{5}, expectedTags: []int64{6},
			expectedQuery: "computers", expectedSbn: "1111111111", err: false,
		},
		{language: []string{"one"}, err: true},
		{publisher: []string{"two"}, err: true},
		{author: []string{"three"}, err: true},
		{category: []string{"four"}, err: true},
		{fileType: []string{"five"}, err: true},
		{tag: []string{"six"}, err: true},
		{language: []string{"0"}, err: true},
		{publisher: []string{"-1"}, err: true},
		{author: []string{"1.5"}, err: true},
		{category: []string{"1E5"}, err: true},
		{fileType: []string{"0"}, err: true},
		{tag: []string{"0"}, err: true},
	}

	for _, tc := range tt {
		t.Run(t.Name(), func(t *testing.T) {
			values := map[string][]string{
				queryParamLanguageFilter:  tc.language,
				queryParamPublisherFilter: tc.publisher,
				queryParamAuthorFilter:    tc.author,
				queryParamCategoryFilter:  tc.category,
				queryParamFileTypeFilter:  tc.fileType,
				queryParamTagFilter:       tc.tag,
				queryParamQueryFilter:     tc.query,
				queryParamSbnFilter:       tc.sbn,
			}

			filter, err := NewFilter(values)
			if tc.err {
				require.Error(t, err)
				assert.ErrorAs(t, err, &errors.ValidationError{})
			} else {
				require.NoError(t, err, "should create filter")
				assert.Equal(t, tc.expectedLanguages, filter.Languages)
				assert.Equal(t, tc.expectedPublishers, filter.Publishers)
				assert.Equal(t, tc.expectedAuthors, filter.Authors)
				assert.Equal(t, tc.expectedCategories, filter.Categories)
				assert.Equal(t, tc.expectedFileTypes, filter.FileTypes)
				assert.Equal(t, tc.expectedTags, filter.Tags)
				assert.Equal(t, tc.expectedQuery, filter.Query)
				assert.Equal(t, tc.expectedSbn, filter.SBN)
			}
		})
	}
}
