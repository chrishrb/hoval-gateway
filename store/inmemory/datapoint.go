package inmemory

import "github.com/chrishrb/hoval-gateway/hoval"

func (s *Store) SetDatapoint(name string, dp *hoval.Datapoint) {
	s.Lock()
	defer s.Unlock()

	s.datapoints[name] = dp
}

func (s *Store) LookupDatapointByName(name string) *hoval.Datapoint {
	s.Lock()
	defer s.Unlock()

	dp, exists := s.datapoints[name]
	if !exists {
		return nil
	}
	return dp
}

func (s *Store) LookupDatapointByIdentifier(fg hoval.FunctionGroup, fn uint8, dpID uint16) *hoval.Datapoint {
	s.Lock()
	defer s.Unlock()

	for _, dp := range s.datapoints {
		if dp.FunctionGroup == fg && dp.FunctionNumber == fn && dp.DatapointID == dpID {
			return dp
		}
	}
	return nil
}
