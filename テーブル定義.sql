drop table if exists `user`;
create table `user` (
  `id` smallint   auto_increment primary key comment '',
  `mail_address` varchar(255) not null   unique key comment '',
  `user_code` smallint     comment '',
  `user_last_name` varchar(32)     comment '',
  `user_first_name` varchar(32)     comment '',
  `passphrase` varchar(128)    unique key comment '',
  `division_code` tinyint     comment '',
  `created_at` datetime not null default current_timestamp   comment '',
  `updated_at` datetime not null default current_timestamp on update current_timestamp   comment ''
) ENGINE = InnoDB DEFAULT CHARSET utf8 comment 'ユーザー情報を格納するテーブル。';
alter table `user` add unique key(`user_last_name`, `user_first_name`);
