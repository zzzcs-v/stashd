package store

import "fmt"

// PipelineOp represents a single operation in a pipeline
type PipelineOp struct {
	Type  string // "set", "get", "delete"
	Key   string
	Value string
	TTL   int // seconds, 0 = no TTL
}

// PipelineResult holds the result of a single pipeline operation
type PipelineResult struct {
	Key   string
	Value string
	OK    bool
	Error string
}

// ExecPipeline executes a slice of operations atomically (best-effort, sequential)
// and returns a result for each op in order.
func (s *Store) ExecPipeline(ops []PipelineOp) []PipelineResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	results := make([]PipelineResult, 0, len(ops))

	for _, op := range ops {
		var res PipelineResult
		res.Key = op.Key

		switch op.Type {
		case "set":
			var ttl int
			if op.TTL > 0 {
				ttl = op.TTL
			}
			s.setLocked(op.Key, op.Value, ttl)
			res.OK = true

		case "get":
			v, ok := s.getLocked(op.Key)
			if ok {
				res.Value = v
				res.OK = true
			} else {
				res.Error = "not found"
			}

		case "delete":
			s.deleteLocked(op.Key)
			res.OK = true

		default:
			res.Error = fmt.Sprintf("unknown op type: %s", op.Type)
		}

		results = append(results, res)
	}

	return results
}
