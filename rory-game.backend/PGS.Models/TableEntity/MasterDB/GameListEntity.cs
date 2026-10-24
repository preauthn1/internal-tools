using target.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    [IsNotAutoKeyTable]
    public class GameListEntity
    {
        [Key]
        public int GameID { get; set; }

        public string GameName { get; set; }

        public string GameCode { get; set; }

        public string ThirdPartyGameCode { get; set; }

        public string LanguageCode { get; set; }


        /// <summary>
        /// 對應 GameTypeEnum
        /// </summary>
        public int GameType { get; set; }

        /// <summary>
        /// 對應 GameVendorTypeEnum
        /// </summary>
        public int GameVendorType { get; set; }

        public int Status { get; set; }

        public bool IsSkipInsertZeroBet { get; set; }

        public bool IsHideAfterAmountInfo { get; set; }


        public DateTime CreateTime { get; set; }


    }
}
