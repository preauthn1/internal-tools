using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.DA.Repository;
using MSC.Models.Attributes;
using target.Utility.Helper;
using MSC.Utility.Utility.ThirdParty;

namespace MSC.BL.Services
{
    public class JobSystemMonitorService : BaseAdminService, IJobSystemMonitorService
    {
        public ApiRequestLogRepository ApiRequestLogRepository { get; set; }

        public RequestForbiddenLogRepository RequestForbiddenLogRepository { get; set; }

        public SystemSettingConfigRepository SystemSettingConfigRepository { get; set; }

        /// <summary>
        /// 檢查API是否有異常
        /// </summary>
        public void DetectPointApiRequestUnusual()
        {
            DateTime now = DateTimeHelper.GetCurrentTime();

            DateTime end = new DateTime(now.Year, now.Month, now.Day, now.Hour , 0 , 0);
            DateTime begin = end.AddHours(-1);

            var queryCount = ApiRequestLogRepository.GetEachApiRequestCountDuration("[dbo].[ApiRequestLogPointAPIError]", begin, end);

            var systemSetting = SystemSettingConfigRepository.GetSingleByCache(SystemSettingConfigEnum.SystemWarning_PointApiSystemWarningConditionPerHour);

            int maxCount = systemSetting.SettingValueOfInt;

            foreach (var item in queryCount)
            {
                //若次數超過指定的限制
                if(item.TotalCount > maxCount)
                {
                    string msg = $"{item.RelativeUrl} , {Environment.NewLine}請求次數為 {item.TotalCount}。已判定為異常";
                    MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
                    {
                        Message = msg,
                        TelegramChannelEnum = TelegramChannelEnum.ogc_common_error
                    });
                }
            }
        }



        /// <summary>
        /// 檢查是否有人被 鎖ip/鎖token 紀錄
        /// </summary>
        public void DetectForbiddenRequestRecord()
        {
            DateTime now = DateTimeHelper.GetCurrentTime();

            DateTime end = new DateTime(now.Year, now.Month, now.Day, now.Hour, 0, 0);
            DateTime begin = end.AddHours(-1);

            var queryLogs = RequestForbiddenLogRepository.GetByFilter(new QueryRequestForbiddenLogRequestViewModel
            {
                CreateTimeBegin = begin,
                CreateTimeEnd = end
            });

            if (queryLogs.IsEmptyOrNull())
                return;

            int count = queryLogs.Count();
            string msg = $" {begin.AddHours(8).ToString("yyyy-MM-dd HH:mm")}時段 {Environment.NewLine}有 {count} 筆異常紀錄 (RequestForbiddenLog)";

            MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
            {
                Message = msg,
                TelegramChannelEnum = TelegramChannelEnum.ogc_common_error
            });
        }


    }
}
