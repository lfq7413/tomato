package err

import "errors"
import "strconv"

// E 组装 json 格式错误信息：
// {"code": 105,"error": "invalid field name: bl!ng"}
func E(code int, msg string) error {
	text := `{"code": ` + strconv.Itoa(code) + `,"error": "` + msg + `"}`
	return errors.New(text)
}
