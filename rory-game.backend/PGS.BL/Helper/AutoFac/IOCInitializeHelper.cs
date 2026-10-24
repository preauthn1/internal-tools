using Autofac;
using target.BL.Interface;
using target.BL.Services;
using target.Models.Enums;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Threading.Tasks;

namespace target.Utility.IOC
{
    public class IOCInitializeHelper
    {       
        public static void InitializeIOC(List<string> csProjectNames)
        {
            // 容器建立者
            var builder = new ContainerBuilder();
            
            
            foreach(string name in csProjectNames)
            {
                //傳入的專案名稱 e.g. game.Service
                var service = Assembly.Load(name);
                builder.RegisterAssemblyTypes(service).AsImplementedInterfaces();
            }

            builder.RegisterType<PGSIntegrationService>().Keyed<IGameVendorIntegrationService>((int)GameVendorTypeEnum.target);
            builder.RegisterType<PGSIntegrationService>().Keyed<IGameVendorIntegrationService>((int)GameVendorTypeEnum.Slot88);
            builder.RegisterType<FunkyIntegrationService>().Keyed<IGameVendorIntegrationService>((int)GameVendorTypeEnum.Funky);
            builder.RegisterType<AGGIntegrationService>().Keyed<IGameVendorIntegrationService>((int)GameVendorTypeEnum.AGG);
            builder.RegisterType<AGGLiveIntegrationService>().Keyed<IGameVendorIntegrationService>((int)GameVendorTypeEnum.AGGLive);
            builder.RegisterType<FishingIntegrationService>().Keyed<IGameVendorIntegrationService>((int)GameVendorTypeEnum.CrazyFishing);

            //建立容器
            var container = builder.Build();
            IOCHelper.SetContainer(container);

        }        
    }
}
