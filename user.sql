create database if not exists `bookms_user`;
use bookms_user;
drop table if exists `bookms_user_user`;

create table `bookms_user_user` (
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
`last_login_time` datetime not null,
`biography` varchar(500) not null default '',
primary key (`id`),
unique key `account` (`account`),
unique key `email` (`email`),
unique key `phone` (`phone`)
) engine=InnoDB auto_increment=2 default charset=utf8mb4;