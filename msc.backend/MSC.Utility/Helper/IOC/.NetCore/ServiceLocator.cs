//using Autofac;
//using Microsoft.AspNetCore.Http;
//using Microsoft.AspNetCore.Http.Features;
//using System;
//using System.Collections.Generic;
//using System.Globalization;
//using System.Linq;
//using System.Reflection;
//using System.Text;
//using System.Threading.Tasks;

//namespace target.IOC
//{

//    /// <summary>
//    /// 這個是 .netCore專用的
//    /// </summary>
//    public static class IOCServiceProvider
//    {
//        private static IServiceProviderProxy diProxy;

//        public static IServiceProviderProxy ServiceProvider => diProxy ?? throw new Exception("You should Initialize the ServiceProvider before using it.");

//        public static void Initialize(IServiceProviderProxy proxy)
//        {
//            diProxy = proxy;
//        }
//    }

//    public interface IServiceProviderProxy
//    {
//        T GetService<T>();
//        object GetService(Type type);
//    }

//    public class HttpContextServiceProviderProxy : IServiceProviderProxy
//    {
//        private readonly IHttpContextAccessor contextAccessor;

//        public HttpContextServiceProviderProxy(IHttpContextAccessor contextAccessor)
//        {
//            this.contextAccessor = contextAccessor;
//        }

//        public T GetService<T>()
//        {
//            return (T)contextAccessor.HttpContext.RequestServices.GetService(typeof(T));
//        }

//        public object GetService(Type type)
//        {
//            return contextAccessor.HttpContext.RequestServices.GetService(type);
//        }

//    }



//}
