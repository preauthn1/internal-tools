using System;
using System.Collections.Generic;
using System.Linq;
using System.ComponentModel;
using MSC.Models;
using MSC.Models.ViewModels;
using MSC.Models.Models;

namespace MSC.Utility.Extensions
{
  
    public static class SignKeyRuleHelper
    {

        public static string GeneratePointTokenSignKey(int merchantID, string memberCode)
        {
            return (merchantID.ToString() + memberCode + SystemConstant.API_Vendor_SaltKey).MD5Hash().Truncate(15);
        }


        public static SignKeyResultModel GenerateVendorApiIntegrationSignKey(Dictionary<string,object> dic, string merchantSecret) 
        {           
            string concatText = DictionaryConvertToBeforeMD5Text(dic);

            concatText += $"&merchantsecret:[REDACTED]";
            var result = new SignKeyResultModel
            {
                BeforeMD5 = concatText,
                AfterMD5 = concatText.MD5Hash(),
                AfterMD5_UTF8 = concatText.MD5Hash()
            };

            return result;
        }


        public static SignKeyResultModel GenerateVendorApiIntegrationSignKey<T>(T requestInfo, string merchantSecret) where T : BaseVendorAPIRequestModel
        {
            Dictionary<string, object> dic = RequestInfoModelConvertToDictionary(requestInfo);
            return GenerateVendorApiIntegrationSignKey(dic , merchantSecret);
        }



        public static SignKeyResultModel GenerateProviderApiIntegrationSignKey(Dictionary<string , object> dic, string providersecret)
        {           
            string concatText = DictionaryConvertToBeforeMD5Text(dic);

            concatText += $"&providersecret:[REDACTED]";
            return new SignKeyResultModel
            {
                BeforeMD5 = concatText,
                AfterMD5 = concatText.MD5Hash()
            };
        }

        public static SignKeyResultModel GenerateProviderApiIntegrationSignKey<T>(T requestInfo, string providersecret) where T : BaseProviderAPIRequestModel
        {
            Dictionary<string, object> dic = RequestInfoModelConvertToDictionary(requestInfo);
            return GenerateProviderApiIntegrationSignKey(dic, providersecret);
        }


        private static string DictionaryConvertToBeforeMD5Text(Dictionary<string, object> dic)
        {
            List<string> keyValues = new List<string>();

            foreach (var item in dic.OrderBy(x => x.Key))
            {
                var name = item.Key.ToLower();
                var value = item.Value;

                keyValues.Add($"{name}:{value}");
            }

            string concatText = string.Join("&", keyValues);
            return concatText;
        }



        private static Dictionary<string, object> RequestInfoModelConvertToDictionary<T>(T requestInfo) where T : class
        {
            Dictionary<string, object> dic = new Dictionary<string, object>();

            foreach (var item in requestInfo.GetType().GetProperties().OrderBy(x => x.Name))
            {
                var name = item.Name.ToLower();
                var value = item.GetValue(requestInfo);

                if (value != null && value.GetType() == typeof(decimal))
                    value = Math.Round((decimal)value, 2).ToString("#0.00");

                dic.Add(name, value);
            }

            return dic;
        }



    }
}
