package server

import "github.com/lfq7413/tomato/livequery/t"

// SubscriptionInfo 订阅对象信息
// 每一个客户端请求对应一个对象
type SubscriptionInfo struct {
	Subscription *Subscription
	SessionToken string
	Fields       []string
}

// Subscription 订阅对象
// 一组 ClassName Hash 对应唯一一个对象
// ClientRequestIDs 记录连接到该对象的所有客户端
type Subscription struct {
	Query            t.M
	ClassName        string
	Hash             string
	ClientRequestIDs map[int][]int
}

// NewSubscription ...
func NewSubscription(className string, query t.M, queryHash string) *Subscription {
	s := &Subscription{}
	s.ClassName = className
	s.Query = query
	s.Hash = queryHash
	s.ClientRequestIDs = map[int][]int{}
	return s
}

// AddClientSubscription 添加连接到该订阅对象的客户端请求
func (s *Subscription) AddClientSubscription(clientID, requestID int) {
	requestIDs := s.ClientRequestIDs[clientID]
	if requestIDs == nil {
		requestIDs = []int{}
	}
	requestIDs = append(requestIDs, requestID)
	s.ClientRequestIDs[clientID] = requestIDs
}

// DeleteClientSubscription 删除连接到该订阅对象的客户端请求
func (s *Subscription) DeleteClientSubscription(clientID, requestID int) {
	requestIDs := s.ClientRequestIDs[clientID]
	if requestIDs == nil {
		return
	}
	index := -1
	for i, id := range requestIDs {
		if id == requestID {
			index = i
			break
		}
	}
	if index < 0 {
		return
	}
	requestIDs[index] = requestIDs[len(requestIDs)-1]
	requestIDs = requestIDs[:len(requestIDs)-1]
	s.ClientRequestIDs[clientID] = requestIDs
	if len(requestIDs) == 0 {
		delete(s.ClientRequestIDs, clientID)
	}
}

// HasSubscribingClient 返回连接到该对象的客户端数量
func (s *Subscription) HasSubscribingClient() bool {
	return len(s.ClientRequestIDs) > 0
}
