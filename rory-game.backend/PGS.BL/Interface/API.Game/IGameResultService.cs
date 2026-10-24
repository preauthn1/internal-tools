using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IGameResultService
    {

        GamePlayResultViewModel UpdateSettlementAmount(GamePlayResultRequestModel model);


        PerGameLogViewModel WritePerGameLog(PerGameLogRequestModel model);


    }
}
