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
using Lck.Utility.Extensions.Helper;

namespace target.BL.Services
{

    public class RoleService : BaseAdminSysService, IRoleService
    {

        public RoleRepository RoleRepository { get; set; }

        public RolePermissionMappingRepository RolePermissionMappingRepository { get; set; }

        public UserRoleMappingRepository UserRoleMappingRepository { get; set; }



        public PagedList<RoleViewModel> GetByPagingFilter(QueryRoleApiRequestParameterModel filter)
        {
            var query = RoleRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<RoleEntity, RoleViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
            {
                var convert = ConvertToViewModels(x);
                return convert;
            })
                .ToList();

            return viewModels;
        }


        public RoleViewModel GetByID(int roleID)
        {
            var queryRole = RoleRepository.GetById(roleID);

            if (queryRole == null)
                throw new TerproException(ApiErrorCode.NotFoundRole);

            var convert = ConvertToViewModels(queryRole);

            var queryPermissions = RolePermissionMappingRepository.GetByFilter(new QueryRolePermissionMappingRequestViewModel
            {
                RoleIDs = roleID.ToListCollection()
            });

            convert.Permissions = queryPermissions.Select(x => x.PermissionID).ToList();
            return convert;
        }


        public void Create(RoleCreateRequestModel model)
        {
            model.RuleFor(x => x.RoleName).Rule_NotEmptyText().Rule_MaxLengthOK(25);

            var newEntity = model.ReturnInjecteObject<RoleEntity>();

            var currentRole = GetCurrentUser();
            newEntity.CreateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateBy = currentRole.LoginName;
            newEntity.CreateBy = currentRole.LoginName;

            var newRoleID = RoleRepository.Create(newEntity);

            if(model.TemplateID.HasValue)
            {
                RolePermissionMappingRepository.CopyPermissionFromRole(newRoleID , model.TemplateID.Value);
            }
        }


        public void Update(RoleUpdateRequestModel model)
        {
            model.RuleFor(x => x.RoleID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.RoleName).Rule_NotEmptyText().Rule_MaxLengthOK(25);

            var queryExist = RoleRepository.GetById(model.RoleID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundRole);

            queryExist.RoleName = model.RoleName;
            queryExist.UpdateTime = DateTimeHelper.GetCurrentTime();
            queryExist.UpdateBy = GetCurrentUser().LoginName;

            RoleRepository.Update(queryExist);

        }


        public void Delete(int roleID)
        {
            var queryExist = RoleRepository.GetById(roleID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundRole);

            //刪除mapping
            RoleRepository.Delete(roleID);

            //刪除該角色的權限
            RolePermissionMappingRepository.DeleteAndInsertNew(roleID , new List<int>());

            //刪除被指派為該角色的設置
            UserRoleMappingRepository.DeleteByRoleID(roleID);

        }



        public void UpdateRolePermissions(RolePermissionUpdateModel model)
        {
            model.RuleFor(x => x.RoleID).Rule_NotZeroNotMinus();
            model.Permissions = model.Permissions.Distinct().ToList();

            var queryExist = RoleRepository.GetById(model.RoleID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundRole);

            RolePermissionMappingRepository.DeleteAndInsertNew(model.RoleID , model.Permissions);

            //清除cache (如果未來開放帳號角色一對多 清除cache會有些問題)
            CacheRemoveHelper.RemoveCache_WhenRolePermissionUpdate(model.RoleID);
        }


        private RoleViewModel ConvertToViewModels(RoleEntity queryRole)
        {
            var convert = queryRole.ReturnInjecteObject<RoleViewModel>();
            return convert;
        }

    }
}
