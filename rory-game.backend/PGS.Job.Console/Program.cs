using target.DA.Repository;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using target.Utility.Utility.ThirdParty;
using Quartz;
using Quartz.Impl;
using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using target.Utility.Extensions;
using target.Utility.IOC;

namespace target.Job.Console
{

    partial class Program
    {

        static async Task Main(string[] args)
        {
            InitializeIOC();
            
            var rep = new CronJobSettingRepository();

            //查出DB的設置
            var queryJobSettings = rep.GetByFilter(new QueryCronJobSettingRequestViewModel
            {
                IsEnable = true
            });

            if (queryJobSettings.IsEmptyOrNull())
                return;

            try
            {
                // Grab the Scheduler instance from the Factory
                StdSchedulerFactory factory = new StdSchedulerFactory();
                IScheduler scheduler = await factory.GetScheduler();

                // and start it off
                await scheduler.Start();


                foreach(var item in queryJobSettings)
                {
                    string jobName = $"Job_{item.JobID}";
                    string triggerName = $"Trigger_{item.JobID}";

                    JobBuilder jobTask = InstanceJobTaskByType((JobTypeEnum)item.JobType);

                    IJobDetail job = jobTask
                       .WithIdentity(jobName, "group1")
                       .UsingJobData("JobID", item.JobID.ToString())
                       .UsingJobData("JobApiUrl", item.RequestUrl)
                       .Build();

                    ITrigger trigger = TriggerBuilder.Create()
                        .WithIdentity(triggerName, "group1")
                        .WithCronSchedule(item.CronExpression)
                        .StartAt(item.StartTimeUTC)
                        .Build();

                    await scheduler.ScheduleJob(job, trigger);
                }


                //而 corn 表示式包括秒、分、時、日期、月份、星期、年(可略)七個參數自由搭配

                //"10,20,25 * * * * ? *"：每分鐘的第10、20、25秒會執行
                //"10 0/5 * * * ?"：每5分鐘的第10秒會執行(ex. 10:00:10 am, 10:05:10 am...)
                //"0 20 10-13 ? * WED,FRI"：每星期三與星期五的 10:20, 11:20, 12:20, 13:20 執行
                //"0 0/30 8-9 5,20 * ?"：每月5號及20號的 8:00, 8:30, 9:00, 9:30 執行


                System.Console.ReadLine();
                System.Console.ReadLine();
                System.Console.ReadLine();
                System.Console.ReadLine();
                System.Console.ReadLine();
                System.Console.ReadLine();

            }
            catch(Exception ex)
            {
                string msg = "JobConsole錯誤啦 " + ex.Message;

                MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
                {
                    Channel = SlackChannelEnum.n_api_gameapi,
                    Message = msg,
                });
            }
          
        }



        private static JobBuilder InstanceJobTaskByType(JobTypeEnum type)
        {
            switch(type)
            {
                case JobTypeEnum.CallWebApi:
                    return JobBuilder.Create<CallWebApiJob>();

                default:
                    throw new TerproException(ApiErrorCode.NotImplement);
            }
        }


        private static void InitializeIOC()
        {
            List<string> projectList = new List<string>
            {
                "target.BL"
            };

            IOCInitializeHelper.InitializeIOC(projectList);
        }


    }
}
