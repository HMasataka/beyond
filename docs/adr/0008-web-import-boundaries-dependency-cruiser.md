# ADR-0008: web の import 境界を dependency-cruiser で強制する

- ステータス: Accepted
- 日付: 2026-06-24

## 背景

ADR 0005 で web の依存ルール（粒度の一方向、共通はページ固有を import しない、ページ間 import 禁止）を定めた。
ただし規約は文章でしかなく、層をまたぐ誤った import を機械的に防げない。レビューの目視に頼ると、ページやコンポーネントが増えるほど抜けやすい。

Go 側は ADR 0007 が同じ問題を depguard で解いており、依存方向を `task lint` と CI で落とせる。web にも対応物が要る。
手段は dependency-cruiser に決めてある。`web/tsconfig.json` の `moduleResolution: bundler` と相対パス import を解決後のパスで判定でき、`import type` も含めて層を検査できるためである。

## 決定

ADR 0005 の依存ルールを dependency-cruiser で強制する。設定は `web/.dependency-cruiser.cjs` に置く。

- dependency-cruiser は web の pnpm devDependency として入れる。nixpkgs に収録がないため flake では供給せず、web の JS ツールチェーンを管理する pnpm に乗せる。
- すべての禁止ルールの severity は `error` とする。warn ではブランチ保護の必須チェックで止められない。
- `tsConfig` で `web/tsconfig.json` を読み、`tsPreCompilationDeps: true` で `import type` も依存として拾う。相対 import を解決後のパスで判定する。
- `includeOnly: "^src/"` で web/src のみを対象にする。`node_modules`・`dist`・生成物 `src/api/schema.ts`・テスト（`*.test.ts(x)` と `src/test/`）は除外する。テストは mock や testing-library の都合で層を跨ぐのが正常で、境界はプロダクションコードに効けば足りる。
- `lint:arch` タスクを新設し `task lint` に組み込む。`lint:web`（Biome・rustywind）は src を直接読み node_modules を要しないため、dependency-cruiser はそこに混ぜず分離する。
- CI（`.github/workflows/ci.yml`）の lint ジョブに web 依存の install を加える。dependency-cruiser は node ツールで node_modules を要するためである。

### 強制するルール

ADR 0005 の依存ルールに 1 対 1 で対応させ、次を禁止する。許可は暗黙とし、禁止以外はすべて許可する。

- `components/**` から `pages/**`・`layouts/**` への import。共通コンポーネントはページ固有・レイアウトに依存しない。
- `layouts/**` から `pages/**` への import。レイアウトは router の親ルートとしてページを包み、直接は import しない。
- `api/**` から `components/**`・`layouts/**`・`pages/**`・`router` への import。api は最下層の共有インフラで、上位層に依存しない。
- ページ間の直接 import。`pages/<A>` から `pages/<B>`（A≠B）を禁止する。正規表現の capture group で同一ページ配下を除外する。
- 粒度の上位方向。`atoms` から `molecules`・`organisms`、`molecules` から `organisms` への import を禁止する。`pages/<page>/` 配下の粒度サブ層にも同じ向きを適用する。現状コードに実例はないが、テンプレートの段階から効かせる。
- 循環参照（`no-circular`）。層ルールとは別軸で禁止する。

### 違反検出の確認手順

ルールが実際に違反を捕まえることは、捨てファイルで確認する。本番ファイルやコミットには違反 import を残さない。

1. `web/src/_probe.ts` のような捨てファイルを作り、禁止 import を 1 件入れる。たとえば api からページを import する（`api-no-upward` に該当）。
2. `pnpm --dir web run depcruise` を実行し、対応するルール名が `error` として出て exit が非ゼロになることを確認する。
3. 捨てファイルを削除し、`git status` がクリーンに戻ることと depcruise が再び 0 violations になることを確認する。

## 結果/影響

- 層をまたぐ誤った依存が `task lint` と CI で落ちる。ADR 0005 の依存ルールが目視に頼らず機械的に守られる。
- ページやコンポーネントが空に近い現状でも、テンプレートの段階から制約が効く。

## 検討した代替案

- eslint-plugin-boundaries。DX は良いが ESLint を増やし、Biome 一本という ADR 0005 の方針を崩す。
- Biome `noRestrictedImports`。パス相対の方向ルールを表現できず、型 import と値 import を区別できない。
- dependency-cruiser の nix 自前パッケージ化。nixpkgs に未収録で、buildNpmPackage の保守コストが見合わない。

## 注意点

- pnpm devDependency が 1 つ増える。node ビルドは nix が供給するため、再現性は lockfile と nix で担保する。ローカルで `task lint:arch` を単体実行するには node_modules の install が要る。
- テストと生成物 `schema.ts` を除外するため、これらの import は境界検査の対象外になる。
- Go 側の対応物は ADR 0007（depguard）である。
