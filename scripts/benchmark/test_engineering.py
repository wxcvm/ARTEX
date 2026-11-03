import json
from pathlib import Path
import tempfile
import unittest

from run_engineering import CASES, MARKER, VERSION, digest, parse_go, verify_evidence


class EngineeringReportTests(unittest.TestCase):
    def test_missing_tests_are_not_success(self):
        self.assertFalse(parse_go("", 20, Path("."))["evidence_complete"])

    def test_evidence_hash_and_path_are_verified(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            path = root / "evidence.json"
            path.write_bytes(b'{"fixture":true}')
            metric = {"http_calls": 1, "evidence": [{"file": str(path), "sha256": digest(path.read_bytes())}]}
            self.assertTrue(verify_evidence(metric, root))
            path.write_bytes(b"changed")
            self.assertFalse(verify_evidence(metric, root))
            self.assertFalse(verify_evidence(metric, root / "nested"))

    def test_missing_duplicate_and_skipped_cases_are_visible(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            path = root / "evidence.json"
            path.write_bytes(b"fixture")
            events = [{"Action": "pass", "Test": f"TestTaskTransition{i}"} for i in range(8)]
            for name in CASES:
                metric = {"case": name, "case_version": VERSION, "http_calls": 1,
                          "evidence": [{"file": str(path), "sha256": digest(path.read_bytes())}]}
                events.append({"Action": "output", "Output": MARKER + json.dumps(metric) + "\n"})
            encode = lambda rows: "\n".join(json.dumps(e) for e in rows)
            self.assertTrue(parse_go(encode(events), 1, root)["evidence_complete"])
            split = events.copy()
            last = split.pop()["Output"]
            split += [{"Action": "output", "Output": last[:77]}, {"Action": "output", "Output": last[77:]}]
            self.assertTrue(parse_go(encode(split), 1, root)["evidence_complete"])
            self.assertFalse(parse_go(encode(events[:-1]), 1, root)["evidence_complete"])
            self.assertFalse(parse_go(encode(events + [events[-1]]), 1, root)["evidence_complete"])
            events.append({"Action": "skip", "Test": "TestTaskTransitionSkipped"})
            self.assertEqual(parse_go(encode(events), 1, root)["nonpassing"][0]["status"], "skip")


if __name__ == "__main__":
    unittest.main()
