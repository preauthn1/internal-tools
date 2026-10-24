using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.Utility.Helper;
using target.Models;
using target.DA.Repository;

namespace target.BL.Interface
{

    public interface IServerMerchantMappingService
    {
        AllServerMerchantMappingViewModel GetAllServerMerchantMappings();

        ReversServerMerchantResultViewModel ReversMerchantMappings(ReversServerMerchantRequestViewModel model);

        ServerMerchantMappingViewModel GetServerMerchantMappings(int serverID);

        void UpdateServerMerchantMappings(int serverID, ServerAvailableMerchantUpdateRequestModel model);

    }
}
