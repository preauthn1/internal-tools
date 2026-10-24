namespace MSC.Models.ViewModels
{
    public class VendorApiQueryMerchantRequestRecordRequestModel : VendorApiCommonRequestModel
    {
        public VendorApiQueryMerchantRecordRequestModel RequestInfo { get; set; }
    }


    public class VendorApiQueryMerchantRecordRequestModel : BaseVendorAPIRequestModel
    {
        public string MerchantRequestID { get; set; }
    }

}
