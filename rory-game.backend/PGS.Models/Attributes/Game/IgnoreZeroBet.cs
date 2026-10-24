using System;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Text.RegularExpressions;
using target.Models.Extensions;

namespace target.Models.Attributes
{

    /// <summary>
    /// 遊戲gameList專用的attribute
    /// 定義此遊戲，是否投注額=0的時候，可以忽略寫入此資料 (寫到另一張表)
    /// </summary>
    public class IgnoreZeroBet : System.Attribute
    {
        public bool IsIgnore { get; set; }

        public IgnoreZeroBet(bool ignore)
        {
            this.IsIgnore = ignore;
        }
    }




}
