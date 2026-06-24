#!/bin/sh
# 棚卸し候補を一覧する読み取り専用スクリプト。詳細は SKILL.md。
# 自動削除・自動 retire はしない。台帳には一切書き込まない。

set -u

days=90

die() { printf 'prune: %s\n' "$1" >&2; exit 1; }

while [ $# -gt 0 ]; do
  case "$1" in
    --days) days="$2"; shift 2 ;;
    *) die "unknown argument: $1" ;;
  esac
done

printf '%s' "$days" | grep -Eq '^[0-9]+$' || die "--days must be a non-negative integer: $days"

# CLI 非依存にルートを解決する。Claude は CLAUDE_PROJECT_DIR、他は git トップレベル。
root="${CLAUDE_PROJECT_DIR:-$(git rev-parse --show-toplevel 2>/dev/null || true)}"
[ -n "$root" ] || die "cannot resolve project root (set CLAUDE_PROJECT_DIR or run inside a git repo)"

ledger_dir="$root/docs/learnings"
[ -d "$ledger_dir" ] || die "ledger dir not found: $ledger_dir"

# cutoff 日付（YYYY-MM-DD）を移植性を保って求める。GNU/BSD どちらの date でも、
# どちらも無ければ古さ判定を degrade（active の古さ判定をスキップ）する。
now=$(date +%s)
cutoff_epoch=$((now - days * 86400))
cutoff=$(date -u -d "@$cutoff_epoch" +%F 2>/dev/null \
  || date -u -r "$cutoff_epoch" +%F 2>/dev/null \
  || true)

# 各台帳から「slug / status / 最新観測 date / 観測数 / distinct 数」を 1 行で抽出する。
rows=$(
  for f in "$ledger_dir"/*.md; do
    [ -f "$f" ] || continue
    # case ... ;; を $() 内に置くと bash 3.2 が parse エラーになるため if を使う。
    base=$(basename "$f")
    if [ "$base" = "README.md" ] || [ "$base" = "_template.md" ]; then
      continue
    fi
    awk '
      function trim(s) { sub(/^[ \t]+/, "", s); sub(/[ \t]+$/, "", s); return s }
      BEGIN { fm = 0; slug = ""; status = ""; occ = 0; distinct = 0; latest = "" }
      /^---[ \t]*$/ { fm++; next }
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
        # date は YYYY-MM-DD の辞書順比較が日付順と一致する。
        if (match(line, /date:[ \t]*[0-9]{4}-[0-9]{2}-[0-9]{2}/)) {
          d = substr(line, RSTART, RLENGTH)
          sub(/date:[ \t]*/, "", d)
          if (latest == "" || d > latest) latest = d
        }
        next
      }
      END {
        if (slug == "" || status == "") exit
        # latest 欠落時は "-" を入れ、フィールド位置の崩れ（空フィールド collapse）を防ぐ。
        if (latest == "") latest = "-"
        printf "%s\t%s\t%s\t%s\t%s\n", slug, status, latest, occ, distinct
      }
    ' "$f"
  done
)

# 候補分類。出力は人が読める一覧で、削除・retire は人が手で行う。
stale_pending=""
dead_active=""
retired=""

# 最新観測 date が cutoff より古いか。
# cutoff 不明（date 計算不可）・latest 欠落は判定不能として偽を返す。
# よって degrade 時は古さに依存する候補（pending・active）が出ず、retired のみ一覧される。
is_old() {
  _latest="$1"
  [ -n "$cutoff" ] || return 1
  [ -n "$_latest" ] || return 1
  [ "$_latest" \< "$cutoff" ]
}

OLDIFS=$IFS
IFS='
'
for row in $rows; do
  IFS='	'
  # shellcheck disable=SC2086
  set -- $row
  IFS=$OLDIFS
  slug="$1"; status="$2"; latest="$3"; occ="$4"; distinct="$5"
  [ "$latest" != "-" ] || latest=""
  case "$status" in
    pending)
      if is_old "$latest"; then
        stale_pending="$stale_pending$slug\t$latest\n"
      fi
      ;;
    active)
      # 卒業しきい値（観測 3・distinct 2）に届かず、かつ古い active は死蔵候補。
      if { [ "$occ" -lt 3 ] || [ "$distinct" -lt 2 ]; } && is_old "$latest"; then
        dead_active="$dead_active$slug\t$latest\n"
      fi
      ;;
    retired)
      retired="$retired$slug\t$latest\n"
      ;;
  esac
  IFS='
'
done
IFS=$OLDIFS

print_section() {
  _title="$1"; _body="$2"
  printf '## %s\n' "$_title"
  if [ -z "$_body" ]; then
    printf '（候補なし）\n\n'
    return
  fi
  printf '%b' "$_body" | while IFS=$(printf '\t') read -r s d; do
    [ -n "$s" ] || continue
    [ -n "$d" ] || d="（観測 date 不明）"
    printf -- '- %s （docs/learnings/%s.md, 最新観測 %s）\n' "$s" "$s" "$d"
  done
  printf '\n'
}

printf '【review-learnings prune】棚卸し候補（基準 %s 日'"'"'。' "$days"
if [ -n "$cutoff" ]; then
  printf 'cutoff %s より古い観測を「古い」と判定）\n\n' "$cutoff"
else
  printf 'date 計算が使えないため active の古さ判定はスキップ）\n\n'
fi

print_section "放置 pending（承認されず古い）" "$stale_pending"
print_section "死蔵 active（伸びず古い・卒業しきい値未到達）" "$dead_active"
print_section "retired（廃止済み）" "$retired"

printf 'retire / 削除は人が手で行う。このスクリプトは一覧のみで自動変更はしない。\n'

exit 0
