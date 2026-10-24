using MSC.Models.Enums;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models.Models
{
    public class IPForbiddenCheckResultModel
    {
  
        public bool IsFobidden { get; set; }

        /// <summary>
        /// 是否為新建立的 禁止紀錄 (後續要寫log)
        /// </summary>
        public bool IsNewForbiddenEvent { get; set; }


        public string ForbiddenCacheKey { get; set; }

    }
}
