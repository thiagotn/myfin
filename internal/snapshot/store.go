package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/thiagotn/investment-analyzer/internal/domain"
)

type Store struct {
	dir string
}

func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create snapshots directory: %w", err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) Save(snap *domain.Snapshot) error {
	path := filepath.Join(s.dir, snap.Month+".json")

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot: %w", err)
	}

	return nil
}

func (s *Store) Load(month string) (*domain.Snapshot, error) {
	path := filepath.Join(s.dir, month+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot: %w", err)
	}

	var snap domain.Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	return &snap, nil
}

func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshots directory: %w", err)
	}

	var months []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			month := entry.Name()[:len(entry.Name())-5]
			months = append(months, month)
		}
	}

	sort.Strings(months)
	return months, nil
}

func (s *Store) LoadRecent(n int) ([]domain.Snapshot, error) {
	months, err := s.List()
	if err != nil {
		return nil, err
	}

	if len(months) < n {
		n = len(months)
	}

	start := len(months) - n
	if start < 0 {
		start = 0
	}

	var snapshots []domain.Snapshot
	for _, month := range months[start:] {
		snap, err := s.Load(month)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, *snap)
	}

	return snapshots, nil
}
