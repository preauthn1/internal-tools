using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.GoogleApi
{
    public class SheetReadModel
    {

        /// <summary>
        /// 第幾行
        /// </summary>
        public int RowIndex { get; set; }


        public List<string> ColumnValues { get; set; }

        
    }
}
