using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Configuration;
using example.invalid;
using Microsoft.Extensions.Configuration;

namespace target.Utility.Utility.ThirdParty
{
    public class MyTelegramUtility
    {

        /// <summary>
        /// 取得config設定的slack token
        /// </summary>
        /// <returns></returns>
        private static string GetConfigTokenSetting()
        {
            string configSetting = GetConfigValueByKey("TelegramBotToken");

            if (!string.IsNullOrWhiteSpace(configSetting))
                return configSetting;

            //查看TG cahtID 可以透過  https://example.invalid/bot{pgs_token}/getUpdates
            return "[TG_BOT_TOKEN_REDACTED]";
        }


        public static void SendMessageToChannelAsync(TelegramMessageParameter parameter)
        {
            //單元測試不發通知
            if (SettingHelper.IsMSTest)
                return;

            //非同步
            Task.Factory.StartNew(() =>
                {
                    SendMessage(parameter);
                });
        }

       
        private static void SendMessage(TelegramMessageParameter parameter)
        {
            string token = [REDACTED];            
            string prefixText = GetPrefixText();

            var newService = new BaseTelegramUtility(token, prefixText);
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
