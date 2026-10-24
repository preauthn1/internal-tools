using Autofac;
using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Threading.Tasks;

namespace target.Utility.IOC
{
    public static class IOCHelper
    {

        public static IContainer Container { get; set; }


        public static Dictionary<string, object> MockInstanceMapping = new Dictionary<string, object>();

        /// <summary>
        /// 如果單元測試會強制把這改成true
        /// </summary>
        private static bool _isMSTestMode = false;


        public static void SetContainer(IContainer container)
        {
            Container = container;
        }

        public static void SetMSTestMode()
        {
            _isMSTestMode = true;
        }

        /// <summary>
        /// 若要使用此函式，Runtime啟動環境 必須在global或intial的時候        
        /// 呼叫 IOCInitializeHelper.InitializeIOC(projectList);
        /// </summary>
        public static T Resolve<T>()
        {
            return Container.Resolve<T>();
        }


        public static object Resolve(Type type)
        {
            if(_isMSTestMode)
            {
                return InstanceMockInterface(type);
            }

            return Container.Resolve(type);
        }


        public static T ResolveByKey<T>(object key)
        {
            if (_isMSTestMode)
            {
                return InstanceMockInterfaceByKey<T>(key);
            }

            return Container.ResolveKeyed<T>(key);
        }



        public static object InstanceMockInterface(Type type)
        {
            if (MockInstanceMapping.ContainsKey(type.Name))
                return MockInstanceMapping[type.Name];

            return Container.Resolve(type);
        }



        public static T InstanceMockInterfaceByKey<T>(object key)
        {
            string name = typeof(T).Name;

            if (MockInstanceMapping.ContainsKey(name))
                return (T)MockInstanceMapping[name];

            return Container.ResolveKeyed<T>(key);
        }


        public static void AddMockInstance<T>(object mockInstance)
        {
            string typeName = typeof(T).Name;

            //提示: 若使用mock，instance service 必須在呼叫 AddMockInstance 之後
            if (MockInstanceMapping.ContainsKey(typeName))
                return;

            MockInstanceMapping.Add(typeName, mockInstance);
        }


    }
}
