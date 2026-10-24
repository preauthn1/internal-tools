using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Web;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Primitives;
using MSC.Utility;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.Models.Enums;
using MSC.Utility.Helper.HttpContext;
using MSC.Utility.IOC;
using MSC.Utility.Extensions;
using MSC.Core.Repository;
using example.invalid;
using NSubstitute;
using MSC.Models.Models;

namespace MSC.Testing
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
            MockFunctions();
        }


        private void MockFunctions()
        {
            var myHttpContext = Substitute.For<IHttpContextAccessor>();

            var allPermissions = new HashSet<string>(ExtensionOfEnum.GetEnumAllValue<AdminPermissionTokenEnum>().Select(z=> z.ToString()).ToList());

            var fakeJwt = MyJwtHelper.GenerateAdminUserJwt(_TEST_DATA.Default.UserID, _TEST_DATA.Default.UserName , 1 , true , new List<int>() , allPermissions);

            var mockHeaders = new HeaderDictionary(new Dictionary<String, StringValues>
                    {
                        { SystemConstant.HeaderKey_Lang ,LanguageEnum.zh_cn.ToString()},
                        { SystemConstant.HeaderKey_AccessToken ,fakeJwt}

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



        protected void Assert_NotEmptyCollection<T>(IEnumerable<T> model)
        {
            if(model.IsEmptyOrNull())
            {
                Assert.Fail("model 為空集合");
            }
        }



        protected void Assert_ExpectExceptionThrow(Action executeFunction , ErrorCode code)
        {
            try
            {
                executeFunction();
            }
            catch (MSCException ex)
            {
                Assert.IsTrue(ex.ErrroCode == code);
            }
            catch (Exception ex)
            {
                Assert.Fail("Undefined");
            }
        }


        private void InitializeIOC()
        {
            List<string> projectList = new List<string>
            {
                "MSC.BL"
            };

            IOCInitializeHelper.InitializeIOC(projectList);
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
                       item.PropertyType.BaseType.GetGenericTypeDefinition() == typeof(IDapperRepository<>))
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