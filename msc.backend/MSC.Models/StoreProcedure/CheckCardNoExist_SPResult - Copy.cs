using MSC.Models.Enums;
using MSC.Models.Models;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{

    public class CommonSPResultModel
    {

        public bool IsSuccess 
        { 
            get
            {
                return this.Code == ErrorCode.Success;
            }
        }

        
        public ErrorCode Code { get; set; }

    }
}
