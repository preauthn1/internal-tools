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
    public class JobRoutineTaskService : BaseAdminService, IJobRoutineTaskService
    {
        public SystemWarningNotifyRepository SystemWarningNotifyRepository { get; set; }


        /// <summary>
        /// 發送slack或TG通知
        /// </summary>
        public void CheckSystemWarningAndSendNotify()
        {
            var queryNotSendNotify = SystemWarningNotifyRepository.GetByFilter(new QuerySystemWarningNotifyRequestViewModel
            {
                IsNotSend = true
            });

            foreach(var item in queryNotSendNotify)
            {
                SystemWarningNotifyTypeEnum type = (SystemWarningNotifyTypeEnum)item.NotifyType;

                switch(type)
                {
                    case SystemWarningNotifyTypeEnum.Slack:
                        SendMessageToSlack(item);
                        break;

                    case SystemWarningNotifyTypeEnum.Telegram:
                        SendMessageToTelegram(item);
                        break;

                    default:
                        throw new MSCException(ErrorCode.NotImplement);
                }

                //更新為已發送
                item.IsSend = true;
                item.SendTime = DateTimeHelper.GetCurrentTime();
                SystemWarningNotifyRepository.Update(item);
            }
        }


        private void SendMessageToTelegram(SystemWarningNotifyEntity model)
        {
            var channel = ExtensionOfEnum.GetEnumByTextValue<TelegramChannelEnum>(model.NotifyChannel, TelegramChannelEnum.ogc_common_error);

            MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
            {
                Message = model.Message,
                TelegramChannelEnum = channel
            });
        }


        private void SendMessageToSlack(SystemWarningNotifyEntity model)
        {
            var channel = ExtensionOfEnum.GetEnumByTextValue<SlackChannelEnum>(model.NotifyChannel, SlackChannelEnum.notify_msc_api_points);

            MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
            {
                Message = model.Message,
                Channel = channel
            });
        }


    }
}
