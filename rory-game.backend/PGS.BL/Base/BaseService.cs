using target.Core.Repository;
using target.Models.Paging;
using target.Utility.IOC;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using target.Utility.Extensions;
using target.Utility.Helper;

namespace target.BL.Services
{
    public class BaseService
    {
        public BaseService()
        {
            InstanceAllRepositoryImplement();
        }


        /// <summary>
        /// 轉換成另一種class的 PagedList 物件
        /// </summary>
        protected PagedList<T_New> ConvertToNewTypePagingObject<T_Original , T_New>(PagedList<T_Original> model)
        {
            var resulk = model.ReturnInjecteObject<PagedList<T_New>>();
            return resulk;
        }


        protected void InstanceAllRepositoryImplement()
        {
            //如果service class內的全域變數
            //1. 定義為 Property (非Field)
            //2. 且命名包含 repository
            //3. 有繼承 IDapperRepository


            var allRep = GetType()
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

            var allService = GetType()
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
