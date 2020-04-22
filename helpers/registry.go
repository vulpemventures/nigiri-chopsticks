package helpers

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/sdomino/scribble"
)

// Registry handles writing/reading json file where are stored info about issued asset
type Registry struct {
	db *scribble.Driver
}

// NewRegistry returns a new Registry of error if the path is not absolute
func NewRegistry(path string) (*Registry, error) {
	r := &Registry{}
	db, err := scribble.New(path, nil)
	if err != nil {
		return nil, err
	}
	r.db = db
	return r, nil
}

// AddEntry adds an entry to the register by previously making sure that
// the incoming entry does not already exist in the registry
func (r *Registry) AddEntry(asset string, issuanceInput map[string]interface{}, contract map[string]interface{}) error {
	entry := map[string]interface{}{}
	err := r.db.Read("registry", asset, &entry)
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		return err
	}

	if len(entry) != 0 {
		return errors.New("Asset already exists on registry")
	}

	entry = map[string]interface{}{
		"asset":         asset,
		"issuance_txin": issuanceInput,
		"contract":      contract,
		"name":          contract["name"].(string),
		"ticker":        contract["ticker"].(string),
	}
	return r.db.Write("registry", asset, entry)
}

// GetEntry returns and entry if it exist in registry or NIL
func (r *Registry) GetEntry(asset string) (interface{}, error) {
	entry := map[string]interface{}{}
	err := r.db.Read("registry", asset, &entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *Registry) GetEntries(assets []interface{}) ([]interface{}, error) {
	entries := []interface{}{}

	if len(assets) == 0 {
		records, err := r.db.ReadAll("registry")
		if err != nil {
			return nil, err
		}

		for _, f := range records {
			var entry interface{}
			json.Unmarshal([]byte(f), &entry)
			entries = append(entries, entry)
		}
	} else {
		for _, asset := range assets {
			var entry interface{}
			r.db.Read("registry", asset.(string), &entry)
			entries = append(entries, entry)
		}
	}

	return entries, nil
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
