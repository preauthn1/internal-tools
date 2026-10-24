using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IGameLaunchService
    {
        string LaunchGameAndGetGamePageUrlByPlayerID(int playerID, int gameID);

        string LaunchGameAndGetGamePageUrl(GameLaunchRequestModel model);

        GameLaunchViewModel LaunchGameAndCheckPlayerExist(GameLaunchRequestModel model);

        LobbyGenerateLinkResultModel GenerateGameLinkForLobby(LobbyApiGenerateLinkParameterModel model);

        GameServerViewModel GetGameServers(GameServerRequestModel model);

    }
}
