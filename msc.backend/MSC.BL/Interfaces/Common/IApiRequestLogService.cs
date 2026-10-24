using MSC.Models.Table;
using System;

namespace MSC.BL.Interface
{
    public interface IApiRequestLogService
    {

        void CreateTryCatch(string tableName, ApiRequestLogEntity model);

    }



}
