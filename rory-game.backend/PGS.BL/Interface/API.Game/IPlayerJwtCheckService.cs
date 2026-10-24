using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.ViewModels.Base;

namespace target.BL.Interface
{

    public interface IPlayerJwtCheckService
    {
        bool CheckIsTokenExist(string jwtHash);

        bool CheckIsLobbyTokenExist(string tokenRawText);

        PlayerTokenJwtModel ParsingLobbyToken(string tokenRawText);
    }
}
