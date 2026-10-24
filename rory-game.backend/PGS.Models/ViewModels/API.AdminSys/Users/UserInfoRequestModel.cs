using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 回應參數
    /// </summary>
    public class UserInfoViewModel
    {
        public int UserID { get; set; }

        public string UserAccount { get; set; }

        public int OrganizationID { get; set; }

        public string OrganizationName { get; set; }

        public string DefaultLang { get; set; }

        public int TimeZoneAreaID { get; set; }

        public string TimeZoneAreaName { get; set; }

        public double TimeZone { get; set; }
        
        /// <summary>
        /// 若為超級管理員的組織 某些畫面功能更強大
        /// </summary>
        public bool IsSuperAdminOrganization { get; set; }

    }


}
