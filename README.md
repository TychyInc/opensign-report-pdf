# OpenSign Report PDF

OpenSign請求書PDF生成ライブラリ

## 概要

OpenSignの利用料金請求書をPDF形式で生成するGoライブラリです。基本料金と従量課金に対応し、日本の請求書フォーマットに準拠したPDFを出力します。

## 特徴

- 基本料金（月額2,000円、100件まで）※カスタマイズ可能
- 従量課金（100件超過分は1件あたり100円）※カスタマイズ可能
- 消費税（10%）自動計算
- 日本語フォント（Noto Sans JP）対応
- A4サイズPDF出力
- テーブル形式の明細表示
- io.ReaderベースのAPI（メモリ効率的）

## インストール

```bash
go get github.com/ryuyama/opensign-report-pdf
```

## ライブラリとしての使用方法

```go
package main

import (
    "io"
    "log"
    "os"
    "time"
    
    opensignreport "github.com/ryuyama/opensign-report-pdf"
)

func main() {
    config := opensignreport.Config{
        CustomerName:  "株式会社サンプル",
        ReceiptNumber: "INV-2024-001",
        IssueDate:     time.Now(),
        MonthlyCount:  150,
        // Optional: Customize fees (defaults: BasicFee=2000, FreeUsageCount=100, UsageFeeRate=100)
        BasicFee:      2000,
        FreeUsageCount: 100,
        UsageFeeRate:  100,
    }
    
    reader, err := opensignreport.GenerateInvoice(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Save to file
    file, err := os.Create("invoice.pdf")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    
    _, err = io.Copy(file, reader)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Invoice generated: invoice.pdf")
}
```

## コマンドラインツール

### ビルド

```bash
go build -o opensign-invoice ./cmd/opensign-invoice
```

### 使用方法

```bash
# 基本的な使用方法
./opensign-invoice -name "株式会社サンプル" -count 150 -receipt-no "INV-2024-001"

# 日付を指定
./opensign-invoice -name "株式会社サンプル" -count 150 -receipt-no "INV-2024-001" -date "2024-01-15"
```

### オプション

- `-name`: 顧客名（必須）
- `-count`: 月間利用件数（必須）
- `-receipt-no`: 請求書番号（必須）
- `-date`: 発行日（YYYY-MM-DD形式、省略時は当日）

## API設定

### Config構造体

- `CustomerName` (必須): 顧客名
- `ReceiptNumber` (必須): 請求書番号
- `IssueDate`: 発行日（省略時は当日）
- `MonthlyCount` (必須): 月間利用件数
- `BasicFee`: 基本料金（デフォルト: 2,000円）
- `FreeUsageCount`: 無料利用件数（デフォルト: 100件）
- `UsageFeeRate`: 従量課金単価（デフォルト: 100円/件）

## 料金計算

- 基本料金: 2,000円（100件まで）※カスタマイズ可能
- 従量課金: 100件を超えた分は1件あたり100円 ※カスタマイズ可能
- 消費税: 10%

### 計算例（デフォルト料金の場合）

- 50件の場合: 2,000円 + 消費税200円 = 2,200円
- 150件の場合: 2,000円 + 5,000円（50件×100円） + 消費税700円 = 7,700円

## フォント設定

このライブラリはNoto Sans JPフォントを使用します。`ttf/NotoSansJP-VariableFont_wght.ttf`にフォントファイルを配置してください。

## 請求書の内容

- 宛名
- 請求書番号
- 発行日
- 請求金額（税込）
- 但し書き: "OpenSign利用料として"
- 内訳
  - 基本利用料
  - 利用料（従量課金分）
  - 小計
  - 消費税（10%）
  - 合計（税込）
- 発行者情報
  - 株式会社ティヒ
  - 住所: 東京都目黒区鷹番3-21-14 303
  - 登録番号: T9011001155105

## ライセンス

MIT License

## 作者

ryuyama
