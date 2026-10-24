using Autofac;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Utility.IOC
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

            //建立容器
            var container = builder.Build();
            IOCHelper.SetContainer(container);

        }        
    }
}
