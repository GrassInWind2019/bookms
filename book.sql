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
