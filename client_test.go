package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// helper function to create a mock server for testing
func newMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestMCPClient_GetConfig(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		expectedConfig := &MCPConfig{
			ID:         "test-id",
			Name:       "Test Config",
			Version:    1,
			Enabled:    true,
			Parameters: map[string]string{"key": "value"},
		}

		server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("Expected GET request, got %s", r.Method)
			}
			if r.URL.Path != "/configs/test-id" {
				t.Errorf("Expected path /configs/test-id, got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(expectedConfig)
		})
		defer server.Close()

		client := NewMCPClient(server.URL)
		config, err := client.GetConfig("test-id")

		if err != nil {
			t.Fatalf("GetConfig failed: %v", err)
		}
		if !reflect.DeepEqual(config, expectedConfig) {
			t.Errorf("Expected config %+v, got %+v", expectedConfig, config)
		}
	})

	t.Run("server error", func(t *testing.T) {
		server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer server.Close()

		client := NewMCPClient(server.URL)
		_, err := client.GetConfig("test-id")

		if err == nil {
			t.Fatal("Expected an error, but got nil")
		}
		// Add more specific error checking if needed
	})
}

func TestMCPClient_CreateConfig(t *testing.T) {
	configToCreate := &MCPConfig{
		Name:       "New Config",
		Version:    1,
		Enabled:    true,
		Parameters: map[string]string{"p1": "v1"},
	}
	expectedCreatedConfig := &MCPConfig{
		ID:         "new-id-from-server", // Server might assign the ID
		Name:       "New Config",
		Version:    1,
		Enabled:    true,
		Parameters: map[string]string{"p1": "v1"},
	}

	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/configs" {
			t.Errorf("Expected path /configs, got %s", r.URL.Path)
		}

		var receivedConfig MCPConfig
		if err := json.NewDecoder(r.Body).Decode(&receivedConfig); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		// For simplicity, not comparing all fields of receivedConfig with configToCreate here,
		// but in a real test, you might want to.

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expectedCreatedConfig)
	})
	defer server.Close()

	client := NewMCPClient(server.URL)
	createdConfig, err := client.CreateConfig(configToCreate)

	if err != nil {
		t.Fatalf("CreateConfig failed: %v", err)
	}
	if !reflect.DeepEqual(createdConfig, expectedCreatedConfig) {
		t.Errorf("Expected created config %+v, got %+v", expectedCreatedConfig, createdConfig)
	}
}

func TestMCPClient_UpdateConfig(t *testing.T) {
	configToUpdate := &MCPConfig{
		Name:    "Updated Config",
		Version: 2,
		Enabled: false,
	}
	expectedUpdatedConfig := &MCPConfig{
		ID:      "update-id",
		Name:    "Updated Config",
		Version: 2,
		Enabled: false,
	}

	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/configs/update-id" {
			t.Errorf("Expected path /configs/update-id, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedUpdatedConfig)
	})
	defer server.Close()

	client := NewMCPClient(server.URL)
	updatedConfig, err := client.UpdateConfig("update-id", configToUpdate)

	if err != nil {
		t.Fatalf("UpdateConfig failed: %v", err)
	}
	if !reflect.DeepEqual(updatedConfig, expectedUpdatedConfig) {
		t.Errorf("Expected updated config %+v, got %+v", expectedUpdatedConfig, updatedConfig)
	}
}

func TestMCPClient_DeleteConfig(t *testing.T) {
	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/configs/delete-id" {
			t.Errorf("Expected path /configs/delete-id, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent) // Or http.StatusOK
	})
	defer server.Close()

	client := NewMCPClient(server.URL)
	err := client.DeleteConfig("delete-id")

	if err != nil {
		t.Fatalf("DeleteConfig failed: %v", err)
	}
}

func TestMCPClient_ListConfigs(t *testing.T) {
	expectedConfigs := []MCPConfig{
		{ID: "id1", Name: "Config 1"},
		{ID: "id2", Name: "Config 2"},
	}

	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/configs" {
			t.Errorf("Expected path /configs, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedConfigs)
	})
	defer server.Close()

	client := NewMCPClient(server.URL)
	configs, err := client.ListConfigs()

	if err != nil {
		t.Fatalf("ListConfigs failed: %v", err)
	}
	if !reflect.DeepEqual(configs, expectedConfigs) {
		t.Errorf("Expected configs %+v, got %+v", expectedConfigs, configs)
	}
}

// Add more tests for:
// - Different server responses (e.g., 400, 404 errors)
// - Malformed JSON responses from server
// - Edge cases in client logic (e.g., empty config ID)
// - Test for makeRequest helper if it grows more complex
// - Test for NewMCPClient initialization
