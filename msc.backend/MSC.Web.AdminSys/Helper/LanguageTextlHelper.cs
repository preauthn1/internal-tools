using Microsoft.AspNetCore.Html;
using Microsoft.AspNetCore.Mvc.Rendering;
using MSC.Models;
using MSC.Models.ViewModels.Base;
using MSC.Utility.Helper;
using MSC.Utility.Helper.HttpContext;
using System;
using System.Collections.Generic;


namespace MSC.Utility.Helper
{
    public static class t
    {        

        public static string Text(string i18nKey)
        {
            var headers = MyHttpContext.Current.Request.Headers;

            if (!headers.ContainsKey(SystemConstant.HeaderKey_Lang))
            {
                MyHttpContext.Current.Request.Headers.Add(SystemConstant.HeaderKey_Lang, "zh-cn");
            }

            string debug = MultiLanguageTextHelper.GetCommon_i18nText(i18nKey);
            return MultiLanguageTextHelper.GetCommon_i18nText(i18nKey);
        }

    }
}