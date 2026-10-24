using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class PlayerLoginLogViewModel
    {

        public long Id { get; set; } 

        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public int GameID { get; set; }

        public string GameName { get; set; }

        public int Status { get; set; }

        public string StatusText { get; set; }


        /// <summary>
        /// 另一張table的鍵值
        /// </summary>
        public string RequestUUID { get; set; }


        public int VendorID { get; set; }

        public string VendorName { get; set; }

        /// <summary>
        /// 是否能透過 externalAPILog反查 詳細資訊 (請參考 GameMerchantCustomizeHelper 定義)
        /// </summary>
        public bool IsLaunchRecordExist { get; set; }

        public int HDType { get; set; }

        public string HDBrowser { get; set; }

        public DateTime CreateTime { get; set; }

        public string CreateIPAddress { get; set; }

        public string Memo { get; set; }


    }



}
