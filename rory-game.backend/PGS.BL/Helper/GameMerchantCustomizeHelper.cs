using target.Models;
using System;
using System.Collections.Generic;
using System.Text;
using target.Utility.Extensions;
using target.Models.Models;
using target.Models.Enums;
using target.Models.Table;
using target.BL.Interface;
using target.Utility.Helper;
using Lck.Utility.Extensions.Helper;

namespace target.BL.Services
{
    public class GameMerchantCustomizeHelper
    {
      
        /// <summary>
        /// 定義哪些遊戲有啟動紀錄(detail)
        /// * 若遊戲商的方式是 先打API，取得網址 => 有紀錄
        /// * 若遊戲商的方式是 組出靜態網址 => 無紀錄
        /// </summary>
        public static bool IsGameLaunchRecordExist(GameVendorTypeEnum gameType)
        {
            var existRecordGameVendors = new List<GameVendorTypeEnum>
            {
                GameVendorTypeEnum.Funky,
            };

            return existRecordGameVendors.Contains(gameType);
        }

        public static bool IsGameLaunchRecordExist(int gameType)
        {
            return IsGameLaunchRecordExist((GameVendorTypeEnum)gameType);
        }

    }
}
