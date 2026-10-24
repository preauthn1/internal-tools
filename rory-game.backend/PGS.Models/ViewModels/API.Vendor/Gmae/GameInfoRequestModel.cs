using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 請求參數
    /// </summary>
    public class GameInfoRequestModel : BaseVendorRequestParameter
    {

    }


    /// <summary>
    /// 回應參數
    /// </summary>
    public class GameInfoViewModel
    {

        /// <summary>
        /// 遊戲編號
        /// </summary>
        public int GameID { get; set; }


        /// <summary>
        /// 遊戲代碼
        /// </summary>
        public string GameCode { get; set; }


        public string GameName { get; set; }


        public string ImageUrl { get; set; }



        /// <summary>
        /// 遊戲分類
        /// </summary>
        public int GameType { get; set; }

        public string GameTypeName { get; set; }

    }


}
