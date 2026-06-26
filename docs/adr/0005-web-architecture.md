# ADR-0005: web のアーキテクチャ（アトミックデザイン + 共通/ページ固有の二層）

- ステータス: Accepted
- 日付: 2026-06-22

## 背景

web は React + TypeScript + Vite で、現状は単一画面のサンプルしかない。
コンポーネントの分類・置き場所・依存方向、ルーティング、データ取得、スタイリングの方針がないと、共通部品とページ固有部品が混ざり、依存が無秩序に広がる。

前提として、OpenAPI から生成した型（`api/schema.ts`）を活かしたい。
また依存は Nix、タスクは Task、整形・lint は Biome 一本という既存方針（ADR 0002・0003）に乗せる。

## 決定

コンポーネントをアトミックデザインの粒度で分類し、スコープで「共通」と「ページ固有」に分ける。
ルーティング・データ取得・スタイリングの技術もあわせて定める。

### ディレクトリ構成

```text
web/src/
├── components/                共通（ページ横断）のコンポーネント
│   ├── atoms/<Component>/
│   ├── molecules/<Component>/
│   └── organisms/<Component>/
├── layouts/                   アプリシェル（AppLayout など）
├── pages/
│   └── <page>/
│       ├── index.tsx          ルートに対応。データ取得・結線
│       ├── atoms|molecules|organisms/   このページ固有のコンポーネント
│       └── api.ts             このページ固有の query hooks
├── api/
│   ├── schema.ts              OpenAPI 生成物（型）
│   ├── client.ts              openapi-fetch クライアント（共通・1 箇所）
│   └── queryClient.ts         TanStack Query 設定
├── router.tsx                 ルート定義
└── main.tsx                   エントリ
```

各コンポーネントはフォルダ単位で持ち、`<Component>.tsx`・`<Component>.test.tsx`・re-export だけの `index.ts` を置く。
ページは `pages/<page>/index.tsx` をルート要素のエントリとして直接使い、バレルは置かない。
`<page>`・`<Component>` は対象ごとの名前を入れるプレースホルダ。

### 依存のルール

- 粒度の依存は下位から上位への一方向（`atoms ← molecules ← organisms ← pages`）。下位は上位を import しない。
- 共通コンポーネントはページ固有コンポーネントを import しない。
- ページ間の直接 import は禁止する。共有したいものは共通へ昇格させる。
- レイアウトは router の親ルートとしてページを包む。

### 昇格

コンポーネントと query hook は、まずページ固有（`pages/<page>/`）として作る。
2 つ目のページで必要になった時点で共通（`components/<粒度>/`・共通 hook）へ移す。最初から共通にはしない。

### 技術選定

- ルーティング: react-router v7（`createBrowserRouter`。データルータ）。親ルートに layout、子に各ページを置き `<Outlet/>` で差し込む。
- HTTP クライアント: openapi-fetch。`schema.ts` を直接消費し、追加生成なしで型付きになる。
- データ取得: TanStack Query を openapi-fetch とともにページローカルの query hooks で使う。共通クライアントと QueryClient だけ `api/` に集約する。
- スタイリング: Tailwind CSS v4（`@tailwindcss/vite`・トークンは `@theme`）。クラス順の整列は rustywind を Task と CI に組み込み、Nix で固定する。

## 結果/影響

- 共通とページ固有が分離し、依存が下位・内向きの一方向に揃う。
- OpenAPI の型が `schema.ts` から openapi-fetch、query hooks まで一気通貫で効く。
- 追加ツール（react-router・TanStack Query・Tailwind・rustywind）はいずれも既存方針（Nix で固定、Task に集約、Biome 一本）に乗せられる。

## 検討した代替案

- ルーティングは TanStack Router（型安全だが土台としては尖りすぎ）、ルータ無し（ページ概念が成立しない）。
- データ取得は route loaders 単独（キャッシュ・ミューテーションが弱い。必要時の併用は可）、素の hook（実アプリで破綻する）。
- スタイリングは CSS Modules（依存ゼロだが構築速度で劣る）、Biome `useSortedClasses`（Tailwind v4 の独自クラス・レスポンシブ未対応、autofix が unsafe）、prettier-plugin-tailwindcss（Prettier 併用で Biome 一本の方針を崩す）。
- レイアウトは templates 層の常設（単純なページで二重化する。配置が複雑・再利用が要るページに限って後付けする）。
- コンポーネントはフラットファイル（複雑化でフォルダと混在する）。
- 層単位のバレル（循環参照・バンドル肥大を招き、Biome の `noBarrelFile`・`noReExportAll` とも相反する。コンポーネント単位の re-export だけ置く）。

## 注意点

- rustywind を 1 つ増やす。Tailwind v4 の独自クラスの並びは、生成された CSS を参照させて学習させる。
