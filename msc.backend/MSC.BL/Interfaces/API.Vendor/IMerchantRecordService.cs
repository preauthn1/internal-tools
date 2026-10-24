using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.DA.Repository;

namespace MSC.BL.Interface
{

    public interface IMerchantRecordService
    {
        MerchantRequestRecordViewModel GetMerchantRequestRecord(VendorApiQueryMerchantRequestRecordRequestModel model);

    }
}
