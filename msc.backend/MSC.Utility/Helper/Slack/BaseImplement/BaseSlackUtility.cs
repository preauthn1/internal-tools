using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Configuration;
using System.ComponentModel;

namespace MSC.Utility.Utility.ThirdParty
{
    /// <summary>
    /// 基底的 BaseSlackUtility 基本上不會再改到
    /// 有需要根據公司邏輯 特別客製的，請修正 MySlackUtility
    /// </summary>
    internal class BaseSlackUtility
    {

        private string _prefixText;
        


        public BaseSlackUtility(string prefixText)
        {
            _prefixText = prefixText;
        }
        
        
        public void SendMessage(SlackMessageParameter parameter)
        {           
            if (string.IsNullOrWhiteSpace(parameter.Message))
                return;
                        
            parameter.Message = parameter.Message.Replace("#", "＃");

            string targetUrl = GetUrlSettingFromEnum(parameter.Channel);
            var para = new RequestParameters
            {
                BaseUrl = targetUrl,
                ApiUrl = APIRequestHelper.No_Url_Keyword,
                Method = HttpMethodType.HttpPost,
                TimeoutMillisecond = 5000,
            };
            
            //訊息加上客製前墜詞 (通常是用來識別環境)
            parameter.Message = _prefixText + parameter.Message;
            
            para.AddParameter("text", parameter.Message);
            
            var response = APIRequestHelper.RequestApiResult<string>(para);
            return;

        }


        private string GetUrlSettingFromEnum(SlackChannelEnum targetChannel)
        {
            var result = GetAttributeFromEnum<DescriptionAttribute>(targetChannel);
            return result?.Description;
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
