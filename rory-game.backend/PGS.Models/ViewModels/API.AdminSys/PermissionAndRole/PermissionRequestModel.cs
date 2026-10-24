using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 請求參數
    /// </summary>
    public class PermissionRequestModel
    {

    }


  
    public class PermissionViewModel
    {

        public int PermissionID { get; set; }

        public string DisplayName { get; set; }

        public string PageUrl { get; set; }


        public bool IsShowInMenu { get; set; }

    }



}
