package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/thiagotn/investment-analyzer/internal/crypto"
	"github.com/thiagotn/investment-analyzer/internal/domain"
)

type Store struct {
	dir        string
	passphrase string
}

func NewStore(dir string, passphrase string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create snapshots directory: %w", err)
	}
	return &Store{dir: dir, passphrase: passphrase}, nil
}

func (s *Store) Save(snap *domain.Snapshot) error {
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	if s.passphrase != "" {
		encData, err := crypto.Encrypt(data, s.passphrase)
		if err != nil {
			return fmt.Errorf("failed to encrypt snapshot: %w", err)
		}

		encPath := filepath.Join(s.dir, snap.Month+".json.enc")
		if err := os.WriteFile(encPath, encData, 0644); err != nil {
			return fmt.Errorf("failed to write encrypted snapshot: %w", err)
		}

		plainPath := filepath.Join(s.dir, snap.Month+".json")
		if err := os.Remove(plainPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove plaintext snapshot: %w", err)
		}

		return nil
	}

	path := filepath.Join(s.dir, snap.Month+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot: %w", err)
	}

	return nil
}

func (s *Store) Load(month string) (*domain.Snapshot, error) {
	var data []byte

	encPath := filepath.Join(s.dir, month+".json.enc")
	if _, err := os.Stat(encPath); err == nil {
		encData, err := os.ReadFile(encPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read encrypted snapshot: %w", err)
		}

		data, err = crypto.Decrypt(encData, s.passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt snapshot: %w", err)
		}
	} else {
		plainPath := filepath.Join(s.dir, month+".json")
		data, err = os.ReadFile(plainPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read snapshot: %w", err)
		}
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

	monthSet := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".json.enc") {
			month := name[:len(name)-9]
			monthSet[month] = true
		} else if strings.HasSuffix(name, ".json") {
			month := name[:len(name)-5]
			monthSet[month] = true
		}
	}

	var months []string
	for month := range monthSet {
		months = append(months, month)
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
