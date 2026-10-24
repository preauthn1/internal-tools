using MSC.DA.Repository;
using MSC.Utility;
using MSC.Models.Table;
using System;
using MSC.BL.Interface;
using MSC.Utility.Helper;

namespace MSC.BL.Services
{
    public class ApiRequestLogService : BaseService , IApiRequestLogService
    {

        public ApiRequestLogRepository ApiRequestLogRepository { get; set; }


        public void CreateTryCatch(string tableName , ApiRequestLogEntity model)
        {
            try
            {
                var now = DateTimeHelper.GetCurrentTime();
                model.CreateTime = now;
                model.CreateTimeCategory = 0;
                model.IPAddress = IPHelper.GetIP();
                model.WebSiteName = SettingHelper.WebSiteName;

                ApiRequestLogRepository.CreateToTargetTable(tableName, model);
            }
            catch(Exception ex)
            {
                string msg = ex.Message;
            }


        }

    }



}
