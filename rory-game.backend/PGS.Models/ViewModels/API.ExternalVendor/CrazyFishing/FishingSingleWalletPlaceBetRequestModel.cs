using Newtonsoft.Json.Linq;
using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class FishingSingleWalletPlaceBetRequestModel
    {

        /// <summary>
        /// 平台方管理的玩家唯一id (我方使用 PlayerID)
        /// </summary>
        public string Account { get; set; }


        public string Token { get; set; }

        /// <summary>
        /// 游戏渠道标识
        /// </summary>
        public string Channel { get; set; }

        /// <summary>
        /// 投注金額
        /// </summary>
        public long Stake { get; set; }

        /// <summary>
        /// 投注單號
        /// </summary>
        public string OrderId { get; set; }


        public string Sign { get; set; }

    }


    public class FishingSingleWalletPlaceBetResponseModel : FishingAPIResponseModel
    {       
        public decimal Balance { get; set; }

    }


}
