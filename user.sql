create database if not exists `bookms_user`;
use bookms_user;
-- drop table if exists `bookms_user_user`;

create table if not exists `bookms_user_user` (
`id` int(11) not null auto_increment,
`account` varchar(30) not null default '',
`nickname` varchar(30) not null default '',
`password` varchar(255) not null default '',
`phone` varchar(15) not null default '',
`email` varchar(100) not null default '',
`role` int(11) not null default '2',
`avatar` varchar(255) not null default '',
`status` int(11) not null default '0',
`create_time` datetime not null,
`last_login_time` datetime default null,
`biography` varchar(500) not null default '',
primary key (`id`),
unique key `account` (`account`),
unique key `email` (`email`),
unique key `phone` (`phone`)
) engine=InnoDB auto_increment=2 default charset=utf8mb4;
-- 管理员账号GrassInWind 密码123
insert into bookms_user.bookms_user_user(id,account,nickname,password,phone,email,role,avatar,status,create_time,last_login_time,biography)
values('1', 'GrassInWind', 'GrassInWind', '3632366636663662366437333033336362363134323664633263323861666534303031373738643835316165202cb962ac59075b964b07152d234b70', '13012345678', 'GrassInWind@sina.cn', '0', '', '0', '2020-02-12 14:50:51', '2020-02-12 03:52:34', '普通用户');

create table if not exists `bookms_user_favorite` (
`id` int(11) not null auto_increment,
`user_id` int(11) not null,
`identify` varchar(100) not null default '',
primary key (`id`),
foreign key (`user_id`) references `bookms_user_user`(`id`)
) engine=InnoDB default charset=utf8;

create table if not exists `bookms_user_score`(
`id` int(11) not null auto_increment,
`identify` varchar(100) not null default '',
`user_id` int(11) not null,
`score` int(11) not null default '0',
`create_time` datetime not null,
primary key (`id`),
foreign key (`user_id`) references `bookms_user_user`(`id`)
)engine=InnoDB default charset=utf8;

create table if not exists `bookms_user_comments`(
`id` int(11) not null auto_increment,
`identify` varchar(100) not null default '',
`user_id` int(11) not null,
`content` varchar(3000) not null default '',
`create_time` datetime not null,
primary key (`id`),
foreign key (`user_id`) references `bookms_user_user`(`id`)
)engine=InnoDB default charset=utf8mb4;
