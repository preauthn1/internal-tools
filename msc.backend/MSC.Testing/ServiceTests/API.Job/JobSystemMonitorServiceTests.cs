using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.Testing
{
    [TestClass]
    public class JobSystemMonitorServiceTests : BaseTest
    {

        public IJobSystemMonitorService ServiceInstance { get; set; }


        [TestMethod]
        public void JobSystemMonitorService_DetectPointApiRequestUnusual_Test()
        {
            ServiceInstance.DetectPointApiRequestUnusual();
            return;
        }

    }
}

