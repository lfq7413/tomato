package livequery

type subscription struct {
	query            M
	className        string
	hash             string
	clientRequestIDs map[int][]int
}

func newSubscription(className string, query M, queryHash string) *subscription {
	s := &subscription{}
	s.className = className
	s.query = query
	s.hash = queryHash
	s.clientRequestIDs = map[int][]int{}
	return s
}

func (s *subscription) addClientSubscription(clientID, requestID int) {
	requestIDs := s.clientRequestIDs[clientID]
	if requestIDs == nil {
		requestIDs = []int{}
	}
	requestIDs = append(requestIDs, requestID)
	s.clientRequestIDs[clientID] = requestIDs
}

func (s *subscription) deleteClientSubscription(clientID, requestID int) {
	requestIDs := s.clientRequestIDs[clientID]
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
	s.clientRequestIDs[clientID] = requestIDs
	if len(requestIDs) == 0 {
		delete(s.clientRequestIDs, clientID)
	}
}

func (s *subscription) hasSubscribingClient() bool {
	return len(s.clientRequestIDs) > 0
}
