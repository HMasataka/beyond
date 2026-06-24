# 0010. レビュー知見を per-file 台帳に貯めて卒業を提案する

- ステータス: Accepted
- 日付: 2026-06-24

## コンテキスト

レビューでは ADR にするほどではない指摘が繰り返し出る。
こうした知見は記録しないと忘れ、同じ指摘を何度も繰り返す。
一方で、出たそばから ADR や lint ルールにすると、まだ一度しか観測していない判断を仕組み化してしまい、過剰な制約を生む。

ADR 未満の知見を貯め、再発が一定回数を超えてから卒業（lint ルール化 / ADR / スタイル文書）させる仕組みが要る。
stage ① では手動 capture と卒業提案フックまでを対象とし、自動抽出（会話ログからの抽出など）は対象外とする。

制約は次のとおり。

- フックは全プロンプトで発火するため、高速で、失敗してもプロンプトを妨げてはならない。
- 台帳の集計は依存を増やさず、開発環境（nix devShell）が無くても動く必要がある。

## 決定

ADR 未満の知見を `docs/learnings/` の台帳に貯め、しきい値を超えたら卒業を提案する。

- 1 知見 1 ファイル `docs/learnings/<slug>.md`。frontmatter（`slug` / `scope` / `status` / `occurrences`）と本文（ルール + 理由）で構成する。
- slug を識別子にする。ファイル名と frontmatter の `slug` を一致させる。
- 観測は `capture.sh` で記録する。既存 slug には `occurrences` を追記し、`(pr, date)` 完全一致はスキップする。同一 PR の別 date は別観測として数える。
- しきい値は「観測 3 件以上 かつ 異なる PR 2 件以上 かつ `status: active`」。同一 PR 内の複数観測は 1 PR と数える。
- しきい値超えを `UserPromptSubmit` フック（`propose.sh`）が検出し、候補 slug を stdout に出してリマインドする。卒業先は lint/CI 化 → ADR → スタイル文書の 3 択を人が選ぶ。スクリプトは選ばない。
- `status` の遷移（`promoted` / `retired`）は人の手動編集のみ。stage ① では自動遷移しない。
- 中立スクリプト（`propose.sh` / `capture.sh`）は POSIX sh + awk のみで書き、jq / node / task / nix に依存しない。`propose.sh` は読み取り専用で全経路 exit 0（fail-safe）。
- skill は ADR 0009 に従い `.agents/skills/review-learnings/` を正本とし、`.claude/skills/review-learnings` はシンボリックリンクにする。
- フック定義を置く `.claude/settings.json` は skill ではないため ADR 0009 の symlink 方針の対象外とし、`.claude` 直下の実体ファイルとして持つ。

## 結果

- 繰り返す指摘が台帳に残り、再発回数が見える。卒業の判断を「観測 3・2 PR」という客観的なしきい値に紐づけられる。
- フックが読み取り専用かつ fail-safe なので、台帳が壊れていてもプロンプトを妨げない。依存が POSIX sh + awk のみなので開発環境が無くても動く。
- 残る制約: 観測の記録は手動 capture に依存する。記録を怠れば台帳は育たない。自動抽出は stage ② に委ねる。卒業先の選択と status 遷移は人が担うため、運用の規律が要る。
- 却下した代替案:
  - CodeRabbit の Learnings 機能に任せる。外部サービスに知見が囲い込まれ、リポジトリから参照・卒業の運用ができない。lint や ADR への卒業という出口とも結びつかないため採らない。
  - 個人のメモリ（エージェントの記憶や個人メモ）に貯める。チームで共有できず、PR 横断の再発回数も数えられない。リポジトリに残らないため採らない。
  - リポジトリ内の単一台帳ファイルに全知見を一括で書く。ファイルが肥大化し、frontmatter ごとの集計や卒業時の status 管理がしにくい。per-file なら slug 単位で扱え、grep ベースの集計も単純になるため、単一台帳は採らない。
