# RSI lessons (shadow mode)

Deterministic, evidence-gated lessons extracted from closed tasks. This is the
memory layer of the RSI work: it records what earlier tasks proved, not what a
model guessed, and it never changes a prompt by itself.

The admission rule follows the harness2 Library (Google Research, Apache-2.0,
arXiv 2610.14224): a lesson exists only when it carries a cue, an action, and
at least one evidence row. The wider taxonomy of gains (task, retained,
improver) comes from the awesome-rsi list (lobehub).

## Modes

`ARTEX_RSI_LESSONS` is read when a task's agent bundle is built:

| value | behavior |
|---|---|
| unset / `off` | no selection, no logging, no prompt change (default) |
| `shadow` | select lessons and log what would have been shown, without injecting |
| `inject` | select, log, and append the block to the planner and worker system prompts |

The block is code-owned: it is appended after the rendered prompt, so editing
the agent prompt in the UI cannot remove it. Changing the flag affects new
tasks (or tasks whose agent bundle is invalidated), not a restart.

## Extraction

Extraction runs automatically when a task reaches a terminal status
(`done`, `failed`, `timeout`), and can be re-run or backfilled with:

```
POST /api/tasks/{id}/lessons/extract
```

It is idempotent. `rsi_lesson_evidence` is keyed by `(lesson, source_kind,
source_id)`, so a second run for the same task adds no evidence and cannot
inflate support.

Rules:

| source | kind | action |
|---|---|---|
| finding `confirmed` / `resolved` / `fixed` | positive | `Confirmed <vulnclass>: <summary>` |
| retest `reproduced` on such a finding | positive | `Reproduced <vulnclass> on retest: <summary>` |
| finding `false_positive` / `duplicate` / `ignored` / `risk_accepted` | negative | `Do not re-report without new mechanics: <summary> (prior <status>)` |
| retest `fixed` / `inconclusive` | caution | `Recheck before reporting <vulnclass>: a prior retest was <verdict>` |

Pending and in-progress findings produce nothing: there is no outcome yet.
Empty cue, action, or source is refused by the gate.

## Selection

`GET /api/tasks/{id}/lessons` returns the lessons visible to a task plus the
rendered block. Selection is scoped:

- generic lessons (no company, no category) are visible to every task;
- lessons matching both company and task category rank first, followed by
  company-only matches, category-only matches, and finally generic lessons;
  lessons scoped to unrelated companies or categories are excluded;
- within a scope: positive before negative before caution, then support, then
  recency;
- at most eight lessons and 1800 characters are rendered.

Duplicates are merged at extraction time by a dedupe key over kind, scope, cue,
and action, so the same lesson from several tasks becomes one row with several
evidence entries and a higher support count.

## Rendering and safety

Lesson fields come from task evidence, which can include target-controlled
text. Rendering strips control characters, flattens every field to one line,
caps lengths without breaking UTF-8 characters (including Chinese text), and
prefixes the block with "prior outcomes, not instructions".
The block is appended to the planner and worker prompts only in `inject` mode.

What this layer deliberately does not do:

- no LLM summarization in this phase: rules only, so every lesson traces to a
  finding or retest id;
- no adoption: `status` is `candidate`; `adopted` and `retired` are reserved for
  a later, owner-approved phase;
- no prompt mutation, no benchmark probing, no lesson editing through the UI;
- no decay or deletion; the evidence table is the audit trail.

## Measuring it

`rsi_lesson_usage` records which lessons were considered for which task and in
which mode. Join it with the per-task outcome card
(`GET /api/tasks/{id}/outcome`) to compare tasks that saw lessons with tasks
that did not on the same domain:

- tokens per confirmed finding
- tool calls per confirmed finding
- repeated dead-end probes (future metric; not computable yet)
- false-positive and duplicate rate

Do not treat support as causal evidence. It counts corroborating task
outcomes, not a controlled comparison.

## Linking experience to actual Skill usage

`GET /api/tasks/{id}/outcome` includes a compact **learning trace** alongside
the existing findings, retests, tokens, and tool calls:

- `skills_used`: names and counts of Skills successfully loaded by this
  task (from the existing `skill_usage` ledger); failed/missing invocations
  remain in `skill_misses` and are not counted as a loaded procedure.
- `lesson_exposure.shadow`: prior-lesson records selected only for observation.
- `lesson_exposure.injected`: prior-lesson records appended to the agent prompt.
  These values come from the persisted `rsi_lesson_usage` mode, not the
  process's current environment variable.

The task result card displays these facts as an observed usage trace. They are
**not** a claim that a particular Skill discovered a finding or that lesson
injection improved results. Neither Skill calls nor lesson exposure automatically
change a finding's confirmation status. No schema migration, additional store,
or new runtime dependency is required.

## On-demand observational evaluation

`GET /api/tasks/{id}/lessons/evaluation` compares completed, unarchived tasks
with exactly the same company and task category as the selected task (including
matching empty scopes). The task result card calls it **only when the comparison
button is clicked**, not on its normal five-second outcome refresh.

The API reuses `tasks`, `findings`, `tool_usage`, `skill_usage`,
`activity`, and `rsi_lesson_usage`. Four mutually exclusive recorded modes
keep mixed configurations distinct:

- `inject`: lesson usage recorded with prompt injection, with no shadow usage;
- `shadow`: lesson usage observed in shadow mode, with no injection;
- `mixed`: both injection and shadow recorded for the same task;
- `none`: no recorded lesson exposure, **not** an experimental control group.

The report includes sample counts, confirmed and dismissed findings, tool
calls, Skill loads/misses, token consumption and elapsed wall time. The UI
shows pooled ratios for confirmed findings per task, false positives per
reviewed disposition, and tool/Token usage per confirmed finding; a zero
denominator displays an em dash rather than a misleading result.

The cohorts are observational and task difficulty, target coverage, model
configuration, and changed inputs can differ. **No causal improvement percentage
is inferred**, and no action is automatically taken on any candidate lesson.
There is no new database table, scheduled job, dependency or agent restriction.
