package helper

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// QueryParams digunakan untuk menangani pagination, sorting, filtering, dan pencarian
type QueryParams struct {
	Page    int               `json:"page"`
	Limit   int               `json:"limit"`
	SortBy  string            `json:"sort_by"`
	Order   string            `json:"order"`
	Keyword string            `json:"keyword"`
	Filters map[string]string `json:"filters"`
}

// ParseQueryParams mengekstrak parameter query dari Fiber context
func ParseQueryParams(ctx *fiber.Ctx) QueryParams {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))    // Default: 1
	limit, _ := strconv.Atoi(ctx.Query("limit", "10")) // Default: 10
	sortBy := ctx.Query("sort_by", "+id")              // Default: +id (ascending)
	keyword := ctx.Query("keyword", "")                // Default: kosong

	// Cek apakah sort_by memiliki tanda "+" atau "-"
	order := "ASC"
	if strings.HasPrefix(sortBy, "-") {
		order = "DESC"
		sortBy = strings.TrimPrefix(sortBy, "-") // Hapus tanda "-"
	} else {
		sortBy = strings.TrimPrefix(sortBy, "+") // Hapus tanda "+"
	}

	// Ambil semua filter dinamis (query params kecuali yang sudah dikenal)
	filters := make(map[string]string)
	for key, values := range ctx.Queries() {
		if key != "page" && key != "limit" && key != "sort_by" && key != "keyword" {
			filters[key] = values
		}
	}

	return QueryParams{
		Page:    page,
		Limit:   limit,
		SortBy:  sortBy,
		Order:   order,
		Keyword: keyword,
		Filters: filters,
	}
}
