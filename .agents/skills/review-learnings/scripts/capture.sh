#!/bin/sh
# 知見台帳に観測を記録する（skill から呼ぶ）。使い方は SKILL.md。

set -eu

slug=""
pr=""
origin=""
id=""
date=""
title=""
rule=""
why=""
scopes=""

die() { printf 'capture: %s\n' "$1" >&2; exit 1; }

while [ $# -gt 0 ]; do
  case "$1" in
    --slug)   slug="$2";   shift 2 ;;
    --pr)     pr="$2";     shift 2 ;;
    --origin) origin="$2"; shift 2 ;;
    --id)     id="$2";     shift 2 ;;
    --date)   date="$2";   shift 2 ;;
    --title)  title="$2";  shift 2 ;;
    --rule)   rule="$2";   shift 2 ;;
    --why)    why="$2";    shift 2 ;;
    --scope)  scopes="$scopes $2"; shift 2 ;;
    *) die "unknown argument: $1" ;;
  esac
done

# 既存台帳との後方互換のため --pr を残す（--origin pr の糖衣）。
if [ -n "$pr" ]; then
  [ -z "$origin" ] || die "--pr and --origin are mutually exclusive"
  origin="pr"
  id="$pr"
fi
[ -n "$origin" ] || die "--pr or --origin is required"

[ -n "$slug" ] || die "--slug is required"
[ -n "$id" ]   || die "--pr or --id is required"
[ -n "$date" ] || die "--date is required"
[ -n "$scopes" ] || die "--scope is required (use --scope general if truly global)"

printf '%s' "$slug" | grep -Eq '^[a-z0-9]+(-[a-z0-9]+)*$' || die "slug must be kebab-case: $slug"
printf '%s' "$date" | grep -Eq '^[0-9]{4}-[0-9]{2}-[0-9]{2}$' || die "date must be YYYY-MM-DD: $date"

case "$origin" in
  pr)
    printf '%s' "$id" | grep -Eq '^[0-9]+$' || die "pr must be a number: $id"
    occ_line="  - { pr: $id, date: $date }"
    ;;
  session)
    printf '%s' "$id" | grep -Eq '^[A-Za-z0-9._-]+$' || die "session id must be safe chars: $id"
    occ_line="  - { origin: session, id: $id, date: $date }"
    ;;
  *) die "origin must be pr or session: $origin" ;;
esac

# CLI 非依存にルートを解決する。Claude は CLAUDE_PROJECT_DIR、他は git トップレベル。
root="${CLAUDE_PROJECT_DIR:-$(git rev-parse --show-toplevel 2>/dev/null || true)}"
[ -n "$root" ] || die "cannot resolve project root (set CLAUDE_PROJECT_DIR or run inside a git repo)"
ledger_dir="$root/docs/learnings"
[ -d "$ledger_dir" ] || die "ledger dir not found: $ledger_dir"

file="$ledger_dir/$slug.md"

if [ -f "$file" ]; then
  dup_re=""
  case "$origin" in
    pr)
      dup_re="^[[:space:]]*-[[:space:]]*\{[[:space:]]*pr:[[:space:]]*$id,[[:space:]]*date:[[:space:]]*$date[[:space:]]*\}"
      ;;
    session)
      dup_re="^[[:space:]]*-[[:space:]]*\{[[:space:]]*origin:[[:space:]]*session,[[:space:]]*id:[[:space:]]*$id,[[:space:]]*date:[[:space:]]*$date[[:space:]]*\}"
      ;;
  esac
  if grep -Eq "$dup_re" "$file"; then
    printf 'capture: skip duplicate (%s=%s, date=%s) for %s\n' "$origin" "$id" "$date" "$slug"
    exit 0
  fi
  tmp="$file.tmp.$$"
  awk -v occ="$occ_line" '
    BEGIN { in_occ = 0; done = 0 }
    /^occurrences:[ \t]*$/ { print; in_occ = 1; next }
    {
      if (in_occ && !done) {
        # occurrences ブロックは "  - ..." の行で続く。それ以外が来たら末尾。
        if ($0 !~ /^[ \t]+-/) {
          print occ
          done = 1
          in_occ = 0
        }
      }
      print
    }
    END { if (in_occ && !done) print occ }
  ' "$file" > "$tmp"
  mv "$tmp" "$file"
  printf 'capture: appended occurrence (%s=%s, date=%s) to %s\n' "$origin" "$id" "$date" "$slug"
  exit 0
fi

[ -n "$title" ] || title="$slug"
[ -n "$rule" ]  || rule="（ルール文を記入）"
[ -n "$why" ]   || why="（理由を記入）"

{
  printf -- '---\n'
  printf 'slug: %s\n' "$slug"
  printf 'scope:\n'
  for s in $scopes; do printf -- '  - %s\n' "$s"; done
  printf 'status: active\n'
  printf 'occurrences:\n'
  printf '%s\n' "$occ_line"
  printf -- '---\n'
  printf '\n'
  printf '# %s\n' "$title"
  printf '\n'
  printf '%s\n' "$rule"
  printf '\n'
  printf '**なぜ:** %s\n' "$why"
} > "$file"

printf 'capture: created %s\n' "$file"
exit 0
