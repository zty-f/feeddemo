-- User表
CREATE TABLE `users`
(
    `id`             bigint(0) NOT NULL AUTO_INCREMENT COMMENT 'PK',
    `username`           varchar(32)  NOT NULL,
    `password`       varchar(200)  NULL DEFAULT NULL,
    `follow_count`   bigint(0)  NULL DEFAULT 0,
    `follower_count` bigint(0)  NULL DEFAULT 0,
    `avatar` varchar(255) DEFAULT NULL COMMENT '头像链接',
    PRIMARY KEY (`id`)
);
ALTER TABLE `feed`.`users`
    ADD UNIQUE INDEX `username&password`(`username`, `password`) USING BTREE COMMENT 'username+password的唯一组合索引';


-- 关注表
CREATE TABLE `relations`
(
    `id`           bigint(0) NOT NULL AUTO_INCREMENT,
    `user_id`      bigint(0) NULL DEFAULT 0,
    `following_id` bigint(0) NULL DEFAULT 0,
    PRIMARY KEY (`id`)
);
ALTER TABLE `feed`.`relations`
    ADD UNIQUE INDEX `user_follow`(`user_id`, `following_id`) USING BTREE COMMENT '关注着和被关注者的id构成唯一索引';

-- 博文表
CREATE TABLE `blogs`
(
    `id`             bigint(0) NOT NULL AUTO_INCREMENT  COMMENT 'blog_id',
    `title`          varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
    `content`        VARCHAR(512)  NULL DEFAULT NULL COMMENT '博文内容',
    `favorite_count` bigint(0) NULL DEFAULT 0 COMMENT '点赞量',
    `comment_count`  bigint(0) NULL DEFAULT 0 COMMENT '评论量',
    `user_id`        bigint(0) NOT NULL COMMENT 'FK reference user id',
    `top`            varchar(24) NOT NULL DEFAULT 0 COMMENT '置顶标志字段',
    `create_time`    datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP (0),
    `update_time`    datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP (0),
    PRIMARY KEY (`id`)
);

-- 点赞表
CREATE TABLE `favorites`
(
    `id`          bigint(0) NOT NULL AUTO_INCREMENT,
    `user_id`     bigint(0) NOT NULL,
    `blog_id`    bigint(0) NOT NULL,
    `is_favorite` tinyint(0) NULL DEFAULT 0,
    PRIMARY KEY (`id`)
);
ALTER TABLE `feed`.`favorites`
    ADD UNIQUE INDEX `user_blog`(`user_id`, `blog_id`) USING BTREE COMMENT 'user_id+blog_id的唯一索引';


-- 评论表
CREATE TABLE `comments`
(
    `id`          bigint(0) NOT NULL AUTO_INCREMENT,
    `user_id`     bigint(0) NOT NULL,
    `blog_id`    bigint(0) NOT NULL,
    `content`     varchar(500) DEFAULT '',
    `create_time`    datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP (0),
    PRIMARY KEY (`id`)
);
ALTER TABLE `feed`.`comments`
    ADD INDEX `blog_id`(`blog_id`) USING BTREE COMMENT 'blogId的普通索引';
