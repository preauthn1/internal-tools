using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Utility;

namespace target.BL.Interface
{

    public interface IExternalApiRequestLogRefService
    {

        void Create(ApiRequestLogTemporaryInfoModel model);

        void CreateSingleWalletApiLog(ApiRequestLogTemporaryInfoModel model);
    }
}
