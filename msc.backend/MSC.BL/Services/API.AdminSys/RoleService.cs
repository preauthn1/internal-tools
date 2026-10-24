using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.BL.Services
{
    public class RoleService : BaseService, IRoleService
    {

        public RoleRepository RoleRepository { get; set; }

        public RolePermissionRepository RolePermissionRepository { get; set; }


        public RolePermissionViewModel GetRolePermissions(int roleID)
        {
            var query = RoleRepository.GetById(roleID);

            if (query == null)
                throw new MSCException(ErrorCode.NotFoundData);

            var queryPermissions = RolePermissionRepository.GetByFilter(new QueryRolePermissionRequestViewModel
            {
                RoleIDs = roleID.ToListCollection()
            });

            List<string> existPermissions = queryPermissions.Select(x => x.PermissionKey).OrderBy(z => z).ToList();

            var allPermissions = ExtensionOfEnum.GetEnumAllValue<AdminPermissionTokenEnum>().Select(z => z.ToString()).ToList();

            return new RolePermissionViewModel
            {
                RoleID = query.RoleID,
                RoleName = query.DisplayName,
                TotalPermissionCount = queryPermissions.Count(),
                Permissions = existPermissions,
                UnassignPermissions = allPermissions.Where(x=> !existPermissions.Contains(x)).OrderBy(z=> z).ToList()
            };               
                
        }


    }
}
