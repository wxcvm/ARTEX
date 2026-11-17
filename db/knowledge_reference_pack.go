package db

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

//go:embed knowledgepacks/cyberstrike-csf-v1.json
var cyberStrikeReferencePackJSON []byte

//go:embed knowledgepacks/cyberstrike-csf-supplement-v1.json
var cyberStrikeSupplementPackJSON []byte

//go:embed knowledgepacks/cyberstrike-csf-third-v1.json
var cyberStrikeThirdPackJSON []byte

//go:embed knowledgepacks/cyberstrike-csf-fourth-v1.json
var cyberStrikeFourthPackJSON []byte

type knowledgeReferencePack struct {
	ID string `json:"id"`
	Framework string `json:"framework"`
	Authority string `json:"authority"`
	Upstream string `json:"upstream"`
	Revision string `json:"revision"`
	UpstreamLicense string `json:"upstream_license"`
	Category string `json:"category"`
	Entries []knowledgeReferenceEntry `json:"entries"`
}

type knowledgeReferenceEntry struct {
	Code string `json:"code"`
	Title string `json:"title"`
	SourcePath string `json:"source_path"`
	SourceBlob string `json:"source_blob"`
	Keywords string `json:"keywords"`
	Body string `json:"body"`
}

func referencePackMarker(pack knowledgeReferencePack) string {
	return "migration.knowledge_reference_pack." + pack.ID
}

func referenceDocumentID(pack knowledgeReferencePack, entry knowledgeReferenceEntry) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte("artex:knowledge-reference:" + pack.ID + ":" + entry.Code)).String()
}

func referenceDocumentContent(pack knowledgeReferencePack, entry knowledgeReferenceEntry) string {
	return fmt.Sprintf("%s\n\n关键词 / Keywords: %s\n\n---\n参考资料 / Reference only — 不代表目标已通过验证或存在漏洞。\n框架: %s (%s)\n来源索引: %s/blob/%s/%s\n源文件 blob: %s\n上游仓库许可: %s；本文为 ARTEX 独立整理的核查参考，不复制上游技能正文或自动执行指令。\n",
		entry.Body, entry.Keywords, pack.Framework, pack.Authority,
		pack.Upstream, pack.Revision, entry.SourcePath, entry.SourceBlob, pack.UpstreamLicense)
}

// SeedCyberStrikeReferencePack installs the shipped reference material once
// into the existing global document library. Existing edits and deletions
// remain authoritative; starting ARTEX again never reinstalls the pack.
func (d *DB) SeedCyberStrikeReferencePack() error {
	for _, raw := range [][]byte{cyberStrikeReferencePackJSON, cyberStrikeSupplementPackJSON, cyberStrikeThirdPackJSON, cyberStrikeFourthPackJSON} {
		var pack knowledgeReferencePack
		if err := json.Unmarshal(raw, &pack); err != nil {
			return fmt.Errorf("decode CyberStrike reference pack: %w", err)
		}
		if err := d.seedKnowledgeReferencePack(pack); err != nil {
			return fmt.Errorf("seed CyberStrike reference pack %s: %w", pack.ID, err)
		}
	}
	return nil
}

func (d *DB) seedKnowledgeReferencePack(pack knowledgeReferencePack) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.Exec(`INSERT INTO settings(key,value) VALUES($1,'true') ON CONFLICT (key) DO NOTHING`, referencePackMarker(pack))
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 0 {
		for _, entry := range pack.Entries {
			// A stable identity and DO NOTHING protect preexisting documents;
			// the marker and all new rows commit as one transaction.
			if _, err := tx.Exec(`INSERT INTO knowledge_documents(id,task_id,title,category,content,source_filename)
VALUES($1,NULL,$2,$3,$4,$5) ON CONFLICT (id) DO NOTHING`,
				referenceDocumentID(pack, entry), entry.Title, pack.Category,
				referenceDocumentContent(pack, entry), entry.SourcePath); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}
