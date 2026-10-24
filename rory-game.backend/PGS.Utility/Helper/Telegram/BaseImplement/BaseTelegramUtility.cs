using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Configuration;
using System.ComponentModel;

namespace target.Utility.Utility.ThirdParty
{
    /// <summary>
    /// 基底的 BaseSlackUtility 基本上不會再改到
    /// 有需要根據公司邏輯 特別客製的，請修正 MySlackUtility
    /// </summary>
    internal class BaseTelegramUtility
    {

        private string _token;
        private string _prefixText;
        

        public BaseTelegramUtility(string token , string prefix = "")
        {
            _token = [REDACTED];
            _prefixText = prefix;
        }

        
        public void SendMessage(TelegramMessageParameter parameter)
        {           
            if (string.IsNullOrWhiteSpace(parameter.Message))
                return;
                        
            parameter.Message = _prefixText + parameter.Message.Replace("#", "＃");

            string chatID = GetAttributeFromEnum<DescriptionAttribute>(parameter.TelegramChannelEnum)?.Description;

            string msg = parameter.Message;

            var para = new RequestParameters
            {
                BaseUrl = "https://example.invalid/",
                ApiUrl = $"bot{_token}/sendMessage?chat_id={chatID}&text={msg}",
                Method = HttpMethodType.HttpGet,
                TimeoutMillisecond = 10000,
            };

            var response = APIRequestHelper.RequestApiResult<string>(para);
            return;
        }

        public static T GetAttributeFromEnum<T>(System.Enum enumVal) where T : System.Attribute
        {
            var type = enumVal.GetType();
            var memInfo = type.GetMember(enumVal.ToString());
            var attributes = memInfo[0].GetCustomAttributes(typeof(T), false);
            return (attributes.Length > 0) ? (T)attributes[0] : null;
        }



    }
}
