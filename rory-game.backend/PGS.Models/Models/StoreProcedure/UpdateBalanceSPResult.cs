using target.Models.Models;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class UpdateBalanceSPResult
    {
        public decimal BeforeAmount { get; set; }


        /// <summary>
        /// 更新後金額 (法幣)
        /// </summary>
        public decimal AfterAmount { get; set; }


        public bool IsSuccess
        {
            get
            {
                var code = (int)this.Code;
                return code == 1 || code == 0;
            }
        }


        public ApiErrorCode Code { get; set; }

    }


    public class UpdateBalanceAndConvertRateResultModel
    {
        public decimal BeforeAmount { get; set; }

        public decimal AfterAmountXXXX { get; set; }

        public decimal AfterGamePoint { get; set; }


    }



}
