using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models
{
    public class AdminUserCacheModelAdminUpdateViewModel
    {
        public int UserID { get; set; }
        public string UserName { get; set; }
        public int OrganizationID { get; set; }
        public bool IsSuperAdminOrg { get; set; }

        public List<int> Roles { get; set; }

        public List<int> AvailableOrganizationIDs { get; set; }

        /// <summary>
        /// 返回userMerchantMapping分配的商戶
        /// </summary>
        public List<int> AvailableMerchantIDs { get; set; }


        /// <summary>
        /// 若是 [超級管理員] 返回全部 [啟用] 的商戶
        /// 若是 [一般人員]   返回可操作的merchantIDs
        /// </summary>
        public List<int> AllEnableMerchantsOrAvailableMerchants { get; set; }


        /// <summary>
        /// 若是超級管理員 回空陣列
        /// 若是一般人員 返回可操作的merchantIDs
        /// </summary>
        public List<int> SuperAdminOrAvailableMerchantIDs
        {
            get
            {
                return this.IsSuperAdminOrg ? new List<int>() : AvailableMerchantIDs;
            }
        }




    }
}
