using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IRoleService
    {
        PagedList<RoleViewModel> GetByPagingFilter(QueryRoleApiRequestParameterModel filter);

        RoleViewModel GetByID(int targetID);

        void Create(RoleCreateRequestModel model);

        void Update(RoleUpdateRequestModel model);

        void Delete(int targetID);

        void UpdateRolePermissions(RolePermissionUpdateModel model);

    }
}
