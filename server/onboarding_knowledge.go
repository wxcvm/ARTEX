package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Autumn-27/artex/agent"
)

// Uses the same read-only handle as the agent tools. Reading never creates a
// task, executes material, or marks an imported claim as independently verified.
func (s *Server) onboardingKnowledge(w http.ResponseWriter, r *http.Request) {
	k := agent.Knowledge()
	if k == nil {
		writeErr(w, http.StatusServiceUnavailable, "知识库尚未加载。桌面版请保留完整 knowledge 目录；外部部署请配置 knowledge_dir 或 ARTEX_KB_DIR 后重启。")
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	if idText := r.URL.Query().Get("id"); idText != "" {
		id, err := strconv.ParseInt(idText, 10, 64)
		if err != nil || id < 1 {
			writeErr(w, 400, "条目编号无效")
			return
		}
		entry, err := k.Get(id)
		if err != nil {
			writeErr(w, 500, "知识库正文读取失败")
			return
		}
		if entry == nil {
			writeErr(w, 404, "条目不存在")
			return
		}
		writeJSON(w, http.StatusOK, entry)
		return
	}
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" || len(query) > 1000 {
		writeErr(w, 400, "请填写关键词（最多 1000 字节）")
		return
	}
	hits, degraded, err := k.Search(query, 30)
	if err != nil {
		writeErr(w, 500, "知识库检索失败，请检查知识库文件")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"hits": hits, "degraded": degraded})
}
