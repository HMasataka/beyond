# ADR-0010: コミットメッセージと PR タイトルを commitlint で強制する

- ステータス: Accepted
- 日付: 2026-06-26

## 背景

このリポジトリのコミットは Conventional Commits に沿った形式で書いてきた。
ただし規約は慣習として共有されているだけで、形式を外れたコミットを止める仕組みはない。
規約が守られているかは人の目に依存し、レビューの負担になる。

このリポジトリはマージ時に merge コミットを残す運用で、PR の各コミットがそのまま `main` に乗る。
そのため履歴の品質はコミットメッセージそのものに左右され、検査の対象も個々のコミットになる。

開発ツールは nix flake の devShell から供給し、web 固有の Node ツールは `web/package.json` で管理してきた。
コミット規約はリポジトリ全体の関心事で、web に限った話ではない。

## 決定

コミットメッセージと PR タイトルを **commitlint** で検査し、ルールは `@commitlint/config-conventional` を正本とする。

- リポジトリルートに commitlint 専用の `package.json` を置き、`@commitlint/cli`、`@commitlint/config-conventional`、`husky` を依存に持つ。web とは独立した Node プロジェクトとして扱う。
- ルールは `commitlint.config.mjs` に書き、`@commitlint/config-conventional` を継承する。
- ローカルでは husky の `commit-msg` フックで、コミットのたびに検査する。
- CI では PR のコミット範囲（`base..head`）を検査する。push 時は範囲が取れないため PR でのみ走らせる。
- PR タイトルは `amannn/action-semantic-pull-request` で Conventional Commits 準拠を検査する。タイトルの編集でも再検査できるよう `edited` イベントも対象にする。

## 結果/影響

- 形式を外れたコミットはローカルのフックで止まり、すり抜けても CI が PR を fail させる。規約の遵守を人の目に頼らずに済む。
- コミット規約のツールがルートの `package.json` に集まり、web の関心事と混ざらない。
- merge コミットは commitlint が既定で無視するため、現在の merge コミット運用のまま検査できる。
- PR タイトルの検査はマージ戦略に依存しない。将来 squash merge へ切り替えてタイトルがコミットになっても、同じ語彙で揃うため移行できる。

## 検討した代替案

- commitlint を `web/package.json` に同居させる案は、CI の pnpm 基盤を流用できる利点はあるが、web 固有でない関心事を web に持ち込むため採らない。
- commitlint を nix flake の devShell で供給する案は、ツール供給方針には沿うが、`@commitlint/config-conventional` の解決に追加の作業が要るため採らない。

## 注意点

- ルートに Node の依存と `node_modules` が増える。フックはローカルの `nix develop` 環境に husky と pnpm がある前提で動く。
