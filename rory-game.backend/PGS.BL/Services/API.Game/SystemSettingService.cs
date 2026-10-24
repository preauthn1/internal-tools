using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.DA.Repository;

namespace target.BL.Services
{

    public class SystemSettingService : BaseService, ISystemSettingService
    {
        public ForbiddenIPAddressRepository ForbiddenIPAddressRepository { get; set; }


        public List<DenyIpAddressSettingViewModel> GetDenyIpAddressList()
        {
            var query = ForbiddenIPAddressRepository.GetAll();

            var viewModels = query.Select(x => new DenyIpAddressSettingViewModel
            {
                IPAddress = x.IPAddress
            })
                .ToList();

            return viewModels;
        }


    }
}
