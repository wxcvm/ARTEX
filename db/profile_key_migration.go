package db

import (
	"context"
	"errors"
	"strconv"
)

var (
	ErrModelKeyMigrationKey = errors.New("configure ARTEX_LLM_SECRET_KEY before migration")
	ErrModelKeyMigrationIntegrity = errors.New("existing model ciphertext cannot be authenticated with the current key")
	ErrModelKeyMigrationChanged = errors.New("model key storage changed; refresh the diagnostic preview")
)

// ModelKeyMigrationResult deliberately returns counts, never key material.
type ModelKeyMigrationResult struct {
	Migrated int `json:"migrated"`
	SkippedMixed int `json:"skipped_mixed"`
}

type pendingModelKey struct {
	id int64
	value string
}

// MigrateModelKeys converts exactly the plaintext-only profiles from the
// caller's preview. Existing ciphertext is authenticated first and mixed rows
// are left untouched. All writes use one transaction and are rolled back on
// failure; original API keys and model settings are preserved.
func (d *DB) MigrateModelKeys(ctx context.Context, expectedIDs []string) (ModelKeyMigrationResult, error) {
	aead, err := profileSecretCipher()
	if err != nil || aead == nil {
		return ModelKeyMigrationResult{}, ErrModelKeyMigrationKey
	}
	expected := make(map[string]bool, len(expectedIDs))
	for _, id := range expectedIDs {
		if id == "" || expected[id] {
			return ModelKeyMigrationResult{}, ErrModelKeyMigrationChanged
		}
		expected[id] = true
	}
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return ModelKeyMigrationResult{}, err
	}
	defer tx.Rollback()

	// Lock the relevant current rows. The CASE avoids loading plaintext from
	// mixed rows, which must be resolved explicitly instead of overwritten.
	rows, err := tx.QueryContext(ctx, `SELECT id,
		CASE WHEN COALESCE(api_key_ciphertext,'') = '' THEN COALESCE(api_key,'') ELSE '' END,
		COALESCE(api_key,'') <> '',
		COALESCE(api_key_ciphertext,'')
		FROM llm_profiles
		WHERE COALESCE(api_key,'') <> '' OR COALESCE(api_key_ciphertext,'') <> ''
		ORDER BY id FOR UPDATE`)
	if err != nil {
		return ModelKeyMigrationResult{}, err
	}
	var pending []pendingModelKey
	var result ModelKeyMigrationResult
	for rows.Next() {
		var id int64
		var plain, encrypted string
		var hasPlain bool
		if err := rows.Scan(&id, &plain, &hasPlain, &encrypted); err != nil {
			rows.Close()
			return ModelKeyMigrationResult{}, err
		}
		if encrypted != "" {
			if _, err := decodeProfileSecret("", encrypted); err != nil {
				rows.Close()
				return ModelKeyMigrationResult{}, ErrModelKeyMigrationIntegrity
			}
			if hasPlain {
				result.SkippedMixed++
			}
			continue
		}
		if !expected[strconv.FormatInt(id, 10)] {
			rows.Close()
			return ModelKeyMigrationResult{}, ErrModelKeyMigrationChanged
		}
		pending = append(pending, pendingModelKey{id: id, value: plain})
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return ModelKeyMigrationResult{}, err
	}
	if len(pending) != len(expected) {
		return ModelKeyMigrationResult{}, ErrModelKeyMigrationChanged
	}
	for _, row := range pending {
		plain, encrypted, err := encodeProfileSecret(row.value)
		if err != nil || plain != "" || encrypted == "" {
			return ModelKeyMigrationResult{}, ErrModelKeyMigrationKey
		}
		// Refresh the existing hint without exposing short complete keys.
		hint := "••••"
		chars := []rune(row.value)
		if len(chars) > 4 {
			hint = "…" + string(chars[len(chars)-4:])
		}
		update, err := tx.ExecContext(ctx, `UPDATE llm_profiles
			SET api_key=NULL, api_key_ciphertext=$1, api_key_hint=$2
			WHERE id=$3 AND COALESCE(api_key_ciphertext,'') = ''`, encrypted, hint, row.id)
		if err != nil {
			return ModelKeyMigrationResult{}, err
		}
		n, err := update.RowsAffected()
		if err != nil {
			return ModelKeyMigrationResult{}, err
		}
		if n != 1 {
			return ModelKeyMigrationResult{}, ErrModelKeyMigrationChanged
		}
		result.Migrated++
	}
	if err := tx.Commit(); err != nil {
		return ModelKeyMigrationResult{}, err
	}
	return result, nil
}
