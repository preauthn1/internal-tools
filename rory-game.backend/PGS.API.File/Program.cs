using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using target.Models;
using target.Utility;

namespace target.API.File
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
                .ConfigureAppConfiguration((hostingContext, config) =>
                {
                    config.AddEnvironmentVariables();
                })
                .ConfigureWebHostDefaults(webBuilder =>
                {
                    var port = SettingHelper.WebHostListenPort ?? "3108";
                    webBuilder.UseUrls($"http://*:{port}");
                    webBuilder.UseStartup<Startup>();
                    webBuilder.UseShutdownTimeout(TimeSpan.FromSeconds(SystemConstant.ApplicationShutdownDelaySeconds));
                });
    }
}
