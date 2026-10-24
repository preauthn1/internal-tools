using game.Entity.Game;
using target.Models.Enums;
using target.Utility;
using target.Utility.IOC;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;

namespace target.Testing
{

    public class TestDataConfig
    {
        public TestSampleDataModel Default => new TestSampleDataModel();

        public TestSampleData_SingleWallet SingleWallet => new TestSampleData_SingleWallet();

        public TestSampleData_Not_SuperUser_Model Not_Super_User => new TestSampleData_Not_SuperUser_Model();

        public TestSampleData_GameVendor ExternalGameVendor => new TestSampleData_GameVendor();
    }



    public class TestSampleDataModel
    {

        public int UserID = 137;

        public string UserName = "unit_test_001";

        public int PlayerID = 1000035;

        public string PlayerName = "unit_test_001";

        public int ServerID = 999999;

        public int GameID = (int)GameListEnum.DFDC_多福多财;

        public int ThirdPartyGameID = 250001;


        public long TradeID = 20200100000031136;

        public int MerchantID = 175001;

        public string MerchantSecret = [REDACTED]";

        public string Language_Text_Test = "GameList_fun_10291";

        public int OrganizationID = 11;


        public string IPAddress = "0.0.0.0";


        public string TransactionNo_PgsBet = "DFDC210127083125449giepp";


        /// <summary>
        /// 單元測試刪測試資料 可能會判斷比對此欄位相同 刪除資料
        /// </summary>
        public string DisplayName__MustDelete = "unitTest建立";

        /// <summary>
        /// 單元測試刪測試資料 可能會判斷比對此欄位相同 刪除資料
        /// </summary>
        public long TradeID_MustDelete = 20190100000071179;

        /// <summary>
        /// 單元測試刪測試資料 可能會判斷比對此欄位相同 刪除資料
        /// </summary>
        public decimal BetAmount_MustDelete = 3377;

    }

    public class TestSampleData_Not_SuperUser_Model
    {

        public int MerchantID = 175010;

        public int PlayerID = 1000217;

        public string PlayerName = "unit_test_003";

    }



    public class TestSampleData_SingleWallet
    {
        public int MerchantID = 175022;

        public string MerchantSecret = [REDACTED]";

        public int PlayerID = 1015906;

        public decimal BalanceAmount = 99999;

        public string PlayerName = "UnitTest_SW001";

        /// <summary>
        /// 這種token沒有過期的問題，純粹只是要能確認是對方的token即可
        /// </summary>
        public string SessionTokenRaw = "ab184bbf-99bf-4ff6-a09d-280e1695eafe";
        public string SessionTokenPGS = "f33f5ad8-d8fd-4856-83bf-baad0fe002a3";

        public string BetRecordNo = "HHF21013114014694149gp3";

    }


    public class TestSampleData_GameVendor
    {
        public int GameID = 290001;
    }




}