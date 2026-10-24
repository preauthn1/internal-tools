using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.DA.Repository;
using target.Models;
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
    public class LanguageTextServiceTests : BaseTest
    {

        public ILanguageTextService ServiceInstance { get; set; }


        [TestMethod]
        public void LanguageTextService_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID , new QueryLanguageTextApiRequestParameterModel());
            Assert_NotEmptyCollection(query.Items);
            return;
        }



        [TestMethod]
        public void GetAllTextByLanguageCode_Test()
        {
            var queryList = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID, new QueryLanguageTextApiRequestParameterModel());

            var queryByCode = ServiceInstance.GetAllTextByLanguageCode(queryList.Items.FirstOrDefault().LanguageCode);

            Assert_NotEmptyCollection(queryByCode.Items);
            return;
        }



        [TestMethod]
        public void BatchUpdate_Test()
        {
            var queryList = ServiceInstance.BatchUpdate(_TEST_DATA.Default.UserID, new LanguageTextBatchUpdateRequestModel
            {
                LanguageCode = _TEST_DATA.Default.Language_Text_Test,
                Items = new List<LanguageTextBatchUpdateEachLanguageRequestModel>
                {
                    new LanguageTextBatchUpdateEachLanguageRequestModel
                    {
                        Language = SystemConstant.DefaultLang,
                        Text = "黃金蟹"
                    }
                }
            });

            return;
        }


    }
}

