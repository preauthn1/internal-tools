# rory-game.backend

原始路径：`rory-game.backend`　|　文件：**1062**　|　主语言：C# 950, JSON 46, YAML 11

## 目录构成

```
PGS.Models/ViewModels/  — 245
PGS.BL/Services/  — 78
PGS.DA/Repository/  — 74
PGS.Models/TableEntity/  — 67
PGS.BL/Interface/  — 66
PGS.Testing/Services/  — 56
PGS.Utility/Helper/  — 43
PGS.Models/Enums/  — 41
PGS.Models/Models/  — 39
PGS.API.AdminSys/Controllers/  — 34
helm-charts/backend-api/  — 20
PGS.API.ExternalVendor/Controllers/  — 19
PGS.API.Game/Controllers/  — 18
PGS.Models/Models.Core/  — 14
```

## 关键文件

- CI/构建：`.gitlab-ci.yml`、`DockerFiles/Dockerfile`
- 服务端主程序：`PGS.API.Vendor/Program.cs`、`PGS.API.Vendor/Startup.cs`、`PGS.Backend.Utility/Program.cs`

## 主要文件

- `PGS.Models/CultureText/i18nKeys.cs`
- `PGS.BL/Services/GameVendorAPI/AGGIntegrationService.cs`
- `PGS.BL/Services/API.Public/GameLaunchService.cs`
- `PGS.BL/Services/GameVendorAPI/FunkyIntegrationService.cs`
- `PGS.BL/Services/GameVendorAPI/AGGLiveIntegrationService.cs`
- `PGS.BL/Services/API.ExternalVendor/GameVendorSingleWalletService.cs`
- `PGS.Web.Modules/SwaggerCustomize/PGSSwaggerQueryParameterIgnoreFilter.cs`
- `PGS.BL/Services/GameVendorAPI/FishingIntegrationService.cs`

## 说明

该工程为平台的一个组件。仅收录文本源码与配置，图片、音频等二进制资源未收录；凭据、密钥、内网地址与域名已按实际内容替换。

