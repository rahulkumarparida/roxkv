package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	master "github.com/rahulkumarparida/roxkv/agents/Master"
	"github.com/rahulkumarparida/roxkv/agents/abstractor"
	"github.com/rahulkumarparida/roxkv/internal/config"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

var StoreInstance *store.MemoryAlloc
var AgentInstance abstractor.Provider

func init() {
	abstractor.SetOnProviderChangeListener(func(newProvider, newModel string) {
		activeCfg := abstractor.GetConfig()
		if provider, err := abstractor.NewProvider(*activeCfg); err == nil {
			AgentInstance = provider
		}
	})
}

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func handlePreflight(w http.ResponseWriter, r *http.Request) bool {
	setCorsHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func WebChatServer(stre *store.MemoryAlloc, provider abstractor.Provider) error {
	StoreInstance = stre
	AgentInstance = provider

	router := mux.NewRouter()
	router.HandleFunc("/healthz", HealthHandler).Methods("GET")
	router.HandleFunc("/api/chat/", callAI).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/chat/config", aiConfigHandler).Methods("GET", "POST", "OPTIONS")

	// LLM Manager Endpoints
	router.HandleFunc("/providers", listProvidersHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/providers", listProvidersHandler).Methods("GET", "OPTIONS")

	router.HandleFunc("/providers/select", selectProviderHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/providers/select", selectProviderHandler).Methods("POST", "OPTIONS")

	router.HandleFunc("/providers/{provider}/models", getProviderModelsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/providers/{provider}/models", getProviderModelsHandler).Methods("GET", "OPTIONS")

	router.HandleFunc("/providers/{provider}/config", getProviderConfigHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/providers/{provider}/config", getProviderConfigHandler).Methods("GET", "OPTIONS")

	router.HandleFunc("/providers/{provider}/config", saveProviderConfigHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/providers/{provider}/config", saveProviderConfigHandler).Methods("POST", "OPTIONS")

	router.HandleFunc("/providers/{provider}/stats", getProviderStatsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/providers/{provider}/stats", getProviderStatsHandler).Methods("GET", "OPTIONS")

	listener, err := net.Listen("tcp", ChatHTTPAddr)
	if err != nil {
		return fmt.Errorf("start AI chat HTTP server on %s: %w", ChatHTTPAddr, err)
	}
	fmt.Println("Listening WebChat API at localhost" + ChatHTTPAddr)
	startupLog("server ports", "AI HTTP "+ChatHTTPAddr)

	httpServer := &http.Server{Handler: router}
	go func() {
		if serveErr := httpServer.Serve(listener); serveErr != nil && serveErr != http.ErrServerClosed {
			log.Printf("AI chat HTTP server stopped: %v", serveErr)
		}
	}()

	return nil
}

var payload struct {
	Query string `json:"query"`
}

func callAI(w http.ResponseWriter, r *http.Request) {
	if handlePreflight(w, r) {
		return
	}

	activeCfg := abstractor.GetConfig()
	if activeCfg == nil || strings.TrimSpace(activeCfg.Provider) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Provider not selected",
			"message": "No active LLM provider is currently selected.",
		})
		return
	}

	// Request Validation: Check API Key if provider requires it
	loadedCfg, _ := config.LoadProvider(activeCfg.Provider)
	apiKey := activeCfg.APIKey
	if apiKey == "" && loadedCfg != nil {
		apiKey = loadedCfg.APIKey
	}

	if abstractor.RequiresAPIKey(activeCfg.Provider) && strings.TrimSpace(apiKey) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "API key required",
			"message": "The selected provider requires an API key before requests can be executed.",
		})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println("Error decoding JSON payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	tcpConn, rw, err := hijacker.Hijack()
	if err != nil || tcpConn == nil {
		fmt.Println("Not able to hijack the connection:", err)
		return
	}
	defer tcpConn.Close()

	// Increment query count for active provider
	abstractor.IncrementQueryCount(activeCfg.Provider)

	client := utils.CreateClient(tcpConn, "admin")
	utils.TotalConnecntions = append(utils.TotalConnecntions, client)

	query := payload.Query
	data := master.MasterAgent(query, StoreInstance, client, AgentInstance)

	rw.WriteString("HTTP/1.1 200 OK\r\n")
	rw.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	rw.WriteString("Cache-Control: no-cache\r\n")
	rw.WriteString("Access-Control-Allow-Origin: *\r\n")
	rw.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(fmt.Sprintf("%v", data))))
	rw.WriteString("\r\n")
	rw.Flush()

	rw.WriteString(fmt.Sprintf("%v", data))
	rw.Flush()
}

func aiConfigHandler(w http.ResponseWriter, r *http.Request) {
	if handlePreflight(w, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		cfg := abstractor.GetConfig()
		json.NewEncoder(w).Encode(cfg)
		return
	}

	if r.Method == http.MethodPost {
		var newConfig abstractor.ProviderConfig
		if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := abstractor.SaveConfig(newConfig); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Also save to per-provider file
		_ = config.SaveProvider(newConfig.Provider, abstractor.MapFromAbstractorConfig(newConfig))

		provider, err := abstractor.NewProvider(newConfig)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		AgentInstance = provider

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}
}

// GET /providers
func listProvidersHandler(w http.ResponseWriter, r *http.Request) {
	if handlePreflight(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	providers := abstractor.GetAvailableProviders()
	json.NewEncoder(w).Encode(map[string]any{
		"providers": providers,
	})
}

// GET /providers/{provider}/models
func getProviderModelsHandler(w http.ResponseWriter, r *http.Request) {
	if handlePreflight(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	providerName := vars["provider"]

	models, exists := abstractor.GetSupportedModels(providerName)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Provider not found",
			"message": fmt.Sprintf("Provider %s is not registered", providerName),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"provider": providerName,
		"models":   models,
	})
}

// GET /providers/{provider}/config
func getProviderConfigHandler(w http.ResponseWriter, r *http.Request) {
	if handlePreflight(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	providerName := vars["provider"]

	cfg, err := config.LoadProvider(providerName)
	if err != nil {
		// Default empty fields if missing
		def := config.ProviderConfig{
			Provider:    providerName,
			Model:       "",
			Endpoint:    "",
			APIKey:      "",
			Temperature: 0.2,
			TopP:        0.9,
			MaxTokens:   400,
			Stream:      false,
			Timeout:     30,
		}
		models, ok := abstractor.GetSupportedModels(providerName)
		if ok && len(models) > 0 {
			def.Model = models[0]
		}
		json.NewEncoder(w).Encode(def)
		return
	}

	json.NewEncoder(w).Encode(cfg)
}

// POST /providers/{provider}/config
func saveProviderConfigHandler(w http.ResponseWriter, r *http.Request) {
	if handlePreflight(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	providerName := vars["provider"]

	var newConfig config.ProviderConfig
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON body"})
		return
	}

	newConfig.Provider = providerName
	if err := config.SaveProvider(providerName, newConfig); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	activeCfg := abstractor.GetConfig()
	if strings.EqualFold(activeCfg.Provider, providerName) {
		absCfg := abstractor.MapToAbstractorConfig(newConfig)
		if absCfg.Model == "" {
			absCfg.Model = activeCfg.Model
		}
		if err := abstractor.SaveConfig(absCfg); err == nil {
			if provider, err := abstractor.NewProvider(absCfg); err == nil {
				AgentInstance = provider
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Configuration saved",
	})
}

// GET /providers/{provider}/stats
func getProviderStatsHandler(w http.ResponseWriter, r *http.Request) {
	if handlePreflight(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	providerName := vars["provider"]

	activeCfg := abstractor.GetConfig()
	loadedCfg, _ := config.LoadProvider(providerName)

	status := abstractor.GetConnectionStatus(providerName, loadedCfg)
	queries := abstractor.GetQueryCount(providerName)

	modelName := activeCfg.Model
	if loadedCfg != nil && loadedCfg.Model != "" {
		modelName = loadedCfg.Model
	}

	json.NewEncoder(w).Encode(map[string]any{
		"totalQueries":     queries,
		"currentModel":     modelName,
		"currentProvider":  activeCfg.Provider,
		"connectionStatus": status,
	})
}

// POST /providers/select
func selectProviderHandler(w http.ResponseWriter, r *http.Request) {
	if handlePreflight(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Provider string `json:"provider"`
		Model    string `json:"model"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON body"})
		return
	}

	pName := strings.ToLower(strings.TrimSpace(req.Provider))
	if pName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Provider is required"})
		return
	}

	models, ok := abstractor.GetSupportedModels(pName)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unknown provider"})
		return
	}

	selectedModel := req.Model
	if selectedModel == "" {
		selectedModel = models[0]
	}

	loadedCfg, err := config.LoadProvider(pName)
	var newCfg config.ProviderConfig
	if err == nil && loadedCfg != nil {
		newCfg = *loadedCfg
	} else {
		newCfg = config.ProviderConfig{
			Provider:    pName,
			Endpoint:    "",
			APIKey:      "",
			Temperature: 0.2,
			TopP:        0.9,
			MaxTokens:   400,
			Stream:      false,
			Timeout:     30,
		}
		if pName == "ollama" {
			newCfg.Endpoint = "http://localhost:11434"
		}
	}
	newCfg.Model = selectedModel
	_ = config.SaveProvider(pName, newCfg)

	absCfg := abstractor.MapToAbstractorConfig(newCfg)
	if err := abstractor.SaveConfig(absCfg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if provider, err := abstractor.NewProvider(absCfg); err == nil {
		AgentInstance = provider
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status":   "success",
		"provider": pName,
		"model":    selectedModel,
	})
}
