#!/bin/sh
# 知見台帳に観測を記録する（skill から呼ぶ）。使い方は SKILL.md。

set -eu

slug=""
pr=""
date=""
title=""
rule=""
why=""
scopes=""

die() { printf 'capture: %s\n' "$1" >&2; exit 1; }

while [ $# -gt 0 ]; do
  case "$1" in
    --slug)  slug="$2";  shift 2 ;;
    --pr)    pr="$2";    shift 2 ;;
    --date)  date="$2";  shift 2 ;;
    --title) title="$2"; shift 2 ;;
    --rule)  rule="$2";  shift 2 ;;
    --why)   why="$2";   shift 2 ;;
    --scope) scopes="$scopes $2"; shift 2 ;;
    *) die "unknown argument: $1" ;;
  esac
done

[ -n "$slug" ] || die "--slug is required"
[ -n "$pr" ]   || die "--pr is required"
[ -n "$date" ] || die "--date is required"
[ -n "$scopes" ] || die "--scope is required (use --scope general if truly global)"

printf '%s' "$slug" | grep -Eq '^[a-z0-9]+(-[a-z0-9]+)*$' || die "slug must be kebab-case: $slug"
printf '%s' "$pr"   | grep -Eq '^[0-9]+$' || die "pr must be a number: $pr"
printf '%s' "$date" | grep -Eq '^[0-9]{4}-[0-9]{2}-[0-9]{2}$' || die "date must be YYYY-MM-DD: $date"

# CLI 非依存にルートを解決する。Claude は CLAUDE_PROJECT_DIR、他は git トップレベル。
root="${CLAUDE_PROJECT_DIR:-$(git rev-parse --show-toplevel 2>/dev/null || true)}"
[ -n "$root" ] || die "cannot resolve project root (set CLAUDE_PROJECT_DIR or run inside a git repo)"
ledger_dir="$root/docs/learnings"
[ -d "$ledger_dir" ] || die "ledger dir not found: $ledger_dir"

file="$ledger_dir/$slug.md"
occ_line="  - { pr: $pr, date: $date }"

if [ -f "$file" ]; then
  if grep -Eq "^[[:space:]]*-[[:space:]]*\{[[:space:]]*pr:[[:space:]]*$pr,[[:space:]]*date:[[:space:]]*$date[[:space:]]*\}" "$file"; then
    printf 'capture: skip duplicate (pr=%s, date=%s) for %s\n' "$pr" "$date" "$slug"
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
  printf 'capture: appended occurrence (pr=%s, date=%s) to %s\n' "$pr" "$date" "$slug"
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
