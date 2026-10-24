using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Models;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IGameServerDistributeRefService
    {

        PlayerLaunchGameServerChooseResultModel ChooseAvailableGameServer(int playerID, int merchantID, int gameID , ServerDistributeOptionalParameter optional_para);

    }
}
