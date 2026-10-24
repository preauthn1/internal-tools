using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IViewPermissionCheckService
    {
        void CheckViewPlayerPermission(int currentUserID, int targerPlayerID);

        void CheckViewUserPermission(int currentUserID, int targetUserID);

        void CheckViewMerchantPermission(int currentUserID, int targetMerchantID);

    }
}
