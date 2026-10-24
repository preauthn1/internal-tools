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
using target.Utility.Helper;
using target.Models;
using target.DA.Repository;

namespace target.BL.Services
{

    public class PermissionService : BaseAdminSysService, IPermissionService
    {
        public PermissionRepository PermissionRepository { get; set; }

        public RolePermissionMappingRepository RolePermissionMappingRepository { get; set; }

        public UserRoleMappingRepository UserRoleMappingRepository { get; set; }


        public UserPermissionViewModel GetAvailablePermissions(int userID)
        {
            var currentUser = GetCurrentOwnPermissionByCache(userID);

            List<int> queryRolePermissionIDs = RolePermissionMappingRepository.GetByFilter(new QueryRolePermissionMappingRequestViewModel
                {
                    RoleIDs = currentUser.Roles
                })
                .Select(x=> x.PermissionID)
                .ToList();

            var queryPermissions = PermissionRepository.GetByIds(queryRolePermissionIDs);
            
            List<UserPermissionItemViewModel> viewModels = ConvertToViewModels(queryPermissions);

            //複製一份一模一樣的
            var allPermissionViewModels = viewModels.Select(x => x.ReturnInjecteObject<UserPermissionItemViewModel>()).ToList();

            //處理巢狀迴圈
            NestedFormatSetValue(viewModels);

            return new UserPermissionViewModel
            {
                NestedPermissions = viewModels.Where(x => x.RootPermissionID == x.PermissionID).ToList(),
                AllPermissions = allPermissionViewModels.Select(x => x.ReturnInjecteObject<PermissionTokenBriefViewModel>()).ToList(),
            };
        }

    

        public UserPermissionViewModel GetPermissionsNestedFormat(int userID)
        {
            var query = PermissionRepository.GetAll();

            var convertBrief = query.Select(x => new PermissionTokenBriefViewModel
                {
                    PageUrl = x.PageUrl,
                    Token = [REDACTED],
                    PermissionID = x.PermissionID,
                    ParentPermissionID = x.ParentPermissionID
                })
                .ToList();

            List<UserPermissionItemViewModel> viewModels = ConvertToViewModels(query);

            //處理巢狀迴圈
            NestedFormatSetValue(viewModels);

            var rootItems = viewModels
                .Where(x => x.RootPermissionID == x.PermissionID)
                .OrderBy(z => z.Sort)
                .ToList();

            return new UserPermissionViewModel
            {
                NestedPermissions = rootItems,
                AllPermissions = convertBrief
            };
        }



        public List<PermissionViewModel> GetPermissions(int userID)
        {
            var query = PermissionRepository.GetAll();

            return query.Select(x => new PermissionViewModel
            {
                DisplayName = x.PermissionName,
                PageUrl = x.PageUrl,
                PermissionID = x.PermissionID,
                IsShowInMenu = x.IsShowInMenu
            })
            .ToList();
        }


        public bool IsPermissionExist(int userID , List<AdminSysPermissionKey> permissionKeys)
        {
            var currentUser = GetCurrentOwnPermissionByCache(userID);

            //主要使用cache來處理
            var myPermissions = UserRoleMappingRepository.GetRolePermissionsByCache(currentUser.Roles);

            //為了避免麻煩 轉小寫判斷
            var tokenText = permissionKeys.Select(x => x.ToString().ToLower()).ToList();

            return myPermissions.Any(z=> tokenText.Contains(z.PermissionToken));
        }


        private static List<UserPermissionItemViewModel> ConvertToViewModels(List<PermissionEntity> query)
        {
            return query.Select(x => new UserPermissionItemViewModel
            {
                Sort = x.Sort,
                DisplayName = x.PermissionName,
                PageUrl = x.PageUrl,
                PermissionID = x.PermissionID,
                PermissionType = x.PermissionType,
                Token = [REDACTED],
                ParentPermissionID = x.ParentPermissionID,
                IsShowInMenu = x.IsShowInMenu
            })
            .OrderBy(z=> z.Sort)
            .ToList();
        }

        private static void NestedFormatSetValue(List<UserPermissionItemViewModel> viewModels)
        {
            Dictionary<int, int> rootMapping = new Dictionary<int, int>();

            //先處理第一層的
            foreach (var item in viewModels)
            {
                if (item.ParentPermissionID == 0)
                {
                    item.RootPermissionID = item.PermissionID;
                    rootMapping.Add(item.PermissionID, item.RootPermissionID);
                }
            }

            int executeTimes = 1;

            while (true)
            {
                //防止無窮迴圈
                if (executeTimes >= 20)
                    break;

                bool isAny = false;
                foreach (var item in viewModels)
                {
                    if (rootMapping.ContainsKey(item.ParentPermissionID) && !rootMapping.ContainsKey(item.PermissionID))
                    {
                        isAny = true;
                        item.RootPermissionID = rootMapping[item.ParentPermissionID];
                        rootMapping.Add(item.PermissionID, item.RootPermissionID);
                    }
                }

                //如果都處理完了 跳出迴圈
                if (!isAny)
                    break;
            }


            var groupBy = viewModels.GroupBy(x => x.RootPermissionID).ToDictionary(x => x.Key, z => z.ToList());

            //不曉得效能是否頗差
            foreach (var item in groupBy)
            {
                var matchItems = viewModels.Where(x => x.RootPermissionID == item.Key).ToList();

                foreach (var sameRootItem in matchItems)
                {
                    sameRootItem.ChildPermissions = matchItems
                        .Where(x => x.ParentPermissionID == sameRootItem.PermissionID)
                        .OrderBy(x => x.Sort)
                        .ToList();
                }
            }


        }
    }
}
