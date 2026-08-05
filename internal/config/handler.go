package config

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// handleCors handles the OPTIONS preflight and sets common CORS headers
func handleCors(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

// SaveProviderHandler handles POST /api/config/provider/{provider}
func SaveProviderHandler(w http.ResponseWriter, r *http.Request) {
	if handleCors(w, r) {
		return
	}
	
	vars := mux.Vars(r)
	providerName := vars["provider"]

	var newConfig ProviderConfig
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Optionally ensure the provider name in the path matches the struct
	newConfig.Provider = providerName

	if err := SaveProvider(providerName, newConfig); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Configuration saved"})
}

// LoadProviderHandler handles GET /api/config/provider/{provider}
func LoadProviderHandler(w http.ResponseWriter, r *http.Request) {
	if handleCors(w, r) {
		return
	}
	
	vars := mux.Vars(r)
	providerName := vars["provider"]

	config, err := LoadProvider(providerName)
	if err != nil {
		// Differentiate between not found and other errors
		if err.Error() == "provider configuration not found: "+providerName {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// DeleteProviderHandler handles DELETE /api/config/provider/{provider}
func DeleteProviderHandler(w http.ResponseWriter, r *http.Request) {
	if handleCors(w, r) {
		return
	}
	
	vars := mux.Vars(r)
	providerName := vars["provider"]

	if err := DeleteProvider(providerName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Configuration deleted"})
}

// ListProvidersHandler handles GET /api/config/providers
func ListProvidersHandler(w http.ResponseWriter, r *http.Request) {
	if handleCors(w, r) {
		return
	}
	
	providers, err := ListConfiguredProviders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"providers": providers,
		"count":     len(providers),
	})
}
