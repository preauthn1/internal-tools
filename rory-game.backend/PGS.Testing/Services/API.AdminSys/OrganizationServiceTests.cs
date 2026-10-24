using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;
using target.Models.CoreLibs;
using target.Models.ViewModels;
using target.DA.Repository;
using target.Models.Enums;

namespace target.Testing
{
    [TestClass]
    public class OrganizationServiceTests : BaseTest
    {

        public IOrganizationService ServiceInstance { get; set; }


        [TestMethod]
        public void Organization_GetPagingByFilter_Test()
        {
            var query = ServiceInstance.GetPagingByFilter(new QueryOrganizationApiRequestParameterModel { PageSize = 10 } );
            Assert_NotEmptyCollection(query.Items);

            var matchFirst = query.Items.FirstOrDefault();
            var querySingle = ServiceInstance.GetByID(matchFirst.OrganizationID);

            Assert.IsNotNull(querySingle);
            return;
        }

    }
}

