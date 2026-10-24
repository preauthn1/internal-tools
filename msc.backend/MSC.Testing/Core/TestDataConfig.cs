
using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;

namespace MSC.Testing
{

    public class TestDataConfig
    {
        public TestSampleDataModel Default => new TestSampleDataModel();

        public TestSampleData_Not_SuperUser_Model Not_Super_User => new TestSampleData_Not_SuperUser_Model();

    }


    public class TestSampleDataModel
    {

        public int UserID = 125002;

        public int RoleID = 1;

        public string UserName = "unit_test_001";

        public string UserPassword = [REDACTED]";



        public int MerchantID = 175001;

        public string MerchantSecret = [REDACTED]";

        public string MemberCode = "unit_test_888";

        public string IPAddress = "0.0.0.0";

        public string CallbackUrl = "http://0.0.0.0";

    }

    public class TestSampleData_Not_SuperUser_Model
    {



    }





}