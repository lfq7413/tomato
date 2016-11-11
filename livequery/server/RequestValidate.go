package server

import (
	"errors"

	"github.com/lfq7413/tomato/livequery/t"
)

// Validate 校验客户端发送的请求是否合法
func Validate(data t.M, op string) error {
	switch op {
	case "general":
		return validateGeneral(data)
	case "connect":
		return validateConnect(data)
	case "subscribe":
		return validateSubscribe(data)
	case "unsubscribe":
		return validateUnsubscribe(data)
	case "update":
		return validateUpdate(data)
	default:
		return errors.New("invalid op")
	}
}

// validateGeneral 校验 op 操作符是否支持
func validateGeneral(data t.M) error {
	if v, ok := data["op"]; ok {
		if op, ok := v.(string); ok {
			if op != "connect" && op != "subscribe" && op != "unsubscribe" && op != "update" {
				return errors.New("op is not in [connect, subscribe, unsubscribe]")
			}
			return nil
		}
		return errors.New("op is not string")
	}
	return errors.New("need op")
}

// validateConnect 校验 connect 请求格式
func validateConnect(data t.M) error {
	if v, ok := data["applicationId"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("applicationId is not string")
		}
	}
	// else {
	// 	return errors.New("need applicationId")
	// }

	if v, ok := data["masterKey"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("masterKey is not string")
		}
	}

	if v, ok := data["clientKey"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("clientKey is not string")
		}
	}

	if v, ok := data["restAPIKey"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("restAPIKey is not string")
		}
	}

	if v, ok := data["javascriptKey"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("javascriptKey is not string")
		}
	}

	if v, ok := data["windowsKey"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("windowsKey is not string")
		}
	}

	if v, ok := data["sessionToken"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("sessionToken is not string")
		}
	}

	return nil
}

// validateSubscribe 校验 subscribe 请求格式
func validateSubscribe(data t.M) error {
	// 必须包含 requestId 字段
	if v, ok := data["requestId"]; ok {
		if _, ok := v.(float64); ok == false {
			return errors.New("requestId is not number")
		}
	} else {
		return errors.New("need requestId")
	}

	// 必须包含 query 字段
	if v, ok := data["query"]; ok {
		if query, ok := v.(map[string]interface{}); ok {
			// 校验 query 字段
			err := validateQuery(query)
			if err != nil {
				return err
			}
		} else {
			return errors.New("query is not object")
		}
	} else {
		return errors.New("need query")
	}

	// sessionToken 字段可选
	if v, ok := data["sessionToken"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("sessionToken is not string")
		}
	}

	return nil
}

// validateUnsubscribe 校验 unsubscribe 请求格式
func validateUnsubscribe(data t.M) error {
	if v, ok := data["requestId"]; ok {
		if _, ok := v.(float64); ok == false {
			return errors.New("requestId is not number")
		}
	} else {
		return errors.New("need requestId")
	}

	return nil
}

// validateUpdate 校验 update 请求格式，可复用 validateSubscribe
func validateUpdate(data t.M) error {
	return validateSubscribe(data)
}

// validateQuery 校验 query 字段
func validateQuery(data t.M) error {
	// 必须包含 className 字段
	if v, ok := data["className"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("className is not string")
		}
	} else {
		return errors.New("need className")
	}

	// 必须包含 where 字段
	if v, ok := data["where"]; ok {
		if _, ok := v.(map[string]interface{}); ok == false {
			return errors.New("where is not object")
		}
	} else {
		return errors.New("need where")
	}

	// fields 字段可选
	if v, ok := data["fields"]; ok {
		if fields, ok := v.([]string); ok {
			if len(fields) < 1 {
				return errors.New("minItems is not 1")
			}
		} else {
			return errors.New("fields is not []string")
		}
	}

	return nil
}
