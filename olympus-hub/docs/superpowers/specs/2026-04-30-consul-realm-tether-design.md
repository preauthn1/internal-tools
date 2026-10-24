# 設計規格：Consul 聖域神聖鏈結 (Consul Realm Tether)

**日期：** 2026-04-30
**狀態：** 已核准

---

## 背景

`HubController::invoke()` 目前執行 artisan 指令時，使用宿主專案 `.env` 中的預設 DB 連線。本功能要在執行指令前，依據 `realmId` 從 Consul KV 取得對應的 DB 連線資訊，並覆蓋 Laravel 的 runtime config，使指令連線至正確的正式環境。

---

## 目標

- 在 `invoke()` 執行 `Artisan::call()` 前，自動切換 DB 連線至對應位面（Realm）
- Consul 位址從宿主專案 `.env` 的 `CONSUL_HTTP_ADDR` 讀取，預設為 `http://consul:8500`
- 僅處理金流位面（Cash Flow Realms）；遊戲位面（Game Realms）留待未來開發

---

## 架構

```
HubController::invoke()
    └── ConsulRealmTetherService::activate($realmId)
            ├── ConsulKvClient::fetch("DB/{REALM}")
            ├── EnvParser::parse($rawValue)
            └── DatabaseConfigOverrider::apply($vars)
```

---

## 元件清單

| 元件 | 路徑 | 職責 |
|---|---|---|
| `ConsulRealmTetherService` | `src/Services/ConsulRealmTetherService.php` | 流程協調：取得 KV → 解析 → 覆蓋 |
| `ConsulKvClient` | `src/Consul/ConsulKvClient.php` | 對 Consul HTTP API 發出 GET 請求，回傳原始字串值 |
| `EnvParser` | `src/Support/EnvParser.php` | 將 `KEY=VALUE\n...` 格式解析成 associative array |
| `DatabaseConfigOverrider` | `src/Support/DatabaseConfigOverrider.php` | 更新 Laravel `config('database.connections.mysql.*')` 並呼叫 `DB::purge()` |

---

## Config

`config/olympus.php` 新增：

```php
'consul' => [
    'base_url' => env('CONSUL_HTTP_ADDR', 'http://consul:8500'),
],
```

---

## Consul KV 規格

- **Key 映射規則：** `realmId` → `DB/{strtoupper($realmId)}`
  - 例：`tpower` → `DB/TPOWER`、`vworld` → `DB/VWORLD`
- **API 端點：** `GET {base_url}/v1/kv/DB/{REALM}?raw`
  - `?raw` 讓 Consul 直接回傳純文字，免除 base64 decode
- **Value 格式（.env 多行字串）：**
  ```
  DB_CONNECTION=mysql
  DB_PORT=3306
  DB_DATABASE=xsgcasino
  DB_HOST_READ=...
  DB_HOST_WRITE=...
  DB_USERNAME_READ=...
  DB_PASSWORD_READ=...
  DB_USERNAME_WRITE=...
  DB_PASSWORD_WRITE=...
  ```

---

## 覆蓋策略

使用 Laravel runtime config 覆蓋，不呼叫 `putenv()`：

```php
config([
    'database.connections.mysql.read.host'      => [$vars['DB_HOST_READ']],
    'database.connections.mysql.read.username'  => $vars['DB_USERNAME_READ'],
    'database.connections.mysql.read.password'  => $vars['DB_PASSWORD_READ'],
    'database.connections.mysql.write.host'     => [$vars['DB_HOST_WRITE']],
    'database.connections.mysql.write.username' => $vars['DB_USERNAME_WRITE'],
    'database.connections.mysql.write.password' => $vars['DB_PASSWORD_WRITE'],
    'database.connections.mysql.database'       => $vars['DB_DATABASE'],
    'database.connections.mysql.port'           => $vars['DB_PORT'],
]);
DB::purge('mysql');
```

> 若宿主 `database.php` 未使用 read/write split，Laravel 會忽略 `read.*`/`write.*` 子鍵，不影響運作。

理由：宿主 artisan commands 應透過 `config('database.*')` 存取 DB，而非直接讀 `env()`。

---

## 錯誤處理

所有 Exception 皆在 `HubController::invoke()` 的 `catch` 統一處理（最外層處理原則），內層只拋轉。

| 情境 | Exception |
|---|---|
| Consul 無法連線 / HTTP 錯誤 | `ConsulConnectionException` |
| Key 不存在（Consul 回傳 404） | `RealmTetherNotFoundException` |
| Value 為空或解析失敗 | `InvalidTetherDataException` |

---

## 範圍限制

- 僅覆蓋 DB 連線（`database.connections.mysql.*`）
- GameRealms（`fon666`, `vplus33-game`）不在本次範圍內，`invoke()` 若收到 GameRealm 的請求，維持現有行為（不做連線切換）
- 不修改宿主專案的 `.env` 檔案，僅覆蓋 runtime config

---

## 測試計畫

| 測試項目 | 類型 |
|---|---|
| `ConsulKvClient` 正常取回 value | Unit（mock HTTP）|
| `ConsulKvClient` 404 時拋出 `RealmTetherNotFoundException` | Unit |
| `ConsulKvClient` 連線失敗時拋出 `ConsulConnectionException` | Unit |
| `EnvParser` 正確解析多行字串 | Unit |
| `EnvParser` 空字串時拋出 `InvalidTetherDataException` | Unit |
| `DatabaseConfigOverrider` 正確更新 config 並呼叫 `DB::purge()` | Unit（mock DB facade）|
| `ConsulRealmTetherService::activate()` 整合流程 | Unit |
| `HubController::invoke()` 成功路徑 | Feature（mock Service）|
| `HubController::invoke()` Consul 失敗時回傳 JSON error | Feature |
