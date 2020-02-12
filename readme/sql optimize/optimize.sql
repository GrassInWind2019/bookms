create database if not exists `bookms`;
use `bookms`;
drop table `bookms_book`;
create table if not exists `bookms_book` (
`id` int(11) not null auto_increment,
`book_name` varchar(200) not null default '',
`identify` varchar(100) not null default '',
`description` varchar(1000) not null default '',
`catalog` varchar(6000) not null default '',
`cover` varchar(1000) not null default '',
`status` int(11) not null default '0',
`sort` int(11) not null default '0',
`create_time` datetime not null,
`doc_count` int(11) not null default '0',
`comment_count` int(11) not null default '0',
`favorite_count` int(11) not null default '0',
`average_score` int(11) not null default '0',
`score_count` int(11) not null default '0',
`comment_people_count` int(11) not null default '0',
`author` varchar(50) not null,
`book_count` int(11) not null default '0',
primary key (`id`),
unique key (`identify`)
)ENGINE=InnoDB auto_increment=2 default charset=utf8mb4;

drop table `bookms_book_category`;
create table if not exists `bookms_book_category` (
`id` int(11) not null auto_increment,
`identify` varchar(100) not null default '',
`category_id` int(11) not null default '0',
primary key (`id`),
key(`identify`),
foreign key(`category_id`) references `bookms_category`(`id`)
)engine=InnoDB auto_increment=2 default charset=utf8;

create table if not exists `bookms_category` (
`id` int(11) not null auto_increment,
`pid` int(11) not null default '0',
`category_name` varchar(30) not null default '',
`description` varchar(500) not null default '',
`icon` varchar(500) not null default '',
`book_count` int(11) not null default '0',
`sort` int(11) not null default '0',
`status` tinyint(1) not null default '0',
primary key (`id`),
unique key `category_name` (`category_name`)
)engine=InnoDB auto_increment=2 default charset=utf8mb4;
-- ALTER TABLE `bookms_category` ADD COLUMN `icon` varchar(500) NOT NULL  DEFAULT '';

create table if not exists `bookms_book_record` (
`id` int(11) not null auto_increment,
`identify` varchar(100) not null default '',
`lend_status` int(11) not null default(0),
`user_id` int(11) not null default(0),
`lend_time` datetime not null,
`return_time` datetime not null,
`lend_count` int(11) not null default(0),
`store_position` varchar(200) not null default '',
primary key (`id`)
)engine=InnoDB auto_increment=1 default charset=utf8;

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
