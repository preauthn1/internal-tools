using Lck.Utility.Extensions.Helper;
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
using System.Reflection;

namespace MSC.BL.Services
{
    public class PermissionService : BaseAdminService, IPermissionService
    {
        public RolePermissionRepository RolePermissionRepository { get; set; }


        /// <summary>
        /// 將所有權限都塞給指定的roleID 基本上是塞給admin
        /// </summary>
        /// <param name="roleID"></param>
        public void AssignAllPermissions(int roleID)
        {
            var allPermissions = GetConstants<AdminPermissionTokenEnum>();

            //超級管理員 不需要有這頁
            allPermissions = allPermissions.Where(z => z != AdminPermissionTokenEnum.Merchant_MerchantSummary.ToString()).ToList();

            var queryExistData = RolePermissionRepository.GetByFilter(new QueryRolePermissionRequestViewModel
            {
                RoleIDs = roleID.ToListCollection()
            });

            List<string> existPermissions = queryExistData.Select(z => z.PermissionKey).ToList();

            var mustInsert = allPermissions
                .Where(z => !existPermissions.Contains(z))
                .Select(z => new RolePermissionEntity
                {
                    PermissionKey = z,
                    RoleID = roleID
                })
                .ToList();

            RolePermissionRepository.CreateMany(mustInsert);
        }


        private List<string> GetConstants<T>()
        {
            FieldInfo[] fieldInfos = typeof(T).GetFields(BindingFlags.Public |
                 BindingFlags.Static | BindingFlags.FlattenHierarchy);

            return fieldInfos.Where(fi => fi.IsLiteral && !fi.IsInitOnly).Select(z => z.GetValue(z).ToString()).ToList();
        }


    }
}
