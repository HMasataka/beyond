---
name: review-learnings
description: >-
  ADR 未満のレビュー知見を docs/learnings/ の台帳に貯め、再発が閾値を超えたら卒業（lint/ADR/style）を提案する。
  「知見に残して」「レビュー知見」「learning 追加」「卒業候補」などの時に使う。
  観測の記録（capture）、台帳の確認（recall）、卒業の検討（graduate）を扱う。
---

# review-learnings

ADR にするほどではないが繰り返し出るレビュー指摘・気づきを台帳に貯め、
再発が閾値を超えたら lint/CI 化・ADR・スタイル文書へ卒業させる仕組み。

台帳の書式と運用は `docs/learnings/README.md`、設計判断は `docs/adr/0010-review-learnings.md`・`docs/adr/0011-review-learnings-auto-extraction.md`・`docs/adr/0012-review-learnings-noise-control.md` を参照する。

## capture（観測を記録する）

知見を観測したら `capture.sh` で記録する。観測元は PR と会話セッションの 2 種類がある。

```sh
# PR 由来（レビュー指摘など）
.agents/skills/review-learnings/scripts/capture.sh \
  --slug <kebab-case> --pr <番号> --date <YYYY-MM-DD> --scope <name> \
  [--scope <name>]... [--title <タイトル>] [--rule <ルール文>] [--why <理由>]

# 会話由来（PR 番号が無い。extract フローで使う）
.agents/skills/review-learnings/scripts/capture.sh \
  --slug <kebab-case> --origin session --id <session-id> --date <YYYY-MM-DD> --scope <name> \
  --status pending
```

- `--scope` は最低 1 つ必須。複数指定で範囲を足す。
- 新規 slug は frontmatter + 本文の雛形でファイルを作る。`--title` `--rule` `--why` を渡して内容を埋める。
- 既存 slug は `occurrences` に 1 件追記する。`(pr, date)` または `(origin, id, date)` 完全一致はスキップする。
- `--pr N` は `--origin pr --id N` の糖衣。会話由来は PR 番号が無いため `--origin session --id <session-id>` で記録する。
- `--status <active|pending>`（既定 `active`）。新規作成時の初期 status を決める。会話由来（extract）は `--status pending` で承認待ちにする。既存 slug への追記では status は変えない。`promoted` / `retired` は capture では設定できない（人手遷移）。
- `date` は観測した日を入れる。

## extract（会話ログから観測を抽出する）

SessionStart で「未抽出の会話がある」と注入されたら、本体エージェントが保留区間を読んで知見を記録する。
フックは LLM を呼ばない。記録係（`record.sh`）が保留キューに区間を残し、抽出の判断はここで行う。

1. `.local/review-learnings/queue` を読む。各行は `<session-id>\t<transcript-path>` のタブ区切り。
2. 各 session のウォーターマーク `.local/review-learnings/watermark/<session-id>`（最後に抽出した行番号。無ければ 0）を読む。
3. ウォーターマークの次の行から transcript jsonl の末尾までを絶対パスで Read する。
4. content が空 / null の行（compaction race 由来）と compaction サマリ行はスキップする。
5. レビュー指摘パターンの知見だけを判断し、`capture.sh --origin session --id <session-id> --date <date> --status pending` で記録する。会話由来は PR 番号が無いため `origin=session` で記録し、承認待ち（`pending`）で入れる。
6. capture が成功した session のみ、ウォーターマーク `.local/review-learnings/watermark/<session-id>` に「読み終えた行番号（transcript の総行数）」を書いて前進させる。

注入は抽出完了を意味しない。capture しなかった session のウォーターマークは前進させない。
同一区間が再注入されても、ウォーターマークが前進していなければ同じ範囲を再度判断するだけで、台帳には `(origin, id, date)` 一致で重複追記されない（冪等）。
1 会話 = 1 観測。同一会話内で同じ知見が複数回出ても origin+id（session 単位）で 1 件に正規化される。

## recall（台帳を確認する）

`docs/learnings/*.md` を読む。slug 単位で 1 ファイル。`status: active` が観測中の知見。

## graduate（卒業を検討する）

`status: active` かつ「観測 3 件以上・異なる観測元 2 つ以上」になったら卒業候補（観測元は PR または会話 session）。
`UserPromptSubmit` フック（`scripts/propose.sh`）がプロンプト時に候補をリマインドする。

卒業先は人が次の 3 択から選ぶ。スクリプトは選ばない。

1. lint / CI 化
2. ADR
3. スタイル文書

卒業させたら台帳の `status` を `promoted`（廃止なら `retired`）に手で書き換える。自動遷移はしない。

## approve（承認待ちを確認する）

会話由来の自動抽出は `pending` で入る。`pending` の知見を人が確認し、採用するなら `status` を `active` に、却下するなら `retired` に手で書き換える。
`pending` のままでは卒業しきい値に算入されない（`propose.sh` は `active` のみ数える）。status の自動遷移はしない。

## prune（棚卸し候補を見る）

`scripts/prune.sh [--days N]`（既定 90）で棚卸し候補を一覧する。読み取り専用で台帳には書き込まない。手動 / skill アクションで呼ぶ（フックではない）。

候補は 3 区分で出る。

- 放置 pending: 承認されないまま最新観測が `--days` より古いもの。
- 死蔵 active: 卒業しきい値に届かず最新観測が `--days` より古いもの。
- retired: 廃止済みのもの。

prune は一覧のみで、retire / 削除は人が手で行う。自動削除・自動 retire はしない。

## フックとの関係

スクリプトは Claude `.claude/settings.json` と codex `.codex/hooks.json` のフックから呼ばれる。全経路 exit 0（fail-safe）。

- `propose.sh`: `UserPromptSubmit`。卒業候補を読み取り専用でリマインドする。
- `record.sh`: `Stop`。会話の保留区間（transcript パスと session id）を `.local/review-learnings/queue` に冪等記録する。`transcript_path` が null / 欠落なら no-op。
- `backfill.sh`: `SessionStart`。未抽出区間が残る session があれば extract を促す文脈を注入する。出力は Claude が `additionalContext`（`--cli claude`）、codex が平文。

## スコープ

ノイズ対策は pending 承認ゲート（approve）と棚卸し補助（prune）で扱う。設計は `docs/adr/0012-review-learnings-noise-control.md` を参照する。
近似重複のマージ・Gemini 対応・新規フックの追加（pending 件数の自動通知など）はスコープ外。設計の全体は `docs/adr/0010-review-learnings.md`・`docs/adr/0011-review-learnings-auto-extraction.md`・`docs/adr/0012-review-learnings-noise-control.md` を参照する。
