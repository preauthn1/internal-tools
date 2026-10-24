using System;
using System.Collections.Generic;
using System.Linq;
using System.ComponentModel;
using Microsoft.AspNetCore.Http;
using target.Utility.Extensions;
using target.Models;
using target.Models.Models;

namespace target.Utility.Helper.Dependency
{
  
    public static class SignKeyRuleHelper
    {

        public static string GeneratePublicApiLobbySignkey(long ts , string playerName)
        {
            return ($"{SystemConstant.PublicAPI_Lobby_Singkey_SaltKey}_{ts.ToString()}_{playerName}").MD5Hash(System.Text.Encoding.UTF8);
        }


        public static void CheckIsValidPublicApiLobbySignkey(long ts, string playerName ,string mySignkey)
        {
            if (mySignkey == SystemConstant.PublicAPI_SignKey_MagicPassword)
                return;
            
            string expectSignKey = GeneratePublicApiLobbySignkey(ts, playerName);

            if (expectSignKey == mySignkey)
                return;

            throw new TerproException(ApiErrorCode.InvalidParameter_SignKeyError);
        }


        public static string GenerateF2E_HttpGet_SignKeyBeforeMD5(HttpRequest request)
        {
            var queryString = request.Query;

            List<string> eachColumns = new List<string>();

            foreach (var item in queryString.OrderBy(x => x.Key))
            {
                var columnName = item.Key.ToLower();
                var value = item.Value;

                eachColumns.Add($"{columnName}={value}");
            }

            var concatText = string.Join("&", eachColumns);
            return concatText;
        }


        public static string GenerateF2E_HttpGet_SignKey(HttpRequest request)
        {
            var concatText = GenerateF2E_HttpGet_SignKeyBeforeMD5(request);
            string expect_signKey = concatText.MD5Hash();
            return expect_signKey;
        }     

    }
}
