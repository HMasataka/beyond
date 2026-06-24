# 知見台帳（review-learnings）

ADR にするほどではないが、レビューで繰り返し出る指摘や設計上の気づきを貯める場所。
再発が一定回数を超えたら、lint ルール化・ADR・スタイル文書のいずれかへ「卒業」させる。
仕組みの設計判断は `docs/adr/0010-review-learnings.md`・`docs/adr/0011-review-learnings-auto-extraction.md`・`docs/adr/0012-review-learnings-noise-control.md` を参照する。

## 書式

知見は 1 件 1 ファイル `docs/learnings/<slug>.md` で持つ。
ファイルは YAML frontmatter と本文からなる。

- `slug`: kebab-case の識別子。ファイル名と一致させる。
- `scope`: 適用範囲のリスト（`api` / `web` / `general` など）。1 行 1 要素のブロックリストで書く。
- `status`: `pending`（承認待ち）/ `active`（観測中）/ `promoted`（卒業済み）/ `retired`（廃止）。会話由来の自動抽出は `pending` で入る。
- `occurrences`: 観測の一覧。1 行 1 件で書く。ブロック内に空行を挟まない（capture の追記が前提とする）。観測元で 2 つの形を持つ。
  - PR 由来: `- { pr: <番号>, date: <YYYY-MM-DD> }`
  - 会話由来: `- { origin: session, id: <session-id>, date: <YYYY-MM-DD> }`（`docs/adr/0011-review-learnings-auto-extraction.md` を参照）

本文は frontmatter の直後に H1 見出し（短いタイトル）を置き、ルール文と `**なぜ:**` の理由を続ける。

### テンプレート

```markdown
---
slug: example-slug
scope:
  - general
status: active
occurrences:
  - { pr: 15, date: 2026-06-21 }
  - { origin: session, id: 0a1b2c3d, date: 2026-06-24 }
---

# 短いタイトル

ルール文を 1〜2 文で書く。

**なぜ:** 守るべき理由（制約・仕様・設計意図）を書く。
```

## capture（観測の記録）

観測は review-learnings skill 経由で記録する。skill は中立スクリプトを呼ぶ。

```sh
# PR 由来
.agents/skills/review-learnings/scripts/capture.sh \
  --slug <kebab-case> --pr <番号> --date <YYYY-MM-DD> --scope <name> \
  [--scope <name>]... [--title <タイトル>] [--rule <ルール文>] [--why <理由>]

# 会話由来（extract フローで使う。PR 番号が無い。承認待ちの pending で記録する）
.agents/skills/review-learnings/scripts/capture.sh \
  --slug <kebab-case> --origin session --id <session-id> --date <YYYY-MM-DD> --scope <name> --status pending
```

- 既存 slug があれば `occurrences` に 1 件追記する。`(pr, date)` または `(origin, id, date)` が完全一致する観測は追記しない。
- 同じ PR でも `date` が異なれば別の観測として数える。
- 会話ログからの自動抽出の流れは skill の extract フローと `docs/adr/0011-review-learnings-auto-extraction.md` を参照する。
- `date` は観測した日を入れる。PR マージ日ではなく、その知見に気づいた日を記録する。

## ライフサイクル（status の遷移）

status は `pending → active → promoted` / `retired` と進む。遷移は人が手で行い、自動遷移はしない。

- `pending`: 承認待ち。会話由来の自動抽出（extract）はここから入る。人が内容を確認し、採用なら `active`、却下なら `retired` へ手で書き換える。
- `active`: 観測中。卒業しきい値の算入対象はこの状態だけ。`pending` はしきい値に数えない。
- `promoted`: 卒業済み（lint / CI 化・ADR・スタイル文書のいずれかへ移した）。
- `retired`: 廃止。

PR 由来のレビュー指摘は人手の記録なので `active` で入る。承認ゲートの設計は `docs/adr/0012-review-learnings-noise-control.md` を参照する。

## 卒業（graduate）

`status: active` の知見が「観測 3 件以上 かつ 異なる観測元 2 つ以上」になったら卒業を検討する。
観測元は PR 番号、または会話の session id で数える。同一 PR・同一 session 内の複数観測は 1 つとして数える。

しきい値を超えると `UserPromptSubmit` フックが候補をリマインドする。
卒業先は次の 3 択から人が選ぶ。スクリプトは選択しない。

1. lint / CI 化: 機械的に検出できるなら lint ルールや CI チェックにする。
2. ADR: 設計判断として残すべきなら ADR を起こす。
3. スタイル文書: 規約・スタイルガイドに書くべきならドキュメント化する。

卒業させたら `status` を `promoted` に、廃止する知見は `retired` に手で書き換える。
status の自動遷移はしない（人の手動編集のみ）。

## prune（棚卸し）

台帳が育つと、承認されないまま放置された `pending`・伸びずに死蔵した `active`・廃止済みの `retired` が溜まる。
`scripts/prune.sh [--days N]`（既定 90）で棚卸し候補を一覧する。読み取り専用で、台帳には書き込まない。

候補は 3 区分で出る。

1. 放置 pending: `pending` のまま最新観測が `--days` より古いもの。
2. 死蔵 active: `active` だが卒業しきい値に届かず、最新観測が `--days` より古いもの。
3. retired: 廃止済みのもの。

prune は一覧だけを出す。retire や削除は人が手で行う。自動削除・自動 retire はしない。
