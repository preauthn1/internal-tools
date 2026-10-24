using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.GoogleApi
{
    public class SheetReqeustParameter
    {
        
        public string SheetIdOrSheetUrl { get; set; }


        /// <summary>
        /// 預設是抓第一個tab的sheet，如果有需要特別指定某個sheet名稱 
        /// </summary>
        public string TargetSheetName { get; set; }

    }
}
