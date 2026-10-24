using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.Models.ViewModels.API.Point;
using MSC.Models.ViewModels.RequestViewModel;

namespace MSC.BL.Interface
{

    public interface ICardConsumeRecordService
    {
        PagedList<CardConsumeRecordViewModel> GetByPagingFilter(int userID ,QueryCardConsumeRecordApiRequestParameterModel filter);

        CardConsumeRecordViewModel GetByID(int userID, int targetID);

    }
}
