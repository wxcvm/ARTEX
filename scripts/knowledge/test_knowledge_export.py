from contextlib import closing
import io
import json
from pathlib import Path
import sqlite3
import tarfile
import tempfile
import unittest
from unittest.mock import patch
import knowledge_export as ex


class ExportTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        self.db, self.tar, self.out = (self.root / n for n in ('knowledge.db', 'source.tar.gz', 'out.json'))
        with closing(sqlite3.connect(self.db)) as c, c:
            c.execute('CREATE TABLE poc (' + ','.join(f + ' TEXT' for f in ex.TEXT_FIELDS) + ',verified INTEGER,asset_target TEXT)')
        self.archive([])

    def row(self, code='a', body='中文正文', path='C:\\old\\redteam\\pocs\\a.txt'):
        with closing(sqlite3.connect(self.db)) as c, c:
            c.execute('INSERT INTO poc (code,title,content,path,verified,asset_target) VALUES (?,?,?,?,?,?)',
                      (code, '中文测试', body, path, 1, 'private-history-should-not-export'))

    def archive(self, files):
        with tarfile.open(self.tar, 'w:gz') as tf:
            for name, body in files:
                data = body if isinstance(body, bytes) else body.encode()
                info = tarfile.TarInfo(name)
                info.size = len(data)
                tf.addfile(info, io.BytesIO(data))

    def run_export(self, out=None):
        ex.export(self.db, self.tar, out or self.out)
        return json.loads((out or self.out).read_text(encoding='utf8'))

    def test_provenance_determinism_and_inputs(self):
        self.row()
        self.archive([('bundle/redteam/pocs/a.txt', '中文正文')])
        hashes = [ex.file_sha(p) for p in (self.db, self.tar)]
        data = self.run_export()
        self.run_export(self.root / 'second.json')
        self.assertEqual(self.out.read_bytes(), (self.root / 'second.json').read_bytes())
        self.assertEqual(hashes, [ex.file_sha(p) for p in (self.db, self.tar)])
        self.assertEqual(data['entries'][0]['attachment']['status'], 'match')
        self.assertTrue(data['entries'][0]['source_verified'])
        self.assertNotIn(b'private-history', self.out.read_bytes())
        self.assertNotIn('C:\\old', self.out.read_text(encoding='utf8'))

    def test_conflict_missing_none(self):
        self.row('a'); self.row('b', path='pocs/missing.txt'); self.row('c', body='', path=None)
        self.archive([('bundle/redteam/pocs/a.txt', '不同正文')])
        data = self.run_export()
        self.assertEqual(data['summary'], dict(total=3, conflicts=1, missing=1, with_content=2))
        self.assertEqual([e['attachment']['status'] for e in data['entries']], ['conflict','missing','none'])
        self.assertEqual(data['entries'][0]['content'], '中文正文')
        self.assertEqual(data['entries'][0]['attachment']['content'], '不同正文')

    def test_duplicate_tar_ambiguous(self):
        self.row()
        self.archive([('one/redteam/pocs/a.txt','a'),('two/redteam/pocs/a.txt','b')])
        a = self.run_export()['entries'][0]['attachment']
        self.assertEqual(a['status'], 'ambiguous'); self.assertEqual(a['content'], '')

    def test_traversal_and_absolute_members_ignored(self):
        self.row()
        self.archive([('../redteam/pocs/a.txt','a'),('/redteam/pocs/a.txt','b')])
        self.assertEqual(self.run_export()['summary']['missing'], 1)
        self.assertEqual(ex.material_path('pocs/../secret'), '')
        self.assertEqual(ex.material_path('/etc/passwd'), '')

    def test_no_overwrite_or_input_overlap(self):
        self.row()
        for path in (self.db, self.tar):
            before = ex.file_sha(path)
            with self.assertRaises(ValueError): self.run_export(path)
            self.assertEqual(before, ex.file_sha(path))
        self.out.write_text('keep')
        with self.assertRaises(ValueError): self.run_export()
        self.assertEqual(self.out.read_text(), 'keep')

    def test_duplicate_codes_rejected_without_output(self):
        self.row(); self.row()
        with self.assertRaises(ValueError): self.run_export()
        self.assertFalse(self.out.exists())

    def test_invalid_utf8_rejected_without_output(self):
        self.row(); self.archive([('bundle/redteam/pocs/a.txt', b'\xff')])
        with self.assertRaises(UnicodeError): self.run_export()
        self.assertFalse(self.out.exists())

    def test_pending_wal_is_not_silently_ignored(self):
        self.row()
        Path(str(self.db) + '-wal').write_bytes(b'pending fixture')
        with self.assertRaises(ValueError): self.run_export()
        self.assertFalse(self.out.exists())

    def test_schema_and_size_limits(self):
        self.row()
        with patch.object(ex, 'MAX_MATERIAL', 2):
            with self.assertRaises(ValueError): self.run_export()
        with patch.object(ex, 'MAX_ENTRIES', 0):
            with self.assertRaises(ValueError): self.run_export()
        with patch.object(ex, 'MAX_JSON', 1):
            with self.assertRaises((ValueError, sqlite3.Error)): self.run_export()
        with closing(sqlite3.connect(self.db)) as c, c: c.execute('DROP TABLE poc')
        with self.assertRaises(ValueError): self.run_export()
        self.assertFalse(self.out.exists())


if __name__ == '__main__':
    unittest.main()
