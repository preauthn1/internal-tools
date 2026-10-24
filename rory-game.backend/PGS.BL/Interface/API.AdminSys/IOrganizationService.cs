using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.ViewModels.API.AdminSys;

namespace target.BL.Interface
{

    public interface IOrganizationService
    {
        PagedList<OrganizationViewModel> GetPagingByFilter(QueryOrganizationApiRequestParameterModel filter);

        OrganizationViewModel GetByID(int targetID);

        void Create(OrganizationCreateRequestModel model);

        void Update(OrganizationUpdateRequestModel model);

        void Delete(int targetID);

    }
}
