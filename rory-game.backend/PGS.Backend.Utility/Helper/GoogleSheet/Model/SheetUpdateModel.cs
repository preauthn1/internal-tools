using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Lck.GoogleApi
{
    public class SheetUpdateModel
    {


        public string SheetId { get; set; }

        
        
        /// <summary>
        /// 依序往右邊寫入值
        /// </summary>
        public List<string> UpdateCellValues { get; set; }
        
    }


    


}
