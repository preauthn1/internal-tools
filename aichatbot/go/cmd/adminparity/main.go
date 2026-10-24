// adminparity 是後台設定和檔案匯入的比對（#39）：從空的資料庫只透過後台 API 做出 eval/ 的森日茶屋設定（admin.BuildViaAPI），
// 匯出到 -o 的檔案，給 go/bin/chat 的 BOT_CONFIG_FILE 用，再用 eval/snapshot.py 和直接匯入的逐題比對：
//
//	go -C go run ./cmd/adminparity -o /tmp/admin.json [-eval ../eval]
//
// 翻譯用假模型，不需要模型服務；資料庫是暫存的，不寫入正式資料。
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"aichatbot/internal/admin"
	"aichatbot/internal/merchant"
	"aichatbot/internal/store"
)

func main() {
	out := flag.String("o", "", "匯出檔的路徑（必填）")
	evalDir := flag.String("eval", "", "森日茶屋的 eval/ 資料夾（預設是工作目錄的 eval/，沒有就用 ../eval/）")
	flag.Parse()
	if *out == "" {
		flag.Usage()
		os.Exit(2)
	}
	if *evalDir == "" {
		*evalDir = filepath.Join("..", "eval")
		if _, err := os.Stat("eval/store.json"); err == nil {
			*evalDir = "eval"
		}
	}
	if err := run(*evalDir, *out); err != nil {
		fmt.Fprintln(os.Stderr, "錯誤：", err)
		os.Exit(1)
	}
}

func run(evalDir, out string) error {
	ctx := context.Background()
	cfg, err := merchant.FromEvalDir(evalDir)
	if err != nil {
		return err
	}
	tmp, err := os.MkdirTemp("", "adminparity-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	db, err := store.Open(filepath.Join(tmp, "bot.sqlite"))
	if err != nil {
		return err
	}
	defer db.Close()
	if err := admin.BuildViaAPI(ctx, db, evalDir, cfg); err != nil {
		return err
	}
	c, err := merchant.Load(ctx, db.SQL())
	if err != nil {
		return err
	}
	data, err := merchant.Export(c)
	if err != nil {
		return err
	}
	if err := os.WriteFile(out, data, 0o644); err != nil {
		return err
	}
	fmt.Println("已匯出：", out)
	return nil
}
