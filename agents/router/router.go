package router

import (
	"sort"
	"strings"

	"github.com/rahulkumarparida/roxkv/agents/abstractor"
	"github.com/rahulkumarparida/roxkv/agents/registry"
)

// DefaultMaxTools is the maximum number of tools returned when routing
// finds matches. It keeps the LLM context lean for small models.
const DefaultMaxTools = 10

// stopWords are common words excluded from query matching to avoid
// noise in the scoring process.
var stopWords = map[string]bool{
	"a": true, "an": true, "the": true, "and": true, "or": true,
	"is": true, "are": true, "was": true, "were": true, "be": true,
	"to": true, "of": true, "in": true, "for": true, "on": true,
	"it": true, "do": true, "does": true, "did": true, "with": true,
	"at": true, "by": true, "from": true, "as": true, "i": true,
	"me": true, "my": true, "we": true, "our": true, "you": true,
	"your": true, "he": true, "she": true, "they": true, "them": true,
	"this": true, "that": true, "what": true, "which": true, "who": true,
	"how": true, "can": true, "will": true, "would": true, "should": true,
	"could": true, "please": true, "tell": true, "about": true,"also": true, "just": true, "like": true, "some": true, "any": true,"so": true, "if": true, "then": true, "than": true, "but": true,
	"because": true, "while": true, "where": true, "when": true, "all": true,
	"anyone": true, "everyone": true, "someone": true, "noone": true,
}

// scoredTool pairs a metadata entry with its computed relevance score.
type scoredTool struct {
	meta  registry.ToolMetadata `json:"meta"`
	score float64               `json:"score"`
}

// RouteTools scores every registered tool against the user query and
// returns only the top-scoring subset. If no tool scores above zero
// the full tool list is returned as a fallback.
//
// The scoring algorithm is purely metadata-driven — no hardcoded
// query→tool mappings exist. Each tool accumulates points from:
//   - Exact keyword matches:       +3 × (priority/10)
//   - Synonym matches:             +2 × (priority/10)
//   - Partial substring matches:   +1
//   - Example phrase word overlap:  +2 × (priority/10)
//   - Category name in query:       +1
func RouteTools(query string, allMeta []registry.ToolMetadata, maxTools int) []abstractor.GenericToolDefinition {
	if maxTools <= 0 {
		maxTools = DefaultMaxTools
	}

	queryLower := strings.ToLower(strings.TrimSpace(query))
	queryWords := tokenize(queryLower)

	if len(queryWords) == 0 {
		return extractTools(allMeta)
	}

	scored := make([]scoredTool, 0, len(allMeta))

	for _, meta := range allMeta {
		score := scoreTool(meta, queryLower, queryWords)
		if score > 0 {
			scored = append(scored, scoredTool{meta: meta, score: score})
		}
	}

	// Fallback: if nothing matched, expose everything so the LLM can
	// still respond (preserves backward compatibility).
	if len(scored) == 0 {
		return extractTools(allMeta)
	}

	// Sort descending by score.
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Cap the result set.
	if len(scored) > maxTools {
		scored = scored[:maxTools]
	}

	tools := make([]abstractor.GenericToolDefinition, 0, len(scored))
	for _, s := range scored {
		tools = append(tools, s.meta.Tool)
	}
	return tools
}

// scoreTool computes the relevance score for one tool against the query.
func scoreTool(meta registry.ToolMetadata, queryLower string, queryWords []string) float64 {
	var score float64
	priorityWeight := float64(meta.Priority) / 10.0
	if priorityWeight < 0.1 {
		priorityWeight = 0.1
	}

	// --- Keyword matching ---
	for _, kw := range meta.Keywords {
		kwLower := strings.ToLower(kw)
		for _, qw := range queryWords {
			if qw == kwLower {
				// Exact keyword match.
				score += 3.0 * priorityWeight
			} else if strings.Contains(qw, kwLower) || strings.Contains(kwLower, qw) {
				// Partial substring match.
				score += 1.0
			}
		}
	}

	// --- Synonym matching ---
	for _, syn := range meta.Synonyms {
		synLower := strings.ToLower(syn)
		for _, qw := range queryWords {
			if qw == synLower {
				score += 2.0 * priorityWeight
			} else if strings.Contains(qw, synLower) || strings.Contains(synLower, qw) {
				score += 1.0
			}
		}
	}

	// --- Example phrase matching ---
	for _, example := range meta.Examples {
		exLower := strings.ToLower(example)
		exWords := tokenize(exLower)
		overlap := 0
		for _, qw := range queryWords {
			for _, ew := range exWords {
				if qw == ew {
					overlap++
					break
				}
			}
		}
		if overlap > 0 {
			score += float64(overlap) * 2.0 * priorityWeight
		}
	}

	// --- Category bonus ---
	catLower := strings.ToLower(meta.Category)
	if strings.Contains(queryLower, catLower) {
		score += 1.0
	}

	return score
}

// tokenize splits text into lowercase words, removing stop-words.
func tokenize(text string) []string {
	raw := strings.Fields(text)
	words := make([]string, 0, len(raw))
	for _, w := range raw {
		// Strip common punctuation.
		w = strings.Trim(w, ".,;:!?\"'`()[]{}")
		w = strings.ToLower(w)
		if w == "" || stopWords[w] {
			continue
		}
		words = append(words, w)
	}
	return words
}

// extractTools converts a metadata slice into an abstractor.GenericToolDefinition slice.
func extractTools(metas []registry.ToolMetadata) []abstractor.GenericToolDefinition {
	tools := make([]abstractor.GenericToolDefinition, 0, len(metas))
	for _, m := range metas {
		tools = append(tools, m.Tool)
	}
	return tools
}
