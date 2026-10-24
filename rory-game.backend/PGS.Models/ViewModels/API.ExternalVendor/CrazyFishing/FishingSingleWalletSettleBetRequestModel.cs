using Newtonsoft.Json.Linq;
using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class FishingSingleWalletSettleBetRequestModel
    {

        /// <summary>
        /// 平台方管理的玩家唯一id (我方使用 PlayerID)
        /// </summary>
        public int Account { get; set; }



        public string Token { get; set; }

        /// <summary>
        /// 游戏渠道标识
        /// </summary>
        public string Channel { get; set; }

        /// <summary>
        /// 結算金額
        /// </summary>
        public long Stake { get; set; }


        /// <summary>
        /// 實際消耗子彈金額
        /// </summary>
        public long Cost { get; set; }


        /// <summary>
        /// 投注單號
        /// </summary>
        public string OrderId { get; set; }



        public string Sign { get; set; }

    }


    public class FishingSingleWalletSettleBetResponseModel : FishingAPIResponseModel
    {
        public decimal Balance { get; set; }
    }


}
