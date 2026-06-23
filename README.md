# beyond

Go 製 API と React 製 web を OpenAPI で繋ぐモノレポのプロジェクトテンプレート。
新規プロジェクトの土台として、開発環境・コード生成・lint・テスト・CI を最初から揃えてある。

## 構成

```text
.
├── api/                   Go API（chi + oapi-codegen）
│   ├── cmd/main/          エントリポイント（HTTP サーバ）
│   └── internal/
│       ├── handler/       OpenAPI から生成した ServerInterface の実装
│       ├── usecase/       ユースケース
│       ├── repository/    データアクセス
│       ├── di/            依存の組み立て
│       ├── migrations/    マイグレーション
│       └── openapi/       OpenAPI からの生成物（編集しない）
├── web/                   React + TypeScript + Vite
│   └── src/
│       ├── api/           API クライアント（openapi-fetch + TanStack Query）と生成物 schema.ts
│       ├── components/    再利用コンポーネント（Atomic Design。atoms/ など）
│       ├── layouts/       ページ共通の枠（AppLayout）
│       ├── pages/         ルートごとの画面とデータ取得（pages/<name>/）
│       ├── router.tsx     React Router のルート定義
│       └── main.tsx       エントリポイント
├── openapi/openapi.yaml   API 仕様。Go と TypeScript の生成の起点
├── compose.yaml           実行時に立ち上げるサービス（MySQL）
├── flake.nix              開発ツール（Nix で固定）
├── Taskfile.yml           タスクランナー
└── docs/adr/              アーキテクチャ決定記録（ADR）
```

依存とサービスの管理方針は [docs/adr/0002](docs/adr/0002-nix-and-docker-compose.md) を参照。
開発ツールは Nix、実行サービスは Docker Compose で分けて管理する。

## 前提

- [Nix](https://nixos.org/)（flakes を有効化）
- [direnv](https://direnv.net/)（任意。devShell の自動読み込み用）
- Docker（`compose.yaml` のサービス起動用）

## セットアップ

direnv を使う場合、リポジトリに入るだけで devShell が読み込まれる。

```sh
direnv allow
```

direnv を使わない場合は、以下で devShell に入る。

```sh
nix develop
```

devShell には Go・Node.js・pnpm・[Task](https://taskfile.dev/)・golangci-lint・Biome・Air が揃う。
web の依存はテスト・ビルド前に取得する。

```sh
pnpm --dir web install
```

## 開発

API と web の開発サーバをまとめて起動する。

```sh
task dev
```

- API: <http://localhost:8080>（Air によるライブリロード。`ADDR` で変更可）
- web: <http://localhost:5173>（Vite）

## 主なタスク

`task` で一覧を表示する。`:api` / `:web` を付けると片側だけ実行する。

| タスク               | 内容                                         |
| -------------------- | -------------------------------------------- |
| `task dev`           | API と web の開発サーバを同時起動            |
| `task gen`           | OpenAPI から Go と TypeScript のコードを生成 |
| `task lint`          | lint（Go: vet + golangci-lint / web: Biome） |
| `task test`          | テスト（Go: `go test` / web: Vitest）        |
| `task build:api`     | API バイナリをビルド                         |
| `task build:web`     | web の本番バンドルをビルド                   |
| `task compose:up`    | コンテナ起動                                 |
| `task compose:down`  | コンテナ停止（ボリュームは残す）             |
| `task compose:reset` | コンテナ停止しボリュームも削除               |

## コード生成

`openapi/openapi.yaml` を単一の仕様として、両側のコードを生成する。

```sh
task gen
```

- Go: 型と chi サーバ（`api/internal/openapi/`）
- TypeScript: 型（`web/src/api/schema.ts`）

API を追加・変更するときは、まず `openapi/openapi.yaml` を編集し、`task gen` で生成し直す。

## サービス（DB）

`compose.yaml` で MySQL を起動できる。

```sh
task compose:up
```

## ADR

アーキテクチャや設計上の意思決定は ADR に残す。
書くタイミングや構成は [docs/adr/README.md](docs/adr/README.md) を参照。

## テンプレートとしての使い方

このリポジトリを GitHub の「Use this template」から複製すると、初回 push 時に
[template-cleanup](.github/workflows/template-cleanup.yml) ワークフローが一度だけ走り、
`beyond` のプレースホルダを新しいリポジトリ名へ書き換える。

- Go モジュールパス `github.com/HMasataka/beyond` → `github.com/<owner>/<repo>`
- データベース名 `beyond-db` → `<repo>-db`
- プロジェクト名 `beyond` / `Beyond` → `<repo>` / `<Repo>`

書き換えの詳細は [.github/template-cleanup/run.sh](.github/template-cleanup/run.sh) を参照。
