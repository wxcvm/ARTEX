#!/usr/bin/env python3
"""Local-only engineering regression. No ARTEX server, target, model key or DB needed."""
import argparse
import collections
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import time
import uuid

ROOT = Path(__file__).resolve().parents[2]
VERSION = "webapi-engineering-v1"
CASES = {"normal", "pagination", "authentication", "rate_limit", "timeout",
         "malformed_json", "large_utf8", "cancel", "duplicates", "incomplete"}
MARKER = "ARTEX_ENGINEERING_METRIC "


def digest(data):
    return hashlib.sha256(data).hexdigest()


def source_identity():
    def git(*args):
        return subprocess.check_output(["git", *args], cwd=ROOT)
    paths = sorted(set(git("ls-files", "-z", "--cached", "--others", "--exclude-standard").split(b"\0")) - {b""})
    tree = hashlib.sha256()
    missing = []
    for raw in paths:
        path = ROOT / os.fsdecode(raw)
        tree.update(raw + b"\0")
        if path.is_file():
            tree.update(digest(path.read_bytes()).encode())
        else:
            tree.update(b"MISSING")
            missing.append(os.fsdecode(raw))
    return {"commit": git("rev-parse", "HEAD").decode().strip(),
            "dirty": bool(git("status", "--porcelain")),
            "working_tree_sha256": tree.hexdigest(), "missing_tracked_files": missing}


def verify_evidence(metric, evidence_root):
    records = metric.get("evidence", [])
    if not records or len(records) != metric.get("http_calls"):
        return False
    seen = set()
    for record in records:
        path = Path(record.get("file", "")).resolve()
        if not path.is_relative_to(evidence_root.resolve()) or not path.is_file() or path in seen:
            return False
        seen.add(path)
        if digest(path.read_bytes()) != record.get("sha256"):
            return False
    return True


def parse_go(output, repeat, evidence_root):
    metrics, events, failures = [], collections.Counter(), []
    buffers = collections.defaultdict(str)
    for line in output.splitlines():
        event = json.loads(line)
        action, name = event.get("Action"), event.get("Test", "")
        if action in {"pass", "fail", "skip"} and name:
            events[action] += 1
            if action != "pass": failures.append({"test": name, "status": action})
        # test2json can split a single t.Log line across several Output events.
        # Reassemble by test identity before decoding the embedded JSON metric.
        key = (event.get("Package", ""), name)
        buffers[key] += event.get("Output", "")
        while "\n" in buffers[key]:
            text, buffers[key] = buffers[key].split("\n", 1)
            if MARKER in text:
                metric = json.loads(text.split(MARKER, 1)[1])
                metric["evidence_verified"] = verify_evidence(metric, evidence_root)
                metrics.append(metric)
    if any(MARKER in tail for tail in buffers.values()):
        raise ValueError("unterminated engineering metric in Go output")
    counts = collections.Counter(m["case"] for m in metrics)
    complete = counts == collections.Counter({name: repeat for name in CASES})
    complete = complete and all(m.get("case_version") == VERSION and m["evidence_verified"] for m in metrics)
    # A skipped or missing lifecycle test is not acceptance. Require all repetitions.
    lifecycle = sum(1 for line in output.splitlines() if (e := json.loads(line)).get("Action") == "pass"
                    and e.get("Test", "").startswith("TestTaskTransition") and "/" not in e["Test"])
    complete = complete and lifecycle >= 8 * repeat
    return {"cases": metrics, "tests": {name: events[name] for name in ("pass", "fail", "skip")}, "nonpassing": failures,
            "evidence_complete": complete, "lifecycle_passes": lifecycle}


def run(command, env, log, timeout=900):
    started = time.monotonic()
    try:
        result = subprocess.run(command, cwd=ROOT, env=env, capture_output=True, timeout=timeout)
        stdout = result.stdout.decode("utf-8", errors="replace")
        log.write_bytes(result.stdout)
        log.with_suffix(".stderr.log").write_bytes(result.stderr)
        return {"command": command, "exit_code": result.returncode, "elapsed_seconds": time.monotonic()-started,
                "log": str(log), "log_sha256": digest(result.stdout)}, stdout
    except (OSError, subprocess.TimeoutExpired) as error:
        return {"command": command, "exit_code": None, "elapsed_seconds": time.monotonic()-started,
                "environment_error": str(error)}, ""


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repeat", type=int, default=20)
    parser.add_argument("--go", default="go")
    parser.add_argument("--node", default="node")
    parser.add_argument("--race", action="store_true")
    parser.add_argument("--output", type=Path, default=ROOT / "benchmark-results" / "engineering")
    args = parser.parse_args()
    if args.repeat < 1: parser.error("--repeat must be positive")
    out = args.output.resolve() / (time.strftime("%Y%m%dT%H%M%S") + "-" + uuid.uuid4().hex[:8])
    evidence = out / "evidence"
    evidence.mkdir(parents=True)
    report = {"schema_version": 1, "case_version": VERSION, "kind": "deterministic_engineering",
              "source": source_identity(), "repeat": args.repeat, "usage_kind": "simulated",
              "real_model_calls": 0, "discovery_score": None, "commands": []}
    env = dict(os.environ, ARTEX_ENGINEERING_EVIDENCE_DIR=str(evidence))
    # Keep this suite DB-independent even if the developer has a DSN configured.
    env.pop("ARTEX_PG_DSN", None)
    command = [args.go, "test", "-json", "-count=" + str(args.repeat), "-timeout=300s"]
    if args.race: command.append("-race")
    command += ["./server", "-run", "^(TestWebAPIEngineering|TestTaskTransition)"]
    result, stdout = run(command, env, out / "go.jsonl")
    report["commands"].append(result)
    try:
        report["go"] = parse_go(stdout, args.repeat, evidence)
    except (ValueError, KeyError, OSError) as error:
        report["go"] = {"evidence_complete": False, "parse_error": str(error)}
    ui_pass = True
    for iteration in range(args.repeat):
        command = [args.node, "--experimental-strip-types", "--test", "--test-reporter=tap",
                   "web/src/lib/task-outcome.test.mjs"]
        result, stdout = run(command, env, out / f"ui-{iteration+1}.tap")
        report["commands"].append(result)
        counts = dict(re.findall(r"^# (tests|pass|fail|skipped) (\d+)\s*$", stdout, re.M))
        result["tests"] = {key: int(value) for key, value in counts.items()}
        if result["exit_code"] != 0 or int(counts.get("pass", "0")) < 18 or int(counts.get("skipped", "0")):
            ui_pass = False
    passed = (all(c["exit_code"] == 0 for c in report["commands"]) and ui_pass
              and report["go"].get("evidence_complete", False)
              and not report["go"].get("nonpassing"))
    report["passed"] = passed
    report["execution_status"] = ("environment_error" if any("environment_error" in c for c in report["commands"])
                                  else "completed" if passed else "incomplete")
    report["elapsed_seconds"] = sum(c["elapsed_seconds"] for c in report["commands"])
    (out / "results.json").write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding="utf-8")
    metrics = report["go"].get("cases", [])
    lines = ["# ARTEX deterministic engineering regression", "",
             f"- Source: `{report['source']['commit']}`; dirty: `{report['source']['dirty']}`",
             f"- Working tree SHA256: `{report['source']['working_tree_sha256']}`",
             f"- Cases: `{VERSION}`; repetitions: {args.repeat}; status: **{report['execution_status']}**",
             f"- Passed: **{passed}**; elapsed: {report['elapsed_seconds']:.3f}s",
             f"- Evidence complete: {report['go'].get('evidence_complete', False)}",
             f"- Simulated model calls: {sum(m['model_calls'] for m in metrics)}; real model calls: 0",
             f"- HTTP calls: {sum(m['http_calls'] for m in metrics)}; tool calls: {sum(m['tool_calls'] for m in metrics)}",
             "- No vulnerability-discovery score. Assessment statuses below belong to synthetic fixtures, not this run.",
             "", "| Case | Records | Synthetic assessment status | Evidence |", "| --- | ---: | --- | --- |"]
    for name in sorted(CASES):
        selected = [m for m in metrics if m["case"] == name]
        lines.append(f"| {name} | {len(selected)}/{args.repeat} | {', '.join(sorted({m['assessment_status'] for m in selected})) or 'unknown'} | {bool(selected) and all(m['evidence_verified'] for m in selected)} |")
    lines += ["", "See results.json for pass/fail/skip, command exit codes, per-case durations, evidence references and hashes."]
    (out / "report.md").write_text("\n".join(lines)+"\n", encoding="utf-8")
    print(json.dumps({"passed": passed, "status": report["execution_status"], "report": str(out / "report.md")}))
    return 0 if passed else 1


if __name__ == "__main__":
    raise SystemExit(main())
