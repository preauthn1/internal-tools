using MSC.Models.Enums;
using MSC.Utility;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace MSC.Web.AdminSys.Helper
{
    public class AdminSysHelper
    {

        public static string WebSiteDisplayName()
        {
            var type = (WebSiteBrandingTypeEnum)SettingHelper.AdminSysWebsiteSetting.WebSiteBrandingType;

            switch(type)
            {
                case WebSiteBrandingTypeEnum.OGC:
                    return "OGC管理系统";

                default:
                    return "MSCard管理系统";
            }
        }


        public static WebSiteBrandingTypeEnum GetWebSiteBrandingType()
        {
            var type = (WebSiteBrandingTypeEnum)SettingHelper.AdminSysWebsiteSetting.WebSiteBrandingType;
            return type;
        }


        public static (string Begin, string End) GetRecently7Days()
        {
            double timeZone = 8;
            DateTime end = DateTimeHelper.GetCurrentTime(timeZone).Date.AddDays(1);
            DateTime begin = end.AddDays(-7);
            return (begin.ToString("yyyy-MM-dd HH:mm"), end.ToString("yyyy-MM-dd HH:mm"));
        }


        public static (string Begin , string End) GetThisMonthDateOnly()
        {
            var result = GetThisMonth();
            return (result.Begin.ToString("yyyy-MM-dd"), result.End.ToString("yyyy-MM-dd"));
        }


        public static (string Begin, string End) GetThisMonthToTomorrow()
        {
            double timeZone = GetDefaultTimeZone();
            DateTime today = DateTimeHelper.GetCurrentTime(timeZone).Date;            
            var result = GetThisMonth();
            return (result.Begin.ToString("yyyy-MM-dd"), today.AddDays(1).ToString("yyyy-MM-dd"));
        }


        public static (DateTime Begin, DateTime End) GetThisMonth()
        {
            double timeZone = GetDefaultTimeZone();
            var this_month = DateTimeHelper.GetCurrentTime(timeZone);
            var this_month_begin = new DateTime(this_month.Year, this_month.Month, 1);
            var this_month_end = this_month_begin.AddMonths(1);
            return (this_month_begin, this_month_end);
        }


        public static double GetDefaultTimeZone()
        {
            return 8.0;
        }


    }
}
