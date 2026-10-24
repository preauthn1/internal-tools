using MSC.Models;
using System;
using System.Collections.Generic;
using System.Text;
using MSC.Utility.Extensions;

namespace MSC.Utility.Helper
{
    public class AppHelper
    {
        public static string GeneratePasswordHash(string password)
        {
            return $"{SystemConstant.Password_SaltKey}_{password}".MD5Hash();
        }

    }
}
