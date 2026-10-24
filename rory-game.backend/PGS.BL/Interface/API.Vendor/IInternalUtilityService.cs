using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IInternalUtilityService
    {
        List<DemoGameCategoryViewModels> GetMerchantServers(int merchantID);

        LobbyApiGetGamesInfoViewModel GetLobbyGames(LobbyApiGetGamesRequestModel model);
    }
}
