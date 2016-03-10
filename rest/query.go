package rest

import (
	"fmt"
	"sort"
	"strings"

	"github.com/lfq7413/tomato/auth"
)

// Query ...
type Query struct {
	auth        *auth.Auth
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
	auth *auth.Auth,
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

	q.getUserAndRoleACL()
	q.validateClientClassCreation()
	q.replaceSelect()
	q.replaceDontSelect()
	q.replaceInQuery()
	q.replaceNotInQuery()
	q.runFind()
	q.runCount()
	q.handleInclude()
	return q.response
}

func (q *Query) getUserAndRoleACL() error {
	return nil
}

func (q *Query) validateClientClassCreation() error {
	return nil
}

func (q *Query) replaceSelect() error {
	return nil
}

func (q *Query) replaceDontSelect() error {
	return nil
}

func (q *Query) replaceInQuery() error {
	return nil
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
