using Autofac;
using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Utility.IOC
{
    public static class IOCHelper
    {

        public static IContainer Container { get; set; }
        

        public static void SetContainer(IContainer container)
        {
            Container = container;
        }


        /// <summary>
        /// 若要使用此函式，Runtime啟動環境 必須在global或intial的時候        
        /// 呼叫 IOCInitializeHelper.InitializeIOC(projectList);
        /// </summary>
        /// <typeparam name="T"></typeparam>
        /// <returns></returns>
        public static T Resolve<T>()
        {
            return Container.Resolve<T>();
        }


        public static T Resolve<T>(T service)
        {
            return Container.Resolve<T>();
        }

        public static object Resolve(Type type)
        {
            return Container.Resolve(type);
        }

    }
}
