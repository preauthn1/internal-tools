using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Lck.GoogleApi
{
    public class SheetReadModel
    {

        public SheetReadModel()
        {
            this.ColumnValues = new List<string>();
        }

        /// <summary>
        /// 第幾行
        /// </summary>
        public int RowIndex { get; set; }


        public List<string> ColumnValues { get; set; }

        
    }
}
