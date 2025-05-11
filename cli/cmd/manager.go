package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Context struct {
	Host string `json:"host"`
}

type State struct {
	Current  string             `json:"current"`
	Contexts map[string]Context `json:"contexts"`
}

var stateFile = filepath.Join(os.Getenv("HOME"), ".daxagrid", "contexts.json")

func SaveContext(name string, host string) error {
	os.MkdirAll(filepath.Dir(stateFile), 0755)

	var state State
	_ = load(&state)

	if state.Contexts == nil {
		state.Contexts = make(map[string]Context)
	}

	state.Contexts[name] = Context{Host: host}
	state.Current = name

	return write(state)
}

func GetCurrentHost() (string, error) {
	var state State
	if err := load(&state); err != nil {
		return "", err
	}

	ctx, ok := state.Contexts[state.Current]
	if !ok {
		return "", errors.New("no active kontext set")
	}

	return ctx.Host, nil
}

func load(state *State) error {
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return nil // empty
	}
	return json.Unmarshal(data, state)
}

func write(state State) error {
	data, _ := json.MarshalIndent(state, "", "  ")
	return os.WriteFile(stateFile, data, 0644)
}
