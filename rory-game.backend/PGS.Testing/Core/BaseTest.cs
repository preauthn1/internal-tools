using game.Entity.Game;
using target.Utility;
using target.Utility.IOC;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Web;
using target.Utility.Extensions;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.Utility.Helper.HttpContext;
using Microsoft.AspNetCore.Http;
using NSubstitute;
using Microsoft.Extensions.Primitives;
using target.Models;
using example.invalid;
using target.Utility.Helper;
using target.Models.Enums;
using target.Core.Repository;
using target.Models.Models;

namespace target.Testing
{

    public class BaseTest
    {

        public static TestDataConfig _TEST_DATA => new TestDataConfig();


        public BaseTest()
        {
            //如果沒呼叫這段 讀config會出問題
            RuntimeHelper.Initialize("Dev");
            InitializeIOC();
            InstanceAllRepositoryImplement();
            MockHttpRequest();
        }


        protected string GeneatePlayerJwtToken(int playerID , int merchantID , string playerName , string sessionToken = "")
        {
            if(string.IsNullOrWhiteSpace(sessionToken))
            {
                sessionToken = [REDACTED];
            }

            var result = MyJwtHelper.GeneratePlayerJwt(
                    playerID,
                    merchantID,
                    playerName, 
                    SystemConstant.DefaultLang,
                    sessionToken,
                    sessionToken
                );

            return result.JwtToken;
        }



        private void MockHttpRequest()
        {
            var fakeJwt = MyJwtHelper.GenerateAdminUserJwt(_TEST_DATA.Default.UserID, _TEST_DATA.Default.UserName, _TEST_DATA.Default.OrganizationID, true, _TEST_DATA.Default.MerchantID.ToListCollection());

            MockHttpHeaderJwtToken(fakeJwt);
        }


        protected void MockHttpHeaderJwtToken(string jwtToken)
        {
            var myHttpContext = Substitute.For<IHttpContextAccessor>();

            var mockHeaders = new HeaderDictionary(new Dictionary<String, StringValues>
                    {
                        { SystemConstant.HeaderKey_Lang ,LanguageEnum.zh_cn.ToString()},
                        { SystemConstant.HeaderKey_AccessToken ,jwtToken}

                    }) as IHeaderDictionary;

            //設定假的實體的Add方法當傳入3,6 回傳 9
            myHttpContext.HttpContext.Request.Headers.Returns(mockHeaders);
            myHttpContext.HttpContext.Connection.RemoteIpAddress.Returns(new MockIPAddress(0));

            MyHttpContext.Configure(myHttpContext);
        }



        /// <summary>
        /// 如果某些 test 只有本機開發才執行的話，可加入此判斷
        /// </summary>
        protected bool IsNotDeveloperComputerName()
        {
            string currentComputerName = System.Environment.MachineName ?? "unknow";

            if (SystemConstant.RD_ComputerNames.Contains(currentComputerName))
            {
                return false;
            }

            return true;
        }


        protected T InstanceRepository<T>()
        {
            return (T)Activator.CreateInstance(typeof(T));
        }

        protected void AssertExpectException(Exception ex, ApiErrorCode code)
        {
            AssertExpectException(ex , new List<ApiErrorCode> { code });
        }

        protected void AssertExpectException(Exception ex , List<ApiErrorCode> codes)
        {
            var terproEx = ex as TerproException;

            if(terproEx == null)
            {
                Assert.Fail("非 TerproException");
            }

            bool isDefine = codes.Any(z => z == terproEx.ErrroCode);

            if(!isDefine)
            {
                Assert.Fail("非預期的errorCodes");
            }
        }

        protected void Assert_NotEmptyCollection<T>(IEnumerable<T> model)
        {
            if(model.IsEmptyOrNull())
            {
                Assert.Fail("model 為空集合");
            }
        }


        private void InitializeIOC()
        {
            List<string> projectList = new List<string>
            {
                "target.BL"
            };

            IOCInitializeHelper.InitializeIOC(projectList);
            IOCHelper.SetMSTestMode();
        }


        protected void InstanceAllRepositoryImplement()
        {
            //如果service class內的全域變數
            //1. 定義為 Property (非Field)
            //2. 且命名包含 repository
            //3. 有繼承 IDapperRepository

            var allService = this.GetType()
                .GetProperties(BindingFlags.NonPublic | BindingFlags.Instance | BindingFlags.Public)
                .Where(x => x.Name.ToLower().Contains("service"))
                .ToList();


            foreach (PropertyInfo item in allService)
            {
                if (item.PropertyType.IsInterface)
                {
                    var instanceValue = IOCHelper.Resolve(item.PropertyType);
                    item.SetValue(this, instanceValue);
                }
            }

            var allRep = this.GetType()
             .GetProperties(BindingFlags.NonPublic | BindingFlags.Instance | BindingFlags.Public)
             .Where(x => x.Name.ToLower().Contains("repository") || x.Name.ToLower().Contains("rep"))
             .ToList();


            foreach (PropertyInfo item in allRep)
            {
                if (item.PropertyType.BaseType != null)
                {
                    if (item.PropertyType.BaseType.IsGenericType &&
                        (
                        item.PropertyType.BaseType.GetGenericTypeDefinition() == typeof(FormatIntegerPrimaryKeyRepository<>) ||
                        item.PropertyType.BaseType.GetGenericTypeDefinition() == typeof(FormatLongPrimaryKeyRepository<>)
                       ))
                    {
                        var instanceValue = Activator.CreateInstance(item.PropertyType);
                        item.SetValue(this, instanceValue);
                    }
                }
            }

        }

    }



    public class MockIPAddress : IPAddress
    {
        public MockIPAddress(long newAddress):base(newAddress)
        {

        }

        public override string ToString()
        {
            return "0.0.0.0";
        }
    }

}