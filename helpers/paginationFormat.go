package helpers

import "event-management/structs"

func PaginationFormat(page int, limit int, totalRows int64, totalPages int, data interface{}) structs.Pagination {
	return structs.Pagination{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Data:       data,
	}
}
