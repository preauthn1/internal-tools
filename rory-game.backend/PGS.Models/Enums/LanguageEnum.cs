using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;

namespace target.Models.Enums
{
    /// <summary>
    /// 支援的語系
    /// </summary>
    [Description("資料庫是存字串_編號本身沒意義")]
    public enum LanguageEnum
    {
        [Description(i18nKeys.Enums_LanguageEnum_zh_cn)]
        zh_cn,

        [Description(i18nKeys.Enums_LanguageEnum_zh_tw)]
        zh_tw,

        [Description(i18nKeys.Enums_LanguageEnum_en_us)]
        en_us,

        [Description(i18nKeys.Enums_LanguageEnum_en_us)]
        vi_vn,

        [Description(i18nKeys.Enums_LanguageEnum_en_us)]
        ms_my,
    }
}
