using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;

namespace target.BL.Interface
{
    public interface IGamePayTableImageService
    {

        PagedList<GamePayTableImageViewModel> GetByPagingFilter(int userID, QueryGamePayTableImageApiRequestParameterModel filter);


        CRUDResultModel CreateOrUpdate(int userID, GamePayTableImageUpdateRequestModel model);


        CRUDResultModel Delete(int userID, int targetID);


    }
}
