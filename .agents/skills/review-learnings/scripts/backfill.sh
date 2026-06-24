#!/bin/sh
# SessionStart backfill 注入。Claude / codex の SessionStart から呼ぶ。詳細は SKILL.md。
# 不変条件: 全経路で exit 0（fail-safe）。状態破損・transcript 失効でもセッション開始を妨げない。
#
# 出力アダプタ: --cli claude は additionalContext JSON、それ以外（codex 等）は平文 stdout。

set -u

cli="codex"
while [ $# -gt 0 ]; do
  case "$1" in
    --cli) cli="${2:-codex}"; shift 2 2>/dev/null || shift ;;
    *) shift ;;
  esac
done

root="${CLAUDE_PROJECT_DIR:-$(git rev-parse --show-toplevel 2>/dev/null || true)}"
[ -n "$root" ] || exit 0

state_dir="$root/.local/review-learnings"
queue="$state_dir/queue"
wm_dir="$state_dir/watermark"

[ -f "$queue" ] || exit 0

# 未抽出区間が残る session を集める。
pending=""
while IFS="$(printf '\t')" read -r sid tpath; do
  [ -n "$sid" ] || continue
  [ -n "$tpath" ] || continue
  printf '%s' "$sid" | grep -Eq '^[A-Za-z0-9._-]+$' || continue
  [ -f "$tpath" ] || continue
  total=$(awk 'END { print NR }' "$tpath" 2>/dev/null || printf '0')
  [ -n "$total" ] || total=0
  wm_file="$wm_dir/$sid"
  if [ -f "$wm_file" ]; then
    mark=$(awk 'NR==1 { print $0 }' "$wm_file" 2>/dev/null || printf '0')
    printf '%s' "$mark" | grep -Eq '^[0-9]+$' || mark=0
  else
    mark=0
  fi
  if [ "$total" -gt "$mark" ]; then
    pending="$pending $sid"
  fi
done < "$queue"

# 前後の空白を除去。未抽出ゼロなら無出力で終了。
pending=$(printf '%s' "$pending" | awk '{ $1=$1; print }')
[ -n "$pending" ] || exit 0

count=$(printf '%s' "$pending" | awk '{ print NF }')
msg=$(printf '【review-learnings】未抽出の会話が %s 件あります。review-learnings skill の extract フローで保留区間（.local/review-learnings/queue）を読み、レビュー指摘パターンの知見を capture.sh --origin session で記録してください。詳細は docs/learnings/README.md。' "$count")

case "$cli" in
  claude)
    # additionalContext へ注入。awk で JSON 文字列としてエスケープする。
    esc=$(printf '%s' "$msg" | awk '{ gsub(/\\/, "\\\\"); gsub(/"/, "\\\""); print }')
    printf '{"hookSpecificOutput":{"hookEventName":"SessionStart","additionalContext":"%s"}}\n' "$esc"
    ;;
  *)
    printf '%s\n' "$msg"
    ;;
esac

exit 0
