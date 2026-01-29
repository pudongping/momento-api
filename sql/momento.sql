-- =============================================
-- 「时光账记」数据库建表语句
-- 数据库版本: MySQL 5.7
-- 字符集: utf8mb4
-- 排序规则: utf8mb4_unicode_ci

-- 规范如下：
-- 1. 时间戳: 所有时间戳字段均为秒级时间戳（INT(11) UNSIGNED）
-- 2. 每张表都会有创建时间（created_at）和更新时间（updated_at）字段，类型均为 INT(11) UNSIGNED NOT NULL DEFAULT 0
-- 3. 每张表除了特殊说明，均会有自增主键 id 字段，并且类型为 BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT
-- 主键命名规范为 “表名单数_id” 比如用户表 表名为 users 主键则为 user_id
-- 4. 如果表中需要自定义排序字段时，统一使用 sort_num 字段，类型为 INT(11) UNSIGNED NOT NULL DEFAULT 0
-- 5. 如果表中需要自定义状态字段时，统一使用 status 字段，类型为 TINYINT(1) NOT NULL DEFAULT 0
-- 6. 如果表中需要用到枚举类型时，统一使用 smallint(1) NOT NULL DEFAULT 0 并且一般起始值不为 0，比如 1-未支付 2-已支付
-- 7. 如果表中需要用到金额类型时，统一使用 decimal(10,2) NOT NULL DEFAULT 0
-- 8. 如果表中需要用到时间、日期类型时，统一使用 int(11) unsigned NOT NULL DEFAULT 0
-- 9. 如果表中需要用到布尔类型时，统一使用 tinyint(1) NOT NULL DEFAULT 1 一般使用 1-是 2-否
-- 10. 软删除时，统一使用 deleted_at 字段，类型为 INT(11) UNSIGNED NOT NULL DEFAULT 0
-- 11. 表字段一般不用 not null，需要设置默认值
-- =============================================

-- =============================================
-- 索引优化说明

-- 1. 所有外键字段都建立了索引
-- 2. 常用查询条件字段建立了索引
-- 3. 时间字段建立了索引，便于按时间范围查询
-- 4. 组合索引遵循最左前缀原则
-- 5. 唯一索引用于保证数据唯一性
-- =============================================

-- 创建数据库
CREATE
DATABASE IF NOT EXISTS `momento` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `momento`;

-- =============================================
-- 1. 用户相关表
-- =============================================

-- 用户表
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `user_id`         bigint(20) UNSIGNED NOT NULL COMMENT '用户唯一ID（雪花算法）',
    `openid`          varchar(255)  NOT NULL DEFAULT '' COMMENT '微信openid',
    `unionid`         varchar(255)  NOT NULL DEFAULT '' COMMENT '微信unionid',
    `nickname`        varchar(50)   NOT NULL DEFAULT '' COMMENT '用户昵称',
    `avatar`          varchar(1024) NOT NULL DEFAULT '' COMMENT '用户头像URL',
    `phone`           varchar(80)   NOT NULL DEFAULT '' COMMENT '手机号',
    `is_disable`      TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否禁用 1-启用 2-禁用',
    `register_ip`     VARCHAR(255)  NOT NULL DEFAULT '' COMMENT '注册时的ip',
    `login_ip`        VARCHAR(255)  NOT NULL DEFAULT '' COMMENT '登录时的ip',
    `last_login_time` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '最近登录的时间（秒级时间戳）',
    `created_at`      int(11) unsigned NOT NULL DEFAULT 0 COMMENT '创建时间（秒级时间戳）',
    `updated_at`      int(11) unsigned NOT NULL DEFAULT 0 COMMENT '更新时间（秒级时间戳）',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `uk_openid` (`openid`),
    KEY               `idx_phone` (`phone`),
    KEY               `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 用户设置表
DROP TABLE IF EXISTS `user_settings`;
CREATE TABLE `user_settings`
(
    `user_setting_id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id`         BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    `background_url`  VARCHAR(1024)  NOT NULL DEFAULT '' COMMENT '首页背景图片URL',
    `budget`          DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '月度预算',
    `created_at`      INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（秒级时间戳）',
    `updated_at`      INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（秒级时间戳）',
    PRIMARY KEY (`user_setting_id`),
    UNIQUE KEY `uk_uid` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户设置表';

-- =============================================
-- 2. 账本相关表
-- =============================================

-- 账本表
DROP TABLE IF EXISTS `account_books`;
CREATE TABLE `account_books`
(
    `book_id`         BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '账本ID',
    `name`            VARCHAR(50) NOT NULL DEFAULT '' COMMENT '账本名称',
    `creator_user_id` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建者用户ID',
    `member_count`    INT(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '成员数量',
    `created_at`      INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（秒级时间戳）',
    `updated_at`      INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（秒级时间戳）',
    PRIMARY KEY (`book_id`),
    KEY               `idx_creator_uid` (`creator_user_id`),
    KEY               `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账本表';

-- 账本成员表
DROP TABLE IF EXISTS `account_book_members`;
CREATE TABLE `account_book_members`
(
    `book_member_id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '账本成员ID',
    `book_id`        BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '账本ID',
    `user_id`        BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    `is_creator`     TINYINT(1) NOT NULL DEFAULT 2 COMMENT '是否为创建者 1-是 2-否',
    `is_default`     TINYINT(1) NOT NULL DEFAULT 2 COMMENT '是否为默认账本 1-是 2-否',
    `status`         ENUM('joined','waiting','rejected') NOT NULL DEFAULT 'joined' COMMENT '成员状态',
    `joined_at`      INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '加入时间（秒级时间戳）',
    `created_at`     INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（秒级时间戳）',
    `updated_at`     INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（秒级时间戳）',
    PRIMARY KEY (`book_member_id`),
    UNIQUE KEY `uk_book_uid` (`book_id`, `user_id`),
    KEY              `idx_uid` (`user_id`),
    KEY              `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账本成员表';

-- 账本邀请表
DROP TABLE IF EXISTS `account_book_invitations`;
CREATE TABLE `account_book_invitations`
(
    `invitation_id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '邀请ID',
    `book_id`       BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '账本ID',
    `inviter_uid`   BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '邀请人UID',
    `target_uid`    BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '被邀请人UID',
    `status`        ENUM('pending','accepted','rejected') NOT NULL DEFAULT 'pending' COMMENT '邀请状态',
    `created_at`    INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（秒级时间戳）',
    `updated_at`    INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（秒级时间戳）',
    PRIMARY KEY (`invitation_id`),
    KEY             `idx_book_id` (`book_id`),
    KEY             `idx_inviter_uid` (`inviter_uid`),
    KEY             `idx_target_uid` (`target_uid`),
    KEY             `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账本邀请表';

-- =============================================
-- 3. 标签相关表
-- =============================================

-- 标签表
DROP TABLE IF EXISTS `tags`;
CREATE TABLE `tags`
(
    `tag_id`     INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '标签ID',
    `user_id`    BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID（系统标签为空）',
    `name`       VARCHAR(20) NOT NULL DEFAULT '' COMMENT '标签名称',
    `color`      VARCHAR(20) NOT NULL DEFAULT '#4CAF50' COMMENT '标签颜色',
    `icon`       VARCHAR(50) NOT NULL DEFAULT 'tag' COMMENT '标签图标',
    `is_system`  TINYINT(1) NOT NULL DEFAULT 2 COMMENT '是否为系统标签 1-是 2-否',
    `type`       ENUM('expense','income') NOT NULL DEFAULT 'expense' COMMENT '标签类型',
    `sort_num`   INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序',
    `created_at` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（秒级时间戳）',
    `updated_at` INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（秒级时间戳）',
    PRIMARY KEY (`tag_id`),
    KEY          `idx_uid_type` (`user_id`, `type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签表';

-- =============================================
-- 4. 交易记录相关表
-- =============================================

-- 交易记录表
DROP TABLE IF EXISTS `transactions`;
CREATE TABLE `transactions`
(
    `transaction_id`    BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '交易记录ID',
    `book_id`           BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '账本ID',
    `user_id`           BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    `type`              ENUM('expense','income') NOT NULL DEFAULT 'expense' COMMENT '交易类型',
    `amount`            DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '金额',
    `tag_id`            INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '标签ID',
    `remark`            VARCHAR(500)   NOT NULL DEFAULT '' COMMENT '备注',
    `transaction_time`  INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '交易时间（秒级时间戳）',
    `is_auto_generated` TINYINT(1) NOT NULL DEFAULT 2 COMMENT '是否自动生成 1-是 2-否',
    `created_at`        INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（秒级时间戳）',
    `updated_at`        INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（秒级时间戳）',
    PRIMARY KEY (`transaction_id`),
    KEY                 `idx_book_id` (`book_id`),
    KEY                 `idx_uid` (`user_id`),
    KEY                 `idx_tag_id` (`tag_id`),
    KEY                 `idx_transaction_time` (`transaction_time`),
    KEY                 `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='交易记录表';

-- =============================================
-- 5. 周期性记账相关表
-- =============================================

-- 周期性记账表
DROP TABLE IF EXISTS `recurring_transactions`;
CREATE TABLE `recurring_transactions`
(
    `recurring_id`         INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '周期性记账ID',
    `book_id`              BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '账本ID',
    `user_id`              BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    `name`                 VARCHAR(50)    NOT NULL DEFAULT '' COMMENT '记账名称',
    `type`                 ENUM('expense','income') NOT NULL DEFAULT 'expense' COMMENT '交易类型',
    `amount`               DECIMAL(10, 2) NOT NULL DEFAULT 0.00 COMMENT '金额',
    `tag_id`               INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '标签ID',
    `remark`               VARCHAR(255)   NOT NULL DEFAULT '' COMMENT '备注',
    `recurring_type`       ENUM('daily','weekly','monthly','quarterly','yearly') NOT NULL DEFAULT 'monthly' COMMENT '周期类型',
    `recurring_hour`       TINYINT(2) UNSIGNED NOT NULL DEFAULT 9 COMMENT '执行小时（0-23）',
    `recurring_minute`     TINYINT(2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '执行分钟（0-59）',
    `recurring_weekday`    TINYINT(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '星期几（0-6，0为周日）',
    `recurring_month`      TINYINT(2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '月份（1-12）',
    `recurring_day`        TINYINT(2) UNSIGNED NOT NULL DEFAULT 0 COMMENT '日期（1-31）',
    `is_recurring_enabled` TINYINT(1) NOT NULL DEFAULT 2 COMMENT '是否启用 1-是 2-否',
    `last_executed_at`     INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '最后执行时间（秒级时间戳）',
    `next_execution_time`  INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '下一次执行时间（秒级时间戳）',
    `created_at`           INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（秒级时间戳）',
    `updated_at`           INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（秒级时间戳）',
    PRIMARY KEY (`recurring_id`),
    KEY                    `idx_book_id` (`book_id`),
    KEY                    `idx_uid` (`user_id`),
    KEY                    `idx_recurring_type` (`recurring_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='周期性记账表';

-- =============================================
-- 6. 节日相关表
-- =============================================

-- 节日表
DROP TABLE IF EXISTS `festivals`;
CREATE TABLE `festivals`
(
    `festival_id`   INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '节日ID',
    `user_id`       BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    `festival_name` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '节日名称',
    `festival_date` INT(8) UNSIGNED NOT NULL DEFAULT 0 COMMENT '节日日期（YYYYMMDD格式）',
    `is_show_home`  TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否在首页显示 1-是 2-否',
    `created_at`    INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（秒级时间戳）',
    `updated_at`    INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（秒级时间戳）',
    PRIMARY KEY (`festival_id`),
    KEY             `idx_uid` (`user_id`),
    KEY             `idx_festival_date` (`festival_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='节日表';

-- =============================================
-- 7. 文件上传相关表
-- =============================================

-- 文件上传记录表
DROP TABLE IF EXISTS `upload_files`;
CREATE TABLE `upload_files`
(
    `file_id`       BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '文件ID',
    `user_id`       BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    `relative_path` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '文件相对路径',
    `absolute_url`  VARCHAR(255) NOT NULL DEFAULT '' COMMENT '文件绝对URL',
    `file_size`     INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '文件大小（字节）',
    `file_type`     VARCHAR(20)  NOT NULL DEFAULT '' COMMENT '文件类型',
    `business_type` VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '业务类型',
    `upload_time`   INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '上传时间（秒级时间戳）',
    `created_at`    INT(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（秒级时间戳）',
    PRIMARY KEY (`file_id`),
    KEY             `idx_uid` (`user_id`),
    KEY             `idx_upload_time` (`upload_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件上传记录表';

-- =============================================
-- 初始化系统标签数据
-- =============================================
INSERT INTO `tags` (`name`, `color`, `icon`, `is_system`, `type`, `sort_num`)
VALUES ('买菜', '#4CAF50', 'cart-filled', 1, 'expense', 1),
       ('房贷', '#2196F3', 'home-filled', 1, 'expense', 2),
       ('孩子', '#FF9800', 'staff-filled', 1, 'expense', 3),
       ('餐饮', '#F44336', 'shop-filled', 1, 'expense', 4),
       ('交通', '#9C27B0', 'star-filled', 1, 'expense', 5),
       ('其他', '#607D8B', 'more-filled', 1, 'expense', 6),
       ('工资', '#4CAF50', 'wallet-filled', 1, 'income', 7),
       ('奖金', '#FF9800', 'gift-filled', 1, 'income', 8),
       ('投资', '#2196F3', 'medal-filled', 1, 'income', 9),
       ('其他收入', '#607D8B', 'flag-filled', 1, 'income', 10);