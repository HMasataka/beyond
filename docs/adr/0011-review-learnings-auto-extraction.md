# 0011. レビュー知見を会話ログから自動抽出する

- ステータス: Accepted
- 日付: 2026-06-24

## コンテキスト

ADR 0010 で導入した知見台帳は、観測の記録を手動 capture に依存する。
記録を怠れば台帳は育たず、卒業のしきい値も意味を持たない。
レビュー指摘や設計上の気づきの多くはエージェントとの会話の中で出るため、その会話ログから観測を拾えれば記録漏れを減らせる。

ただし会話ログからの抽出には次の制約がある。

- フックは全プロンプト・全ターンで発火するため、高速で、失敗してもセッションを妨げてはならない。
- 知見かどうかの判断には文脈理解が要る。決定的なスクリプトでは判断できない。
- 会話ログ（transcript）の所在やフックの種類は CLI ごとに異なる。Claude Code と codex の両方で動かしたい。
- 会話由来の観測には PR 番号が無い。ADR 0010 の `{ pr, date }` スキーマでは記録できない。

## 決定

会話ログからの自動抽出を「記録係フック」と「本体エージェントによる抽出」に分離する。

- フックは LLM を呼ばない。記録係 `record.sh` は `Stop` で発火し、抽出すべき会話区間の所在（session id と transcript パス）を保留キュー `.local/review-learnings/queue` に冪等記録するだけにする。判断は本体エージェントが担う。
- 抽出は本体エージェントが行う。`SessionStart` フック `backfill.sh` が「未抽出の会話がある」と文脈注入し、本体が review-learnings skill の extract フローで保留区間の transcript を読み、知見を判断して `capture.sh` で記録する。別プロセスの LLM 起動はしない。
- ウォーターマークは最後に抽出した位置（transcript の行番号）を session 別に `.local/review-learnings/watermark/<session-id>` へ持つ。各回は前回以降の新区間だけを対象にし、重複を防ぐ。ウォーターマークは capture 成功でのみ前進させる（注入は抽出完了を意味しない）。
- 会話由来の観測スキーマは `{ origin: session, id: <session-id>, date }` とする。既存の `{ pr, date }` は後方互換で維持する。`propose.sh` の distinct カウントを「pr 値 もしくは origin+id」をキーに一般化し、会話由来観測もしきい値に算入できるようにする。1 会話 = 1 観測（origin+id で正規化）。
- フックは `Stop` ベースに統一する。codex には `SessionEnd` が無いため、区間記録の役割を `Stop` に寄せて Claude / codex で共通化する。`Stop` は毎ターン発火するため記録係は冪等にする。文脈注入は `SessionStart` 経由で、Claude は `additionalContext`、codex は平文を出力アダプタで切り替える。
- 状態ファイルは `.local/review-learnings/` 配下に置く。`.local/` は gitignore 済み・markdownlint 除外済みで、各環境ローカルに留まる。テンプレ同梱対象はスクリプトとフック定義（`.claude/settings.json` / `.codex/hooks.json`）のみとする。
- codex の transcript は graceful degradation とする。codex の `transcript_path` が null / 欠落なら `record.sh` は no-op（exit 0）で、codex の会話自動抽出だけスキップする。提案フックと手動 capture は動く。
- ルート解決は CLI で異なる。Claude はフック定義で `$CLAUDE_PROJECT_DIR` を使い、codex は同等のプレースホルダが無いためフック定義で `git rev-parse --show-toplevel` を使う。スクリプト内部のルート解決も `CLAUDE_PROJECT_DIR` 優先・git fallback で揃える。
- スクリプトは ADR 0010 と同じく POSIX sh + awk のみで書き、jq / node / task / nix に依存しない。全フックスクリプトは全経路 exit 0（fail-safe）。

## 結果

- 会話で出た知見が記録漏れしにくくなる。会話由来観測が origin / id で台帳に残り、卒業しきい値にも算入される。
- フックが薄い記録係に徹し全経路 exit 0 なので、状態破損・書込不可・transcript 失効でもセッションを妨げない。同一 `record.sh` / `backfill.sh` を Claude と codex で共有し、出力アダプタだけが分岐する。
- 既存の手動 capture と提案フック（ADR 0010・`{ pr, date }` 台帳・symlink）は無改変で動く。
- 残る制約: 抽出の質は本体エージェントの判断に依存する。ノイズ対策（承認待ち・定期 prune）は本 ADR の対象外とし、別の決定で扱う。codex 環境が手元に無いため、codex の `transcript_path` の実体と append-only 性は実機未検証で、degrade 時は会話抽出のみスキップされる。backfill の対象は保留キューに記録された未抽出 session に限り、過去全件は遡及しない。
- 却下した代替案:
  - フックの中で LLM CLI（`claude -p` 等）を起動して抽出させる。フックは高速・fail-safe であるべきで、LLM 起動はタイムアウトと失敗のリスクを持ち込む。別プロセス起動はセッションのコストと複雑さを増やし、codex 移植も難しくなるため採らない。
  - データ保全目的で `PreCompact` フックを追加する。Claude Code の transcript jsonl は append-only で、compaction が起きても過去行は削除・書き換えされない。resume も同一ファイルを再利用する。よって session 別ウォーターマークは compaction を跨いでも有効で、`SessionStart` の照合だけで圧縮前を含む全未抽出区間を拾える。フックを増やすと codex 移植が複雑化するだけで compaction race も防げないため、`PreCompact` は配線しない。
  - 単一グローバルのウォーターマークにする。別 session・別プロジェクトの区間を取り違える。session 別に持てば誤りなく区間を切り出せるため、グローバルウォーターマークは採らない。
