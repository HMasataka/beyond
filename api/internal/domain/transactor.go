package domain

import "context"

// Transactor は usecase が依存するトランザクション境界。
// シグネチャを stdlib のみに保ち、domain を database/sql に依存させない。
type Transactor interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
