# 概要

マークダウンからMariaDB（MySQL）のDDLを生成する。

## インストール

``` cmd
go get -u github.com/saihe/markdowntoddl/
```

## 起動

``` cmd
markdowntoddl マークダウンファイル
```
  
``` cmd
例）
markdowntoddl tables.md
```

## マークダウン記述方法

### h1見出し

成果物のファイル名を記載する。

### h2見出し

テーブル名を記載する。  
内容を記載した場合、テーブルコメントになる。

### h3見出し

* テーブル名_columnsと記載する。  
表形式でテーブル定義を記載する。  
テーブルヘッダー（カラム）は以下の通り。  
`name|type|not null|default|key|extra|comment`

| | |
---|---
name|カラム物理名
type|データ型（サイズ）
not null|not null制約の有無（y/記載なし）
default|デフォルト値
key|キー制約（p/u/記載なし）
extra|その他制約
comment|カラムコメント

* テーブル名_extraと記載する。  
番号なしリストで記載する。
※現在は複合ユニークキーのみ対応
カラム名はバッククォートで囲う。
