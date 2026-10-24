using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using target.BL.Services;
using target.Models;
using target.Utility;
using target.Web.Modules.Filter;

namespace target.API.Game
{
    public class Program
    {
        public static void Main(string[] args)
        {
            CreateHostBuilder(args).Build().Run();
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>
            Host.CreateDefaultBuilder(args)
                .ConfigureAppConfiguration((hostingContext, configuration) =>
                {
                    //如果沒呼叫這段 讀config會出問題
                    RuntimeHelper.Initialize(hostingContext.HostingEnvironment.EnvironmentName);

                })
                 .ConfigureServices((hostBuilderContext, services) =>
                 {
                     services.AddHostedService<LifetimeEventsHostedService>();

                     //如果config配置有開啟背景同步作業
                     if (SettingHelper.ApiRequestLogBackgroundSyncEnable)
                     {
                         services.AddHostedService<HostRoutineSyncLogService>();
                     }

                     services.Configure<HostOptions>(option =>
                     {
                         option.ShutdownTimeout = TimeSpan.FromSeconds(SystemConstant.ApplicationShutdownDelaySeconds);
                     });
                 })
                .ConfigureWebHostDefaults(webBuilder =>
                {
                    var port = SettingHelper.WebHostListenPort ?? SystemConstant.Port_API_Game.ToString();
                    webBuilder.UseUrls($"http://*:{port}");
                    webBuilder.UseStartup<Startup>();
                    webBuilder.UseShutdownTimeout(TimeSpan.FromSeconds(SystemConstant.ApplicationShutdownDelaySeconds));
                });
    }
}
