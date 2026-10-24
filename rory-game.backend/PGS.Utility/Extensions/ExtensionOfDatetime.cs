using System;
using System.Collections.Generic;
using System.Globalization;
using System.Threading;


namespace target.Utility.Extensions
{
    /// <summary>
    /// Datetime 的擴充方法
    /// </summary>
    public static class ExtensionOfDatetime
    {
        public static string ToISO8601(this DateTime time , double timeZone = 0)
        {
            return time.ToString("yyyy-MM-ddTHH:mm:ssZ");
        }

        public static string ToDateTimeMinuteStringOrEmpty(this DateTime? time)
        {
            if (!time.HasValue)
                return string.Empty;

            return time.Value.ToString("yyyy-MM-dd HH:mm");
        }


        public static string ToDateTimeStringOrEmpty(this DateTime? time)
        {
            if (!time.HasValue)
                return string.Empty;

            return time.Value.ToString("yyyy-MM-dd HH:mm:ss");
        }

        public static string ToDateStringOrEmpty(this DateTime? time)
        {
            if (!time.HasValue)
                return string.Empty;

            return time.Value.ToString("yyyy-MM-dd");
        }

        /// <summary>
        /// 取得這個月首日
        /// </summary>
        public static DateTime ToMonthFirstDay(this DateTime time)
        {
            return ToMonthOneDay(time, 1);
        }

        /// <summary>
        /// 取得這個月某一日
        /// </summary>
        public static DateTime ToMonthOneDay(this DateTime time,int day)
        {
            return new DateTime(time.Year, time.Month, day);
        }

        /// <summary>
        /// 取得這個月最後一日
        /// </summary>
        public static DateTime ToMonthLastDay(this DateTime time)
        {
            DateTime resultTime = new DateTime(time.Year, time.Month, DateTime.DaysInMonth(time.Year, time.Month), 23, 59, 59, 999);
            return resultTime;
        }

        public static DateTime To_23_59_59(this DateTime time)
        {
            return new DateTime(time.Year, time.Month, time.Day, 23, 59, 59 , 999);
        }

        public static DateTime To_00_00_00(this DateTime time)
        {
            return new DateTime(time.Year, time.Month, time.Day, 0, 0, 0);
        }

        public static DateTime To_00_00(this DateTime time)
        {
            return new DateTime(time.Year, time.Month, time.Day, time.Hour, 0, 0);
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
        /// <param name="unixTimeStamp"></param>
        /// <returns></returns>
        public static DateTime UnixTimeStampToDateTime(this double unixTimeStamp)
        {
            // Unix timestamp is seconds past epoch
            System.DateTime dtDateTime = new DateTime(1970, 1, 1, 0, 0, 0, 0, System.DateTimeKind.Utc);
            dtDateTime = dtDateTime.AddSeconds(unixTimeStamp);
            return dtDateTime;
        }
    }
    

}