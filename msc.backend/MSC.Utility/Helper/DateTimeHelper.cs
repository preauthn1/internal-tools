using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Collections.Specialized;
using MSC.Utility.Extensions;
using example.invalid;
using example.invalid;
using Microsoft.AspNetCore.Http;
using System.Reflection;
using System.Runtime.Versioning;

namespace MSC.Utility.Helper
{

    public static class DateTimeHelper
    {
        public const string ISO_8601_Format = "yyyy-MM-ddTHH:mm:ssZ";


        public static DateTime GetDefaultTime(double timezone = 0)
        {
            return new DateTime(1900, 1 ,1);
        }

        public static DateTime GetCurrentTime(double timezone = 0)
        {
            return DateTime.UtcNow.AddHours(timezone);
        }

        public static long GetCurrentTimeTimestamp(double timezone = 0)
        {
            DateTime now = GetCurrentTime();
            return now.CovertToUnixTimestamp();
        }

        public static TimeZoneAreaModel GetTimeZoneToday(double timezone)
        {
            var now = DateTime.UtcNow.AddHours(timezone);
            var begin = now.To_00_00_00();
            var end = now.To_23_59_59();

            return new TimeZoneAreaModel
            {
                TimeZone = timezone,
                TodayBegin = begin,
                TodayEnd = end,
                TodayBeginUTC = begin.AddHours(-1 * timezone),
                TodayEndUTC = end.AddHours(-1 * timezone),
            };
        }

        public static string ToISO_8601(DateTime utcTime)
        {
            var convert = utcTime.ToString(ISO_8601_Format);
            return convert;
        }

    }


    public class TimeZoneAreaModel
    {
        public double TimeZone { get; set; }

        public DateTime TodayBegin { get; set; }

        public DateTime TodayEnd { get; set; }

        public DateTime TodayBeginUTC { get; set; }

        public DateTime TodayEndUTC { get; set; }

    }
}
