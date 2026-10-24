using MSC.Models.Paging;
using MSC.Utility.IOC;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using MSC.Core.Repository;
using MSC.DA.Repository;
using Lck.Utility.Extensions.Helper;
using MSC.Models.Models;

namespace MSC.BL.Services
{
    public class BaseService
    {

        public BaseService()
        {
            InstanceAllRepositoryImplement();
        }
        

        protected void InstanceAllRepositoryImplement()
        {
            //如果service class內的全域變數
            //1. 定義為 Property (非Field)
            //2. 且命名包含 repository
            //3. 有繼承 IDapperRepository


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
                    else
                    {
                        var instanceValue = Activator.CreateInstance(item.PropertyType);
                        item.SetValue(this, instanceValue);
                    }
                }
            }

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
        }


        protected T InstanceRepository<T>()
        {
            return (T)Activator.CreateInstance(typeof(T));
        }
              
    }
}
