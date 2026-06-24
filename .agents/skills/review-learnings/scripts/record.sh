#!/bin/sh
# フック記録係。Stop（Claude / codex）から呼ぶ。詳細は SKILL.md。
# 不変条件: 全経路で exit 0（fail-safe）。状態破損・書込不可・transcript 欠落でもセッションを妨げない。

set -u

# フック stdin は単一行 JSON 前提。改行を畳んでから解析する。
payload=$(cat 2>/dev/null | tr '\n\r' '  ' || true)

# トップレベルの文字列フィールドを取り出す（jq 非依存）。null / 欠落は空文字。
json_str() {
  printf '%s' "$payload" | awk -v key="$1" '
    {
      # "key" : "value" を拾う。value は最初の非エスケープ " まで。
      re = "\"" key "\"[ \t]*:[ \t]*\""
      if (match($0, re)) {
        rest = substr($0, RSTART + RLENGTH)
        out = ""
        n = length(rest)
        i = 1
        while (i <= n) {
          c = substr(rest, i, 1)
          if (c == "\\") { out = out substr(rest, i, 2); i += 2; continue }
          if (c == "\"") break
          out = out c
          i++
        }
        print out
        exit
      }
    }
  '
}

session_id=$(json_str session_id)
transcript_path=$(json_str transcript_path)

# transcript が無ければ何もしない（codex で null/欠落のときの graceful degradation）。
[ -n "$session_id" ] || exit 0
[ -n "$transcript_path" ] || exit 0

# session id・transcript_path は安全な値のみ受理（パス・キー注入、tab 区切り破壊を防ぐ）。
printf '%s' "$session_id" | grep -Eq '^[A-Za-z0-9._-]+$' || exit 0
printf '%s' "$transcript_path" | grep -Eq '^/[^[:cntrl:]]*$' || exit 0

root="${CLAUDE_PROJECT_DIR:-$(git rev-parse --show-toplevel 2>/dev/null || true)}"
[ -n "$root" ] || exit 0

state_dir="$root/.local/review-learnings"
queue="$state_dir/queue"

mkdir -p "$state_dir" 2>/dev/null || exit 0

# 冪等記録: 既存 session 行は transcript_path を更新、無ければ追記。重複行は作らない。
tmp="$queue.tmp.$$"
new_line="$session_id	$transcript_path"
if [ -f "$queue" ]; then
  awk -F '\t' -v sid="$session_id" -v line="$new_line" '
    BEGIN { replaced = 0 }
    $1 == sid { if (!replaced) { print line; replaced = 1 } ; next }
    { print }
    END { if (!replaced) print line }
  ' "$queue" > "$tmp" 2>/dev/null || exit 0
  mv "$tmp" "$queue" 2>/dev/null || { rm -f "$tmp" 2>/dev/null; exit 0; }
else
  { printf '%s\n' "$new_line" > "$queue"; } 2>/dev/null || exit 0
fi

exit 0
