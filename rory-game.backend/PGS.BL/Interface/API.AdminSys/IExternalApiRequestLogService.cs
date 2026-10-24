using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;

namespace target.BL.Interface
{
    public interface IExternalApiRequestLogService
    {

        PagedList<ExternalApiRequestLogViewModel> GetByPagingFilter(int userID, QueryExternalApiRequestLogApiRequestParameterModel filter);

        ExternalApiRequestLogViewModel GetByRequestID(string requestID);

        ExternalApiRequestLogViewModel GetByID(long id);


    }
}
