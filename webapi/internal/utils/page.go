package utils

func GetPageStart(page, pageSize int) int {
	if page >= 1 {
		page = page - 1
	}
	return page * pageSize
}
