package rest

import (
	"fmt"
	"sort"
	"strings"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/utils"
)

// Query ...
type Query struct {
	auth        *Auth
	className   string
	where       map[string]interface{}
	findOptions map[string]interface{}
	response    map[string]interface{}
	doCount     bool
	include     [][]string
	keys        []string
}

// NewQuery ...
func NewQuery(
	auth *Auth,
	className string,
	where map[string]interface{},
	options map[string]interface{},
) *Query {
	query := &Query{
		auth:        auth,
		className:   className,
		where:       where,
		findOptions: map[string]interface{}{},
		response:    nil,
		doCount:     false,
		include:     [][]string{},
		keys:        []string{},
	}

	if auth.IsMaster == false {
		if auth.User != nil {
			query.findOptions["acl"] = []string{auth.User.ID}
		} else {
			query.findOptions["acl"] = nil
		}
		if className == "_Session" {
			if query.findOptions["acl"] == nil {
				// TODO session 无效
			}
			user := map[string]string{"__type": "Pointer", "className": "_User", "objectId": auth.User.ID}
			and := []interface{}{where, user}
			query.where = map[string]interface{}{"$and": and}
		}
	}

	for k, v := range options {
		switch k {
		case "keys":
			if s, ok := v.(string); ok {
				query.keys = strings.Split(s, ",")
				query.keys = append(query.keys, "objectId", "createdAt", "updatedAt")
			}
		case "count":
			query.doCount = true
		case "skip":
			query.findOptions["skip"] = v
		case "limit":
			query.findOptions["limit"] = v
		case "order":
			if s, ok := v.(string); ok {
				fields := strings.Split(s, ",")
				sortMap := map[string]int{}
				for _, v := range fields {
					if strings.HasPrefix(v, "-") {
						sortMap[v[1:]] = -1
					} else {
						sortMap[v] = 1
					}
				}
				query.findOptions["sort"] = sortMap
			}
		case "include":
			if s, ok := v.(string); ok { // v = "user.session,name.friend"
				paths := strings.Split(s, ",") // paths = ["user.session","name.friend"]
				pathSet := []string{}
				for _, path := range paths {
					parts := strings.Split(path, ".") // parts = ["user","session"]
					for lenght := 1; lenght <= len(parts); lenght++ {
						pathSet = append(pathSet, strings.Join(parts[0:lenght], "."))
					} // pathSet = ["user","user.session"]
				} // pathSet = ["user","user.session","name","name.friend"]
				sort.Strings(pathSet) // pathSet = ["name","name.friend","user","user.session"]
				for _, set := range pathSet {
					query.include = append(query.include, strings.Split(set, "."))
				} // query.include = [["name"],["name","friend"],["user"],["user","seeeion"]]
			}
		default:
		}
	}

	return query
}

// Execute ...
func (q *Query) Execute() map[string]interface{} {

	fmt.Println("keys       ", q.keys)
	fmt.Println("doCount    ", q.doCount)
	fmt.Println("findOptions", q.findOptions)
	fmt.Println("include    ", q.include)

	q.buildRestWhere()
	q.runFind()
	q.runCount()
	q.handleInclude()
	return q.response
}

func (q *Query) buildRestWhere() error {
	q.getUserAndRoleACL()
	q.validateClientClassCreation()
	q.replaceSelect()
	q.replaceDontSelect()
	q.replaceInQuery()
	q.replaceNotInQuery()
	return nil
}

func (q *Query) getUserAndRoleACL() error {
	if q.auth.IsMaster || q.auth.User == nil {
		return nil
	}
	roles := q.auth.GetUserRoles()
	roles = append(roles, q.auth.User.ID)
	q.findOptions["acl"] = roles
	return nil
}

func (q *Query) validateClientClassCreation() error {
	sysClass := []string{"_User", "_Installation", "_Role", "_Session", "_Product"}
	if config.TConfig.AllowClientClassCreation {
		return nil
	}
	if q.auth.IsMaster {
		return nil
	}
	for _, v := range sysClass {
		if v == q.className {
			return nil
		}
	}
	if orm.CollectionExists(q.className) {
		return nil
	}
	// TODO 无法操作不存在的表
	return nil
}

func (q *Query) replaceSelect() error {
	return nil
}

func (q *Query) replaceDontSelect() error {
	return nil
}

func (q *Query) replaceInQuery() error {
	inQueryObject := findObjectWithKey(q.where, "$inQuery")
	if inQueryObject == nil {
		return nil
	}
	// 必须包含两个 key ： where className
	inQueryValue := utils.MapInterface(inQueryObject["$inQuery"])
	if inQueryValue == nil ||
		inQueryValue["where"] == nil ||
		inQueryValue["className"] == nil {
		// TODO $inQuery 用法不正确
		return nil
	}

	values := []interface{}{}
	response := Find(
		q.auth,
		utils.String(inQueryValue["className"]),
		utils.MapInterface(inQueryValue["where"]),
		map[string]interface{}{})
	// 组装查询到的对象
	if response != nil &&
		response["results"] != nil &&
		utils.SliceInterface(response["results"]) != nil &&
		len(utils.SliceInterface(response["results"])) > 0 {
		for _, v := range utils.SliceInterface(response["results"]) {
			result := utils.MapInterface(v)
			pointer := map[string]interface{}{
				"__type":    "Pointer",
				"className": inQueryValue["className"],
				"objectId":  result["objectId"],
			}
			values = append(values, pointer)
		}
	}
	delete(inQueryObject, "$inQuery")
	if inQueryObject["$in"] != nil &&
		utils.SliceInterface(inQueryObject["$in"]) != nil {
		in := utils.SliceInterface(inQueryObject["$in"])
		in = utils.AppendInterface(in, values)
	} else {
		inQueryObject["$in"] = values
	}

	return q.replaceInQuery()
}

func (q *Query) replaceNotInQuery() error {
	return nil
}

func (q *Query) runFind() error {
	return nil
}

func (q *Query) runCount() error {
	return nil
}

func (q *Query) handleInclude() error {
	return nil
}

// 查找带有指定 key 的对象
func findObjectWithKey(root interface{}, key string) map[string]interface{} {
	if s := utils.SliceInterface(root); s != nil {
		for _, v := range s {
			answer := findObjectWithKey(v, key)
			if answer != nil {
				return answer
			}
		}
	}
	if m := utils.MapInterface(root); m != nil {
		if m[key] != nil {
			return m
		}
		for _, v := range m {
			answer := findObjectWithKey(v, key)
			if answer != nil {
				return answer
			}
		}
	}
	return nil
}
