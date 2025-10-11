package udict

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/93aakash/urban_dict/internal/udict/database"
	"github.com/93aakash/urban_dict/internal/udict/model"
)

const baseURL = "https://api.urbandictionary.com/v0/define?term="

var (
	ErrDefNotFound = errors.New("no definition found")
)

type DictResponse struct {
	List []model.Definition `json:"List"`
}

type Dict struct {
	client *http.Client
	db     *database.DB
}

func New() (*Dict, error) {
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: &http.Transport{},
	}

	dbPath, err := GetDatabasePath()
	if err != nil {
		return nil, fmt.Errorf("failed to configure %s: %w", AppName, err)
	}

	db, err := database.NewDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create dict store: %w", err)
	}

	if err := db.InitDB(); err != nil {
		return nil, fmt.Errorf("failed to initialize dict: %w", err)
	}

	return &Dict{client: client, db: db}, nil
}

func (d *Dict) Get(query string) ([]model.Definition, error) {
	definitions, err := d.db.GetDefinitions(query)
	if err == nil && len(definitions) > 0 {
		return definitions, nil
	}

	definitions, err = d.fetchDefinitions(query)
	if err != nil {
		slog.Debug("failed to fetch definition from API", "error", err)
		return nil, err
	}
	if len(definitions) == 0 {
		return nil, ErrDefNotFound
	}

	if err := d.db.InsertDefinitions(definitions); err != nil {
		slog.Debug("failed to add definition to dict", "error", err)
		return nil, fmt.Errorf("failed to insert definition into dict: %w", err)
	}

	return definitions, nil
}

func (d *Dict) Delete(query string) error {
	if err := d.db.DeleteDefinition(query); err != nil {
		slog.Debug("failed to delete definition from dict", "error", err)
		return fmt.Errorf("failed to add definition to dict: %w", err)
	}

	return nil
}

// fetchDefinitions retrieves the definition of a word from the Urban Dictionary API
func (d *Dict) fetchDefinitions(query string) ([]model.Definition, error) {
	// Encode the query to make sure it's URL-safe
	encodedTerm := url.QueryEscape(query)

	// Build the API URL
	apiURL := fmt.Sprintf("%s%s", baseURL, encodedTerm)

	resp, err := d.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Urban Dictionary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from Urban Dictionary", resp.StatusCode)
	}

	var dictResp DictResponse
	if err := json.NewDecoder(resp.Body).Decode(&dictResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return dictResp.List, nil
}
