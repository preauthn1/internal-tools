using System;
using System.Collections.Generic;
using System.Configuration;
using example.invalid;
using System.Security.Cryptography;
using System.Text;
using System.Text.RegularExpressions;
using System.Linq;
using target.Utility.Extensions;
using target.Models;
using target.Models.Enums;

namespace target.Utility.Helper
{
    /// <summary>
    /// 通用型的正規表達式 或 格式定義
    /// </summary>
    public class RegexRuleHelper
    {
        
        public static bool IsSupportLanguage(string inputLang)
        {
            if (string.IsNullOrWhiteSpace(inputLang))
                return false;

            inputLang = inputLang.ToLower().Replace("-", "_");

            var allLang = ExtensionOfEnum.GetEnumAllValue<LanguageEnum>();

            foreach(var myEnum in allLang)
            {
                string langCode = myEnum.ToString().ToLower();

                if (langCode == inputLang)
                    return true;
            }

            return false;
        }


        /// <summary>
        /// 是否有超過小數點兩位
        /// </summary>
        public static bool IsDecimalPointOverThanPlacesTwo(decimal value)
        {
            decimal modValue = value * 1000 % 10;

            if (modValue != 0)            
                return true;         

            return false;
        }


        public static bool IsUrlFormat(string input)
        {
            if (string.IsNullOrWhiteSpace(input))
                return false;
            
            if(input.Contains("http://") ||
               input.Contains("https://"))
            {
                return true;
            }

            return false;
        }


        public static bool IsIPAddressFormat(string input)
        {
            if (string.IsNullOrWhiteSpace(input))
                return false;

            //有句點符號 且 拆解後只剩下數字
            if (input.Split('.').Length <= 2)
                return false;

            string tempReplace = input.Replace(".", "");
            
            if (tempReplace.IsNumeric() == false)
                return false;

            return true;
        }


        public static bool IsQQAccount(string input)
        {
            if (string.IsNullOrWhiteSpace(input))
                return false;

            if (Regex.IsMatch(input, "^[1-9][0-9]{5,11}$") == false)            
                return false;            

            return true;
        }


        public static bool IsWeChatAccount(string input)
        {
            if (string.IsNullOrWhiteSpace(input))
                return false;

            if (Regex.IsMatch(input, "^[a-zA-Z]{1}[-_a-zA-Z0-9]{5,19}$") == false)            
                return false;            

            return true;
        }



        public static bool IsBankCard(string input)
        {
            if (string.IsNullOrWhiteSpace(input))
                return false;
            
            Match match = Regex.Match(input, @"^\d{10,19}$");
            return match.Success;
        }



        /// <summary>
        /// 是否包含特殊符號
        /// </summary>
        public static bool IsAnySpecialCharacters(string input)
        {
            string replaceText = input.ReplaceAllSpecialChar();
            
            if (replaceText.Length != input.Length)
                return true;

            return false;
        }
        
        
        /// <summary>
        /// 是否為中國手機格式
        /// </summary>
        public static bool IsChinaPhoneNumberFormat(string input)
        {
            if (string.IsNullOrWhiteSpace(input))
                return false;

            if (input.Length <= 8) //字數不得太少
                return false;
            
            Match match = Regex.Match(input, @"^(13[0-9]|14[57]|15[0-3,5-9]|165|166|167|17[0135678]|18[0-9]|19[1,8-9])\d{8}$");
            return match.Success;
        }
        

        public static bool IsContainNumbers(string input)
        {
            Regex regex = new Regex(@"\d+");
            return regex.IsMatch(input);
        }


        public static bool IsAccountFormat(string input)
        {
            if (string.IsNullOrWhiteSpace(input))
                return false;

            if (input.Length < SystemConstant.PlayerAccountMinLength) //字數不得太少
                return false;

            if (input.Length > SystemConstant.PlayerAccountMaxLength)
                return false;

            //英文數字、並允許部分特殊符號
            var isMatch = Regex.IsMatch(input, @"^[a-zA-Z0-9!@\-_`.+,]*$");

            return isMatch;
        }
        

        /// <summary>
        /// 只有英文跟數字
        /// </summary>
        public static bool IsLetterAndNumberOnly(string input)
        {
            Match match = Regex.Match(input, @"^[a-zA-Z0-9]+$");
            return match.Success;
        }
        

        public static bool IsPasswordFormat(string input)
        {
            if (string.IsNullOrWhiteSpace(input))
                return false;

            if (input.Length < SystemConstant.PlayerPasswordMinLength) //字數不得太少
                return false;

            if(input.Length > SystemConstant.PlayerPasswordMaxLength)
                return false;

            return true;
        }


        /// <summary>
        /// 检测是否为Email格式
        /// </summary>
        public static bool IsValidEmail(string strEmail)
        {
            return Regex.IsMatch(strEmail, @"^[\w\.]+([-]\w+)*@[A-Za-z0-9-_]+[\.][A-Za-z0-9-_]+$");
        }


        /// <summary>
        /// 检测是否为有效域名
        /// </summary>
        public static bool IsValidDomain(string host)
        {
            if (host.IndexOf(".") == -1)
                return false;

            return new Regex(@"^\d+$").IsMatch(host.Replace(".", string.Empty)) ? false : true;
        }


        /// <summary>
        /// 验证是否为正整数
        /// </summary>
        public static bool IsInt(string str)
        {
            return Regex.IsMatch(str, @"^[0-9]*$");
        }


        /// <summary>
        /// 判断是否 英文字母、數字、特殊符號
        /// </summary>
        public static bool IsLetterOrNumber(string str)
        {
            return Regex.IsMatch(str, @"^[A-Za-z0-9]*$");
        }




    }
}
