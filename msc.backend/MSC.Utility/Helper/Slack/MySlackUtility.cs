using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Configuration;
using example.invalid;
using Microsoft.Extensions.Configuration;

namespace MSC.Utility.Utility.ThirdParty
{
    public class MySlackUtility
    {
               
        public static void SendMessageToChannelAsync(SlackMessageParameter parameter)
        {
            //非同步
            Task.Factory.StartNew(() =>
                {
                    SendMessage(parameter);
                });
        }

        public static void SendMessageToChannelNotAsync(SlackMessageParameter parameter)
        {
            //同步
            SendMessage(parameter);
        }


        private static void SendMessage(SlackMessageParameter parameter)
        {        
            string prefixText = GetPrefixText();

            var newService = new BaseSlackUtility(prefixText);
            newService.SendMessage(parameter);
        }


        private static string GetPrefixText()
        {
            string siteName = GetConfigValueByKey("WebSiteName");

            //可根據專案客製需求，不同環境要帶不同的前綴文字 提供識別
            if (IsDev())
                return $"【Debug環境】{siteName} ";

            return $" `{siteName}` ";
        }


        private static string GetConfigValueByKey(string keyName)
        {
            return BaseSettingHelper.GetSettingValue(keyName);
        }


        private static bool IsDev()
        {
#if DEBUG
            return true;
#else
            return false;
#endif
        }



    }
}
