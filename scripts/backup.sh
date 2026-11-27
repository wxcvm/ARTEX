#!/usr/bin/env bash
# Backup the ARTEX database (and optionally the data directory) with rotation.
#
# Usage:
#   scripts/backup.sh [--with-data] [--keep N] [--output-dir DIR]
#
#   --with-data    also archive the bind-mounted data/ directory (worker
#                  artifacts, evidence, skills state) as a tar.gz
#   --keep N       number of db dumps to keep (default 14; data archives rotate
#                  to the same count)
#   --output-dir   where to write (default backups/, relative to the repo root)
#
# The database dump runs inside the postgres compose service, so no psql client
# is needed on the host. For isolated, empty-database restore and separately
# protected configuration/key backups, see the README Backups section.
set -euo pipefail
cd "$(dirname "$0")/.."

usage() { sed -n '2,15p' "$0" | sed 's/^# \?//'; }

KEEP=14
OUTPUT_DIR="backups"
WITH_DATA=0
while [ $# -gt 0 ]; do
  case "$1" in
    --keep) KEEP="${2:?}"; shift 2 ;;
    --output-dir) OUTPUT_DIR="${2:?}"; shift 2 ;;
    --with-data) WITH_DATA=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
done
if ! [[ "$KEEP" =~ ^[1-9][0-9]{0,8}$ ]]; then
  echo "--keep must be an integer between 1 and 999999999 (without leading zeros)" >&2
  exit 2
fi

# Compose reads .env for credentials; source it for the dump command too.
if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  . ./.env
  set +a
fi
PGUSER_NAME="${POSTGRES_USER:-artex}"
PGDATABASE_NAME="${POSTGRES_DB:-artex}"

# Dumps may include credentials and private task data. Restrict new files
# and directories even when the caller has a permissive umask.
umask 077
mkdir -p "$OUTPUT_DIR"
case "$OUTPUT_DIR" in /*) ;; *) OUTPUT_DIR="./$OUTPUT_DIR" ;; esac
TS=$(date -u +%Y%m%dT%H%M%SZ)
DB_OUT="$OUTPUT_DIR/artex-db-$TS.sql.gz"
# Stage privately on the destination filesystem. Publish only validated files,
# and never overwrite a completed backup with the same timestamp.
STAGE=$(mktemp -d "$OUTPUT_DIR/.artex-backup.XXXXXXXX")
trap 'rm -rf -- "$STAGE"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
DB_TMP="$STAGE/database.sql.gz"
DATA_OUT="$OUTPUT_DIR/artex-data-$TS.tar.gz"
DATA_TMP="$STAGE/data.tar.gz"
if [ -e "$DB_OUT" ] || { [ "$WITH_DATA" -eq 1 ] && [ -e "$DATA_OUT" ]; }; then
  echo "refusing to overwrite an existing backup for $TS" >&2
  exit 1
fi

if ! docker compose exec -T postgres pg_dump -U "$PGUSER_NAME" "$PGDATABASE_NAME" | gzip > "$DB_TMP"; then
  echo "pg_dump failed (is the stack running?)" >&2
  exit 1
fi
gzip -t "$DB_TMP"

if [ "$WITH_DATA" -eq 1 ]; then
  if ! tar -czf "$DATA_TMP" data 2>"$STAGE/tar.err"; then
    echo "data archive failed (check file access and whether data is changing):" >&2
    head -5 "$STAGE/tar.err" >&2
    exit 1
  fi
  gzip -t "$DATA_TMP"
fi

# A same-filesystem hard link publishes a complete file and fails if a file
# already occupies the destination. Avoid GNU-only ln -T so BSD/macOS work too.
publish_backup() {
  local source="$1" destination="$2" nested
  if [ -e "$destination" ] || [ -L "$destination" ]; then
    echo "refusing to overwrite existing backup: $destination" >&2
    return 1
  fi
  if ! ln "$source" "$destination"; then
    return 1
  fi
  # BSD ln treats an existing directory as a target directory. Detect that
  # case (including a raced-in symlink to a directory) and remove only the
  # hard link we may have created inside it.
  if [ -d "$destination" ] || [ -L "$destination" ] ||
     [ ! -f "$destination" ] || [ ! "$source" -ef "$destination" ]; then
    nested="$destination/$(basename "$source")"
    if [ -f "$nested" ] && [ "$source" -ef "$nested" ]; then
      rm -f -- "$nested"
    fi
    echo "backup destination is not a regular file: $destination" >&2
    return 1
  fi
}
if ! publish_backup "$DB_TMP" "$DB_OUT"; then
  echo "could not publish database backup" >&2
  exit 1
fi
if [ "$WITH_DATA" -eq 1 ] && ! publish_backup "$DATA_TMP" "$DATA_OUT"; then
  rm -f -- "$DB_OUT"
  echo "could not publish data backup" >&2
  exit 1
fi
echo "wrote $DB_OUT ($(du -h "$DB_OUT" | cut -f1))"
if [ "$WITH_DATA" -eq 1 ]; then
  echo "wrote $DATA_OUT ($(du -h "$DATA_OUT" | cut -f1))"
fi

prune() {
  ls -1t "$OUTPUT_DIR"/$1 2>/dev/null | tail -n +$((KEEP + 1)) | while IFS= read -r old; do
    echo "pruning $old"
    rm -f "$old"
  done
}
prune 'artex-db-*.sql.gz'
[ "$WITH_DATA" -eq 1 ] && prune 'artex-data-*.tar.gz'
echo "backup complete ($KEEP newest kept)"
