package db_connect

import (
	"context"

	"gorm.io/gorm"
)

type clusterMockImpl struct {
	db *gorm.DB
}

// Initialize implements DBCLuster.
func (c *clusterMockImpl) Initialize() error {
	return nil
}

// Read implements DBCLuster.
func (c *clusterMockImpl) Read(ctx context.Context) *gorm.DB {
	return c.db.WithContext(ctx)
}

// Write implements DBCLuster.
func (c *clusterMockImpl) Write(ctx context.Context) *gorm.DB {
	return c.db.WithContext(ctx)
}

func NewMockDBCluster(db *gorm.DB) DBCLuster {
	return &clusterMockImpl{
		db: db,
	}
}
