# 遊戲平台 Prototype — target 轉帳模式對接

使用 Laravel + Vue Starter Kit（Fortify 認證）建立的簡易遊戲平台，對接 target 遊戲商的**轉帳模式**。

## 功能與遊戲流程（轉帳模式）

- **會員系統**：註冊（E-Mail／帳號／密碼／確認密碼）、帳號登入；註冊即發初始金 **10,000**。
- **遊戲大廳**（公開首頁 `/`）：顯示遊戲商遊戲列表；登入後右上顯示會員名稱與平台餘額。
- **進入遊戲**（#3）：點遊戲於新分頁開啟——伺服端先確保遊戲商玩家綁定（平台帳號 ≠ 遊戲商 playerName，格式 `{prefix}_{亂數}_{user_id}`，不外流平台帳號）、**把平台餘額全額轉入**遊戲商錢包，再 302 導向遊戲啟動連結。
- **離開遊戲自動結算**（#5 轉出）：轉帳模式下遊戲商不會回呼「遊戲關閉」，故以**玩家回到大廳＝停止遊戲**為信號，自動把遊戲商錢包餘額**全額轉回平台**並以 toast 告知。
- **遊戲紀錄同步**（#4）：見下方「遊戲記錄同步 Command」；玩家於 `/records` 檢視自己的紀錄。
- **錢包紀錄（帳本）**：每筆餘額異動（初始金／轉入／轉出／退回）寫入 `wallet_transactions` 帳本並關聯遊戲商；玩家於 `/wallet` 檢視。
- **轉帳一致性與對帳**：進出遊戲的轉帳以 `game_transfers` 記錄流水／狀態（`reference` 唯一、冪等）；跨系統逾時或遊戲商回「狀態不合法」時不直接失敗，改由排程 `game:reconcile-transfers` 依遊戲商 `/transfer` 對帳補償，`game:audit-wallet-sessions` 稽核錢包工作階段，無法自動恢復者轉人工（`manual_review`）。

> 多遊戲商可插拔：服務層以 `App\Services\Games` 的 provider 介面 + driver 設計，target 只是其中一個 driver；遊戲商設定存於 `game_providers` 資料表（非 `config/games.php`），由 `GameProviderManager` 依 `driver` + `settings` 解析。新增遊戲商 = 新增一個 `implements GameProvider` 的 driver + 在 `game_providers` 新增一列（`key`／`name`／`driver`／`settings`）。

## 專案結構

本專案將 **Docker 環境**與 **Laravel 專案**分成兩個目錄：

```
developr_exam/
├── docker/          # Docker Compose 環境（不含應用程式碼）
│   ├── docker-compose.yml
│   ├── .env(.example)   # 容器層變數：MySQL 帳密、對外 port
│   ├── php/             # 自建 PHP 8.4-fpm image（含 Node 22、Composer）
│   └── nginx/           # nginx 設定
└── src/             # Laravel 專案（交付的原始碼）
    └── .env(.example)   # 應用層變數：DB 連線、target 對接資訊
```

技術棧：**Laravel 13 + PHP 8.4 + MySQL 8.0 + Redis 7**、Inertia + Vue 3 + Tailwind。

---

## 環境需求

本專案完全容器化，**主機端只需要 Docker**，不需要另外安裝 PHP / Node / MySQL。

| 需求 | 版本 |
|------|------|
| Docker Engine | 24+（開發使用 29.x） |
| Docker Compose | v2+（開發使用 v5.2） |

對外使用的 port（可於 `docker/.env` 調整）：

| Port | 服務 | 用途 |
|------|------|------|
| `80` | nginx | 應用程式入口 → http://game-lobby.test（需設定 hosts）或 http://localhost |
| `3306` | MySQL | 資料庫對外連線（供 GUI 工具使用） |
| `6379` | Redis | cache / session / queue（供 GUI 工具使用） |
| `5173` | Vite | 前端 dev server（HMR，開發時才需要） |

容器服務一覽：

| 服務 | 說明 |
|------|------|
| `app` | PHP 8.4-fpm，執行 Laravel（image 內含 Node 22 + Composer + phpredis） |
| `nginx` | 反向代理，root 指向 `src/public` |
| `mysql` | MySQL 8.0，資料以 named volume 持久化 |
| `redis` | Redis 7，作為 **cache / session / queue** 的 driver，資料以 named volume 持久化 |
| `scheduler` | 常駐執行 `php artisan schedule:work`，觸發排程任務 |
| `node` | Vite dev server（開發用，非必要啟動） |

> **快取/工作階段/佇列均使用 Redis**（`CACHE_STORE`、`SESSION_DRIVER`、`QUEUE_CONNECTION` 皆為 `redis`）

---

## 安裝步驟

> 以下指令請在 **`docker/` 目錄**下執行（`docker compose` 會讀取該目錄的設定）。
> 由於 `vendor/`、`node_modules/` 不納入版控，clone 後需自行安裝依賴。

### 1. 準備環境變數

```bash
# 容器層變數（MySQL 帳密、port）
cp docker/.env.example docker/.env

# 應用層變數（DB 連線、target 對接資訊）
cp src/.env.example src/.env
```

接著編輯這兩個檔案，填入密碼並**保持一致**（詳見下方「資料庫設定說明」）：

- `docker/.env`：`MYSQL_ROOT_PASSWORD`、`MYSQL_PASSWORD`
- `src/.env`：`DB_PASSWORD`（需等於 `docker/.env` 的 `MYSQL_PASSWORD`）、`PGS_MERCHANT_ID`、`PGS_MERCHANT_SECRET`

### 2. 建置並啟動容器

```bash
cd docker
docker compose up -d --build
```

### 3. 安裝依賴、初始化應用

```bash
# PHP 依賴
docker compose exec app composer install

# 產生 APP_KEY
docker compose exec app php artisan key:generate

# 建立資料表並匯入遊戲商設定（game_providers；target 由 src/.env 的 PGS_* 種入）
docker compose exec app php artisan migrate --seed

# 前端依賴與建置
docker compose exec app npm install
docker compose exec app npm run build
```

### 4. 設定遊戲網域（hosts）

應用對外埠為 **80**，並以 `APP_URL=http://game-lobby.test` 作為正式網址。
`game-lobby.test` 需解析到本機，請於主機的 hosts 檔加入一筆記錄（**需系統管理員權限**）：

```
0.0.0.0    game-lobby.test
```

- Windows（以系統管理員開啟 PowerShell）：
  ```powershell
  Add-Content -Path "$env:windir\System32\drivers\etc\hosts" -Value "`n127.0.0.1`tgame-lobby.test"
  ```
- macOS / Linux：
  ```bash
  echo "0.0.0.0  game-lobby.test" | sudo tee -a /etc/hosts
  ```

> nginx `server_name _` 接受任何 host，故也可直接用 **http://localhost** 開啟（免設定 hosts）；
> 若埠 80 已被其他程式（IIS、Skype 等）占用，可改 `docker/.env` 的 `NGINX_PORT` 後重建 nginx。

### 5. 開啟應用

瀏覽器開啟 **http://game-lobby.test**（或 http://localhost），即可看到平台首頁；`/register`、`/login` 為 Fortify 提供的註冊 / 登入頁。

### （選用）前端開發模式

需要 HMR 熱重載時，啟動 `node` 服務跑 Vite dev server（取代 `npm run build`）。
`node` 服務放在 `dev` profile，預設不啟動，需明確指定 profile：

```bash
docker compose --profile dev up -d node
```

> 啟動 dev server 後，Laravel 會改用 Vite dev server 提供前端資產（`public/hot` 存在時）。
> 若要回到靜態資產模式，停掉 node 服務即可：`docker compose stop node`。

---

## 資料庫設定說明

MySQL 跑在 `mysql` 容器內，**資料庫與帳號會在容器首次啟動時依 `docker/.env` 自動建立**，無需手動建 database。

### 變數對應關係

`docker/.env`（給 MySQL 容器）與 `src/.env`（給 Laravel）必須對應一致：

| Laravel（`src/.env`） | MySQL 容器（`docker/.env`） | 預設值 |
|----------------------|----------------------------|--------|
| `DB_HOST=mysql` | —（容器服務名） | `mysql` |
| `DB_PORT=3306` | —（容器內埠） | `3306` |
| `DB_DATABASE=developr_exam` | `MYSQL_DATABASE` | `developr_exam` |
| `DB_USERNAME=exam` | `MYSQL_USER` | `exam` |
| `DB_PASSWORD=[REDACTED] | `MYSQL_PASSWORD` | **需自行設定且兩邊一致** |

> `DB_HOST` 為 **`mysql`**（容器服務名）而非 `0.0.0.0`——因為 Laravel 是從 `app` 容器連到 `mysql` 容器。

### 從主機端連線（GUI 工具）

若要用 TablePlus / DBeaver 等工具連線，走對外映射的 port：

- Host：`0.0.0.0`　Port：`3306`
- User：`exam`（或 root）　Password：`docker/.env` 所設
- Database：`developr_exam`

### 重建資料庫

```bash
# 清空並重跑所有 migration，並重新種入遊戲商設定
docker compose exec app php artisan migrate:fresh --seed

# 若要連 MySQL 資料一起清掉（含 volume）
docker compose down -v && docker compose up -d
```

---

## 如何執行遊戲記錄同步 Command

平台以 Artisan Command 定時向 target 拉取玩家遊戲記錄並存入本地 `game_records` 資料表，玩家登入後可於「遊戲紀錄」頁（`/records`）查看自己的記錄。

指令以**每個遊戲商的游標（上次同步到的時間）做增量同步**：因遊戲商 `/game-records` API 對查詢時間範圍有上限（≤1 小時），指令會自動把區間**切成 ≤1 小時的片段並逐頁拉取**，再依 `game_provider_accounts` 綁定把 `playerName` 對回平台會員；對不到的玩家（非本平台）會略過。以 `unique(provider, record_id)` + `updateOrCreate` 確保**重複執行不會產生重複資料**。

### 手動執行

```bash
# 增量同步（自游標或預設回看時間到現在）
docker compose exec app php artisan game:sync-records

# 指定遊戲商與時間區間（ISO8601 UTC）
docker compose exec app php artisan game:sync-records --provider=target \
  --since=2026-07-17T08:00:00Z --until=2026-07-17T09:00:00Z
```

指令結束會輸出各遊戲商的摘要：`fetched`（取得筆數）/ `matched`（對到會員）/ `upserted`（寫入）/ `skipped`（略過）。

### 自動排程執行

排程頻率定義於 `src/routes/console.php`，例如：

```php
use App\Models\LogEntry;
use Illuminate\Support\Facades\Schedule;

Schedule::command('game:sync-records')->everyFiveMinutes()->withoutOverlapping();
Schedule::command('game:reconcile-transfers')->everyMinute()->withoutOverlapping();
Schedule::command('game:audit-wallet-sessions')->everyFiveMinutes()->withoutOverlapping();
Schedule::command('model:prune', ['--model' => [LogEntry::class]])->daily()->withoutOverlapping();
```

除了記錄同步，另有轉帳一致性與日誌保養相關排程：

- `game:reconcile-transfers`（每分鐘）：對未結案的轉帳依遊戲商 `/transfer` 對帳補償，逾重試上限轉人工。
- `game:audit-wallet-sessions`（每 5 分鐘）：稽核玩家錢包工作階段，處理殘留在遊戲商錢包的餘額。
- `model:prune`（每日）：清除逾保留期的 DB 日誌（`logs` 資料表），保留天數由 `LOG_DATABASE_RETENTION_DAYS` 控制（預設 30 天）。

實際的定時觸發由 **`scheduler` 容器**負責——它常駐執行 `php artisan schedule:work`，會每分鐘檢查並執行到期的排程任務，**無需在主機設定 cron**。

確認排程已註冊：

```bash
docker compose exec app php artisan schedule:list
```

觀察排程容器的執行輸出：

```bash
docker compose logs -f scheduler
```

---

## 常用指令附錄

```bash
# 服務狀態
docker compose ps

# 進入 app 容器
docker compose exec app bash

# 查看應用日誌（另同時寫入 DB：logs 資料表，LOG_STACK=single,database；
# DB 日誌由每日 model:prune 依 LOG_DATABASE_RETENTION_DAYS（預設 30 天）自動清理）
docker compose exec app php artisan pail    # 即時日誌
docker compose logs -f app                  # 容器日誌

# 列出所有路由（含 Fortify 認證路由）
docker compose exec app php artisan route:list

# 停止 / 移除
docker compose down
```

> **開發工具**：登入後可於 `/dev/provider-tester`（側欄「API 測試」）互動式呼叫遊戲商各端點（games／playerExists／walletAmount／gameLink／transferIn/Out／transferRecord／logoutPlayer／gameRecords）並檢視原始回應與錯誤碼，方便驗證對接。此頁**僅在 `local`／`testing` 環境開放**，正式環境回 404。
