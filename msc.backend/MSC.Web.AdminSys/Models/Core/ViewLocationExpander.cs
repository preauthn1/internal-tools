using Microsoft.AspNetCore.Mvc.Razor;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace MSC.Web.AdminSys.Models.Core
{
    public class ViewLocationExpander : IViewLocationExpander
    {

        /// <summary>
        /// 自訂的views路徑擴充
        /// </summary>
        public IEnumerable<string> ExpandViewLocations(ViewLocationExpanderContext context, IEnumerable<string> viewLocations)
        {
            //{2} is area, {1} is controller,{0} is the action
            string[] locations = new string[] 
            { 
                "/Views/_InternalUtility/{1}/{0}.cshtml" ,
            };

            return locations.Union(viewLocations);          //Add mvc default locations after ours
        }


        public void PopulateValues(ViewLocationExpanderContext context)
        {
            context.Values["customviewlocation"] = nameof(ViewLocationExpander);
        }
    }

}
