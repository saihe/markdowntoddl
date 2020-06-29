drop table if exists `user`;
create table `user` (
  `id` smallint   auto_increment primary key comment '連番',
  `mail_address` varchar(255) not null   unique key comment 'メールアドレス。',
  `user_code` smallint     comment '社員番号。　入退社で、同じ社員番号でも違う社員が存在することを想定している。',
  `user_last_name` varchar(32)     comment '姓。',
  `user_first_name` varchar(32)     comment '名。',
  `passphrase` varchar(128)    unique key comment 'メールアドレスとパスワードからハッシュ化。',
  `division_code` tinyint     comment '所属コード。',
  `created_at` datetime not null default current_timestamp   comment '作成日時。',
  `updated_at` datetime not null default current_timestamp on update current_timestamp   comment '更新日時。'
) ENGINE = InnoDB DEFAULT CHARSET utf8 comment 'ユーザー情報を格納するテーブル。';
