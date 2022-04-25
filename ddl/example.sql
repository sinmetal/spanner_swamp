CREATE TABLE Users (
  UserID STRING(32) NOT NULL,
  MailAddress STRING(320) NOT NULL,
  FirstName STRING(256) NOT NULL,
  LastName STRING(256) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY (UserID);

CREATE INDEX UsersByMailAddress
  ON Users (
    MailAddress,
);

CREATE TABLE UserSecretInfos (
  UserID STRING(32) NOT NULL,
  Address STRING(1024),
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(UserID),
  INTERLEAVE IN PARENT Users ON DELETE CASCADE;

# 全UserのAccessLogを並べるのは不安がある
CREATE TABLE UserAccessLogs1 (
  LastAccess TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  UserID STRING(32) NOT NULL,
) PRIMARY KEY (LastAccess DESC, UserID);

# UserごとのAccessLogなら大丈夫だろう
CREATE TABLE UserAccessLogs2 (
  UserID STRING(32) NOT NULL,
  LastAccess TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY (UserID, LastAccess DESC);

# 全UserのAccessLogを並べたいなら、Shardingする必要が出てくることがある
CREATE TABLE UserAccessLogs3 (
  ShardID INT64 NOT NULL,
  LastAccess TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  UserID STRING(32) NOT NULL,
) PRIMARY KEY (ShardID, LastAccess DESC, UserID);

CREATE TABLE Shops (
  ShopID STRING(32) NOT NULL,
  ShopName STRING(128) NOT NULL,
  Status STRING(16) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true)
) PRIMARY KEY (ShopID);

CREATE TABLE Items (
  ShopID STRING(32) NOT NULL,
  CategoryID STRING(16) NOT NULL,
  ItemID STRING(32) NOT NULL,
  ItemName STRING(128) NOT NULL,
  Status STRING(16) NOT NULL,
  Price INT64 NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true)
) PRIMARY KEY (ShopID, CategoryID, ItemID),
  INTERLEAVE IN PARENT Shops ON DELETE CASCADE;

CREATE INDEX ItemsByStatusAndUpdatedAtDesc
ON Items (
  Status,
  UpdatedAt DESC
);

# Itemsは更新頻度が低いので、UpdatedAtにIndexを貼っている
CREATE INDEX ItemsByUpdatedAtDesc
ON Items (
  UpdatedAt DESC
);

# Usersの子どもにしているが、Shopsの子どもにするもしくは親なしという選択肢もある
CREATE TABLE Orders (
  UserID STRING(32) NOT NULL,
  ShopID STRING(32) NOT NULL,
  OrderID STRING(32) NOT NULL,
  Price INT64 NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true)
) PRIMARY KEY (UserID, ShopID, OrderID),
  INTERLEAVE IN PARENT Users ON DELETE CASCADE;

CREATE INDEX OrdersByShopIDAndUserID
ON Orders (
    ShopID,
    UserID,
    OrderID
);

CREATE INDEX OrdersByUserIDAndUpdatedAtDesc
ON Orders (
    UserID,
    UpdatedAt DESC
);

CREATE TABLE OrderDetails (
  UserID STRING(32) NOT NULL,
  ShopID STRING(32) NOT NULL,
  OrderID STRING(32) NOT NULL,
  CategoryID STRING(16) NOT NULL,
  ItemID STRING(32) NOT NULL,
  OrderDetailID STRING(32) NOT NULL,
  Price INT64 NOT NULL,
  Quantity INT64 NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
  UpdatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true)
) PRIMARY KEY (UserID, ShopID, OrderID, CategoryID, ItemID, OrderDetailID),
  INTERLEAVE IN PARENT Orders ON DELETE CASCADE;
