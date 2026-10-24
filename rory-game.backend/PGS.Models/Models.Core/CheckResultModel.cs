using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.CoreLibs
{
    public class CheckResultModel
    {
        public CheckResultModel()
        {

        }
        
        public CheckResultModel(string msg)
        {
            Message = msg;
        }


        public bool IsSuccess { get; set; }

        public string Message { get; set; }


    }
}
