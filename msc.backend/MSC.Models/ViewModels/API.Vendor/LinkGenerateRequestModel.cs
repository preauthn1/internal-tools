using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class LinkGenerateRequestModel : VendorApiCommonRequestModel
    {

        public LinkGenerateRequestInfoModel RequestInfo { get; set; }

    }




    public class LinkGenerateRequestInfoModel : BaseVendorAPIRequestModel
    {

        /// <summary>
        /// 請求號 或 訂單號
        /// </summary>
        public string RequestID { get; set; }

        public string MemberCode { get; set; }

        public string ReturnUrl { get; set; }

        public string NotifyUrl { get; set; }

        public string Lang { get; set; }

        public decimal Amount { get; set; }

    }


}
