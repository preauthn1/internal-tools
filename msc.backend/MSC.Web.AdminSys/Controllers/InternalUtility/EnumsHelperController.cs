using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Reflection;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.Web.AdminSys.Models;
using MSC.Utility.Extensions;

namespace MSC.Web.AdminSys.Controllers
{
    public class EnumsHelperController : Controller
    {       

        public IActionResult ConvertToExcelFormat(string enumName = "")
        {
            if (enumName.IsNullOrWhiteSpace())
                return Content("請帶入參數 enumName");

            var allEnums = GetAllEnums();

            List<EnumsToExcelFormatViewModel> values = new List<EnumsToExcelFormatViewModel>();
            foreach (Type type in allEnums)
            {
                if (type.Name.ToLower() == enumName.ToLower())
                {
                    foreach (var enumValue in System.Enum.GetValues(type))
                    {
                        values.Add(new EnumsToExcelFormatViewModel
                        {
                            Category = "Enums",
                            ColumnKeyName = $"Enums_{type.Name}_{enumValue.ToString()}"
                        });
                    }
                    break;
                }
            }

            return PartialView(values);
        }


        private List<Type> GetAllEnums()
        {
            List<Type> all = new List<Type>();
            var listOfEnums1 = GetTypesFromSpecificNamespace(Assembly.Load("MSC.Models"), "MSC.Models.Enums");
            var listOfEnums2 = GetTypesFromSpecificNamespace(Assembly.Load("MSC.Models"), "MSC.Models.Models");

            all.AddRange(listOfEnums1);
            all.AddRange(listOfEnums2);
            return all;
        }


        private List<Type> GetTypesFromSpecificNamespace(Assembly assembly, string nameSpace)
        {
            return assembly
                .GetTypes()
                .Where(t => String.Equals(t.Namespace, nameSpace, StringComparison.Ordinal) && !t.IsClass)
                .ToList();
        }


    }
}
