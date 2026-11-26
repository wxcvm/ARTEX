"""Exercise backup protection using a disposable fake Docker executable."""
import gzip
import os
from pathlib import Path
import shutil
import stat
import subprocess
import tempfile
import unittest


class BackupTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        (self.root / "scripts").mkdir()
        shutil.copyfile(Path(__file__).with_name("backup.sh"),
                        self.root / "scripts" / "backup.sh")
        (self.root / "bin").mkdir()
        docker = self.root / "bin" / "docker"
        docker.write_text("#!/bin/sh\nprintf 'synthetic dump\\n'\n")
        docker.chmod(0o700)
        (self.root / "data").mkdir()
        (self.root / "data" / "fixture.txt").write_text("synthetic")
        self.env = dict(os.environ)
        self.env["PATH"] = str(self.root / "bin") + os.pathsep + self.env["PATH"]

    def backup(self, *args):
        return subprocess.run(
            ["bash", "-c", 'umask 022; exec bash scripts/backup.sh "$@"', "backup", *args],
            cwd=self.root, env=self.env, capture_output=True, text=True)

    def test_private_outputs_and_positive_retention(self):
        result = self.backup("--with-data", "--keep", "1")
        self.assertEqual(result.returncode, 0, result.stderr)
        out = self.root / "backups"
        self.assertEqual(stat.S_IMODE(out.stat().st_mode), 0o700)
        files = sorted(out.iterdir())
        self.assertEqual(len(files), 2)
        for f in files:
            self.assertEqual(stat.S_IMODE(f.stat().st_mode), 0o600)
        dump = next(out.glob("artex-db-*.sql.gz"))
        with gzip.open(dump, "rt") as stream:
            self.assertEqual(stream.read(), "synthetic dump\n")

    def test_invalid_retention_does_not_touch_existing_backup(self):
        out = self.root / "backups"
        out.mkdir()
        existing = out / "artex-db-existing.sql.gz"
        existing.write_bytes(b"fixture")
        for keep in ("0", "00", "-1", "08", "9999999999999999999999", "bad"):
            with self.subTest(keep=keep):
                result = self.backup("--keep", keep)
                self.assertEqual(result.returncode, 2, result.stderr)
                self.assertEqual(list(out.iterdir()), [existing])
                self.assertEqual(existing.read_bytes(), b"fixture")


if __name__ == "__main__":
    unittest.main()
