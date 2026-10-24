# API Contract 模板

> **使用時機**：Phase 2 — API Contract 協商階段
> **填寫者**：PM Agent（初稿），FE + BE Agent（審閱修改）
> **使用方式**：複製此模板，替換所有 `[placeholder]` 為實際內容

---

## API Contract — [Feature Name]

### 版本資訊

| 項目 | 內容 |
|------|------|
| **Contract 版本** | v1.0 |
| **狀態** | 🔄 協商中 / 🔓 已鎖定 |
| **建立者** | PM |
| **建立日期** | [YYYY-MM-DD] |
| **鎖定日期** | [YYYY-MM-DD]（鎖定後填寫） |
| **對應需求文件** | [Feature Name 後端需求文件] |

---

### 共用約定

#### 基礎 URL
```
[環境] [URL]
開發環境：http://localhost:8000/api/v1
測試環境：[待定]
```

#### 認證方式
- **方式**：Laravel Sanctum Bearer Token
- **Header**：`Authorization: Bearer {token}`
- **Token 有效期限**：[待定義]

#### 統一回應格式

**成功回應**：
```json
{
    "success": true,
    "data": { },
    "message": "操作成功"
}
```

**列表回應（含分頁）**：
```json
{
    "success": true,
    "data": [ ],
    "meta": {
        "current_page": 1,
        "per_page": 15,
        "total": 50,
        "last_page": 4
    },
    "message": "操作成功"
}
```

**錯誤回應**：
```json
{
    "success": false,
    "message": "錯誤描述",
    "errors": {
        "field_name": ["錯誤訊息 1", "錯誤訊息 2"]
    }
}
```

#### 共用參數規範

| 參數類型 | 格式 | 範例 |
|---------|------|------|
| 分頁 | `?page={n}&per_page={n}` | `?page=1&per_page=15` |
| 排序 | `?sort_by={field}&sort_order={asc\|desc}` | `?sort_by=created_at&sort_order=desc` |
| 時間格式 | ISO 8601 | `2026-03-24T12:00:00Z` |
| ID 格式 | integer (auto-increment) | `1`, `42` |

#### HTTP 狀態碼對照

| 狀態碼 | 用途 | 觸發條件 |
|--------|------|---------|
| 200 | 成功取得/更新 | GET、PUT、PATCH 成功 |
| 201 | 成功建立 | POST 建立資源成功 |
| 204 | 成功刪除 | DELETE 成功（無回應主體） |
| 400 | 請求格式錯誤 | JSON 格式錯誤、缺少必要 Header |
| 401 | 未認證 | Token 缺失、過期、無效 |
| 403 | 無權限 | Token 有效但角色/權限不足 |
| 404 | 資源不存在 | 請求的資源 ID 不存在 |
| 422 | 驗證失敗 | 輸入資料不符合驗證規則 |
| 429 | 請求過於頻繁 | 超過 Rate Limit |
| 500 | 伺服器錯誤 | 非預期的伺服器端錯誤 |

---

### 端點清單

| # | 方法 | 端點 | 描述 | 認證 | 角色限制 |
|---|------|------|------|------|---------|
| 1 | [METHOD] | /api/v1/[resource] | [描述] | 是/否 | [角色] |
| 2 | | | | | |
| 3 | | | | | |

---

### 端點詳細規格

---

#### API-001：[端點描述]

**基本資訊**
- **端點**：`[METHOD] /api/v1/[resource]`
- **描述**：[功能描述]
- **認證**：是 / 否
- **角色限制**：[角色列表]
- **對應 AC**：AC-BE-[XXX]

**請求 Header**
| Header | 值 | 必填 |
|--------|---|------|
| Authorization | Bearer {token} | 是 |
| Content-Type | application/json | 是 |
| Accept | application/json | 是 |

**請求參數**（URL Parameters，若適用）
| 參數 | 型別 | 必填 | 說明 |
|------|------|------|------|
| [param] | [type] | [是/否] | [說明] |

**請求 Body**（POST/PUT/PATCH）
```json
{
    "field_a": "string — 必填，[說明]",
    "field_b": "number — 選填，[說明，預設值：X]",
    "field_c": "boolean — 選填，[說明，預設值：false]"
}
```

**請求 Body 驗證規則**
| 欄位 | 型別 | 必填 | 驗證規則 | 錯誤訊息 |
|------|------|------|---------|---------|
| field_a | string | 是 | max:255 | field_a 不可超過 255 字元 |
| field_b | integer | 否 | min:0, max:9999 | field_b 必須在 0-9999 之間 |

**成功回應**
- **狀態碼**：[200 / 201]
```json
{
    "success": true,
    "data": {
        "id": 1,
        "field_a": "value",
        "field_b": 100,
        "created_at": "2026-03-24T12:00:00Z",
        "updated_at": "2026-03-24T12:00:00Z"
    },
    "message": "[成功訊息]"
}
```

**錯誤回應**
| 狀態碼 | 情境 | 回應範例 |
|--------|------|---------|
| 401 | 未附帶 Token | `{ "success": false, "message": "未授權" }` |
| 403 | 角色權限不足 | `{ "success": false, "message": "無權限執行此操作" }` |
| 404 | 資源不存在 | `{ "success": false, "message": "找不到指定的資源" }` |
| 422 | 驗證失敗 | `{ "success": false, "message": "驗證失敗", "errors": { "field_a": ["..."] } }` |

**備註**
- [商業邏輯注意事項]
- [效能考量]
- [其他注意事項]

---

#### API-002：[端點描述]

（同上格式，每個端點一個區塊）

---

### FE 型別定義建議

供前端 Agent 參考，建立對應的 TypeScript 型別：

```typescript
// types/[feature].ts

// 回應格式
interface ApiResponse<T> {
  success: boolean
  data: T
  message: string
}

interface PaginatedResponse<T> {
  success: boolean
  data: T[]
  meta: {
    current_page: number
    per_page: number
    total: number
    last_page: number
  }
  message: string
}

// 領域模型
interface [ResourceName] {
  id: number
  // ... 依 API 回應定義
  created_at: string
  updated_at: string
}

// 請求 Body
interface Create[ResourceName]Request {
  // ... 依 API 請求 Body 定義
}

interface Update[ResourceName]Request {
  // ... 依 API 請求 Body 定義
}
```

---

### MSW Mock 建議

供前端 Agent 在 Phase 3 並行開發時使用：

```typescript
// mocks/handlers/[feature].ts
import { http, HttpResponse } from 'msw'

export const [feature]Handlers = [
  // API-001
  http.[method]('/api/v1/[resource]', async ({ request }) => {
    // 依 Contract 定義的回應格式
    return HttpResponse.json({
      success: true,
      data: { /* mock data */ },
      message: '操作成功',
    })
  }),
]
```

---

### 變更紀錄

| 版本 | 日期 | 變更內容 | 提出者 | 決定者 |
|------|------|----------|--------|--------|
| v1.0 | [YYYY-MM-DD] | 初版建立 | PM | PM |
| | | | | |

---

### 審閱簽核

| 角色 | 狀態 | 日期 | 備註 |
|------|------|------|------|
| PM | ✅ 確認 | [日期] | |
| FE | ⬜ 待審閱 / ✅ 確認 / ❌ 有異議 | | |
| BE | ⬜ 待審閱 / ✅ 確認 / ❌ 有異議 | | |
| QA | ⬜ 待審閱 / ✅ 確認 | | |
