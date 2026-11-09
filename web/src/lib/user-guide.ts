// The repository Markdown file is the single source for both GitHub docs and /help.
// Keep this module free of React and Node APIs so it can be tested with node --test.

export interface GuideSection {
  id: string;
  title: string;
  content: string;
}

export function guideSlug(title: string): string {
  // Match GitHub's heading anchors for this guide (Chinese text is preserved).
  return title.toLowerCase().trim()
    .replace(/[^\p{L}\p{N}\s_-]/gu, "")
    .replace(/\s/g, "-");
}

export function parseGuide(markdown: string): { introduction: string; sections: GuideSection[] } {
  const headers = [...markdown.matchAll(/^## (.+)\r?$/gm)];
  const introduction = markdown.slice(0, headers[0]?.index ?? markdown.length)
    .replace(/^# [^\n]*\r?\n/, "").trim();
  const sections = headers.map((header, index) => {
    const title = header[1].trim();
    const content = markdown.slice(
      header.index + header[0].length,
      headers[index + 1]?.index ?? markdown.length,
    ).trim();
    return { id: guideSlug(title), title, content };
  });
  return { introduction, sections };
}

// Resolve links relative to docs/ in the source repository, not /help in
// the running frontend. GitHub URLs open the source; #anchors stay in-app.
export function guideLinks(markdown: string): string {
  return markdown.replace(/(!?\[[^\]\n]*\]\()([^)]+)(\))/g, (whole, start: string, link: string, end: string) => {
    if (/^(?:[a-z][a-z\d+.-]*:|#|\/)/i.test(link)) return whole;
    const [relativePath, fragment] = link.split("#", 2);
    const path = ["docs"];
    for (const component of relativePath.split("/")) {
      if (component === "." || component === "") continue;
      if (component === "..") path.pop();
      else path.push(component);
    }
    const target = path.join("/");
    const image = start.startsWith("!");
    const base = image
      ? "https://raw.githubusercontent.com/wxcvm/ARTEX/main/"
      : "https://github.com/wxcvm/ARTEX/blob/main/";
    return start + base + target + (fragment ? "#" + fragment : "") + end;
  });
}
