using MSC.Models.Table;
using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class CardInputRequestModel
    {

        public string SignKey { get; set; }

        public long TS { get; set; }

        public string MerchantToken { get; set; }

        public string CardNo { get; set; }

        public string CardPassword { get; set; }

    }


    public class TriggerNotifyAgainRequestModel
    {
        public int MerchantID { get; set; }

        public string MerchantRequestID { get; set; }

        /// <summary>
        /// 若為true，即使已經 status=成功，也會觸發通知
        /// </summary>
        public bool IsSkipStatusCheck { get; set; }

    }



    public class CardInputFakeTestRequestModel
    {

        public int MerchantID { get; set; }


        /// <summary>
        /// 若要透過此API測試，必須先從db撈一筆未使用的 LinkRequest支付請求單號 做為測試使用
        /// </summary>
        public string MerchantRequestID { get; set; }


    }

    public class CardInputFakeTestTargetCardNumberRequestModel
    {
        public int MerchantID { get; set; }

        /// <summary>
        /// 若要透過此API測試，必須先從db撈一筆未使用的 LinkRequest支付請求單號 做為測試使用
        /// </summary>
        public string MerchantRequestID { get; set; }


        public string CardNo { get; set; }


        public string CardPassword { get; set; }


        public string MagicPassword { get; set; }
    }




    public class SignKeyPreviewRequestModel
    {
        public object RequestInfo { get; set; }

        public string MyMerchantSecret { get; set; }

    }



}
