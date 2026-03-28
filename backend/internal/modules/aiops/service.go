package aiops

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type lexiconTerm struct {
	Term     string
	Severity int
	Category string
}

type Service struct {
	db      *sql.DB
	mu      sync.RWMutex
	terms   []lexiconTerm
	loaded  time.Time
	version int
}

func NewService(db *sql.DB) *Service {
	svc := &Service{db: db}
	_ = svc.ReloadLexicon()
	return svc
}

func (s *Service) ReloadLexicon() error {
	path := filepath.Join("internal", "modules", "aiops", "data", "lexicon_vi.txt")
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open lexicon failed: %w", err)
	}
	defer file.Close()

	terms := make([]lexiconTerm, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}
		term := strings.TrimSpace(parts[0])
		sev := parseSeverity(strings.TrimSpace(parts[1]))
		cat := strings.TrimSpace(parts[2])
		if term == "" || sev <= 0 {
			continue
		}
		terms = append(terms, lexiconTerm{Term: term, Severity: sev, Category: cat})
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan lexicon failed: %w", err)
	}

	s.mu.Lock()
	s.terms = terms
	s.loaded = time.Now()
	s.version++
	s.mu.Unlock()

	return nil
}

func (s *Service) Summary() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]interface{}{
		"termCount": len(s.terms),
		"loadedAt":  s.loaded.Format(time.RFC3339),
		"version":   s.version,
	}
}

func (s *Service) Analyze(req analyzeRequest) (string, int, []termScore) {
	norm := normalizeText(req.Text)
	s.mu.RLock()
	defer s.mu.RUnlock()

	hits := make([]termScore, 0)
	total := 0
	for _, term := range s.terms {
		count := strings.Count(norm, term.Term)
		if count <= 0 {
			continue
		}
		score := term.Severity * count
		total += score
		hits = append(hits, termScore{
			Term:     term.Term,
			Severity: term.Severity,
			Category: term.Category,
			Hits:     count,
		})
	}

	risk := classifyRisk(total)
	return risk, total, hits
}

func (s *Service) AnalyzeAndFlag(ctx context.Context, req analyzeRequest) (map[string]interface{}, error) {
	risk, score, hits := s.Analyze(req)
	if risk == "none" {
		return map[string]interface{}{
			"riskLevel": risk,
			"score":     score,
			"hits":      hits,
			"flagged":   false,
		}, nil
	}

	reasonCode := "aiops_risk_" + risk
	reasonDetail := summarizeHits(hits)

	_, err := s.db.ExecContext(ctx,
		"INSERT INTO flags (source_type, source_id, risk_level, reason_code, reason_detail, status) VALUES (?, ?, ?, ?, ?, 'queued')",
		req.SourceType,
		req.SourceID,
		risk,
		reasonCode,
		reasonDetail,
	)
	if err != nil {
		return nil, fmt.Errorf("insert aiops flag failed: %w", err)
	}

	return map[string]interface{}{
		"riskLevel": risk,
		"score":     score,
		"hits":      hits,
		"flagged":   true,
	}, nil
}

func (s *Service) EnrichedQueue(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, source_type, source_id, risk_level, reason_code, reason_detail, status, created_at FROM flags WHERE status = 'queued' ORDER BY risk_level DESC, id ASC LIMIT 200",
	)
	if err != nil {
		return nil, fmt.Errorf("query enriched queue failed: %w", err)
	}
	defer rows.Close()

	items := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, sourceID uint64
		var sourceType, riskLevel, reasonCode, reasonDetail, status, createdAt string
		if err := rows.Scan(&id, &sourceType, &sourceID, &riskLevel, &reasonCode, &reasonDetail, &status, &createdAt); err != nil {
			return nil, fmt.Errorf("scan enriched queue failed: %w", err)
		}
		items = append(items, map[string]interface{}{
			"id":         id,
			"sourceType": sourceType,
			"sourceId":   sourceID,
			"riskLevel":  riskLevel,
			"reasonCode": reasonCode,
			"reasonDetail": reasonDetail,
			"status":     status,
			"createdAt":  createdAt,
			"handoffHint": handoffHint(riskLevel, reasonCode),
		})
	}
	return items, nil
}

func parseSeverity(raw string) int {
	switch raw {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	default:
		return 0
	}
}

func normalizeText(input string) string {
	text := strings.ToLower(input)
	replacer := strings.NewReplacer("\n", " ", "\t", " ", ".", " ", ",", " ", ";", " ", ":", " ", "!", " ", "?", " ")
	return replacer.Replace(text)
}

func classifyRisk(score int) string {
	if score >= 8 {
		return "red"
	}
	if score >= 3 {
		return "yellow"
	}
	return "none"
}

func summarizeHits(hits []termScore) string {
	parts := make([]string, 0, len(hits))
	for _, hit := range hits {
		parts = append(parts, fmt.Sprintf("%s(x%d,s%d)", hit.Term, hit.Hits, hit.Severity))
	}
	return strings.Join(parts, "; ")
}

func handoffHint(riskLevel, reasonCode string) string {
	if riskLevel == "red" {
		return "Immediate moderator review required; consider temporary hide"
	}
	if strings.Contains(reasonCode, "defamation") {
		return "Cross-check for defamatory statements and evidence requirement"
	}
	return "Queue for normal moderation SLA"
}
