package helper

import (
	"math"
	"strings"

	"gorm.io/gorm"
)

// PaginatedResult menyimpan hasil query dengan metadata pagination
type PaginatedResult struct {
	Records    interface{}
	TotalCount int64
	PageCount  int
}

// ApplyFiltersAndPagination menerapkan filter, sorting, dan pagination ke query database
func ApplyFiltersAndPagination(db *gorm.DB, model interface{}, query QueryParams, searchableFields []string) PaginatedResult {
	// Hitung total data sebelum pagination
	var totalCount int64
	db.Model(model).Count(&totalCount)

	// Hitung offset
	offset := (query.Page - 1) * query.Limit

	// Query dengan filtering & sorting
	dbQuery := db.Model(model)

	// Pencarian global (jika ada keyword dan field tersedia)
	if query.Keyword != "" && len(searchableFields) > 0 {
		var searchConditions []string
		var searchValues []interface{}

		// Buat query pencarian dinamis berdasarkan field yang tersedia
		for _, field := range searchableFields {
			searchConditions = append(searchConditions, field+" LIKE ?")
			searchValues = append(searchValues, "%"+query.Keyword+"%")
		}

		// Gabungkan kondisi dengan "OR"
		searchQuery := strings.Join(searchConditions, " OR ")
		dbQuery = dbQuery.Where(searchQuery, searchValues...)
	}

	// Filter berdasarkan field yang dinamis
	for field, value := range query.Filters {
		dbQuery = dbQuery.Where(field+" = ?", value)
	}

	// Ambil data dengan pagination
	dbQuery = dbQuery.Order(query.SortBy + " " + query.Order).Limit(query.Limit).Offset(offset)

	// Eksekusi query
	dbQuery.Find(model)

	// Hitung total halaman
	pageCount := int(math.Ceil(float64(totalCount) / float64(query.Limit)))

	return PaginatedResult{
		Records:    model,
		TotalCount: totalCount,
		PageCount:  pageCount,
	}
}
