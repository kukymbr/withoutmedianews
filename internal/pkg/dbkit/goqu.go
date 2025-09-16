package dbkit

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/kukymbr/withoutmedianews/internal/domain"
	"go.uber.org/zap"
)

func GoquScanStruct(ctx context.Context, ds *goqu.SelectDataset, target any, logger *zap.Logger) error {
	return goquScanLogged(ctx, ds, target, logger, func(ctx context.Context, target any) error {
		ok, err := ds.ScanStructContext(ctx, target)
		if err != nil {
			return err
		}

		if !ok {
			return domain.ErrNotFound
		}

		return nil
	})
}

func GoquScanStructs(ctx context.Context, ds *goqu.SelectDataset, target any, logger *zap.Logger) error {
	return goquScanLogged(ctx, ds, target, logger, ds.ScanStructsContext)
}

func goquScanLogged(
	ctx context.Context,
	ds *goqu.SelectDataset,
	target any,
	logger *zap.Logger,
	scanFunc func(context.Context, any) error,
) error {
	if err := scanFunc(ctx, target); err != nil {
		query, args, _ := ds.ToSQL()

		logger.Debug(
			"failed scan query",
			zap.Error(err),
			zap.String("query", query),
			zap.Any("args", args),
		)

		return err
	}

	return nil
}
