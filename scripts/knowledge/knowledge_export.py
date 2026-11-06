#!/usr/bin/env python3
"""Export a read-only knowledge.db + tar into a provenance-preserving preview.
No scripts are run, files extracted, historical task tables read, or network used.
"""
import argparse
from contextlib import closing
import hashlib
import json
from pathlib import Path
import sqlite3
import tarfile

MAX_ENTRIES = 10000
MAX_MATERIAL = 2 << 20
MAX_JSON = 32 << 20
TEXT_FIELDS = ('code', 'title', 'kind', 'category', 'component', 'versions', 'cve',
               'tags', 'source', 'source_url', 'description', 'usage', 'content', 'path')


def sha(data):
    return hashlib.sha256(data).hexdigest()


def file_sha(path):
    with path.open('rb') as f:
        return hashlib.file_digest(f, 'sha256').hexdigest()


def material_path(raw):
    parts = raw.replace('\\', '/').split('/')
    try:
        parts = parts[parts.index('pocs'):]
    except ValueError:
        return ''
    if len(parts) < 2 or any(p in ('', '.', '..') or ':' in p or '\x00' in p for p in parts):
        return ''
    return '/'.join(parts)


def archive_path(raw):
    # Archive member names are relative POSIX names, never host filesystem paths.
    parts = raw.split('/')
    if any(p in ('', '.', '..') or '\\' in p or ':' in p or '\x00' in p for p in parts):
        return ''
    for i in range(len(parts) - 2):
        if parts[i:i + 2] == ['redteam', 'pocs']:
            return '/'.join(parts[i + 1:])
    return ''


def export(database, archive, output):
    database, archive, output = map(lambda p: Path(p).resolve(), (database, archive, output))
    if output in (database, archive) or output.exists():
        raise ValueError('output must be a new file distinct from both inputs')
    for suffix in ('-wal', '-journal'):
        sidecar = Path(str(database) + suffix)
        if sidecar.exists() and sidecar.stat().st_size:
            raise ValueError('export requires a closed snapshot without a pending journal/WAL')
    before = file_sha(database)
    entries = []
    total_bytes = 0
    with closing(sqlite3.connect(database.as_uri() + '?mode=ro&immutable=1', uri=True)) as con:
        con.execute('PRAGMA query_only=ON')
        con.execute('PRAGMA trusted_schema=OFF')
        con.setlimit(sqlite3.SQLITE_LIMIT_LENGTH, MAX_JSON)
        schema = con.execute("SELECT type,sql FROM sqlite_master WHERE name='poc'").fetchone()
        if not schema or schema[0] != 'table' or 'VIRTUAL TABLE' in schema[1].upper():
            raise ValueError('poc must be an ordinary table')
        columns = {r[1] for r in con.execute('PRAGMA table_info(poc)')}
        if not set(TEXT_FIELDS + ('verified',)).issubset(columns):
            raise ValueError('missing required poc columns')
        seen = set()
        fields = TEXT_FIELDS + ('verified',)
        for row in con.execute('SELECT ' + ','.join(fields) + ' FROM poc LIMIT ?', (MAX_ENTRIES + 1,)):
            if len(entries) >= MAX_ENTRIES:
                raise ValueError('too many entries')
            values = dict(zip(fields, row))
            for field in TEXT_FIELDS:
                value = values[field]
                if value is None:
                    value = ''
                if not isinstance(value, str):
                    raise ValueError('non-text material field')
                size = len(value.encode('utf-8', errors='strict'))
                if size > MAX_MATERIAL:
                    raise ValueError('material exceeds 2 MiB')
                total_bytes += size
                values[field] = value
            if total_bytes > MAX_JSON:
                raise ValueError('catalog exceeds 32 MiB')
            code = values['code']
            if not code or code in seen or not values['title'].strip():
                raise ValueError('missing/duplicate code or missing title')
            seen.add(code)
            verified = values.pop('verified')
            if verified not in (None, 0, 1):
                raise ValueError('invalid source verified flag')
            raw_path = values.pop('path')
            path = material_path(raw_path)
            values.update(id=sha((before + '\0' + code).encode())[:32],
                          source_verified=verified == 1,
                          content_sha256=sha(values['content'].encode()),
                          attachment={'path': path, 'status': 'missing' if raw_path else 'none',
                                      'content': '', 'sha256': ''})
            entries.append(values)
    wanted = {e['attachment']['path'] for e in entries if e['attachment']['path']}
    matches = {}
    # Streaming avoids extraction, tar path writes and loading unrelated histories.
    with tarfile.open(archive, 'r|*') as tf:
        for member in tf:
            path = archive_path(member.name)
            if path not in wanted:
                continue
            if path in matches or not member.isfile():
                matches[path] = None  # ambiguous or non-regular; never pick an arbitrary body
                continue
            if member.size > MAX_MATERIAL or member.size < 0:
                raise ValueError('attachment exceeds 2 MiB')
            with tf.extractfile(member) as f:
                raw = f.read(MAX_MATERIAL + 1)
            if len(raw) != member.size:
                raise ValueError('incomplete attachment')
            body = raw.decode('utf-8', errors='strict')
            total_bytes += len(raw)
            if total_bytes > MAX_JSON:
                raise ValueError('catalog exceeds 32 MiB')
            matches[path] = (body, sha(raw))
    summary = {'total': len(entries), 'conflicts': 0, 'missing': 0, 'with_content': 0}
    for e in entries:
        a = e['attachment']
        if a['path'] in matches:
            match = matches[a['path']]
            if match is None:
                a['status'] = 'ambiguous'
            else:
                a['content'], a['sha256'] = match
                a['status'] = 'match' if a['content'] == e['content'] else 'conflict'
        summary['conflicts'] += a['status'] == 'conflict'
        summary['missing'] += a['status'] in ('missing', 'ambiguous')
        summary['with_content'] += bool(e['content'])
    result = {'version': 1, 'source': {'name': database.name, 'sha256': before},
              'summary': summary, 'entries': sorted(entries, key=lambda e: e['code'])}
    payload = (json.dumps(result, ensure_ascii=False, sort_keys=True, separators=(',', ':')) + '\n').encode('utf-8')
    if len(payload) > MAX_JSON:
        raise ValueError('encoded catalog exceeds 32 MiB')
    if file_sha(database) != before:
        raise ValueError('database changed during export')
    # All validation finishes before exclusive creation. Never overwrite an input/output.
    with output.open('xb') as f:
        try:
            f.write(payload)
        except BaseException:
            f.close()
            output.unlink()
            raise
    return summary


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--database', required=True)
    p.add_argument('--archive', required=True)
    p.add_argument('--output', required=True)
    args = p.parse_args()
    try:
        summary = export(args.database, args.archive, args.output)
    except (OSError, ValueError, sqlite3.Error, tarfile.TarError) as exc:
        p.exit(1, 'Export failed: ' + str(exc) + '\n')
    print(json.dumps(summary))


if __name__ == '__main__':
    main()
