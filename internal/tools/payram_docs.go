package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/payram/payram-analytics-mcp-server/internal/protocol"
)

// payramDocsTool provides search and section retrieval over local PayRam docs.
// Actions:
//   - search: query markdown corpus, optional category filter, limit results
//   - get_section: return a specific section (by path and optional heading) or whole file
type payramDocsTool struct {
	sections       []docSection
	sectionsByPath map[string][]docSection // path -> sections
	files          map[string]string       // path -> full content
}

// docSection represents a single heading + content block within a markdown file.
type docSection struct {
	Path     string
	Heading  string
	Body     string
	Category string
	Tags     []string
}

// PayramDocs builds the docs tool, indexing markdown under docs/payram-docs by default.
func PayramDocs() *payramDocsTool {
	root := strings.TrimSpace(os.Getenv("PAYRAM_DOCS_ROOT"))
	if root == "" {
		root = filepath.Join("docs", "payram-docs")
	}

	sections, byPath, files := indexDocs(root)
	applyTopics(sections)

	return &payramDocsTool{
		sections:       sections,
		sectionsByPath: byPath,
		files:          files,
	}
}

// Descriptor describes the tool.
func (t *payramDocsTool) Descriptor() protocol.ToolDescriptor {
	return protocol.ToolDescriptor{
		Name:        "payram_docs",
		Description: "Search PayRam docs and return sections. Categories: faqs, features, onboarding-guide. Actions: search, get_section, list_index. Keywords are boosted when they match headings or curated topic tags (analytics, payouts, hot wallet, payment links, multi-brand, multi-currency, customer deposit wallets, API/webhooks, config/deployment, debug).",
		InputSchema: &protocol.JSONSchema{
			Type: "object",
			Properties: map[string]protocol.JSONSchema{
				"action": {
					Type:        "string",
					Enum:        []string{"search", "get_section", "list_index"},
					Description: "Action to perform",
				},
				"query": {
					Type:        "string",
					Description: "Search query (required for search)",
				},
				"category": {
					Type:        "string",
					Description: "Optional category filter (faqs, features, onboarding-guide)",
				},
				"limit": {
					Type:        "integer",
					Description: "Max results (default 3, max 10)",
				},
				"path": {
					Type:        "string",
					Description: "Doc path relative to docs/payram-docs (required for get_section)",
				},
				"heading": {
					Type:        "string",
					Description: "Heading within the file (optional; if omitted returns whole file)",
				},
			},
			Required: []string{"action"},
		},
	}
}

type docsArgs struct {
	Action   string `json:"action"`
	Query    string `json:"query"`
	Category string `json:"category"`
	Limit    int    `json:"limit"`
	Path     string `json:"path"`
	Heading  string `json:"heading"`
}

// Invoke routes search and section fetch.
func (t *payramDocsTool) Invoke(ctx context.Context, raw json.RawMessage) (protocol.CallResult, *protocol.ResponseError) {
	_ = ctx
	var args docsArgs
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &args); err != nil {
			return protocol.CallResult{}, &protocol.ResponseError{Code: -32602, Message: "invalid arguments"}
		}
	}

	switch args.Action {
	case "search":
		if strings.TrimSpace(args.Query) == "" {
			return protocol.CallResult{}, &protocol.ResponseError{Code: -32602, Message: "query is required for search"}
		}
		limit := args.Limit
		if limit <= 0 {
			limit = 3
		}
		if limit > 10 {
			limit = 10
		}
		return t.search(args.Query, args.Category, limit)
	case "get_section":
		if strings.TrimSpace(args.Path) == "" {
			return protocol.CallResult{}, &protocol.ResponseError{Code: -32602, Message: "path is required for get_section"}
		}
		return t.getSection(args.Path, args.Heading)
	case "list_index":
		return t.listIndex(), nil
	default:
		return protocol.CallResult{}, &protocol.ResponseError{Code: -32602, Message: "action must be search or get_section"}
	}
}

// search performs a simple keyword match over headings and bodies.
func (t *payramDocsTool) search(query, category string, limit int) (protocol.CallResult, *protocol.ResponseError) {
	q := strings.ToLower(strings.TrimSpace(query))
	words := strings.Fields(q)
	if len(words) == 0 {
		return protocol.CallResult{}, &protocol.ResponseError{Code: -32602, Message: "empty query"}
	}
	cat := strings.TrimSpace(strings.ToLower(category))

	type hit struct {
		sec   docSection
		score int
	}

	hits := make([]hit, 0)
	for _, sec := range t.sections {
		if cat != "" && strings.ToLower(sec.Category) != cat {
			continue
		}
		hScore := scoreSection(sec, words)
		if hScore > 0 {
			hits = append(hits, hit{sec: sec, score: hScore})
		}
	}

	if len(hits) == 0 {
		return protocol.CallResult{Content: []protocol.ContentPart{{Type: "text", Text: "No results."}}}, nil
	}

	sort.Slice(hits, func(i, j int) bool {
		if hits[i].score == hits[j].score {
			return hits[i].sec.Path < hits[j].sec.Path
		}
		return hits[i].score > hits[j].score
	})

	if len(hits) > limit {
		hits = hits[:limit]
	}

	var b strings.Builder
	b.WriteString("Results:\n")
	for i, h := range hits {
		excerpt := trimExcerpt(h.sec.Body, 320)
		fmtPath := h.sec.Path
		if h.sec.Heading != "" {
			fmtPath += "#" + h.sec.Heading
		}
		b.WriteString(fmt.Sprintf("%d) [%s] (%s)\n", i+1, fmtPath, h.sec.Category))
		b.WriteString(excerpt)
		b.WriteString("\n\n")
	}

	return protocol.CallResult{Content: []protocol.ContentPart{{Type: "text", Text: strings.TrimSpace(b.String())}}}, nil
}

// getSection returns a specific section or whole file.
func (t *payramDocsTool) getSection(path, heading string) (protocol.CallResult, *protocol.ResponseError) {
	norm := filepath.ToSlash(strings.TrimSpace(path))
	norm = strings.TrimPrefix(norm, "./")

	sections, ok := t.sectionsByPath[norm]
	if !ok {
		return protocol.CallResult{}, &protocol.ResponseError{Code: -32004, Message: "path not found"}
	}

	if strings.TrimSpace(heading) == "" {
		full := strings.TrimSpace(t.files[norm])
		if full == "" {
			return protocol.CallResult{}, &protocol.ResponseError{Code: -32004, Message: "content not found"}
		}
		return protocol.CallResult{Content: []protocol.ContentPart{{Type: "text", Text: full}}}, nil
	}

	target := strings.ToLower(strings.TrimSpace(heading))
	for _, sec := range sections {
		if strings.ToLower(sec.Heading) == target {
			text := fmt.Sprintf("%s\n\n%s", sec.Heading, sec.Body)
			return protocol.CallResult{Content: []protocol.ContentPart{{Type: "text", Text: strings.TrimSpace(text)}}}, nil
		}
	}

	return protocol.CallResult{}, &protocol.ResponseError{Code: -32004, Message: "heading not found in file"}
}

// listIndex returns available categories, topics, and per-file headings (truncated).
func (t *payramDocsTool) listIndex() protocol.CallResult {
	cats := make(map[string]struct{})
	fileHeadings := make(map[string][]string)
	for _, sec := range t.sections {
		cats[sec.Category] = struct{}{}
		hs := fileHeadings[sec.Path]
		if len(hs) < 8 { // cap to avoid huge payloads
			fileHeadings[sec.Path] = append(hs, sec.Heading)
		}
	}

	// sort categories
	catList := make([]string, 0, len(cats))
	for c := range cats {
		catList = append(catList, c)
	}
	sort.Strings(catList)

	// topics map
	topics := topicMap()

	var b strings.Builder
	b.WriteString("Categories: ")
	b.WriteString(strings.Join(catList, ", "))
	b.WriteString("\n\nTopics (keywords -> paths):\n")
	// invert topics for display
	keys := make([]string, 0, len(topics))
	for p := range topics {
		keys = append(keys, p)
	}
	sort.Strings(keys)
	for _, p := range keys {
		b.WriteString(fmt.Sprintf("- %s : %s\n", p, strings.Join(topics[p], ", ")))
	}

	b.WriteString("\nHeadings (truncated per file):\n")
	paths := make([]string, 0, len(fileHeadings))
	for p := range fileHeadings {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		hs := fileHeadings[p]
		b.WriteString(fmt.Sprintf("- %s: %s\n", p, strings.Join(hs, "; ")))
	}

	return protocol.CallResult{Content: []protocol.ContentPart{{Type: "text", Text: strings.TrimSpace(b.String())}}}
}

// scoreSection gives weight to heading matches and body matches.
func scoreSection(sec docSection, words []string) int {
	head := strings.ToLower(sec.Heading)
	body := strings.ToLower(sec.Body)
	tags := strings.ToLower(strings.Join(sec.Tags, " "))
	score := 0
	for _, w := range words {
		if w == "" {
			continue
		}
		if strings.Contains(tags, w) {
			score += 5
		}
		if strings.Contains(head, w) {
			score += 3
		}
		if strings.Contains(body, w) {
			score += 1
		}
	}
	return score
}

// indexDocs walks root, parsing markdown files into sections.
func indexDocs(root string) ([]docSection, map[string][]docSection, map[string]string) {
	sections := make([]docSection, 0)
	byPath := make(map[string][]docSection)
	files := make(map[string]string)

	root = filepath.Clean(root)

	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		content := string(data)

		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}
		rel = filepath.ToSlash(rel)

		cats := strings.Split(rel, "/")
		category := "docs"
		if len(cats) > 1 {
			category = cats[0]
		}

		secs := parseSections(rel, category, content)
		sections = append(sections, secs...)
		byPath[rel] = secs
		files[rel] = strings.TrimSpace(content)
		return nil
	})

	return sections, byPath, files
}

// applyTopics attaches curated topic tags to sections based on their path.
func applyTopics(sections []docSection) {
	pathTopics := topicMap()
	for i := range sections {
		p := strings.ToLower(sections[i].Path)
		if tags, ok := pathTopics[p]; ok {
			sections[i].Tags = append(sections[i].Tags, tags...)
		}
	}
}

// topicMap defines curated topic tags for known docs.
func topicMap() map[string][]string {
	return map[string][]string{
		"onboarding-guide/hot-wallet-setup.md":               {"hot wallet", "sweep", "custody"},
		"onboarding-guide/funds-sweeping.md":                 {"funds sweep", "automation", "consolidation"},
		"onboarding-guide/testing-payment-links.md":          {"payment links", "testing", "checkout"},
		"features/payment-links.md":                          {"payment links", "checkout", "invoice"},
		"features/payouts.md":                                {"payouts", "withdraw", "settlement"},
		"features/analytics-and-reporting.md":                {"analytics", "reports", "dashboard"},
		"features/user-management.md":                        {"user roles", "rbac", "permissions"},
		"features/multi-brand-setup.md":                      {"multi-brand", "white label"},
		"features/multi-currency-and-multi-chain-support.md": {"multi-currency", "multi-chain", "tokens", "networks"},
		"features/customer-deposit-wallets.md":               {"deposit wallet", "customer wallet"},
		"features/fiat-onramp.md":                            {"on-ramp", "fiat", "buy crypto"},
		"faqs/api-integration-faqs.md":                       {"api", "webhook", "integration"},
		"faqs/configuration-faqs.md":                         {"config", "yaml", "env"},
		"faqs/deployment-faqs.md":                            {"deployment", "server requirements", "install"},
		"faqs/debug-faqs.md":                                 {"debug", "errors", "failures"},
		"faqs/fund-management-faqs.md":                       {"fund management", "settlement", "withdraw"},
		"faqs/referral-faqs.md":                              {"referral", "affiliate"},
		"faqs/general-faqs.md":                               {"overview", "general"},
	}
}

// parseSections breaks a markdown file into heading-based sections.
func parseSections(path, category, content string) []docSection {
	lines := strings.Split(content, "\n")
	currentHeading := "Introduction"
	var buf []string
	sections := make([]docSection, 0)

	flush := func() {
		text := strings.TrimSpace(strings.Join(buf, "\n"))
		if text == "" {
			buf = buf[:0]
			return
		}
		sections = append(sections, docSection{
			Path:     path,
			Heading:  currentHeading,
			Body:     text,
			Category: category,
		})
		buf = buf[:0]
	}

	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if h := parseHeading(trim); h != "" {
			flush()
			currentHeading = h
			continue
		}
		buf = append(buf, line)
	}
	flush()
	return sections
}

// parseHeading returns heading text if the line is a markdown heading.
func parseHeading(line string) string {
	if !strings.HasPrefix(line, "#") {
		return ""
	}
	// Count leading #'s then trim
	idx := 0
	for idx < len(line) && line[idx] == '#' {
		idx++
	}
	if idx == 0 || idx == len(line) {
		return ""
	}
	return strings.TrimSpace(line[idx:])
}

// trimExcerpt shortens body text to the requested length (runes).
func trimExcerpt(s string, max int) string {
	s = strings.TrimSpace(s)
	if max <= 0 {
		return s
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max]) + "..."
}
