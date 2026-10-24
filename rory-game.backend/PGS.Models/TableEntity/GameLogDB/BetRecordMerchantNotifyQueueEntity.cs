using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    public class BetRecordMerchantNotifyQueueEntity
    {
        [Key]
        public long Id { get; set; }

        /// <summary>
        /// 請參考 BetRecordMerchantNotifyQueueTypeEnum
        /// </summary>
        public int Type { get; set; }
        public int GameID { get; set; }
        public int MerchantID { get; set; }
        public int PlayerID { get; set; }
        public string BetRecordNoPGS { get; set; }
        public DateTime NextNotifyTime { get; set; }
        public int RetryTimes { get; set; }
        public bool IsSuccess { get; set; }
        public string Memo { get; set; }
        public string DebugInfo { get; set; }
        public DateTime CreateTime { get; set; }
        public string LastRequestUUID { get; set; }


    }
}
