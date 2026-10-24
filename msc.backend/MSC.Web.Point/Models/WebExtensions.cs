using MSC.Utility.Helper;
using System;

namespace MSC.Web.Point
{

    public class t
    {
        public static string Text(string i18nKeys)
        {
            return MultiLanguageTextHelper.GetCommon_i18nText(i18nKeys);
        }
    }


    public class WebExtensions
    {
        public string RequestId { get; set; }

        public bool ShowRequestId => !string.IsNullOrEmpty(RequestId);
    }
}
