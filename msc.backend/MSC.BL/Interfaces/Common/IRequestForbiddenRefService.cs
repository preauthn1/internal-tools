using MSC.DA.Repository;
using MSC.Models.Table;
using MSC.Utility;
using MSC.Utility.Helper;
using MSC.Utility.Utility;
using MSC.Utility.Extensions;
using System;
using System.Collections.Generic;
using MSC.Models.Models;
using MSC.Models.Enums;

namespace MSC.BL.Services
{

    public interface IRequestForbiddenRefService
    {

        public void AddForbiddenCounterTimes(RequestForbiddenTypeEnum type, string forbiddenKeyOrIP);


        public IPForbiddenCheckResultModel CheckForbiddenCounterTooManyTimes(RequestForbiddenTypeEnum type, string forbiddenKeyOrIP, int maxTimes);



    }
}
