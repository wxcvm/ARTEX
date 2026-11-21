package db

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

// SARIFSourceID is an opaque, non-secret identifier for this database's
// reports. It survives exports and restarts; it is not proof of authorization.
func (d *DB) SARIFSourceID() (string, error) {
	return storedSARIFSourceID(d)
}

type sarifSourceSettings interface {
	GetSetting(string) (string, bool, error)
	InsertSettingIfAbsent(string, string) (bool, error)
}

const sarifSourceKey = "report.sarif_source_id"

func storedSARIFSourceID(store sarifSourceSettings) (string, error) {
	id, found, err := store.GetSetting(sarifSourceKey)
	if err != nil {
		return "", err
	}
	if !found {
		var bytes [16]byte
		if _, err := rand.Read(bytes[:]); err != nil {
			return "", err
		}
		if _, err := store.InsertSettingIfAbsent(sarifSourceKey, hex.EncodeToString(bytes[:])); err != nil {
			return "", err
		}
		// Two exporters may race. Use the persisted winner in both reports.
		id, found, err = store.GetSetting(sarifSourceKey)
		if err != nil {
			return "", err
		}
	}
	decoded, err := hex.DecodeString(id)
	if !found || err != nil || len(decoded) != 16 {
		return "", errors.New("invalid SARIF report source identifier")
	}
	return hex.EncodeToString(decoded), nil
}
