package dictionary

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
)
type Entry struct {
	Definition string `json:"definition"`
}



type Dictionary struct {
	filepath string
}

func New(filepath string) *Dictionary {
	return &Dictionary{filepath: filepath}
}

func (d *Dictionary) load() (map[string]Entry, error) {
	var data map[string]Entry
	file, err := os.ReadFile(d.filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(map[string]Entry), nil
		}
		return nil, err
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *Dictionary) save(data map[string]Entry) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(d.filepath, file, 0666)
}

func (d *Dictionary) Add(word string, definition string) error {
	data, err := d.load()
	if err != nil {
		return err
	}
	data[word] = Entry{Definition: definition}
	return d.save(data)
}

func (d *Dictionary) Get(word string) (Entry, bool, error) {
	entries, err := d.load()
	if err != nil {
		return Entry{}, false, err
	}
	entry, found := entries[word]
	return entry, found, nil
}

func (d *Dictionary) Remove(word string) error {
	data, err := d.load()
	if err != nil {
		return err
	}
	_, exists := data[word]
	if !exists {
		return errors.New("word not found")
	}
	delete(data, word)
	return d.save(data)
}

func (d *Dictionary) List() ([]string, error) {
	data, err := d.load()
	if err != nil {
		return nil, err
	}
	var words []string
	for word := range data {
		words = append(words, word)
	}
	sort.Strings(words)
	return words, nil
}