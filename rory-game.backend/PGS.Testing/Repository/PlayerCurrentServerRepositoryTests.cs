using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class PlayerCurrentServerRepositoryTests : BaseTest
    {

        public PlayerCurrentServerRepository PlayerCurrentServerRepository { get; set; }


        [TestMethod]
        public void PlayerCurrentServerRepository_PagingByFilter_Test()
        {
            var query = PlayerCurrentServerRepository.GetPagingByFilter(new QueryPlayerCurrentServerRequestViewModel
            {
                PageSize = 10
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }

    }
}

