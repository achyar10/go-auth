package helper

// GenerateMetadata untuk response metadata pagination
func GenerateMetadata(query QueryParams, totalCount int64, pageCount int) map[string]interface{} {
	return map[string]interface{}{
		"page":        query.Page,
		"per_page":    query.Limit,
		"page_count":  pageCount,
		"total_count": totalCount,
	}
}
