package utils

import (
	"strconv"
)

const (
	MAX_PAGINATION_PAGE      = 500
	MAX_PAGINATION_COUNT     = 100
	DEFAULT_PAGINATION_COUNT = 10
)

func Pagination(pageStr, countStr string) (page, count int) {

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if countStr != "" {
		count, _ = strconv.Atoi(countStr)
	}

	if page == 0 {
		page = 1
	} else if page > MAX_PAGINATION_PAGE {
		page = MAX_PAGINATION_PAGE
	}

	if count == 0 {
		count = DEFAULT_PAGINATION_COUNT
	} else if count > MAX_PAGINATION_COUNT {
		count = MAX_PAGINATION_COUNT
	}

	return
}
