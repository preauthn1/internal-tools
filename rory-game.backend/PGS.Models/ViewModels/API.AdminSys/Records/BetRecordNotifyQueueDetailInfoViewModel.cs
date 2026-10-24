using target.Models.Enums;
using target.Models.Table;
using System;
using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class BetRecordNotifyQueueDetailInfoViewModel
    {
        public long Id { get; set; }

        /// <summary>
        /// BetRecordMerchantNotifyQueueTypeEnum
        /// </summary>
        public int Type { get; set; }
        public string TypeName { get; set; }
        public string BetRecordNoPGS { get; set; }
        public string LastRequestUUID { get; set; }
        public int RetryTimes { get; set; }
        public bool IsSuccess { get; set; }
        public List<SingleWalletMerchantApiRequestLogEntity> Logs {get;set;}

    }






}
