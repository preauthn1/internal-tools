using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class GameListViewModel
    {

        public int GameID { get; set; }

        public string GameName { get; set; }

        public string GameCode { get; set; }

        /// <summary>
        /// 對應 GameTypeEnum
        /// </summary>
        public int GameType { get; set; }

     
        public string GameTypeText { get; set; }

        /// <summary>
        /// 對應 GameVendorTypeEnum
        /// </summary>
        public int GameVendorType { get; set; }


        public string GameVendorTypeText { get; set; }

        public int Status { get; set; }
        
        public bool IsSkipInsertZeroBet { get; set; }

        public DateTime CreateTime { get; set; }

        public string ThirdPartyGameCode { get; set; }

        public bool IsThirdPartyGame { get; set; }


        /// <summary>
        /// 列舉是否有定義
        /// </summary>
        public bool IsGameListEnumDefined { get; set; }


        /// <summary>
        /// 是否語系檔有定義好文字 (只判斷PGS自製遊戲)
        /// </summary>
        public bool IsMultiLanguageDefined { get; set; }


        public string GameImageUrl { get; set; }

    }



}
