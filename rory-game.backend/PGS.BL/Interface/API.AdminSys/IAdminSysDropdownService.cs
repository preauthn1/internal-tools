using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.ViewModels.API.AdminSys;

namespace target.BL.Interface
{

    public interface IAdminSysDropdownService
    {
        List<DropdownValueViewModel> GetGames(int userID);

        List<DropdownValueViewModel> GetGameVendorTypes();

        List<DropdownValueViewModel> GetGameTypes();

        List<DropdownValueViewModel> GetInternalTestPlayerAccounts();

        List<DropdownValueViewModel> GetServers(int userID, QueryServerListApiRequestParameterModel filter);

        List<DropdownValueViewModel> GetBetRecordStatus();

        List<DropdownValueViewModel> GetAvailableRoles(int userID);

        List<DropdownValueViewModel> GetRoomLevels();

        List<DropdownValueViewModel> GetAvailableMerchants(int userID , QueryDropdownMerchantRequestViewModel filter);

        List<DropdownValueViewModel> GetInternalMerchants(int userID);

        List<DropdownValueViewModel> GetDynamicQueryItems();

        List<DropdownValueViewModel> GetMerchantTypes();

        List<DropdownValueViewModel> GetTransactionTypes();

        List<DropdownValueViewModel> GetAvailableOrgs(int userID);

        List<DropdownValueViewModel> GetThirdPartyVendors();

        List<DropdownValueViewModel> GetTimeZoneAreas();

    }
}
