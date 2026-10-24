using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using target.Utility.Extensions;
using target.Models;
using target.Models.Models;
using target.Utility.Helper;
using System;
using System.Linq;
using target.Utility;
using target.Utility.IOC;
using target.BL.Interface;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.Logging;
using example.invalid;
using target.Models.Table;
using target.Utility.Helper.Dependency;
using target.Utility.Utility.ThirdParty;
using Microsoft.AspNetCore.Diagnostics;
using System.Data.SqlClient;
using System.Collections.Generic;
using target.Utility.Extensions.Web;
using Microsoft.Extensions.Hosting;
using System.Threading;
using System.Threading.Tasks;

namespace target.Web.Modules.Filter
{

    /// <summary>
    /// 針對k8s graceful shutdown 的處理
    /// </summary>
    public class LifetimeEventsHostedService : IHostedService
    {
        private readonly IHostApplicationLifetime _hostApplicationLifetime;


        public LifetimeEventsHostedService(IHostApplicationLifetime hostApplicationLifetime)
        {
            _hostApplicationLifetime = hostApplicationLifetime;
        }


        public Task StartAsync(CancellationToken cancellationToken)
        {
            _hostApplicationLifetime.ApplicationStarted.Register(OnStarted);
            _hostApplicationLifetime.ApplicationStopping.Register(OnStopping);
            _hostApplicationLifetime.ApplicationStopped.Register(OnStopped);

            return Task.CompletedTask;
        }


        public Task StopAsync(CancellationToken cancellationToken)
        {
            return Task.CompletedTask;
        }

        private void OnStopped()
        {

        }

        private void OnStopping()
        {         
            Thread.Sleep(SystemConstant.ApplicationShutdownDelaySeconds * 1000);
        }

        private void OnStarted()
        {

        }
    }


}
