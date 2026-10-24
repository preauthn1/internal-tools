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
using Newtonsoft.Json.Linq;

namespace MSC.BL.Services
{
    public class WebHookNotifyService : BaseAdminService, IWebHookNotifyService
    {
      
        /// <summary>
        /// 發送slack或TG通知
        /// </summary>
        public void AzureNotifySQLWarning(JObject model)
        {
            string json = "【AzureNotifySQLWarning】" +  model.ConvertToJsonString().Truncate(1000);

            MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
            {
                Channel = SlackChannelEnum.notify_msc_api_job,
                Message = json
            });

            MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
            {
                TelegramChannelEnum = TelegramChannelEnum.ogc_common_error,
                Message = json
            });
        }

    }
}
