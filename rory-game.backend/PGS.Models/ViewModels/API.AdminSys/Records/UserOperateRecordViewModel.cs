using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{


    public class UserOperateRecordViewModel
    {
        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public int OperateType { get; set; }

        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public string Memo { get; set; }

        public string UpdateBy { get; set; }

        public DateTime CreateTime { get; set; }

        public string CreateIPAddress { get; set; }

    }



}
