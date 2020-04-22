package helpers

import (
	"encoding/json"
	"errors"
	"os"
)

// Registry handles writing/reading json file where are stored info about issued asset
type Registry struct {
	path string
}

// NewRegistry returns a new Registry of error if the path is not absolute
func NewRegistry(path string) *Registry {
	r := &Registry{path}
	return r
}

// AddEntry adds an entry to the register by previously making sure that
// the incoming entry does not already exist in the registry
func (r *Registry) AddEntry(asset string, issuanceInput map[string]interface{}, contract map[string]interface{}) error {
	registry, err := r.load()
	if err != nil {
		return err
	}

	_, ok := registry[asset]
	if ok {
		return errors.New("Asset already exists on registry")
	}
	entry := map[string]interface{}{
		"asset":         asset,
		"issuance_txin": issuanceInput,
		"contract":      contract,
		"name":          contract["name"].(string),
		"ticker":        contract["ticker"].(string),
	}
	registry[asset] = entry
	return r.save(registry)
}

func (r *Registry) load() (map[string]interface{}, error) {
	file, err := os.OpenFile("registry.json", os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	data := map[string]interface{}{}

	decoder.Decode(&data)

	return data, nil
}

func (r *Registry) save(payload map[string]interface{}) error {
	file, err := os.OpenFile("registry.json", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "   ")
	encoder.Encode(payload)
	return nil
}
