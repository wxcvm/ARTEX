import assert from "node:assert/strict";
import test from "node:test";
import { skillLabel, toolLabel, knowledgeCategory, knowledgeKind } from "./resource-labels.ts";

test("shipped resources have Chinese display labels without rewriting their identifiers", () => {
  assert.equal(skillLabel("playwright-cli"), "Playwright 浏览器自动化");
  assert.equal(toolLabel("kb_search"), "检索知识库");
  assert.equal(knowledgeCategory("info-leak"), "信息泄露");
  assert.equal(knowledgeKind("script"), "脚本");
});

test("unknown keys and custom tools retain their original names", () => {
  assert.equal(skillLabel("my-custom-skill"), "my-custom-skill");
  assert.equal(toolLabel("kb_search", false), "kb_search");
  assert.equal(toolLabel("custom_tool"), "custom_tool");
  assert.equal(knowledgeCategory("new-category"), "new-category");
  assert.equal(knowledgeKind(""), "");
  assert.equal(skillLabel("constructor"), "constructor");
});
