using System;

namespace MSC.Models.ViewModels
{
    public class MerchantCreateRequestModel
    {
        public string DisplayName { get; set; }

        public int Status { get; set; }

        public string Memo { get; set; }

        public string Currency { get; set; }

        public int MaxTimesOfRequestLinkPerMember { get; set; }

        public string AvailableAmounts { get; set; }

        public string CardInputPage_DomainUrl { get; set; }

    }


    public class MerchantUpdateRequestModel
    {
        public int MerchantID { get; set; }

        public string DisplayName { get; set; }

        public int Status { get; set; }

        public string Memo { get; set; }

        public int MaxTimesOfRequestLinkPerMember { get; set; }

        public string AvailableAmounts { get; set; }

        public string CardInputPage_DomainUrl { get; set; }
    }

}
