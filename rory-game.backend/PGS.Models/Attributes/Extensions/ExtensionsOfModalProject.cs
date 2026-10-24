using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Extensions
{
    /// <summary>
    /// model專案 才會用到的extension
    /// 正常狀況是放在 Helper/Extension專案
    /// </summary>
    internal static class ExtensionsOfModalProject
    {
        /// <summary>
        /// 判斷string是否為數字
        /// </summary>
        public static bool IsNumeric(this string s)
        {
            double retNum;
            bool isNum = Double.TryParse(Convert.ToString(s), System.Globalization.NumberStyles.Any, System.Globalization.NumberFormatInfo.InvariantInfo, out retNum);
            return isNum;
        }


        /// <summary>
        /// 注意裡面的時間是 +0時區的
        /// </summary>
        public static long CovertToUnixTimestamp(this DateTime utcTime)
        {
            return (long)utcTime.Subtract(new DateTime(1970, 1, 1, 0, 0, 0)).TotalSeconds;
        }


        /// <summary>
        /// 返回拿到的時間是 +0時區的
        /// </summary>
        public static DateTime UnixTimeStampToDateTime(this double unixTimeStamp)
        {
            // Unix timestamp is seconds past epoch
            System.DateTime dtDateTime = new DateTime(1970, 1, 1, 0, 0, 0, 0, System.DateTimeKind.Utc);
            dtDateTime = dtDateTime.AddSeconds(unixTimeStamp);
            return dtDateTime;
        }

    }
}
