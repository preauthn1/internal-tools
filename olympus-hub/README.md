# olympus-hub

> 本目录原始 README 见文末，本节为快照补充信息。

原始路径：`olympus-hub`　|　文件 **29**　|　主语言：PHP 18, JSON 2, XML 1

## 目录构成

```
src/Exceptions/  — 3
docs/plans/  — 3
resources/views/  — 3
src/Support/  — 2
docs/superpowers/  — 2
routes/  — 1
src/  — 1
src/Consul/  — 1
src/Services/  — 1
src/Http/  — 1
src/Console/  — 1
tests/Consul/  — 1
```

---

## 原始说明

# ⚡ Olympus Hub (奧林帕斯樞紐) ⚡

> **「歡迎來到數位世界的頂端。從這裡俯瞰，所有的 Production 領域皆為臣民。」**

Olympus Hub 是一個專為 Laravel 開發環境設計的終極管理面板 (Dashboard)。它不提供任何外部類別，而是透過一個 Web UI，讓你在不同 Production 環境之間進行觀測與操作。

> ⚠️ **重要提示**：本套件必須配合 **Consul Key-Value Store** 才能運作。如果您尚未建立 Consul 環境，請參考下方的「啟動 Consul 服務」章節。

## 🏛️ 主要特點

- **開發專用 (Dev-Only)**：僅在 `local` 或 `testing` 環境中啟動，確保生產環境安全。
- **直觀界面 (Dashboard UI)**：仿照 Laravel Telescope 的設計，提供獨立的 Web 控制面板。
- **零干擾 (Non-Intrusive)**：不暴露任何 Class 給外部調用，維持程式碼的簡潔。
- **環境切換 (Environment Switching)**：(即將登場) 在介面上直接切換與查看不同的 Production 資料。

## 🛠️ 安裝方式

在你的 Laravel 專案中，透過 Composer 安裝：

```bash
composer require internel/olympus-hub --dev
```

安裝後，請發布設定檔以進行自定義配置：

```bash
php artisan vendor:publish --tag=olympus-config
```

## 🐳 啟動 Consul 服務 (Start Consul Service)

本套件提供了一個預設的 Shell 腳本來簡化 Docker 容器的啟動流程。請在**宿主機 (Host Machine)** 執行以下指令：

```bash
# 1. 設定認證 Token 環境變數 (請替換為您的 Token)
export INITIAL_MANAGEMENT_TOKEN=[REDACTED]

# 2. 執行啟動腳本
sh vendor/internel/olympus-hub/scripts/start-consul.sh
```

> **注意**：啟動後，請同步將此 Token 設定於主專案 `.env` 的 `CONSUL_HTTP_TOKEN` 中，以便系統進行身分驗證。

## 📥 初始資料匯入 (Initial Data Import)

啟動後，您可以使用 Artisan 指令批次匯入 `realms/` 目錄下的配置：

```bash
# 匯入單一環境 (對應 realms/dev.env)
php artisan olympus:import dev

# 匯入 realms/ 目錄下所有 .env 檔案
php artisan olympus:import
```

此指令會掃描目錄，將檔名轉換為全大寫的 Consul Key（例如 `dev.env` -> `DB/DEV`）並寫入 Consul KV。

> **資料內容說明**：匯入的 `.env` 檔案應包含該環境（Realm）的資料庫連線資訊，例如：
> ```env
> DB_CONNECTION=mysql
> DB_HOST=0.0.0.0
> DB_DATABASE=production_db
> DB_USERNAME=user
> DB_PASSWORD=[REDACTED]
> ```
> 當您在 UI 觸發切換時，系統會讀取這些資訊並動態覆寫主專案的資料庫配置。

## 📜 系統配置 (Configuration)

### 1. 環境變數 (.env)
主專案可以透過以下變數覆蓋預設連線配置：
- `CONSUL_HTTP_ADDR`: Consul 的連線位址（例如：`http://localhost:8500`）。
- `CONSUL_HTTP_TOKEN`: 存取 Consul 的認證 Token (ACL Token)。請使用您在啟動 Docker 時設定的 Token。

### 2. 快捷指令 (Commands)
在 `config/olympus.php` 中，您可以自定義在 UI 上顯示的快捷指令。
- **第一層**：分類 ID (如 `game` 代表遊戲商，`finance` 代表金流商)。
- **第二層**：顯示分類名稱。
- **第三層**：`artisan 指令 => 顯示名稱`。

### 3. 進階設定
- `realms_path`: 本地環境配置檔案 (`.env`) 的存放路徑，預設為專案根目錄下的 `realms/`。請直接在 `config/olympus.php` 中修改。

## ⚙️ 技術運行機制 (Technical Mechanism)

Olympus Hub 整合 **Consul Key-Value Store** 作為分散式配置中心，實現動態環境切換：

1.  **配置存儲**：環境配置存儲於 Consul KV 的 `DB/` 路徑（例如 `DB/TPOWER`）。
2.  **資料格式**：Key 的內容必須符合標準的 `.env` 檔案格式。
3.  **動態啟動流程**：
    *   使用者在 Web UI 觸發「活化 (Activate)」動作。
    *   系統透過 Consul API 獲取配置字串並解析為鍵值對。
    *   動態覆寫 Laravel 的 `database.connections` 配置。
4.  **本地環境保護**：若環境 ID 為 `localhost`，系統將直接沿用本地 `.env` 配置。

## ⚖️ 授權協議

MIT - 眾神是寬宏大量的。

