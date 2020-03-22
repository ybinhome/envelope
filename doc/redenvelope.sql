

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;



-- ----------------------------
-- Table structure for red_envelope_goods
-- ----------------------------
drop table IF EXISTS `red_envelope_goods`;
create TABLE `red_envelope_goods` (
                                      `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
                                      `envelope_no` varchar(32) NOT NULL COMMENT '红包编号,红包唯一标识 ',
                                      `envelope_type` tinyint(2) NOT NULL COMMENT '红包类型：普通红包，碰运气红包,过期红包',
                                      `username` varchar(64) DEFAULT NULL COMMENT '用户名称',
                                      `user_id` varchar(40) NOT NULL COMMENT '用户编号, 红包所属用户 ',
                                      `blessing` varchar(64) DEFAULT NULL COMMENT '祝福语',
                                      `amount` decimal(30,6) unsigned NOT NULL DEFAULT '0.000000' COMMENT '红包总金额',
                                      `amount_one` decimal(30,6) unsigned NOT NULL DEFAULT '0.000000' COMMENT '单个红包金额，碰运气红包无效',
                                      `quantity` int(10) unsigned NOT NULL COMMENT '红包总数量 ',
                                      `remain_amount` decimal(30,6) unsigned NOT NULL DEFAULT '0.000000' COMMENT '红包剩余金额额',
                                      `remain_quantity` int(10) unsigned NOT NULL COMMENT '红包剩余数量 ',
                                      `expired_at` datetime(3) NOT NULL COMMENT '过期时间',
                                      `status` tinyint(2) NOT NULL COMMENT '红包/订单状态：0 创建、1 发布启用、2过期、3失效',
                                      `order_type` tinyint(2) NOT NULL COMMENT '订单类型：发布单、退款单 ',
                                      `pay_status` tinyint(2) NOT NULL COMMENT '支付状态：未支付，支付中，已支付，支付失败 ',
                                      `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
                                      `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON update CURRENT_TIMESTAMP(3) COMMENT '更新时间',
                                      `origin_envelope_no` varchar(32) NULL DEFAULT '' COMMENT '原订单号' AFTER `updated_at`,
                                      PRIMARY KEY (`id`) USING BTREE,
                                      UNIQUE KEY `envelope_no_idx` (`envelope_no`) USING BTREE,
                                      KEY `id_user_idx` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1273 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for red_envelope_item
-- ----------------------------
drop table IF EXISTS `red_envelope_item`;
create TABLE `red_envelope_item` (
                                     `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
                                     `item_no` varchar(32) NOT NULL COMMENT '红包订单详情编号 ',
                                     `envelope_no` varchar(32) NOT NULL COMMENT '红包编号,红包唯一标识 ',
                                     `recv_username` varchar(64) DEFAULT NULL COMMENT '红包接收者用户名称',
                                     `recv_user_id` varchar(40) NOT NULL COMMENT '红包接收者用户编号 ',
                                     `amount` decimal(30,6) unsigned NOT NULL DEFAULT '0.000000' COMMENT '收到金额',
                                     `quantity` int(10) unsigned NOT NULL COMMENT '收到数量：对于收红包来说是1 ',
                                     `remain_amount` decimal(30,6) unsigned NOT NULL DEFAULT '0.000000' COMMENT '收到后红包剩余金额',
                                     `account_no` varchar(32) NOT NULL COMMENT '红包接收者账户ID',
                                     `pay_status` tinyint(2) NOT NULL COMMENT '支付状态：未支付，支付中，已支付，支付失败 ',
                                     `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
                                     `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON update CURRENT_TIMESTAMP(3) COMMENT '更新时间',
                                     `desc` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
                                     PRIMARY KEY (`id`) USING BTREE,
                                     UNIQUE KEY `item_no_idx` (`item_no`) USING BTREE,
                                     KEY `envelope_no_idx` (`envelope_no`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=141 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;


SET FOREIGN_KEY_CHECKS = 1;
