

##  red_envelope_goods

| 字段            | 字段名称       | 数据类型                 | 说明                                              |
| --------------- | -------------- | ------------------------ | ------------------------------------------------- |
| id              | 自增ID         | bigint(20)               | 自增ID                                            |
| envelope_no     | 红包编号       | varchar(32)              | 红包编号,红包唯一标识                             |
| envelope_type   | 红包类型       | tinyint(2)               | 红包类型：普通红包，碰运气红包,过期红包           |
| username        | 用户名称       | varchar(64)              | 用户名称                                          |
| user_id         | 用户编号       | varchar(40)              | 用户编号, 红包所属用户                            |
| blessing        | 祝福语         | varchar(64)              | 祝福语                                            |
| amount          | 红包总金额     | decimal(30,6)   unsigned | 红包总金额                                        |
| amount_one      | 单个红包金额   | decimal(30,6)   unsigned | 单个红包金额，碰运气红包无效                      |
| quantity        | 红包总数量     | int(10) unsigned         | 红包总数量                                        |
| remain_amount   | 红包剩余金额额 | decimal(30,6)   unsigned | 红包剩余金额额                                    |
| remain_quantity | 红包剩余数量   | int(10) unsigned         | 红包剩余数量                                      |
| expired_at      | 过期时间       | datetime(3)              | 过期时间                                          |
| status          | 红包/订单状态  | tinyint(2)               | 红包/订单状态：0 创建、1   发布启用、2过期、3失效 |
| order_type      | 订单类型       | tinyint(2)               | 订单类型：发布单、退款单                          |
| pay_status      | 支付状态       | tinyint(2)               | 支付状态：未支付，支付中，已支付，支付失败        |
| created_at      | 创建时间       | datetime(3)              | 创建时间                                          |
| updated_at      | 更新时间       | datetime(3)              | 更新时间                                          |



## red_envelope_item

| 字段          | 字段名称           | 数据类型                 | 说明                                       |
| ------------- | ------------------ | ------------------------ | ------------------------------------------ |
| id            | 自增ID             | bigint(20)               | 自增ID                                     |
| item_no       | 红包订单详情编号   | bigint(20)               | 红包订单详情编号                           |
| envelope_no   | 红包编号           | varchar(32)              | 红包编号,红包唯一标识                      |
| recv_username | 红包接收者用户名称 | varchar(64)              | 红包接收者用户名称                         |
| recv_user_id  | 红包接收者用户编号 | varchar(40)              | 红包接收者用户编号                         |
| amount        | 收到金额           | decimal(30,6)   unsigned | 收到金额                                   |
| quantity      | 收到数量           | int(10) unsigned         | 收到数量：对于收红包来说是1                |
| remain_amount | 收到后红包剩余金额 | decimal(30,6)   unsigned | 收到后红包剩余金额                         |
| account_no    | 红包接收者账户ID   | varchar(32)              | 红包接收者账户ID                           |
| pay_status    | 支付状态           | tinyint(2)               | 支付状态：未支付，支付中，已支付，支付失败 |
| created_at    | 创建时间           | datetime(3)              | 创建时间                                   |
| updated_at    | 更新时间           | datetime(3)              | 更新时间                                   |