# ADR-0006: マイグレーションツールに goose を採用する

- ステータス: Accepted
- 日付: 2026-06-23

## 背景

API は MySQL（`compose.yaml`）を使う。スキーマの変更を管理するマイグレーションツールが要る。
候補は goose・golang-migrate・atlas。

依存は Nix で固定し（ADR 0002）、操作は Taskfile に集約する（ADR 0003）方針に乗せる。

## 決定

マイグレーションツールに goose を採用する。

- プレーン SQL の up/down に加え、必要なら Go で書くマイグレーションも扱える。
- 単一バイナリで、`embed.FS` でアプリに同梱でき、CLI からもアプリ起動時からも適用できる。将来 `di` 経由でマイグレーションを走らせる選択肢を残せる。
- 定義は `internal/migrations/` に置き、`Taskfile.yml` に適用・差し戻し・作成のタスクを用意する。ツールは Nix で固定する（`nixpkgs#goose`）。

## 結果/影響

- スキーマの変更が `internal/migrations/` に集約され、適用方法が Taskfile に揃う。
- SQL と Go の両方でマイグレーションを書けるため、データ移行など SQL だけでは難しい処理にも対応できる。
- `embed.FS` でバイナリに同梱でき、配布物だけでマイグレーションを適用できる。

## 検討した代替案

- golang-migrate: プレーン SQL・CLI で手堅いが、Go マイグレーションやアプリ埋め込みの柔軟性で goose を採る。
- atlas: 宣言的スキーマと diff で高機能だが、重く独自色が強い。テンプレートの土台にはオーバースペック。
