using MSC.DA.Repository;
using MSC.Models.Table;
using MSC.Utility;
using MSC.Utility.Helper;
using MSC.Utility.Utility;
using MSC.Utility.Extensions;
using System;
using System.Collections.Generic;


namespace MSC.BL.Services
{

    public interface IExternalApiRequestLogRefService
    {

        public void Create(ApiRequestLogTemporaryInfoModel model);

    }
}
