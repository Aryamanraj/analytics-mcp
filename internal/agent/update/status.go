package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// UpdateStatus captures persisted update state.
type UpdateStatus struct {
	CurrentVersion      string    `json:"current_version"`
	PreviousVersion     string    `json:"previous_version"`
	LastSuccessVersion  string    `json:"last_success_version"`
	LastSuccessAt       time.Time `json:"last_success_at"`
	LastAttemptVersion  string    `json:"last_attempt_version"`
	LastAttemptAt       time.Time `json:"last_attempt_at"`
	LastErrorCode       string    `json:"last_error_code"`
	LastErrorMessage    string    `json:"last_error_message"`
	LastErrorAt         time.Time `json:"last_error_at"`
	InProgress          bool      `json:"in_progress"`
	InProgressStartedAt time.Time `json:"in_progress_started_at"`
}

// LoadStatus loads persisted status, returning a zero value when missing.
func LoadStatus() (UpdateStatus, error) {
	if err := EnsureBaseDirs(); err != nil {
		return UpdateStatus{}, err
	}

	path := statusPath()
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return UpdateStatus{}, nil
		}
		return UpdateStatus{}, err
	}

	var st UpdateStatus
	if err := json.Unmarshal(raw, &st); err != nil {
		return UpdateStatus{}, err
	}
	return st, nil
}

// SaveStatus persists the provided status atomically.
func SaveStatus(st UpdateStatus) error {
	if err := EnsureBaseDirs(); err != nil {
		return err
	}

	path := statusPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	tmp := path + ".tmp"
	raw, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}

	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	return nil
}

// MarkAttempt sets in-progress markers.
func (s *UpdateStatus) MarkAttempt() {
	s.InProgress = true
	s.InProgressStartedAt = time.Now()
	s.LastAttemptVersion = ""
	s.LastAttemptAt = time.Now()
	s.LastErrorCode = ""
	s.LastErrorMessage = ""
	s.LastErrorAt = time.Time{}
}

// MarkSuccess records a successful update.
func (s *UpdateStatus) MarkSuccess(current, previous string) {
	s.CurrentVersion = current
	s.PreviousVersion = previous
	s.LastSuccessVersion = current
	s.LastSuccessAt = time.Now()
	s.InProgress = false
}

// MarkFailure records a failed update attempt.
func (s *UpdateStatus) MarkFailure(code, msg string) {
	s.LastErrorCode = code
	s.LastErrorMessage = msg
	s.LastErrorAt = time.Now()
	s.InProgress = false
}

func statusPath() string {
	return filepath.Join(StateDir(), "update_status.json")
}
