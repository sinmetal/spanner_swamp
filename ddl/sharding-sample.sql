CREATE TABLE ShardingSamples (
  ShardingSampleID STRING(32) NOT NULL,
  ShardID INT64 NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY (ShardingSampleID);

CREATE INDEX ShardingSamplesByShardIDAndUpdatedAtDesc
ON ShardingSamples (
  ShardID,
  UpdatedAt DESC
);