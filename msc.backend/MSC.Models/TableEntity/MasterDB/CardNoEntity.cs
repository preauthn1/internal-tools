using MSC.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{

    [IsNotAutoKeyTable]
    public class CardNoEntity
    {
        public string CardNo { get; set; }

    }
}
