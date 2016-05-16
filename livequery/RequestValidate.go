package livequery

import "errors"

func validate(data M, op string) error {
	switch op {
	case "general":
		return validateGeneral(data)
	case "connect":
		return validateConnect(data)
	case "subscribe":
		return validateSubscribe(data)
	case "unsubscribe":
		return validateUnsubscribe(data)
	default:
		return errors.New("invalid op")
	}
}

func validateGeneral(data M) error {
	if v, ok := data["op"]; ok {
		if op, ok := v.(string); ok {
			if op != "connect" && op != "subscribe" && op != "unsubscribe" {
				return errors.New("op is not in [connect, subscribe, unsubscribe]")
			}
			return nil
		}
		return errors.New("op is not string")
	}
	return errors.New("need op")
}

func validateConnect(data M) error {
	if v, ok := data["applicationId"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("applicationId is not string")
		}
	} else {
		return errors.New("need applicationId")
	}

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

	if v, ok := data["sessionToken"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("sessionToken is not string")
		}
	}

	return nil
}

func validateSubscribe(data M) error {
	if v, ok := data["requestId"]; ok {
		if _, ok := v.(float64); ok == false {
			return errors.New("requestId is not number")
		}
	} else {
		return errors.New("need requestId")
	}

	if v, ok := data["query"]; ok {
		if query, ok := v.(map[string]interface{}); ok {
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

	if v, ok := data["sessionToken"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("sessionToken is not string")
		}
	}

	return nil
}

func validateUnsubscribe(data M) error {
	if v, ok := data["requestId"]; ok {
		if _, ok := v.(float64); ok == false {
			return errors.New("requestId is not number")
		}
	} else {
		return errors.New("need requestId")
	}

	return nil
}

func validateQuery(data M) error {
	if v, ok := data["className"]; ok {
		if _, ok := v.(string); ok == false {
			return errors.New("className is not string")
		}
	} else {
		return errors.New("need className")
	}

	if v, ok := data["where"]; ok {
		if _, ok := v.(map[string]interface{}); ok == false {
			return errors.New("where is not object")
		}
	} else {
		return errors.New("need where")
	}

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
