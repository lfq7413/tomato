package server

import "github.com/lfq7413/tomato/livequery/t"

// SubscriptionInfo ...
type SubscriptionInfo struct {
	Subscription *Subscription
	SessionToken string
	Fields       []string
}

// Subscription ...
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

// AddClientSubscription ...
func (s *Subscription) AddClientSubscription(clientID, requestID int) {
	requestIDs := s.ClientRequestIDs[clientID]
	if requestIDs == nil {
		requestIDs = []int{}
	}
	requestIDs = append(requestIDs, requestID)
	s.ClientRequestIDs[clientID] = requestIDs
}

// DeleteClientSubscription ...
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

// HasSubscribingClient ...
func (s *Subscription) HasSubscribingClient() bool {
	return len(s.ClientRequestIDs) > 0
}
