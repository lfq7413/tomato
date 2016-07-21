package client

import (
	"regexp"
	"strconv"
	"strings"
)

// FromString 从字符串中解析 SDK 名称及版本信息
func FromString(version string) map[string]string {
	versionRE := `([-a-zA-Z]+)([0-9\.]+)`
	re := regexp.MustCompile(versionRE)
	match := re.FindStringSubmatch(strings.ToLower(version))
	if match != nil && len(match) == 3 {
		return map[string]string{
			"sdk":     match[1],
			"version": match[2],
		}
	}
	return map[string]string{
		"sdk":     "",
		"version": "",
	}
}

// SupportsForwardDelete 是否支持字段删除
func SupportsForwardDelete(clientSDK map[string]string) bool {
	compatibleSDK := map[string]string{
		"js": ">=1.9.0",
	}
	return compatible(compatibleSDK, clientSDK)
}

// compatible 检测 SDK 兼容性
// compatibleSDK 兼容的 SDK 版本
// clientSDK 客户端 SDK 版本
func compatible(compatibleSDK, clientSDK map[string]string) bool {
	// REST API, or custom SDK
	if clientSDK == nil || len(clientSDK) == 0 || clientSDK["sdk"] == "" && clientSDK["version"] == "" {
		return true
	}

	clientVersion := clientSDK["version"]
	compatiblityVersion := compatibleSDK[clientSDK["sdk"]]

	// 客户端版本类型不在兼容性检测列表中
	if compatiblityVersion == "" {
		return true
	}

	return satisfies(clientVersion, compatiblityVersion)
}

// satisfies 给定的 clientVersion 是否满足 compatiblityVersion 条件
func satisfies(clientVersion, compatiblityVersion string) bool {
	// 支持 6 种比较操作符，">=" 先于 ">"
	ops := []string{">=", "<=", ">", "<", "=", ""}
	var option string
	for _, op := range ops {
		if strings.HasPrefix(compatiblityVersion, op) {
			option = op
			break
		}
	}

	compatiblityVersion = compatiblityVersion[len(option):]
	clientVersionList := strings.Split(clientVersion, ".")
	compatiblityVersionList := strings.Split(compatiblityVersion, ".")

	// 使用 0 补齐长度：(1.2   3.5.1) >> (1.2.0   3.5.1)
	length := len(clientVersionList)
	if len(compatiblityVersionList) > length {
		length = len(compatiblityVersionList)
	}
	if len(clientVersionList) < length {
		for i := 0; i < length-len(clientVersionList); i++ {
			clientVersionList = append(clientVersionList, "0")
		}
	}
	if len(compatiblityVersionList) < length {
		for i := 0; i < length-len(compatiblityVersionList); i++ {
			compatiblityVersionList = append(compatiblityVersionList, "0")
		}
	}

	// 比较版本号
	switch option {
	case ">=":
		return greater(clientVersionList, compatiblityVersionList, length) || equal(clientVersionList, compatiblityVersionList, length)
	case "<=":
		return less(clientVersionList, compatiblityVersionList, length) || equal(clientVersionList, compatiblityVersionList, length)
	case ">":
		return greater(clientVersionList, compatiblityVersionList, length)
	case "<":
		return less(clientVersionList, compatiblityVersionList, length)
	case "=":
		return equal(clientVersionList, compatiblityVersionList, length)
	case "":
		return equal(clientVersionList, compatiblityVersionList, length)
	}
	return false
}

func greater(clientVersionList, compatiblityVersionList []string, length int) bool {
	for i := 0; i < length; i++ {
		if toInt(clientVersionList[i]) > toInt(compatiblityVersionList[i]) {
			return true
		}
		if toInt(clientVersionList[i]) < toInt(compatiblityVersionList[i]) {
			return false
		}
	}
	return false
}

func less(clientVersionList, compatiblityVersionList []string, length int) bool {
	for i := 0; i < length; i++ {
		if toInt(clientVersionList[i]) < toInt(compatiblityVersionList[i]) {
			return true
		}
		if toInt(clientVersionList[i]) > toInt(compatiblityVersionList[i]) {
			return false
		}
	}
	return false
}

func equal(clientVersionList, compatiblityVersionList []string, length int) bool {
	for i := 0; i < length; i++ {
		if toInt(clientVersionList[i]) != toInt(compatiblityVersionList[i]) {
			return false
		}
	}
	return true
}

func toInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return 0
	}
	return int(i)
}
