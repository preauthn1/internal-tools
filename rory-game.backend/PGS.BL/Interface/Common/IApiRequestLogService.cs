using target.Models.Table;
using System;

namespace target.BL.Interface
{
    public interface IApiRequestLogService
    {
        void CreateTryCatch(string tableName, ApiRequestLogEntity model);

        void WriteToLocalFile(ApiRequestLogEntity model);

        void ReadLogFileAndSyncToRemote();
    }



}
