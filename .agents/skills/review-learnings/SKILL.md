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

台帳の書式と運用は `docs/learnings/README.md`、設計判断は `docs/adr/0010-review-learnings.md` を参照する。

## capture（観測を記録する）

知見を観測したら `capture.sh` で記録する。

```sh
.agents/skills/review-learnings/scripts/capture.sh \
  --slug <kebab-case> --pr <番号> --date <YYYY-MM-DD> --scope <name> \
  [--scope <name>]... [--title <タイトル>] [--rule <ルール文>] [--why <理由>]
```

- `--scope` は最低 1 つ必須。複数指定で範囲を足す。
- 新規 slug は frontmatter + 本文の雛形でファイルを作る。`--title` `--rule` `--why` を渡して内容を埋める。
- 既存 slug は `occurrences` に 1 件追記する。`(pr, date)` 完全一致はスキップする。
- `date` は観測した日を入れる。

## recall（台帳を確認する）

`docs/learnings/*.md` を読む。slug 単位で 1 ファイル。`status: active` が観測中の知見。

## graduate（卒業を検討する）

`status: active` かつ「観測 3 件以上・異なる PR 2 件以上」になったら卒業候補。
`UserPromptSubmit` フック（`scripts/propose.sh`）がプロンプト時に候補をリマインドする。

卒業先は人が次の 3 択から選ぶ。スクリプトは選ばない。

1. lint / CI 化
2. ADR
3. スタイル文書

卒業させたら台帳の `status` を `promoted`（廃止なら `retired`）に手で書き換える。自動遷移はしない。

## フックとの関係

`scripts/propose.sh` は `.claude/settings.json` の `UserPromptSubmit` フックから呼ばれる読み取り専用スクリプト。
台帳に書き込まず、候補があれば短いリマインダを出し、無ければ無出力で終わる（全経路 exit 0）。

## スコープ

会話ログからの自動抽出・ノイズ対策・外部レビュー連携は未実装（`docs/adr/0010-review-learnings.md` を参照）。
