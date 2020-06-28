package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// カラム
type columns struct {
	Name          string
	Type          string
	NotNull       string
	Default       string
	Key           string
	AutoIncrement string
	Comment       string
}

// エントリーポイント
func main() {
	fmt.Printf("読み込みファイル：[%s]", os.Args[1])
	// ファイル読み込み
	readLines := readFile(os.Args[1])
	// DDL作成
	outputFileName, ddl := parseMarkdownToDDL(readLines)
	// ファイル書き出し
	fmt.Println(outputFileName)
	fmt.Printf("%v", ddl)
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
		fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", filePath, err)
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
	tableComment := ""
	// カラム定義
	tableCounter := 0

	for rowCounter, line := range lines {
		// h1見出し
		if strings.HasPrefix(line, "# ") {
			outputFileName = strings.TrimSpace(line[1:])
			continue
		}
		// h2見出し
		if strings.HasPrefix(line, "## ") {
			h1Flag = false
			h2Flag = true
			tableCounter = 0
			ddl = append(ddl, fmt.Sprintf("drop table if exists `%s`;", strings.TrimSpace(line[2:])))
			ddl = append(ddl, fmt.Sprintf("create table `%s` (", strings.TrimSpace(line[2:])))
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
		if 2 > tableCounter {
			tableCounter++
			continue
		}
		// h3見出し内容
		if 0 < len(strings.TrimSpace(line)) {
			col := parseTalbeToColumn(line)
			ddl = append(ddl, fmt.Sprintf("%s %s %s %s %s %s %s", col.Name, col.Type, col.NotNull, col.Default, col.AutoIncrement, col.Key, col.Comment))
			if rowCounter+1 == len(lines) || 1 > len(strings.TrimSpace(lines[rowCounter+1])) {
				// create table 終了
				if 0 < len(tableComment) {
					ddl = append(ddl, fmt.Sprintf(") ENGINE = InnoDB DEFAULT CHARSET utf8 comment `%s`;", tableComment))
				} else {
					ddl = append(ddl, fmt.Sprintf(") ENGINE = InnoDB DEFAULT CHARSET utf8;"))
				}
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

	colComment := ""
	if 0 < len(strings.TrimSpace(splited[6])) {
		colComment = fmt.Sprintf("comment `%s`", strings.TrimSpace(splited[6]))
	}

	return columns{
		Name:          colName,
		Type:          colType,
		NotNull:       colNotNull,
		Default:       colDefault,
		Key:           colKey,
		AutoIncrement: colAutoIncrement,
		Comment:       colComment,
	}
}
