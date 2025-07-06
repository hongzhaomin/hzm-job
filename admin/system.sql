
-- 用户表
CREATE TABLE `user`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_name`   varchar(128) NOT NULL DEFAULT '' COMMENT '用户名/手机号',
    `password`    varchar(128) NOT NULL DEFAULT '' COMMENT '密码',
    `nick_name`   varchar(128) NOT NULL DEFAULT '' COMMENT '昵称',
    `email`       varchar(128) NOT NULL DEFAULT '' COMMENT '邮件',
    `head_img`    varchar(128) NOT NULL DEFAULT '' COMMENT '头像',
    `sex`         tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '性别：0-女；1-男',
    `profession`  varchar(64)  NOT NULL DEFAULT '' COMMENT '职业',
    `province`    varchar(64)  NOT NULL DEFAULT '' COMMENT '省',
    `city`        varchar(64)  NOT NULL DEFAULT '' COMMENT '市',
    `county`      varchar(64)  NOT NULL DEFAULT '' COMMENT '区',
    `valid`       tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_username` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';

INSERT INTO `user` (`user_name`, `password`, `nick_name`, `email`, `head_img`, `sex`, `profession`, `province`, `city`, `county`) VALUES ('15856982413', '1', '煕风', 'test@163.com', '', 1, '程序员', '安徽省', '合肥市', '蜀山区');

-- 菜单表
CREATE TABLE `menu`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `pid`         bigint(20) NOT NULL DEFAULT '0' COMMENT '父级id',
    `type`        tinyint(2) NOT NULL DEFAULT '0' COMMENT '菜单类型：0-菜单；1-按钮；2-目录',
    `perm_tag`    varchar(128) NOT NULL DEFAULT '' COMMENT '菜单权限标识',
    `title`       varchar(128) NOT NULL DEFAULT '' COMMENT '名称',
    `href`        varchar(128) NOT NULL DEFAULT '' COMMENT '路由地址',
    `icon`        varchar(128) NOT NULL DEFAULT '' COMMENT 'icon图标',
    `target`      varchar(32)  NOT NULL DEFAULT '_self' COMMENT '同a标签的target属性',
    `order_num`   int(11)  NOT NULL DEFAULT '0' COMMENT '排序数字',
    `valid`       tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_permtag` (`perm_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='菜单表';

INSERT INTO menu (pid, type, perm_tag, title, href, icon, target, order_num) VALUES (0, 2, 'sys:manage', '系统管理', '', 'fa fa-th-large', '_self', 30);
INSERT INTO menu (pid, type, perm_tag, title, href, icon, target, order_num) VALUES (1, 0, 'sys:menu', '菜单管理', '/menu', 'fa fa-window-maximize', '_self', 33);
INSERT INTO menu (pid, type, perm_tag, title, href, icon, target, order_num) VALUES (1, 0, 'sys:role', '角色管理', '/role', 'fa fa-lock', '_self', 32);
INSERT INTO menu (pid, type, perm_tag, title, href, icon, target, order_num) VALUES (1, 0, 'sys:user', '用户管理', '/user', 'fa fa-user', '_self', 31);
INSERT INTO menu (pid, type, perm_tag, title, href, icon, target, order_num) VALUES (0, 2, 'programmer', '码农客栈', '', 'fa fa-laptop', '_self', 1);
INSERT INTO menu (pid, type, perm_tag, title, href, icon, target, order_num) VALUES (5, 0, 'json:convert', 'Json转换', '/json', 'fa fa-code', '_self', 2);

-- 角色表
CREATE TABLE `role`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(128) NOT NULL DEFAULT '' COMMENT '角色名称',
    `role_tag`    varchar(128) NOT NULL DEFAULT '' COMMENT '角色标识',
    `valid`       tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_roletag` (`role_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色表';

-- 角色菜单关系表
CREATE TABLE `role_menu`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `perm_tag`    varchar(128) NOT NULL DEFAULT '' COMMENT '菜单权限标识',
    `role_tag`    varchar(128) NOT NULL DEFAULT '' COMMENT '角色标识',
    `valid`       tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_roletag_permtag` (`role_tag`, `perm_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色菜单关系表';

-- 用户角色关系表
CREATE TABLE `user_role`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`     bigint(20) NOT NULL DEFAULT '0' COMMENT '用户id',
    `role_tag`    varchar(128) NOT NULL DEFAULT '' COMMENT '角色标识',
    `valid`       tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_userid_roletag` (`user_id`, `role_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户角色关系表';


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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4 COMMENT='执行器表';

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4 COMMENT='执行器节点表';

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
    `valid`           tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '是否可用：1-可用；0-不可用',
    `create_time`     datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`     datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY               `idx_executorid` (`executor_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4 COMMENT='任务表';
