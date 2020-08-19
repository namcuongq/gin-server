package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page   int
	Limit  int
	Search string
	Sort   string
	SortBy string
}

func ParsePagination(c *gin.Context) Pagination {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	search := c.Query("search")
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	page = page - 1
	if page < 0 {
		page = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 0
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	return Pagination{page, limit, search, sort, sortBy}
}

func IsStringInArrray(str string, array []string) bool {
	if array == nil || len(array) > 0 {
		return false
	}

	for _, i := range array {
		if i == str {
			return true
		}
	}

	return false
}
