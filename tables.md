# テーブル定義

テーブル一覧　論理名（物理名）

* [ユーザー（user）](#user)

## user

ユーザー情報を格納するテーブル。

### columns

name|type|not null|default|key|auto increment|comment
---|---|---|---|---|---|---
id|smallint|||primary key|y|連番
mail_address|varchar(255)|y||unique key||メールアドレス。
user_code|smallint||||||社員番号。　入退社で、同じ社員番号でも違う社員が存在することを想定している。
user_last_name|varchar(32)|||||姓。
user_first_name|varchar(32)|||||名。
passphrase|varchar(128)||||unique key||メールアドレスとパスワードからハッシュ化。
division_code|tinyint||||||所属コード。
created_at|datetime|yes|current_timestamp|||作成日時。
updated_at|datetime|yes|current_timestamp on update current_timestamp|||更新日時。
