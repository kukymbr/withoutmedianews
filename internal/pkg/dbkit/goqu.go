package dbkit

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

func GoquScanStructs(ctx context.Context, ds *goqu.SelectDataset, target any, logger *zap.Logger) error {
	if err := ds.ScanStructsContext(ctx, target); err != nil {
		query, args, _ := ds.ToSQL()

		logger.Debug(
			"scan structs error",
			zap.Error(err),
			zap.String("query", query),
			zap.Any("args", args),
		)

		return err
	}

	return nil
}
