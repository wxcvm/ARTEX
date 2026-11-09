import { readFileSync } from "node:fs";
import { resolve } from "node:path";

import { UserGuide } from "@/components/user-guide";

// Next.js renders this page at build time for static export. The same source
// file also powers the README link: no duplicated guide to maintain.
export default function HelpPage() {
  const markdown = readFileSync(resolve(process.cwd(), "../docs/USER_GUIDE_ZH.md"), "utf8");
  return <UserGuide markdown={markdown} />;
}
