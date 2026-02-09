package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
)

// Create a unique ID for the helm chart source
func (h HelmChartSource) ToBase64Id() (string, error) {
	jsonDataBytes, err := json.Marshal(h)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(jsonDataBytes)
	return str, nil
}

// Convert the unique ID back to the helm chart source
func Base64StringToHelmChart(base64Data string) (HelmChartSource, error) {
	dataBytes, err := base64.StdEncoding.DecodeString(base64Data)
	fmt.Printf("%s\n", dataBytes)
	if err != nil {
		return HelmChartSource{}, err
	}
	var h HelmChartSource
	err = json.Unmarshal(dataBytes, &h)
	fmt.Printf("%v\n", h)
	if err != nil {
		return HelmChartSource{}, err
	}
	return h, nil
}

type Status int

const (
	StatusSuccess    Status = 0
	StatusNotFound   Status = 1
	StatusInProgress Status = 2
)

type StoredResult struct {
	Value  HelmChartAnalysis
	Status Status
}

type ResultStore struct {
	mu    sync.Mutex
	store map[string]StoredResult
}

func NewResultStore() *ResultStore {
	resultStore := &ResultStore{
		store: map[string]StoredResult{},
	}
	return resultStore
}

func (r *ResultStore) Get(key string) (HelmChartAnalysis, Status) {
	r.mu.Lock()
	res, ok := r.store[key]
	r.mu.Unlock()
	if !ok {
		return HelmChartAnalysis{}, StatusNotFound
	}
	if res.Status == StatusInProgress {
		return HelmChartAnalysis{}, StatusInProgress
	}

	// TODO: Use a damn library for deep copying next time...
	orig := res.Value
	newCopy := HelmChartAnalysis{
		RepoURL:  orig.RepoURL,
		ChartRef: orig.ChartRef,
	}
	newCopy.Images = make([]ImageAnalysis, len(orig.Images))
	copy(newCopy.Images, orig.Images)
	return newCopy, StatusSuccess
}

func (r *ResultStore) Put(key string, v HelmChartAnalysis) {
	r.mu.Lock()

	newCopy := HelmChartAnalysis{
		RepoURL:  v.RepoURL,
		ChartRef: v.ChartRef,
	}
	newCopy.Images = make([]ImageAnalysis, len(v.Images))
	copy(newCopy.Images, v.Images)
	r.store[key] = StoredResult{
		Value:  newCopy,
		Status: StatusSuccess,
	}
	fmt.Printf("in the put operation: %v\n", r.store[key])
	r.mu.Unlock()
}

func (r *ResultStore) SetPending(key string) {
	r.mu.Lock()
	r.store[key] = StoredResult{
		Status: StatusInProgress,
	}
	r.mu.Unlock()
}

func (r *ResultStore) UnsetPending(key string) {
	r.mu.Lock()
	delete(r.store, key)
	r.mu.Unlock()
}
