package utils

import "strconv"

// GraphQLレスポンスでIDを扱う際に用いる関数
// GOのIDのデフォルトはuint型だが、GraphQLのデフォルトはstringである。
// そのため、IDはGOの内部ではuint型で扱い、クライアントへのレスポンスはstring型で扱う
func ConvertIdToString(id uint) string{
	return strconv.FormatUint(uint64(id), 10)
}

// GO内部でIDを扱う際に用いる関数
// GOのIDのデフォルトはuint型だが、GraphQLのデフォルトはstringである。
// そのため、IDはGOの内部ではuint型で扱い、クライアントへのレスポンスはstring型で扱う
func ConvertIdToUint(id string)(uint, error) {
	// 64ビットの整数として解析
	u64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(u64), nil
}