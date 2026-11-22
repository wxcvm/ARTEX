package db

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"

	"github.com/Autumn-27/artex/kb"
	"github.com/google/uuid"
)

var ErrKnowledgeHandoffConflict = errors.New("a document with the same handoff name has different content")
var ErrKnowledgeHandoffScope = errors.New("select an ARTEX task before uploading a Purge handoff")

func insertKnowledgeDocumentTx(ctx context.Context, tx *sql.Tx, taskID int64, doc KnowledgeDocument) (KnowledgeDocument, error) {
	return scanKnowledgeDocument(tx.QueryRowContext(ctx, `INSERT INTO knowledge_documents(id,task_id,title,category,content,source_filename)
VALUES($1,NULLIF($2::bigint,0),$3,$4,$5,$6) RETURNING `+knowledgeDocumentCols,
		uuid.NewString(), taskID, doc.Title, doc.Category, doc.Content, doc.SourceFilename))
}

// SavePurgeHandoffDocument reuses one existing document in this task only.
// The filename carries a digest prefix; comparing full canonical content
// avoids silently accepting an edited document or a prefix collision.
func (d *DB) SavePurgeHandoffDocument(ctx context.Context, taskID int64, doc KnowledgeDocument) (KnowledgeDocument, bool, error) {
	if taskID <= 0 {
		return KnowledgeDocument{}, false, ErrKnowledgeHandoffScope
	}
	handoff, err := kb.ParsePurgeHandoff(doc.SourceFilename, doc.Content)
	if err != nil || handoff == nil {
		return KnowledgeDocument{}, false, kb.ErrInvalidPurgeHandoff
	}
	doc.SourceFilename, doc.Content = handoff.Filename, handoff.Content

	var out KnowledgeDocument
	var reused bool
	err = d.WithEvidenceTx(ctx, func(tx *sql.Tx) error {
		// This existing FOR UPDATE path serializes same-task imports and
		// existing document edits, and preserves archive/read-only fencing.
		if err := LockTaskEvidenceTx(tx, taskID); err != nil {
			return err
		}
		// Older ARTEX releases kept the browser's copy suffix verbatim.
		namePattern := "^" + regexp.QuoteMeta(strings.TrimSuffix(doc.SourceFilename, ".md")) + `( \([1-9][0-9]{0,8}\))?\.md$`
		rows, err := tx.QueryContext(ctx, `SELECT `+knowledgeDocumentCols+` FROM knowledge_documents
WHERE task_id=$1 AND lower(source_filename) ~ $2 ORDER BY created_at,id`, taskID, namePattern)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			prior, err := scanKnowledgeDocument(rows)
			if err != nil {
				return err
			}
			if kb.CanonicalPurgeHandoffText(prior.Content) != doc.Content {
				return ErrKnowledgeHandoffConflict
			}
			if !reused {
				out, reused = prior, true
			}
		}
		if err := rows.Err(); err != nil {
			return err
		}
		if reused {
			return nil
		}
		rows.Close()
		out, err = insertKnowledgeDocumentTx(ctx, tx, taskID, doc)
		return err
	})
	if err != nil {
		return KnowledgeDocument{}, false, err
	}
	return out, reused, nil
}
