#!/bin/sh
# UserPromptSubmit フックから呼ぶ読み取り専用スクリプト。詳細は SKILL.md。
# 不変条件: 全経路で exit 0（fail-safe）。台帳には書き込まない。

set -u

# CLI 非依存にルートを解決する。Claude は CLAUDE_PROJECT_DIR、他は git トップレベル。
root="${CLAUDE_PROJECT_DIR:-$(git rev-parse --show-toplevel 2>/dev/null || true)}"
[ -n "$root" ] || exit 0

ledger_dir="$root/docs/learnings"
[ -d "$ledger_dir" ] || exit 0

candidates=$(
  for f in "$ledger_dir"/*.md; do
    [ -f "$f" ] || continue
    case "$(basename "$f")" in
      README.md|_template.md) continue ;;
    esac
    awk '
      function trim(s) { sub(/^[ \t]+/, "", s); sub(/[ \t]+$/, "", s); return s }
      BEGIN { fm = 0; slug = ""; status = ""; occ = 0; distinct = 0 }
      /^---[ \t]*$/ {
        fm++
        next
      }
      fm == 1 && /^slug:/ {
        v = $0; sub(/^slug:[ \t]*/, "", v); slug = trim(v); next
      }
      fm == 1 && /^status:/ {
        v = $0; sub(/^status:[ \t]*/, "", v); status = trim(v); next
      }
      fm == 1 && /^[ \t]*-[ \t]*\{/ {
        line = $0
        # distinct キーは PR 由来なら "pr:N"、会話由来なら "session:<id>"。
        key = ""
        if (match(line, /origin:[ \t]*session/)) {
          if (match(line, /id:[ \t]*[A-Za-z0-9._-]+/)) {
            tok = substr(line, RSTART, RLENGTH)
            sub(/id:[ \t]*/, "", tok); key = "session:" tok
          }
        } else if (match(line, /pr:[ \t]*[0-9]+/)) {
          tok = substr(line, RSTART, RLENGTH)
          sub(/pr:[ \t]*/, "", tok); key = "pr:" tok
        }
        if (key != "") {
          occ++
          if (!(key in seen)) { seen[key] = 1; distinct++ }
        }
        next
      }
      END {
        if (slug == "" || status != "active") exit
        if (occ >= 3 && distinct >= 2) print slug
      }
    ' "$f"
  done
)

[ -n "$candidates" ] || exit 0

shown=$(printf '%s\n' "$candidates" | grep -v '^$' | head -n 5)
count=$(printf '%s\n' "$candidates" | grep -vc '^$')
[ "$count" -gt 0 ] || exit 0

printf '【review-learnings】卒業候補の知見が %s 件あります（観測 3 以上・異なる観測元 2 以上・active）。\n' "$count"
printf '%s\n' "$shown" | while IFS= read -r s; do
  [ -n "$s" ] && printf -- '- %s （docs/learnings/%s.md）\n' "$s" "$s"
done
printf 'lint/CI 化 → ADR → スタイル文書の順で卒業を検討してください。詳細は docs/learnings/README.md。\n'

exit 0
