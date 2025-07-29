-- 用户表
CREATE TABLE `hzm_user`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_name`   varchar(128) NOT NULL COMMENT '用户名',
    `password`    varchar(128) NOT NULL COMMENT '密码',
    `role`        tinyint(2)  DEFAULT '0' COMMENT '角色：0-管理员；1-普通用户',
    `email`       varchar(128)          DEFAULT NULL COMMENT '邮件',
    `valid`       tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_username` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 用户数据权限表(执行器维度)
CREATE TABLE `hzm_user_data_permission`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`     bigint(20) NOT NULL COMMENT '用户id',
    `executor_id` bigint(20) DEFAULT NULL COMMENT '执行器id',
    `valid`       tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_userid` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户数据权限表';

-- 执行器表
CREATE TABLE `hzm_executor`
(
    `id`            bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `name`          varchar(128) NOT NULL DEFAULT '' COMMENT '执行器名称',
    `app_key`       varchar(128) NOT NULL DEFAULT '' COMMENT '执行器标识',
    `registry_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '注册方式：0-自动；1-手动',
    `valid`         tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='执行器表';

-- 执行器节点表
CREATE TABLE `hzm_executor_node`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `executor_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '执行器id',
    `address`     varchar(128) NOT NULL DEFAULT '' COMMENT '节点地址',
    `status`      tinyint(2) NOT NULL DEFAULT '0' COMMENT '节点状态：0-离线；1-在线',
    `valid`       tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_executorid` (`executor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='执行器节点表';

-- 任务表
CREATE TABLE `hzm_job`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `executor_id`     bigint(20) NOT NULL DEFAULT '0' COMMENT '执行器id',
    `name`            varchar(128) NOT NULL DEFAULT '' COMMENT '任务名称',
    `schedule_type`   tinyint(2) NOT NULL DEFAULT '0' COMMENT '调度类型：1-cron表达式；2-极简表达式',
    `schedule_value`  varchar(32)  NOT NULL DEFAULT '' COMMENT '调度值：如果scheduleType是1，则为cron表达式；如果scheduleType是2，则为极简表达式值',
    `description`     varchar(128) NOT NULL DEFAULT '' COMMENT '任务描述',
    `parameters`      varchar(512) NOT NULL DEFAULT '' COMMENT '任务参数',
    `head`            varchar(128) NOT NULL DEFAULT '' COMMENT '负责人',
    `status`          tinyint(2) NOT NULL DEFAULT '0' COMMENT '任务状态：0-未启动；1-已启动',
    `router_strategy` tinyint(2) NOT NULL DEFAULT '0' COMMENT '路由策略：0-轮询；1-随机；2-故障转移',
    `cron_entry_id`   int(11) DEFAULT NULL COMMENT '注册到cron中的id',
    `valid`           tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time`     datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`     datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY               `idx_executorid` (`executor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务表';

-- 任务日志表
CREATE TABLE `hzm_job_log`
(
    `id`                    bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `job_id`                bigint(20) NOT NULL DEFAULT '0' COMMENT '任务id',
    `executor_id`           bigint(20) NOT NULL DEFAULT '0' COMMENT '执行器id',
    `executor_node_address` varchar(128)      DEFAULT NULL COMMENT '执行器节点地址',
    `parameters`            varchar(512)      DEFAULT NULL COMMENT '任务参数',
    `schedule_time`         datetime          DEFAULT NULL COMMENT '任务调度日志时间',
    `status`                tinyint(2) NOT NULL DEFAULT '0' COMMENT '任务调度日志状态：0-待调度；1-任务执行中；2-任务结束',
    `handle_code`           int(11) DEFAULT NULL COMMENT '任务结果编码',
    `handle_msg`            varchar(512)      DEFAULT NULL COMMENT '任务结果消息',
    `finish_time`           datetime          DEFAULT NULL COMMENT '任务完成时间',
    `valid`                 tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time`           datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`           datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY                     `idx_jobid` (`job_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务日志表';
