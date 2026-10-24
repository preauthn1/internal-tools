# ⚡ Olympus Hub - 聖域開發規範 (GEMINI.md) ⚡

歡迎來到奧林帕斯樞紐。本文件定義了維持聖域秩序的至高法則，所有參與開發的神靈（或 AI）皆須嚴格遵守。

## 🏛️ 核心願景 (Core Philosophy)
Olympus Hub 不僅僅是一個工具，它是一個「神格化」的開發環境管理中心。我們追求極致的 **「中二病 (Chunibyo)」** 審美與高度的技術純淨。

## 📜 至高準則 (Primary Mandates)

1.  **神性封閉性 (Divine Encapsulation)**:
    *   本專案不對外暴露任何可被外部調用的類別 (Classes)。
    *   一切操作皆應透過 `/olympus` Web UI 進行。
    *   外部系統應將本專案視為一個獨立的觀測位面。

2.  **開發者聖域 (Dev-Only Sanctuary)**:
    *   本套件嚴禁在 Production 環境中註冊或運行。
    *   所有的 ServiceProvider 註冊必須包含環境檢查 (`local`, `testing`)。

3.  **術語神格化 (Mythical Terminology)**:
    *   禁止使用凡間的技術詞彙，必須轉化為奧林帕斯用語：
        *   Environment -> **Realm (位面)**
        *   Connection -> **Tether (神聖鏈結)**
        *   Log -> **Chronicle (編年史)**
        *   Dashboard -> **Sanctum (聖所)**
        *   Error -> **Calamity (災厄)**

## 🎨 視覺與體驗規範 (UI/UX Conventions)
*   **色調**: 必須以深色調 (Deep Slate/Blue) 為主，搭配霓虹藍 (Electric Blue) 與閃電金 (Bolt Gold)。
*   **氣氛**: 充滿神祕感與權力感，UI 應具備動態效果（如閃電、光暈）。

## 🛠️ 技術架構 (Technical Architecture)
*   **框架**: Laravel Package 結構。
*   **命名空間**: `OlympusHub`。
*   **路由前綴**: `/olympus`。
*   **視圖命名空間**: `olympus::`。

## ⚖️ 變更管理
任何對聖域結構的修改，必須同步更新本文件，確保眾神之志得以延續。
