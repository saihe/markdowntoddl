package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// カラム
type columns struct {
	Name          string
	Type          string
	NotNull       string
	Default       string
	Key           string
	AutoIncrement string
	Extra         string
	Comment       string
}

// 現在日時取得
func now() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

// エントリーポイント
func main() {
	fmt.Printf("%s 処理開始\n", now())
	// ファイル読み込み
	fmt.Printf("読み込みファイル：[%s]\n", os.Args[1])
	readLines := readFile(os.Args[1])
	// DDL作成
	outputFileName, ddl := parseMarkdownToDDL(readLines)
	// ファイル書き出し
	fmt.Printf("書き出しファイル：[%s]\n", outputFileName)
	writeFile(outputFileName, ddl)
	fmt.Printf("%s 処理終了\n", now())
}

// ファイル読み込み
func readFile(filePath string) []string {
	// オープン・クローズ
	fp, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	// 読み込み
	scanner := bufio.NewScanner(fp)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// エラー
	if serr := scanner.Err(); serr != nil {
		fmt.Fprintf(os.Stderr, "ファイル読み込みエラー（ファイル：[%s]、エラー: [%v]）\n", filePath, err)
	}

	return lines
}

// マークダウンからDDL作成
func parseMarkdownToDDL(lines []string) (string, []string) {
	// 返却値
	outputFileName := ""
	ddl := []string{}
	// テーブル定義
	h1Flag := true
	h2Flag := false
	tableName := ""
	tableComment := ""
	// カラム定義
	tableCounter := 0
	columnFlag := false
	extraFlag := false

	for rowCounter, line := range lines {
		// h1見出し
		if strings.HasPrefix(line, "# ") {
			outputFileName = strings.TrimSpace(line[1:]) + ".sql"
			continue
		}
		// h2見出し
		if strings.HasPrefix(line, "## ") {
			h1Flag = false
			h2Flag = true
			tableCounter = 0
			tableName = strings.TrimSpace(line[2:])
			ddl = append(ddl, fmt.Sprintf("drop table if exists `%s`;", tableName))
			ddl = append(ddl, fmt.Sprintf("create table `%s` (", tableName))
			continue
		}
		// h2見出しまで無視
		if h1Flag {
			continue
		}
		// h2見出し後の空行無視
		if 1 > len(strings.TrimSpace(line)) {
			continue
		}

		// h3見出し
		if strings.HasPrefix(line, "### ") {
			h2Flag = false
			if strings.HasSuffix(line, "_columns") {
				columnFlag = true
				extraFlag = false
			}
			if strings.HasSuffix(line, "_extra") {
				columnFlag = false
				extraFlag = true
			}
			continue
		}
		// h3見出しまで無視
		if h2Flag {
			// h2見出し内容
			if 0 < len(strings.TrimSpace(line)) {
				tableComment = tableComment + line
			}
			continue
		}
		// h3見出し後の空行無視
		if 1 > tableCounter && 1 > len(strings.TrimSpace(line)) {
			continue
		}
		// h3見出し内容表戦闘２行無視
		if columnFlag && 2 > tableCounter {
			tableCounter++
			continue
		}
		// h3見出し内容
		if 0 < len(strings.TrimSpace(line)) {
			if columnFlag {
				col := parseTalbeToColumn(line)
				if rowCounter+1 == len(lines) || 1 > len(strings.TrimSpace(lines[rowCounter+1])) {
					ddl = append(ddl, fmt.Sprintf("%s %s %s %s %s %s %s", col.Name, col.Type, col.NotNull, col.Default, col.AutoIncrement, col.Key, col.Comment))
					// create table 終了
					if 0 < len(tableComment) {
						ddl = append(ddl, fmt.Sprintf(") ENGINE = InnoDB DEFAULT CHARSET utf8 comment '%s';", tableComment))
					} else {
						ddl = append(ddl, fmt.Sprintf(") ENGINE = InnoDB DEFAULT CHARSET utf8;"))
					}
				} else {
					ddl = append(ddl, fmt.Sprintf("%s %s %s %s %s %s %s,", col.Name, col.Type, col.NotNull, col.Default, col.AutoIncrement, col.Key, col.Comment))
				}
			}
			if extraFlag {
				ddl = append(ddl, fmt.Sprintf("alter table `%s` add %s;", tableName, line[2:]))
			}
		}
	}

	return outputFileName, ddl
}

// 表の行からカラム定義を抽出
func parseTalbeToColumn(line string) columns {
	splited := strings.Split(line, "|")

	colName := ""
	if 0 < len(strings.TrimSpace(splited[0])) {
		colName = fmt.Sprintf("`%s`", strings.TrimSpace(splited[0]))
	}

	colType := ""
	if 0 < len(strings.TrimSpace(splited[1])) {
		colType = strings.TrimSpace(splited[1])
	}

	colNotNull := ""
	if 0 < len(strings.TrimSpace(splited[2])) {
		colNotNull = "not null"
	}

	colDefault := ""
	if 0 < len(strings.TrimSpace(splited[3])) {
		colDefault = fmt.Sprintf("default %s", strings.TrimSpace(splited[3]))
	}

	colKey := ""
	if 0 < len(strings.TrimSpace(splited[4])) {
		colKey = strings.TrimSpace(splited[4])
	}

	colAutoIncrement := ""
	if 0 < len(strings.TrimSpace(splited[5])) {
		colAutoIncrement = "auto_increment"
	}

	colExtra := ""
	if 0 < len(strings.TrimSpace(splited[6])) {
		colExtra = strings.TrimSpace(splited[6])
	}

	colComment := ""
	if 0 < len(strings.TrimSpace(splited[7])) {
		colComment = fmt.Sprintf("comment '%s'", strings.TrimSpace(splited[7]))
	}

	return columns{
		Name:          colName,
		Type:          colType,
		NotNull:       colNotNull,
		Default:       colDefault,
		Key:           colKey,
		AutoIncrement: colAutoIncrement,
		Extra:         colExtra,
		Comment:       colComment,
	}
}

// ファイル書き出し
func writeFile(outputFileName string, lines []string) {
	// オープン・クローズ
	fp, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	for _, line := range lines {
		var err error
		if strings.HasPrefix(line, "`") {
			_, err = fp.WriteString("  " + strings.TrimSpace(line) + "\n")
		} else {
			_, err = fp.WriteString(strings.TrimSpace(line) + "\n")
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ファイル書き込みエラー（ファイル：[%s]、エラー：[%v]）\n", outputFileName, err)
		}
	}
}
