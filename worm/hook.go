package worm

import "context"

type BeforeInsertInterface interface {
	BeforeInsert(ctx context.Context)
}

type AfterInsertInterface interface {
	AfterInsert(ctx context.Context)
}

type BeforeUpdateInterface interface {
	BeforeUpdate(ctx context.Context)
}

type AfterUpdateInterface interface {
	AfterUpdate(ctx context.Context)
}

type BeforeDeleteInterface interface {
	BeforeDelete(ctx context.Context)
}

type AfterDeleteInterface interface {
	AfterDelete(ctx context.Context)
}

type BeforeQueryInterface interface {
	BeforeQuery(ctx context.Context)
}

type AfterQueryInterface interface {
	AfterQuery(ctx context.Context)
}