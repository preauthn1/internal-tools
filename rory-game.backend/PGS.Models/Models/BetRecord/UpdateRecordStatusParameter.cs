using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models
{
    public class UpdateRecordStatusParameter
    {

        public UpdateRecordStatusParameter()
        {
            this.StatusRules = new UpdateRecordStatusCheckRules();
        }


        public string TransactionNoPGS { get; set; }

        public int TargetStatus { get; set; }

        public string Memo { get; set; }


        public UpdateRecordAmountParameter UpdateAmountValue { get; set; }

        public UpdateRecordStatusCheckRules StatusRules { get; set; }

    }


    public class UpdateRecordAmountParameter
    {
        public decimal ValidAmount { get; set; }

        public decimal PayoutAmount { get; set; }
    }



    public class UpdateRecordStatusCheckRules
    {
        /// <summary>
        /// 限定允許的注單狀態
        /// </summary>
        public List<int> AvailableStatus { get; set; }

        /// <summary>
        /// 不允許的注單狀態
        /// </summary>
        public List<int> NotAvailableStatus { get; set; }
    }




}
