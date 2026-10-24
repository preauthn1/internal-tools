using target.Utility;
using target.Utility.Utility;
using target.Utility.Utility.ThirdParty;
using Quartz;
using System;
using System.Threading.Tasks;

namespace target.Job.Console
{
    public class CallWebApiJob : IJob
    {
        public async Task Execute(IJobExecutionContext context)
        {
            JobDataMap dataMap = context.JobDetail.JobDataMap;

            var apiUrl = dataMap.GetString("JobApiUrl");

            MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
            {
                Message = "CallWebApiJob_" + DateTime.Now.ToString("HH:mm:ss"),
                TelegramChannelEnum = TelegramChannelEnum.pgs_common_error
            });

            string gitSHA1 = RuntimeHelper.GetGitSHA1();
            string pipeLineID = RuntimeHelper.Get_Gitlab_Pipeline_ID();

            var para = new RequestParameters
            {
                BaseUrl = apiUrl,
                ApiUrl = APIRequestHelper.No_Url_Keyword,
                Method = HttpMethodType.HttpGet,
            };

            string requestSource = $"CallWebApiJob_{gitSHA1}_{pipeLineID}";

            para.AddHeaderData("HDCode", requestSource);

            var result = APIRequestHelper.RequestApiResult<string>(para);
            return;
        }
    }

}
