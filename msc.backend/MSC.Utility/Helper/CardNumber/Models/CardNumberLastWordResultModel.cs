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
using MSC.Utility.Helper.HttpContext;

namespace MSC.Utility
{

    public class CardNumberLastWordResultModel
    {
        public List<int> Numbers { get; set; }

        /// <summary>
        /// 有多少組數字
        /// </summary>
        public int TotalNumbers { get; set; }


        /// <summary>
        /// 原始字串 倒數第二數字
        /// </summary>
        public int OriginalLast2Word { get; set; }


        /// <summary>
        /// 原始字串 最後一個數字
        /// </summary>
        public int OriginalLastWord { get; set; }


        /// <summary>
        /// 倒數第二碼數字
        /// </summary>
        public int Last2Word { get; set; }

        /// <summary>
        /// 最後一碼數字
        /// </summary>
        public int LastWord { get; set; }

    }


}
