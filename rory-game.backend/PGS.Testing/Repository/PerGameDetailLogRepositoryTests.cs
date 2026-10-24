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
    public class PerGameDetailLogRepositoryTests : BaseTest
    {

        public PerGameDetailLogRepository PerGameDetailLogRepository { get; set; }


        [TestMethod]
        public void UpdateCurrentAmount_Test()
        {
            var query = PerGameDetailLogRepository.GetByFilter(new QueryDetailLogRequestViewModel
            {
                PlayerIDs = _TEST_DATA.Default.PlayerID.ToListCollection()
            });

            if(query.IsNotEmptyNull())
            {
                long perGameDetailID = query.FirstOrDefault().Id;
                PerGameDetailLogRepository.UpdateCurrentAmount(perGameDetailID , 54321);
            }

            return;
        }

    }
}

