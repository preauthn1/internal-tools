using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class MerchantRequestProcessDetailInfoRequestModel
    {

        public string MerchantRequestID { get; set; }

        public int MerchantID { get; set; }

    }


    public class MerchantRequestProcessDetailInfoViewModel
    {
        public RequestStepInfoViewModel StepRequestLink { get; set; }

        public RequestStepInfoViewModel StepCardConsume { get; set; }

        public RequestStepInfoViewModel StepNotify { get; set; }

        
        public string MerchantRequestID { get; set; }

        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public string MemberCode { get; set; }

        public string CardNo { get; set; }

        public string ReturnUrl { get; set; }

        public string NotifyUrl { get; set; }




        public decimal PointAmount { get; set; }

        public DateTime CreateTime { get; set; }
    }



    public class RequestStepInfoViewModel
    {
        public DateTime? LastTime { get; set; }

        public bool IsSuccess { get; set; }

    }


}
