package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

const dataFile = "engine_data.json"

type Engine struct {
	data map[string]json.RawMessage
	mu   sync.RWMutex
}

func NewEngine() *Engine {
	e := &Engine{
		data: make(map[string]json.RawMessage),
	}
	e.loadData()
	return e
}

func (e *Engine) loadData() {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return // File doesn't exist or can't be read, start with empty data
	}
	json.Unmarshal(data, &e.data)
}

func (e *Engine) saveData() error {
	data, err := json.Marshal(e.data)
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func (e *Engine) Set(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.data[key] = jsonData
	return e.saveData()
}

func (e *Engine) Get(key string, value interface{}) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	jsonData, ok := e.data[key]
	if !ok {
		return errors.New("key not found")
	}
	return json.Unmarshal(jsonData, value)
}

func (e *Engine) Delete(key string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.data, key)
	return e.saveData()
}

func (e *Engine) DeleteAll() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.data = make(map[string]json.RawMessage)
	return e.saveData()
}

