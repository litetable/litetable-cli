package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	litetable2 "github.com/litetable/litetable-cli/internal/litetable"
	"github.com/litetable/litetable-cli/internal/server"
	"net/http"
	"strings"
)

type litetable interface {
	CreateFamilies(ctx context.Context, p *server.CreateFamilyParams) error
	Read(ctx context.Context, p *server.ReadParams) (map[string]*litetable2.Row, error)
	Write(ctx context.Context, p *server.WriteParams) (map[string]*litetable2.Row, error)
	Delete(ctx context.Context, p *server.DeleteParams) error
}

const (
	queryCreate = "CREATE"
	queryRead   = "READ"
	queryWrite  = "WRITE"
	queryDelete = "DELETE"
)

type payload struct {
	Type       string             `json:"type"`
	ReadType   string             `json:"readType"`
	Key        string             `json:"key"`
	Family     string             `json:"family"`
	Qualifiers []server.Qualifier `json:"qualifiers"`
	Latest     int                `json:"latest"`
	Families   []string           `json:"families"`
}

type handler struct {
	server litetable
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	// Decode the JSON payload
	var p payload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to decode JSON payload: %v", err),
		})
		return
	}

	// Check if the query type is empty
	if p.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Query type must be specified",
		})
		return
	}

	if p.Type == queryRead {
		data, err := h.handleReadQuery(r.Context(), &p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}

		if err = json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Failed to encode response: %v", err),
			})
			return
		}
	}

	if p.Type == queryWrite {
		data, err := h.handleWriteQuery(r.Context(), &p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}

		if err = json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Failed to encode response: %v", err),
			})
			return
		}
	}

	if p.Type == queryCreate {
		if err := h.handleCreateFamilyQuery(r.Context(), &p); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}

		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Families created successfully",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Failed to encode response: %v", err),
			})
			return
		}
	}

	if p.Type == queryDelete {
		if err := h.handleDeleteQuery(r.Context(), &p); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("%v", err),
			})
			return
		}

		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Deleted successfully",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Failed to encode response: %v", err),
			})
			return
		}
	}
}

func (h *handler) handleDeleteQuery(ctx context.Context, p *payload) error {
	params := &server.DeleteParams{
		Key:    p.Key,
		Family: p.Family,
	}

	// qualifiers are optional, so only add them if they are present
	if len(p.Qualifiers) > 0 {
		for _, q := range p.Qualifiers {
			params.Qualifiers = append(params.Qualifiers, q.Name)
		}
	}

	return h.server.Delete(ctx, params)
}

func (h *handler) handleCreateFamilyQuery(ctx context.Context, p *payload) error {
	params := &server.CreateFamilyParams{
		Families: p.Families,
	}

	return h.server.CreateFamilies(ctx, params)
}

func (h *handler) handleReadQuery(ctx context.Context, p *payload) (any, error) {

	params := &server.ReadParams{
		Key:        p.Key,
		QueryType:  server.Read,
		Family:     p.Family,
		Qualifiers: []string{},
		Latest:     int32(p.Latest),
	}

	if p.ReadType == "prefix" {
		params.QueryType = server.ReadPrefix
	}

	if p.ReadType == "regex" {
		params.QueryType = server.ReadRegex
		// Don't add wildcards if user's input already contains regex patterns
		if strings.ContainsAny(p.Key, ".*+?^$[](){}|\\") {
			params.Key = p.Key // Keep the regex as-is
		} else {
			// Only add wildcards if it's a simple string search
			params.Key = fmt.Sprintf(".*%s.*", p.Key)
		}
	}

	// for Read we only need to add the qualifiers
	if len(p.Qualifiers) > 0 {
		for _, q := range p.Qualifiers {
			params.Qualifiers = append(params.Qualifiers, q.Name)
		}
	}

	return h.server.Read(ctx, params)

}

func (h *handler) handleWriteQuery(ctx context.Context, p *payload) (any, error) {
	params := &server.WriteParams{
		Key:        p.Key,
		Family:     p.Family,
		Qualifiers: p.Qualifiers,
	}

	return h.server.Write(ctx, params)
}
