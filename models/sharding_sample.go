package models

import "time"

const ShardingSampleTable string = "ShardingSamples"

type ShardingSample struct {
	ShardingSampleID string
	ShardID          int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (v *ShardingSample) Table() string {
	return ShardingSampleTable
}
