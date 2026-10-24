namespace MSC.Models.ViewModels
{
    public class CardProviderCardGenerateRequestModel 
    {

        public CardProviderCardGenerateRequestDetailModel RequestInfo { get; set; }

        public string SignKey { get; set; }

    }


    public class CardProviderCardGenerateRequestDetailModel : BaseProviderAPIRequestModel
    {
        public string TransactionNo { get; set; }

        public decimal Point { get; set; }

        public string Memo { get; set; }
    }



    public class ProviderCardGenerateResultViewModel
    {
        public string SystemOrderNo { get; set; }

        public string CardNo { get; set; }

        public string CardPassword { get; set; }

        public decimal Point { get; set; }

        public string Currency { get; set; }
    }


}
