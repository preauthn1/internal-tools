using System;
using System.Collections.Generic;
using System.Security.Claims;
using System.Text;
using System.Security.Cryptography;
using target.Utility.Extensions;

namespace target.Utility
{
    
    public class BaseJWTTokenHelper
    {
        public static string key = "8sywhyuri4ensti7";

        public static string Generate_JWT_Token(object myObject)
        {
            if (myObject == null)
                return null;

            var header = new
            {
                alg = "AES",
                typ = "JWT"
            };

            string headerText = header.ConvertToJsonString().Base64Encode();
            string payloadText = myObject.ConvertToJsonString().Base64Encode();

            string encrySignature = GenerateSignature(headerText , payloadText);
            return string.Format("{0}.{1}.{2}", headerText, payloadText, encrySignature);
        }

        private static string GenerateSignature(string header , string payload)
        {
            string concat = string.Format("{0}.{1}", header, payload);
            string encrySignature = EncryptAES.Encrypt(concat, key);
            return encrySignature;
        }

        public static T GetObjectFromJWT<T>(string jwtToken)
        {
            var splitJWT = jwtToken.Split('.');

            //解析失敗1. 不合法格式
            if (splitJWT.Length <= 2)
                return default(T);

            string payload = splitJWT[1].Base64Decode();

            string checkSignature = GenerateSignature(splitJWT[0], splitJWT[1]);
            
            //解析失敗2. 資料解密不同 被串改
            if (checkSignature != splitJWT[2])
                return default(T);            

            if (typeof(T) == typeof(string))
                return (dynamic)payload;

            return payload.JsonConvertToModel<T>();
        }


    }
}
