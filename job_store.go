package main

import "sync"

// JobStore is a store for jobs
type JobStore struct {
	sync.RWMutex
	m map[string]Job
}

// NewJobStore create a new job store
func NewJobStore() *JobStore {
	return &JobStore{m: make(map[string]Job)}
}

// Get a job
func (js *JobStore) Get(k string) Job {
	js.RLock()
	j := js.m[k]
	js.RUnlock()
	return j
}

// Save a new job
func (js *JobStore) Save(j Job) {
	js.Lock()
	js.m[j.ID] = j
	js.Unlock()
}
