
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for account
-- ----------------------------
-- | 字段          | 字段名称     | 数据类型 | 说明                                                         |
-- | ------------- | ------------ | -------- | ------------------------------------------------------------ |
-- | id            | 自增ID       | bigint   |                                                              |
-- | account_no    | 账户编号     | char 32  | 账户唯一标识，该表中唯一                                     |
-- | account_name  | 账户名称     | varchar  | 账户对应的名称或者命名，比如xxx积分、xxx零钱                 |
-- | account_type  | 账户类型     | int      | 账户类型，用来区分不同类型的账户：积分账户、会员卡账户、钱包账户、红包账户 |
-- | currency_code | 货币类型编码 | char 3   | 货币类型编码：CNY人民币，EUR欧元，USD美元                    |
-- | user_no       | 用户编号     | varchar  | 用户标识, 账户所属用户                                       |
-- | username      | 用户名称     | varchar  |                                                              |
-- | balance       | 可用余额     | decimal  |                                                              |
-- | status        | 账户状态     | int      | 账户状态：0账户初始化，1启用，2停用                          |
-- | created_at    | 账户创建时间 | time     |                                                              |
-- | Updated_at    | 最后更新时间 | time     |                                                              |
-- ----------------------------
-- ----------------------------
-- Table structure for account
-- ----------------------------
DROP TABLE IF EXISTS `account`;
CREATE TABLE `account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '账户ID',
  `account_no` varchar(32) NOT NULL COMMENT '账户编号,账户唯一标识 ',
  `account_name` varchar(64) NOT NULL COMMENT '账户名称,用来说明账户的简短描述,账户对应的名称或者命名，比如xxx积分、xxx零钱',
  `account_type` tinyint(2) NOT NULL COMMENT '账户类型，用来区分不同类型的账户：积分账户、会员卡账户、钱包账户、红包账户',
  `currency_code` char(3) NOT NULL DEFAULT 'CNY' COMMENT '货币类型编码：CNY人民币，EUR欧元，USD美元 。。。',
  `user_id` varchar(40) NOT NULL COMMENT '用户编号, 账户所属用户 ',
  `username` varchar(64) DEFAULT NULL COMMENT '用户名称',
  `balance` decimal(30,6) unsigned NOT NULL DEFAULT '0.000000' COMMENT '账户可用余额',
  `status` tinyint(2) NOT NULL COMMENT '账户状态，账户状态：0账户初始化，1启用，2停用 ',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `account_no_idx` (`account_no`) USING BTREE,
  KEY `id_user_idx` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=171 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for account_log
-- ----------------------------
-- | 字段              | 字段名称     | 数据类型    | 说明                                                         |
-- | ----------------- | ------------ | ----------- | ------------------------------------------------------------ |
-- | id                | 自增ID       |             |                                                              |
-- | log_no            | 流水编号     | char 32     | 全局不重复字符或数字，唯一性标识                             |
-- | account_no        | 账户编号     | char 32     | 账户标识                                                     |
-- | target_account_no | 目标账户编号 | char 32     |                                                              |
-- | user_no           | 用户编号     | varchar     | 用户标识                                                     |
-- | target_user_no    | 目标用户编号 | varchar     | 交易对象的用户ID，资金去向的账户的所有用户ID                 |
-- | change_type       | 流水交易类型 | int         | 流水交易发生的类型，便于区分和统计：红包转出、转账收入'，类型，0 创建账户，>0 为收入类型，<0 为支出类型', |
-- | change_flag       | 交易变化标识 | tinyint     | 交易变化标识，用来区分是出账还是入账：-1 出账 0 账户创建 1为进账 |
-- | amount            | 交易金额     | decimal     | 该交易涉及的金额                                             |
-- | balance           | 交易后余额   | decimal     | 该交易后的余额                                               |
-- | desc              | 交易备注     | varchar 128 | 交易描述                                                     |
-- | created_at        | 交易时间     | time        | 账户变动发生的时间
-- ----------------------------

-- ----------------------------
-- Table structure for account_log
-- ----------------------------
DROP TABLE IF EXISTS `account_log`;
CREATE TABLE `account_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `trade_no` varchar(32) NOT NULL COMMENT '交易单号 全局不重复字符或数字，唯一性标识 ',
  `log_no` varchar(32) NOT NULL COMMENT '流水编号 全局不重复字符或数字，唯一性标识 ',
  `account_no` varchar(32) NOT NULL COMMENT '账户编号 账户ID',
  `target_account_no` varchar(32) NOT NULL COMMENT '账户编号 账户ID',
  `user_id` varchar(40) NOT NULL COMMENT '用户编号',
  `username` varchar(64) NOT NULL COMMENT '用户名称',
  `target_user_id` varchar(40) NOT NULL COMMENT '目标用户编号',
  `target_username` varchar(64) NOT NULL COMMENT '目标用户名称',
  `amount` decimal(30,6) NOT NULL DEFAULT '0.000000' COMMENT '交易金额,该交易涉及的金额 ',
  `balance` decimal(30,6) unsigned NOT NULL DEFAULT '0.000000' COMMENT '交易后余额,该交易后的余额 ',
  `change_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '流水交易类型，0 创建账户，>0 为收入类型，<0 为支出类型，自定义',
  `change_flag` tinyint(2) NOT NULL DEFAULT '0' COMMENT '交易变化标识：-1 出账 1为进账，枚举',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '交易状态：',
  `decs` varchar(128) NOT NULL COMMENT '交易描述 ',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `id_log_no_idx` (`log_no`) USING BTREE,
  UNIQUE INDEX `idx_tradeno_logno`(`trade_no`, `log_no`) USING BTREE,
  KEY `id_user_idx` (`user_id`) USING BTREE,
  KEY `id_account_idx` (`account_no`) USING BTREE,
  KEY `id_trade_idx` (`trade_no`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1282 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;
