using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;

namespace target.BL.Interface
{
    public interface IVendorMaintenanceService
    {

        PagedList<VendorMaintenanceViewModel> GetByPagingFilter(int userID, QueryVendorMaintenanceApiRequestParameterModel filter);



        CRUDResultModel Create(int userID, VendorMaintenanceCreateViewModel model);


        CRUDResultModel Delete(int userID, int targetID);


    }
}
