package report

import (
	"fmt"
	"strings"

	"github.com/Autumn-27/artex/db"
)

// ScreenshotPath uses numeric indexes instead of user-controlled filenames.
func ScreenshotPath(findingID int64, index int, mime string) string {
	ext := "png"
	if mime == "image/jpeg" {
		ext = "jpg"
	}
	return fmt.Sprintf("evidence/%d/screenshots/screenshot-%d.%s", findingID, index+1, ext)
}

func findingScreenshotsMarkdown(f *db.DBFinding, attachments bool, locale, heading string) string {
	if len(f.Screenshots) == 0 {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "\n%s %s\n\n", heading, reportText(locale, "截图证据", "Screenshot evidence"))
	for i, shot := range f.Screenshots {
		caption := strings.Join(strings.Fields(shot.Caption), " ")
		if caption == "" {
			caption = reportText(locale, "无图注", "No caption")
		}
		fmt.Fprintf(&b, "%d. %s: %s\n", i+1, reportText(locale, "图注", "Caption"), caption)
		if shot.BindingID != "" {
			fmt.Fprintf(&b, "   %s #%s\n", reportText(locale, "关联证据", "Linked evidence"), shot.BindingID)
		}
		if attachments {
			fmt.Fprintf(&b, "   [%s](%s)\n", reportText(locale, "查看截图", "View screenshot"), ScreenshotPath(f.ID, i, shot.MIME))
		}
	}
	if !attachments {
		b.WriteString(reportText(locale, "\n图片可在 Word 或证据 ZIP 导出中查看。\n", "\nImages are included in Word or evidence ZIP exports.\n"))
	}
	b.WriteString("\n")
	return b.String()
}
