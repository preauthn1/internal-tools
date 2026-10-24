using Microsoft.OpenApi.Any;
using Microsoft.OpenApi.Models;
using target.Models;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.ViewModels;
using target.Utility;
using target.Utility.Extensions;
using target.Utility.Helper.HttpContext;
using Swashbuckle.AspNetCore.SwaggerGen;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;


namespace target.API.AdminSys
{

    public static class PGSSwaggerSpecialCaseSetting
    {

        /// <summary>
        /// 請求參數，若名稱符合、swagger UI上會顯示此處定義的欄位說明
        /// </summary>
        public static Dictionary<string, ApiDescriptionModel> APIQueryParameterDescriptionDefinedCN = new Dictionary<string, ApiDescriptionModel>()
        {

            { "orderByColumn" ,  new ApiDescriptionModel("用哪個欄位排序" ,"") },
            { "isASC" ,  new ApiDescriptionModel("正排序、倒排序" ,"") },
            { "signKey" ,  new ApiDescriptionModel("加密驗簽，需將Ts以及部分欄位加密，以檢查防止API偽造" ,"") },
            { "ts" ,  new ApiDescriptionModel("時間戳記timestamp" ,"") },
            { "merchantSecretKey" , new ApiDescriptionModel("运营商密钥" ,"`游戏商预先提供`") },
            { "merchantID" , new ApiDescriptionModel("运营商編號" ,"`游戏商预先提供`") },
            { "playerName" , new ApiDescriptionModel("玩家账号" ," `#格式限制請參考文件說明`") },
            { "playerUniqueID" , new ApiDescriptionModel("玩家唯一識別碼 (系统ID或唯一值帐号)" ,"`#帐号格式限制請參考文件說明`") },
            { "startTime" , new ApiDescriptionModel("查询起始时间 (+0时区) 格式 yyyy-MM-dd HH:mm:ss" ,"") },
            { "endTime" , new ApiDescriptionModel("查询结束时间 (+0时区) 格式 yyyy-MM-dd HH:mm:ss `时间区间不得大于31天` " ,"") },
            { "pageIndex" , new ApiDescriptionModel("查询页数" ,"`从1开始`") },
            { "pageSize" , new ApiDescriptionModel("每页几笔" ,"`限制每次查询10000笔`") },
            { "totalPage" , new ApiDescriptionModel("总页数" ,"") },
            { "totalCount" , new ApiDescriptionModel("总笔数 " ,"") },
            { "gameID" , new ApiDescriptionModel("游戏编号" ,"开放游戏以 `4.6.1 API (Get game list) 為準`") },
            { "roundID" , new ApiDescriptionModel("游戏局号(棋牌游戏)" ,"") },
            { "betID" , new ApiDescriptionModel("投注单号" ,"") },
            { "roundStartTime" , new ApiDescriptionModel("游戏局起始时间(棋牌游戏)" ,"") },
            { "roundEndTime" ,new ApiDescriptionModel( "游戏局结束时间(棋牌游戏)" ,"") },
            { "currency" , new ApiDescriptionModel("商户货币" ,"") },
            { "gameName" ,new ApiDescriptionModel( "游戏名称" ,"") },
            { "isExist" ,new ApiDescriptionModel( "资料是否已存在" ,"") },
            { "status" , new ApiDescriptionModel("状态" ,"") },
            { "betTime" , new ApiDescriptionModel("下注时间" ,"") },
            { "betAmount" , new ApiDescriptionModel("投注金额" ,"") },
            { "payoutAmount" , new ApiDescriptionModel("派彩金额" ,"") },
            { "validAmount" ,new ApiDescriptionModel( "有效流水" ,"") },
            { "winLossAmount" ,new ApiDescriptionModel( "输赢金额(payoutAmount - betAmount - taxAmount)" ,"") },
            { "taxAmount" , new ApiDescriptionModel("税收金额" ,"") },
            { "settlementTime" ,new ApiDescriptionModel( "结算时间" ,"") },
            { "totalWinLossAmount" , new ApiDescriptionModel("总输嬴" ,"") },
            { "totalValidAmount" , new ApiDescriptionModel("总有效流水" ,"") },
            { "totalPayoutAmount" , new ApiDescriptionModel("总派彩金额" ,"") },
            { "totalBetAmount" ,new ApiDescriptionModel( "总投注金额" ,"") },
            { "totalTaxAmount" ,new ApiDescriptionModel( "总税收金额" ,"") },
            { "link" ,new ApiDescriptionModel( "网址/链结" ,"") },
            { "gameCode" ,new ApiDescriptionModel( "游戏代码" ,"") },
            { "imageUrl" , new ApiDescriptionModel("图片链结" ,"") },
            { "gameType" , new ApiDescriptionModel("游戏分类" ,"") },
            { "gameTypeName" , new ApiDescriptionModel("游戏分类名称" ,"") },
            { "roomLevelFlagID" ,new ApiDescriptionModel( "暫無作用" ,"") },
            { "lang" , new ApiDescriptionModel("语系代码" ,"`(参考文件附录列表)`") },
            { "amount" , new ApiDescriptionModel("交易金额" ," `#格式限制請參考文件說明`") },
            { "sessionToken" , new ApiDescriptionModel("作为玩家每次进入游戏的生命周期，并提供运营商API检查身份合法性，需传入长度15~80字串 (建议使用uuid或hash)" ,"") },
            { "merchantOrderNo" , new ApiDescriptionModel("运营商自定转账单号 (*需唯一)，长度限制50字元" ,"") },
            { "transferType" , new ApiDescriptionModel("转账类型" ,"") },
            { "playerID" , new ApiDescriptionModel("玩家ID" ,"") },
            { "beforeAmount" , new ApiDescriptionModel("交易前金额" ,"") },
            { "afterAmount" ,new ApiDescriptionModel( "交易后金额" ,"") },
            { "createTime" , new ApiDescriptionModel("建立时间" ,"") },
            { "pgsTransactionNo" ,new ApiDescriptionModel( "[游戏商] 系统转账单号" ,"") },
            { "pGSTransactionNo" ,new ApiDescriptionModel( "[游戏商] 系统转账单号" ,"") },
            { "merchantTransactionNo" ,new ApiDescriptionModel( "[营运商] 转账单号" ,"") },
            { "transactionNo" , new ApiDescriptionModel("交易单号" ,"") },
        };

        /// <summary>
        /// 請求參數，若名稱符合、swagger UI上會顯示此處定義的欄位說明
        /// </summary>
        public static Dictionary<string, ApiDescriptionModel> APIQueryParameterDescriptionDefinedEN = new Dictionary<string, ApiDescriptionModel>()
        {

            { "orderByColumn" ,  new ApiDescriptionModel("sort by" ,"") },
            { "isASC" ,  new ApiDescriptionModel("sort direction (ascending、descending)" ,"") },
            { "signKey" ,  new ApiDescriptionModel("加密驗簽，需將Ts以及部分欄位加密，以檢查防止API偽造" ,"") },
            { "ts" ,  new ApiDescriptionModel("timestamp" ,"") },
            { "merchantID" , new ApiDescriptionModel("operator's API ID" ,"`provides in advance`") },
            { "merchantSecretKey" , new ApiDescriptionModel("operator's secret key" ,"`provides in advance`") },
            { "playerName" , new ApiDescriptionModel("player account" ," `#reference rules in document`") },
            { "playerUniqueID" , new ApiDescriptionModel("Player’s Unique ID (systemID or Unique Account)" ,"`#reference rules in document`") },
            { "startTime" , new ApiDescriptionModel("Start time of this query (UTC +0) format: yyyy-MM-dd HH:mm:ss" ,"") },
            { "endTime" , new ApiDescriptionModel("End time of this query (UTC +0) format: yyyy-MM-dd HH:mm:ss ` time interval cannot exceed 31 days` " ,"") },
            { "pageIndex" , new ApiDescriptionModel("pageIndex" ,"`Start from 1`") },
            { "pageSize" , new ApiDescriptionModel("pageSize" ,"`(maximum: 10000)`") },
            { "totalPage" , new ApiDescriptionModel("total page" ,"") },
            { "totalCount" , new ApiDescriptionModel("total count " ,"") },
            { "gameID" , new ApiDescriptionModel("ID of game" ,"please reference `4.6.1 API (Get game list).`") },
            { "roundID" , new ApiDescriptionModel("game round ID (board game)" ,"") },
            { "betID" , new ApiDescriptionModel("unique bet ID (bet order number)" ,"") },
            { "roundStartTime" , new ApiDescriptionModel("Start time of this round(board game)" ,"") },
            { "roundEndTime" ,new ApiDescriptionModel( "End time of this round(board game)" ,"") },
            { "currency" , new ApiDescriptionModel("currency" ,"") },
            { "gameName" ,new ApiDescriptionModel( "game display nmae" ,"") },
            { "isExist" ,new ApiDescriptionModel( "If the account exists or not" ,"") },
            { "status" , new ApiDescriptionModel("status" ,"") },
            { "betTime" , new ApiDescriptionModel("bet time" ,"") },
            { "betAmount" , new ApiDescriptionModel("bet amount" ,"") },
            { "payoutAmount" , new ApiDescriptionModel("payout amount" ,"") },
            { "validAmount" ,new ApiDescriptionModel( "valid amount" ,"") },
            { "winLossAmount" ,new ApiDescriptionModel( "winLoss amount(payoutAmount - betAmount - taxAmount)" ,"") },
            { "taxAmount" , new ApiDescriptionModel("tax amount" ,"") },
            { "settlementTime" ,new ApiDescriptionModel( "settlement time" ,"") },
            { "totalWinLossAmount" , new ApiDescriptionModel("Total WinLoss Amount" ,"") },
            { "totalValidAmount" , new ApiDescriptionModel("Total Valid Amount" ,"") },
            { "totalPayoutAmount" , new ApiDescriptionModel("Total Payout Amount" ,"") },
            { "totalBetAmount" ,new ApiDescriptionModel( "Total Bet Amount" ,"") },
            { "totalTaxAmount" ,new ApiDescriptionModel( "Total Tax Amount" ,"") },
            { "link" ,new ApiDescriptionModel( "link url" ,"") },
            { "gameCode" ,new ApiDescriptionModel( "game Code" ,"") },
            { "imageUrl" , new ApiDescriptionModel("image Url" ,"") },
            { "gameType" , new ApiDescriptionModel("game Type" ,"") },
            { "gameTypeName" , new ApiDescriptionModel("game Type Name" ,"") },
            { "roomLevelFlagID" ,new ApiDescriptionModel( "" ,"*`NOT AVAILABLE FOR NOW`") },
            { "lang" , new ApiDescriptionModel("Language Code" ,"`(#reference codes in document)`") },
            { "amount" , new ApiDescriptionModel("transaction amount" ," `#reference rules in document`") },
            { "sessionToken" , new ApiDescriptionModel("For Operator API identifying players and entering game life circle everytime, need to be made by 15~80 string length(uuid or hash) " ,"") },
            { "merchantOrderNo" , new ApiDescriptionModel("Operator’s order ID, should be unique, the limitation is 50 characters." ,"") },
            { "transferType" , new ApiDescriptionModel("transfer Type" ,"") },
            { "playerID" , new ApiDescriptionModel("player ID" ,"") },
            { "beforeAmount" , new ApiDescriptionModel("Amount before this transaction(transferMode only)" ,"") },
            { "afterAmount" ,new ApiDescriptionModel( "Amount after this transaction" ,"") },
            { "createTime" , new ApiDescriptionModel("create Time" ,"") },
            { "pgsTransactionNo" ,new ApiDescriptionModel( "[Game Provider] Transaction number" ,"") },
            { "pGSTransactionNo" ,new ApiDescriptionModel( "[Game Provider] Transaction number" ,"") },
            { "merchantTransactionNo" ,new ApiDescriptionModel( "[Operator] Transaction number of transfer" ,"") },
            { "transactionNo" , new ApiDescriptionModel("Transaction number" ,"") },
        };


        /// <summary>
        /// 回應參數，若名稱符合、swagger UI上會顯示此處定義的範例值
        /// </summary>
        public static Dictionary<string, string> OutputParameterSampleValueDefined = new Dictionary<string, string>()
        {
            { "link" , "https://example.invalid" },
            { "url" , "https://example.invalid" },
            { "imageUrl" , "https://example.invalid" },
            { "betID" , "1052102551" },            
            { "lang" , "zh-tw" },          
            { "defaultLanguage" , "zh-tw" },
            { "defaultLang" , "zh-tw" },
            { "currency" , "CNY" },
            { "memo" , "備註欄位" },
            { "createBy" , "adminUser" },
            { "updateBy" , "adminUser" },
            { "userName" , "用戶名／用戶帳號" },
            { "displayName" , "顯示名稱" },
            { "password" , "輸入您的密碼" },
            { "confirmPassword" , "再次確認您的密碼" },
            { "apiSecretKey" , "8acc72e29907" },
            { "apiDoaminUrl" , "http://0.0.0.0" },
            { "apiValidatePlayerSessionUrl" , "target/validatePlayerSession" },
            { "ipAddress" , "0.0.0.0" },
            { "userGameDetail" , "10001:2000:-10:20180401,10002:1000:-10:20180402" },
            { "serverName" , "多福多財伺服器_1.0.1" },
            { "merchantName" , "PGS娛樂投注平台" },
            { "organizationName" , "PGS天馬團隊" },
            { "gameName" , "DFDZ多福多財" },
            { "gameCode" , "DFDZ" },
            { "gameTypeName" , "Slots" },
            { "orderNo" , "2019010100751241" },
            { "merchantOrderNo" , "AZY201906553641" },
            { "playerName" , "player_account" },
            { "transactionNo" , "202001011234567" },
            { "merchantSecretKey" , "4c269eb2aa8f690598df" },            
        };


        /// <summary>
        /// 回應參數，若名稱符合、swagger UI上會顯示此處定義的範例值
        /// </summary>
        public static Dictionary<string, long> OutputParameterSampleNumberValueDefined = new Dictionary<string, long>()
        {
            { "merchantID" , 175001 },
            { "pageIndex" , 1 },
            { "pageSize" , 10 },
            { "totalPage" , 5 },
            { "gameID" , 1023 },
            { "serverID" , 102300410 },
            { "tradeID" , 20191200000095331},
            { "gameType" , 1000 },
            { "playerID" , 1000003},
            { "faceID" , 15},
            { "gameTimeSecond" , 360},     
        };


        /// <summary>
        /// 回應參數，若名稱符合、swagger UI上會顯示此處定義的範例值
        /// </summary>
        public static Dictionary<string, decimal> OutputParameterSampleDecimalValueDefined = new Dictionary<string, decimal>()
        {
            { "betAmount" , 10.0M},
            { "validAmount" , 10.0M},
            { "payoutAmount" , 20.0M},
            { "currentAmount" , 5000.75M},
            { "updateAmount" , 2000.10M},
            { "flowAmount" , 100.50M},
            { "taxAmount" , 5.5M},
            { "amount" , 100.0M},
        };


        /// <summary>
        /// 這些變數不會產生文件 因為會有無窮迴圈
        /// </summary>
        public static HashSet<string> OutputParameterIgnoreOutputDocumentColumns = new HashSet<string>()
        {
            "childPermissions"
        };


        public static HashSet<string> SkipOutputRequestBodyList = new HashSet<string>();

        public static HashSet<string> SkipOutputResponseColumnWhenPagedList = new HashSet<string>();

        public static Dictionary<string, ApiDescriptionModel> GetParameterDescriptionByWebsite()
        {
            bool isEN = IsLanguageEN();

            if(isEN)
            {
                return APIQueryParameterDescriptionDefinedEN;
            }
            else
            {
                return APIQueryParameterDescriptionDefinedCN;
            }
        }

        public static bool IsVendorAPI()
        {
            return SettingHelper.WebSiteName.Contains("target.API.Vendor"); 
        }


        public static bool IsLanguageEN()
        {
            var value = MyHttpContext.Current.Request.Cookies[SystemConstant.HeaderKey_Lang];

            if (string.IsNullOrWhiteSpace(value))
                return false;

            return value == LanguageEnum.en_us.ToString();
        }



        public static void AddToPagedIgnoreDictionary(string name)
        {
            name = name.ToLower();

            if (SkipOutputResponseColumnWhenPagedList.Contains(name))
                return;

            SkipOutputResponseColumnWhenPagedList.Add(name);
        }



        public static void AddToRequestBodyIgnoreDictionary(string name)
        {
            name = name.ToLower();

            if (SkipOutputRequestBodyList.Contains(name))
                return;

            SkipOutputRequestBodyList.Add(name);
        }


        public static bool IsExistInSkipOutputResponseColumnWhenPagedList(string className , string columnName)
        {
            string keyName = $"{className.ToLower()}_{columnName.ToLower()}";

            if (SkipOutputResponseColumnWhenPagedList.Contains(keyName))
                return true;

            return false;
        }


        public static bool IsExistInSkipOutputRequestBodyList(string className, string columnName)
        {
            string keyName = $"{className.ToLower()}_{columnName.ToLower()}";

            if (SkipOutputRequestBodyList.Contains(keyName))
                return true;

            return false;
        }


        public static bool IsOutputNumberColumnsExist(string columnName)
        {
            if (OutputParameterSampleNumberValueDefined.ContainsKey(columnName))
                return true;

            return false;
        }

        public static bool IsOutputDecimalColumnsExist(string columnName)
        {
            if (OutputParameterSampleDecimalValueDefined.ContainsKey(columnName))
                return true;

            return false;
        }



        public static bool IsOutputKeyExist(string columnName)
        {
            if (OutputParameterSampleValueDefined.ContainsKey(columnName))
                return true;

            return false;
        }

    }


    public class ApiDescriptionModel
    {

        public ApiDescriptionModel(string desc , string memo = "")
        {
            this.Description = desc;
            this.ExtraMemo = memo;
        }


        public string Description { get; set; }

        /// <summary>
        /// 额外说明
        /// </summary>
        public string ExtraMemo { get; set; }


        public string FullDescription 
        { 
            get
            {
                return this.Description + " "+ this.ExtraMemo;
            }                
        }


    }


    public class SwaggerApiDescriptionEN : Attribute
    {
        public SwaggerApiDescriptionEN(string desc)
        {
            this.Description = desc;
        }

        public string Description { get; set; }
    }

    public class SwaggerApiDescriptionCN : Attribute
    {
        public SwaggerApiDescriptionCN(string desc)
        {
            this.Description = desc;
        }

        public string Description { get; set; }
    }

}
