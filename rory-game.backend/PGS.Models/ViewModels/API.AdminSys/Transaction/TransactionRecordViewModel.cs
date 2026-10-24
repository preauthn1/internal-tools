using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class TransactionRecordViewModel
    {
        /// <summary>
        /// 說明
        /// </summary>
        public int TransactionType { get; set; }
        /// <summary>
        /// 錢包模式
        /// </summary>
        public int WalletMode { get; set; }
        /// <summary>
        /// 用戶ID
        /// </summary>
        public int PlayerID { get; set; }
        /// <summary>
        /// 帳號
        /// </summary>
        public string PlayerName { get; set; }
        /// <summary>
        /// 商戶ID
        /// </summary>
        public int MerchantID { get; set; }
        /// <summary>
        /// 商戶名稱
        /// </summary>
        public string MerchantName { get; set; }
        /// <summary>
        /// 當前數量
        /// </summary>
        public decimal CurrentAmount { get; set; }
        /// <summary>
        /// 更新數量
        /// </summary>
        public decimal UpdateAmount { get; set; }
        /// <summary>
        /// 更新後數量
        /// </summary>
        public decimal AfterAmount { get; set; }
        /// <summary>
        /// IP地址
        /// </summary>
        public string IPAddress { get; set; }
        /// <summary>
        /// 創建時間
        /// </summary>
        public DateTime CreateTime { get; set; }
        /// <summary>
        /// 遊戲ID
        /// </summary>
        public int GameID { get; set; }
        /// <summary>
        /// 遊戲名稱
        /// </summary>
        public string GameName { get; set; }



        public string TransactionTypeText { get; set; }
    }

}
