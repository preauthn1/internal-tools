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

    public interface IPlayerGameRecordService
    {

        PagedList<PlayerGameRecordViewModel> GetByPagingFilter(int userID, QueryPlayerGameRecordApiRequestParameterModel filter);

    }
}
