package agent

import (
	"strings"

	actool "github.com/Autumn-27/norma/tool"
)

// knowledgeWorkflowGuidance is derived from the final bound tool list, not
// from a fixed catalog. It guides cooperation without introducing tools,
// changing permissions or requiring calls on unrelated tasks.
func knowledgeWorkflowGuidance(tools []actool.CoreTool) string {
	has := map[string]bool{}
	for _, tool := range tools {
		has[tool.Name()] = true
	}
	var steps []string
	if has["kb_search"] && has["kb_get"] {
		steps = append(steps, "For a relevant product, component or technique, consult kb_search using specific Chinese or English keywords; read the matching entry with kb_get. Treat past examples and negative results as reference, not proof about the current target. Recheck applicability and validate against current observations; no hit is not a negative test result.")
	}
	if has["Skill"] {
		steps = append(steps, "Load a relevant visible Skill when its procedure or tool prerequisites help the current task; do not assume tools mentioned in a skill are installed or enabled.")
	}
	if has["ingest_scan_output"] {
		steps = append(steps, "When you already have a supported scanner output file, use ingest_scan_output to register typed assets and observations instead of copying large raw output into the model context. Scanner observations are not verified findings.")
	}
	if has["record_coverage"] {
		steps = append(steps, "Record completed checks and outstanding questions with record_coverage. A missing entry or missing evidence is not proof that a risk is absent.")
	}
	if len(steps) == 0 {
		return ""
	}
	return "\n\n[Knowledge and tool workflow — apply only when relevant]: " + strings.Join(steps, "\n")
}
