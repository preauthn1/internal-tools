# 發現與決策 (Olympus Hub)

## 需求
- 提供不同 Production 環境的連線管理。
- 中二病 (Chunibyo) 命名與風格。
- Laravel 套件形式，Web UI 介面，開發環境專用。
- 不暴露內部類別給外部使用。

## 研究發現
- **專案環境 1 (apis)**: Laravel 13, PHP 8.4。
- **專案環境 2 (RD-Studio-01)**: Laravel 5.5, PHP 7.3。
- **Docker 配置**: 使用 `php-fpm8` 容器處理 `apis`。
- **Composer 整合**: 透過 `path` repository 搭配 `symlink: true` 可以實現無縫本地開發。

## 技術決策
| 決策 | 理由 |
|------|------|
| 支援 PHP 7.1+ | 確保相容於 `RD-Studio-01` (PHP 7.3)。 |
| 支援 Laravel 5.5 - 13 | 覆蓋所有目標專案的版本範圍。 |
| ServiceProvider 環境過濾 | 防止 Web UI 在生產環境被誤啟動。 |

## 遇到的問題
| 問題 | 解決方案 |
|------|---------|
| `illuminate/support` 版本衝突 | 擴展 composer.json 中的版本限制。 |
| Docker 內 composer 指令遺失 | 確認宿主機與容器的執行環境，使用 `docker exec`。 |

## 資源
- [Laravel Package Development Guide](https://example.invalid/docs/master/packages)
- [Composer Path Repositories](https://example.invalid/doc/05-repositories.md#path)

---
*每執行2次查看/瀏覽器/搜尋操作後更新此檔案*
