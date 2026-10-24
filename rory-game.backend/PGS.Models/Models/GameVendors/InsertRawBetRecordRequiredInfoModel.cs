using target.Models.Enums;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{


    public class InsertRawBetRecordRequiredInfoModel
    {

        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public int MerchantID { get; set; }

        public int GameID { get; set; }

        /// <summary>
        /// 投注金額 (法幣)
        /// </summary>
        public decimal BetAmount { get; set; }


        /// <summary>
        /// 投注金額 (遊戲幣)
        /// </summary>
        public decimal BetGamePoint { get; set; }


        /// <summary>
        /// PGS自訂注單號
        /// </summary>
        public string BetRecordNoPGS { get; set; }


        /// <summary>
        /// Debug用 可以跟external logs 匹配用
        /// </summary>
        public string DebugProcessUUID { get; set; }


        public decimal ConvertRate { get; set; }


    }


    public class SettlementBetRecordRequiredInfoModel
    {

        public int MerchantID { get; set; }


        /// <summary>
        /// PGS自訂注單號
        /// </summary>
        public string BetRecordNoPGS { get; set; }


        /// <summary>
        /// Debug用 可以跟external logs 匹配用
        /// </summary>
        public string DebugExtraMemo { get; set; }



        public decimal PayoutAmount { get; set; }



        public decimal ValidAmount { get; set; }


        /// <summary>
        /// 目前只有AGG需要這欄位
        /// </summary>
        public decimal JackpotAmount { get; set; }

        

        /// <summary>
        /// 是否為後台手動操作
        /// </summary>
        public bool IsAdminSysManualOperate { get; set; }

        /// <summary>
        /// 如果手動結算要記錄操作人
        /// </summary>
        public int? UserID { get; set; }

    }







}
