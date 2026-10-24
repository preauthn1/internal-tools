# 金邊相容性處理

將原封包作轉換成金邊可處理的封包內容

### 設定檔參數

| File              | Section       | Name           | Description                |
|:-----------------:|:-------------:|:--------------:|:---------------------------|
| `gameserver.dat`  | `gameserver`  | `lianyunid`    | 聯運ID(舊的，可相容)          |
| `serverinfo.dat`  | `gameserver`  | `lianyunid`    | 聯運ID(新的，以此為主)        |
| `serverinfo.dat`  | `gameserver`  | `platform`     | 營運平台 (金邊: `PHNUM`)     |
| `serverinfo.dat`  | `server`      | `realaddr`     | 實際對外網址                 |
