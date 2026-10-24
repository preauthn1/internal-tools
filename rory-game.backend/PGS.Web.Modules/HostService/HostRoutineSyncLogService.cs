using target.DA.Repository;
using target.Utility;
using target.Models.Table;
using System;
using target.BL.Interface;
using target.Utility.Helper;
using System.Collections.Generic;
using System.Linq;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using System.Threading;
using System.Threading.Tasks;
using target.Utility.Extensions;
using target.Utility.IOC;

namespace target.BL.Services
{
    public class HostRoutineSyncLogService : IHostedService, IDisposable
    {

        private Timer timer;
        private int number;


        public void Dispose()
        {
            timer?.Dispose();
        }

        public Task StartAsync(CancellationToken cancellationToken)
        {
            timer = new Timer(o => 
            {
                Interlocked.Increment(ref number);

                var service = IOCHelper.Resolve<IApiRequestLogService>();
                service.ReadLogFileAndSyncToRemote();
            },
            null,
            TimeSpan.Zero,
            TimeSpan.FromSeconds(15));

            return Task.CompletedTask;
        }


        public Task StopAsync(CancellationToken cancellationToken)
        {
            return Task.CompletedTask;
        }

    }



}
