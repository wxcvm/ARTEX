package db

import (
	"context"
	"strconv"
)

// ModelKeyStatus reports storage coverage and authentication with the currently
// configured application key. It deliberately contains no key material.
type ModelKeyStatus struct {
	KeyState      string            `json:"key_state"`
	Total         int               `json:"total"`
	Empty         int               `json:"empty"`
	PlaintextOnly int               `json:"plaintext_only"`
	EncryptedOnly int               `json:"encrypted_only"`
	Mixed         int               `json:"mixed"`
	Verified      int               `json:"verified"`
	Failed        int               `json:"failed"`
	Unchecked     int               `json:"unchecked"`
	Profiles      []ModelKeyProfile `json:"profiles"`
}

type ModelKeyProfile struct {
	ID        string `json:"id"`
	Storage   string `json:"storage"`
	Integrity string `json:"integrity"`
}

func newModelKeyStatus() ModelKeyStatus {
	aead, err := profileSecretCipher()
	state := "missing"
	if err != nil {
		state = "invalid"
	} else if aead != nil {
		state = "configured"
	}
	return ModelKeyStatus{KeyState: state, Profiles: []ModelKeyProfile{}}
}

func (s *ModelKeyStatus) observe(id int64, hasPlaintext bool, encrypted string) {
	p := ModelKeyProfile{ID: strconv.FormatInt(id, 10), Integrity: "not_applicable"}
	s.Total++
	switch {
	case hasPlaintext && encrypted != "":
		p.Storage = "mixed"
		s.Mixed++
	case hasPlaintext:
		p.Storage = "plaintext_only"
		s.PlaintextOnly++
	case encrypted != "":
		p.Storage = "encrypted_only"
		s.EncryptedOnly++
	default:
		p.Storage = "empty"
		s.Empty++
	}
	if encrypted != "" {
		switch {
		case s.KeyState != "configured":
			p.Integrity = "unchecked"
			s.Unchecked++
		default:
			// Decrypted bytes are discarded. Neither the original key nor its
			// ciphertext or decoder errors are included in the result.
			if _, err := decodeProfileSecret("", encrypted); err != nil {
				p.Integrity = "failed"
				s.Failed++
			} else {
				p.Integrity = "verified"
				s.Verified++
			}
		}
	}
	s.Profiles = append(s.Profiles, p)
}

// ModelKeyDiagnostics is a read-only snapshot. PostgreSQL returns a boolean for
// legacy plaintext presence, not the plaintext value itself.
func (d *DB) ModelKeyDiagnostics(ctx context.Context) (ModelKeyStatus, error) {
	status := newModelKeyStatus()
	rows, err := d.QueryContext(ctx, `SELECT id, COALESCE(api_key,'') <> '', COALESCE(api_key_ciphertext,'') FROM llm_profiles ORDER BY id`)
	if err != nil {
		return ModelKeyStatus{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var hasPlaintext bool
		var encrypted string
		if err := rows.Scan(&id, &hasPlaintext, &encrypted); err != nil {
			return ModelKeyStatus{}, err
		}
		status.observe(id, hasPlaintext, encrypted)
	}
	if err := rows.Err(); err != nil {
		return ModelKeyStatus{}, err
	}
	return status, nil
}
