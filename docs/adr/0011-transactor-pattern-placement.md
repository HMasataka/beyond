# ADR-0011: transactor パターンで層をまたぐトランザクション境界を置く

- ステータス: Accepted
- 日付: 2026-06-26

## 背景

複数の repository をまたぐ操作を 1 つの DB トランザクションにまとめたい。
オニオンアーキテクチャ（ADR 0004）と depguard による依存方向の強制（ADR 0007）の下では、トランザクション境界をどの層に置き、tx の有無判定（同じ ctx 内ですでに tx があればそれを使い、なければ開始する）をどこに集約するかの規約がない。

ここで配置を一意に縛るのが depguard の制約である。

- `usecase` は `domain` 以外の内部層を deny する。トランザクション境界を usecase から呼ぶには、その境界が `domain` になければならない。
- `infra` は `domain`・`usecase`・`repository`・`handler` を deny する。tx の有無判定を担う具象は `database/sql` を扱うため infra に置くしかないが、infra は `domain` を import できない。
- `repository` の deny に `infra` は含まれない。repository が `infra/database` の型に依存するのは許可される。

トランザクション境界の interface を usecase が、実行口の供給元を repository がそれぞれ参照する以上、両者を 1 つの層にまとめると上のいずれかの deny に触れる。

## 決定

transactor パターンを採用し、3 つの interface と 1 つの具象を層ごとに分けて置く。

- `domain.Transactor`（`api/internal/domain/transactor.go`）。usecase が依存する境界で、`Do(ctx, fn) error` のみを持つ。シグネチャは stdlib（`context`）だけで構成し、domain を `database/sql` に依存させない。
- `database.DBTX` / `database.DBProvider` と具象 `database.TxManager`（`api/internal/infra/database/transaction.go`）。`DBTX` は `*sql.DB` と `*sql.Tx` の双方が満たす最小の実行口、`DBProvider` はそれを供給する `Executor(ctx) DBTX` を持つ。repository はこれに依存する。

`TxManager` は単一の具象で、`Do` と `Executor` の両方を実装する。tx の有無判定をこの 1 箇所に集約し、`Do` が ctx に載せた tx を `Executor` が取り出して返す。infra は domain を import できないため、`TxManager` は `domain.Transactor` を import せず、同一シグネチャの `Do` を持つことで構造的に満たす。

`di` は合成ルートとして全層を import できる。単一の `TxManager` を生成し、usecase へは `domain.Transactor` として、repository へは `database.DBProvider` として渡す。

```go
txManager := database.NewTxManager(db)

var transactor domain.Transactor = txManager
var provider database.DBProvider = txManager

userRepository := repository.NewUserRepository(provider)
userUsecase := usecase.NewUserUsecase(transactor, userRepository)
```

これにより、境界の interface（`domain.Transactor`）と実行口の interface（`database.DBProvider`）が別の層に分かれていても、結線時には同じ具象を指す。

## 結果/影響

- トランザクション境界が `domain.Transactor` として usecase から見え、実行口が `database.DBProvider` として repository から見える。どの層が何に依存するかが型で明示される。
- tx の有無による分岐が `TxManager` に 1 箇所だけ存在する。repository は `Executor(ctx)` を呼ぶだけで、tx 中かどうかを意識しない。
- テストでは責務を分けて差し替えられる。usecase は `domain.Transactor` のモックで境界の呼び出しを検証し、repository は `TxManager` 本実装に実 DB を渡して実行口を検証する。

## 検討した代替案

- 3 interface すべてを 1 パッケージに平置きする案も検討したが採らない。`Transactor` を domain 以外に置くと usecase から参照できず（usecase は domain 以外を deny）、domain に置くと `DBTX`・`DBProvider` が `database/sql` を持ち込んで domain を技術詳細に依存させ、かつ具象 `TxManager` を infra に置けなくなる（infra は domain を import できない）。depguard と衝突する。

## 注意点

- infra は domain を import できないため、`TxManager` が `domain.Transactor` を満たすことは同一シグネチャによる構造的充足に頼る。合成ルートの di が全層を import できるため、di に `var _ domain.Transactor = (*database.TxManager)(nil)` を置き、シグネチャ drift をコンパイル時に検出する。
- tx は ctx 経由で暗黙に運ばれる。`Do` が張った tx を `Executor` が ctx から取り出す流れは型に現れず、規約として理解する必要がある。

関連 ADR: [0004](0004-api-onion-architecture.md)、[0007](0007-depguard-architecture-enforcement.md)。
