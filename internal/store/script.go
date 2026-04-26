package store

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ScriptEngine provides a minimal Lua-like scripting interface for atomic multi-step operations.
// Commands are expressed as a simple line-based DSL: SET key value, GET key, DEL key, INCR key.

type ScriptResult struct {
	Outputs []string
	Error   string
}

type ScriptEngine struct {
	store *Store
}

func NewScriptEngine(s *Store) *ScriptEngine {
	return &ScriptEngine{store: s}
}

// Exec runs a multi-line script atomically under the store's write lock.
// Each line is one command. Blank lines and lines starting with # are ignored.
func (e *ScriptEngine) Exec(script string) ScriptResult {
	e.store.mu.Lock()
	defer e.store.mu.Unlock()

	var outputs []string
	lines := strings.Split(script, "\n")

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out, err := e.execLine(line)
		if err != nil {
			return ScriptResult{Outputs: outputs, Error: fmt.Sprintf("script error on %q: %v", line, err)}
		}
		outputs = append(outputs, out)
	}
	return ScriptResult{Outputs: outputs}
}

func (e *ScriptEngine) execLine(line string) (string, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", nil
	}
	cmd := strings.ToUpper(parts[0])
	switch cmd {
	case "SET":
		if len(parts) < 3 {
			return "", errors.New("SET requires key and value")
		}
		e.store.setLocked(parts[1], parts[2], 0)
		return "OK", nil
	case "GET":
		if len(parts) < 2 {
			return "", errors.New("GET requires key")
		}
		v, ok := e.store.getLocked(parts[1])
		if !ok {
			return "(nil)", nil
		}
		return v, nil
	case "DEL":
		if len(parts) < 2 {
			return "", errors.New("DEL requires key")
		}
		e.store.deleteLocked(parts[1])
		return "OK", nil
	case "INCR":
		if len(parts) < 2 {
			return "", errors.New("INCR requires key")
		}
		v, _ := e.store.getLocked(parts[1])
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil && v != "" {
			return "", fmt.Errorf("value at %s is not an integer", parts[1])
		}
		n++
		e.store.setLocked(parts[1], strconv.FormatInt(n, 10), 0)
		return strconv.FormatInt(n, 10), nil
	default:
		return "", fmt.Errorf("unknown command: %s", cmd)
	}
}
