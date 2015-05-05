package dashboard

import (
	"encoding/json"
	"net/http"
)

// HandleHistory returns an object in JSON with the following members.
//
//     - "Clock": {number} is the current clock value (starting from 0).
//
//     - "TotalHits": {number} is the total number of hits in the past 24
//     hours.
//
//     - "Directions": {array} is the collection of directions.  Each item is
//     an object with the following members:
//
//         - "blocked-values": {array} is the collection of blocked values.
//
//         - "name": {string} is the name of the direction.
func (s *Server) HandleHistory(w http.ResponseWriter, r *http.Request) {
	type HistoryData struct {
		Clock      int32
		TotalHits  uint64
		Directions []interface{}
	}

	data := HistoryData{
		Clock:      s.counter.Clock.GetTime(),
		TotalHits:  s.counter.Count.Count(),
		Directions: make([]interface{}, 0, len(s.conf.Directions)),
	}

	for iDirection := range s.conf.Directions {
		direction := &s.conf.Directions[iDirection]
		direction.Store.CleanUp(data.Clock)

		dirData := map[string]interface{}{
			"name":           direction.Name,
			"blocked-values": direction.BlockedValues(),
		}

		data.Directions = append(data.Directions, dirData)
	}

	json, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
