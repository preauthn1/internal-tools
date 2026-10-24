using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IThirdPartyGameLogicService
    {
        Slot88MD5InfoResultModel GetPlayerInformationByMD5(QuerySlot88MD5RequestModel model);
    }
}
