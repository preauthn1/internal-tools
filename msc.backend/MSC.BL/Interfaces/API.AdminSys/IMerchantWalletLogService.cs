using MSC.BL.Interface;
using MSC.Models;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.Paging;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.API.Point;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.BL.Interface
{
    public interface IMerchantWalletLogService
    {
        PagedList<MerchantWalletLogViewModel> GetByPagingFilter(int userID, QueryMerchantWalletLogApiRequestParameterModel filter);

        MerchantWalletLogViewModel GetByID(int targetID);

    }
}
