package splitdumper

//Breaker is a halting check function for Dump
//Returning true will cause Dump to stop
//Breaker should not modify the map, and it should be idempotent
//It should not store external states, and the same input should give the same output
//While it's possible to store external states, it will not be supported
//The limit is the user defined limit of times a site should be seen before stopping
//The destination is the last dumped destination
//The record is the current state of record(destination=>seen counter)
type Breaker func(limit int, destination string, record map[string]int) bool

//ConservativeBreaker breaks when all known destination meet the limit, and the current destination exceeds it
func ConservativeBreaker() Breaker {
	return func(limit int, destination string, record map[string]int) bool {
		ready := true
		for _, count := range record {
			if count < limit {
				ready = false
				break
			}
		}
		if ready && record[destination] > limit {
			return true
		}
		return false
	}
}

//EagerBreaker breaks when all known destination meet the limit
func EagerBreaker() Breaker {
	return func(limit int, destination string, record map[string]int) bool {
		halt := true
		for _, count := range record {
			if count < limit {
				halt = false
				break
			}
		}
		return halt
	}
}
