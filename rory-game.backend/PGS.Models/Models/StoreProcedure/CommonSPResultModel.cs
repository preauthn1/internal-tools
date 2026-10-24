using target.Models.Models;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class CommonSPResultModel
    {
        public bool IsSuccess
        {
            get
            {
                var code = (int)this.Code;
                return code == 1 || code == 0;
            }
        }


        public ApiErrorCode Code { get; set; }

    }
}
