using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Utility.Helper
{
    /// <summary>
    /// 數字越小 排序越前面
    /// 沒特別指定的話 預設值就是10開始
    /// </summary>
    public class ExcelColumnSequenceAttribute : System.Attribute
    {
        public int Sequence { get; set; }


        public ExcelColumnSequenceAttribute()
        {

        }

        public ExcelColumnSequenceAttribute(int sequence)
        {
            Sequence = sequence;
        }

    }
}
