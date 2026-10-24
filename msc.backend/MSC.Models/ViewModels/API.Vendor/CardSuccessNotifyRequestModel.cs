using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class CardSuccessNotifyRequestModel : VendorApiCommonRequestModel
    {

        public CardSuccessNotifyRequestInfo_V1_Model RequestInfo { get; set; }

    }


    public class CardSuccessNotifyRequestInfo_V1_Model : BaseVendorAPIRequestModel
    {
        public string MSCOrderID { get; set; }

        public string MerchantRequestID { get; set; }

        public string CardNo { get; set; }

        public decimal Point { get; set; }

        public string MemberCode { get; set; }

        /// <summary>
        /// 成功輸入卡號時間
        /// </summary>
        public string CreateTime { get; set; }

        public string IPAddress { get; set; }

    }


}
