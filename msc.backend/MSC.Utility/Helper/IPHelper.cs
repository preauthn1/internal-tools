using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Collections.Specialized;
using MSC.Utility.Extensions;
using example.invalid;
using example.invalid;
using Microsoft.AspNetCore.Http;
using MSC.Utility.Helper.HttpContext;
using MSC.Models;

namespace MSC.Utility
{
    public static class IPHelper
    {
      
        public static string GetIP()
        {
            try
            {
                string k8s_ip = MyHttpContext.Current.Request.Headers["X-Real-IP"];

                if (!string.IsNullOrWhiteSpace(k8s_ip))
                    return k8s_ip;


                return MyHttpContext.Current.Connection.RemoteIpAddress.ToString();
            }
            catch
            {

            }

            return string.Empty;
            
        }


        public static bool IsCompanyIPAddress()
        {
            string ip = GetIP();
            return ip == SystemConstant.CompanyIPAddress || ip == "::1" || ip == "0.0.0.0";
        }


    }

}
