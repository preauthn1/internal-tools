using Lck.Utility.Extensions.Helper;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using MSC.Utility.Extensions;
using Lck.Utility.Extensions.Cache;
using MSC.Utility;
using MSC.Models.Models;

namespace Lck.Utility.Extensions.Helper
{
   
    /// <summary>
    /// 純粹針對 MemoryCache
    /// </summary>
    public class CacheHelper
    {
        protected static TimeSpan _defaultCacheTime
        {
            get
            {
                return new TimeSpan(0, 2, 0);
            }
        }


        public static T GetCacheOrQuery<T>(string key, Func<T> func, int expiredMinute)
        {
            return GetCacheOrQuery(key, func, new TimeSpan(0, expiredMinute, 0));
        }


        public static T GetCacheOrQuery<T>(string key, Func<T> func, TimeSpan? expiredTime = null)
        {
            //若無啟用cache
            if (IsCacheActive() == false)
            {
                return func();
            }

            T cacheValue = GetCacheValue<T>(key);

            if (cacheValue != null)
                return cacheValue;

            T queryResult = func();

            //null的話 就不用cache存了
            if (queryResult == null)
                return queryResult;

            var cloneResult = queryResult.ConvertToJsonString().JsonConvertToModel<T>();
            SetCacheValue(key, cloneResult, expiredTime);

            return queryResult;
        }


        public static T GetCacheValue<T>(string key)
        {
            if (IsCacheActive() == false)
                return default(T);

            var cacheService = InstanceCacheService();

            T cacheValue = cacheService.Get<T>(key);
            var cloneResult = cacheValue.ConvertToJsonString().JsonConvertToModel<T>();
            return cloneResult;
        }


        public static List<string> GetCacheKeysByKeyword(string keyword)
        {
            if (IsCacheActive() == false)
                return new List<string>();

            var cacheService = InstanceCacheService();
            return cacheService.GetCacheKeysByKeyword(keyword); ;
        }



        public static bool SetCacheValue(string key, object value, TimeSpan? expiredTime)
        {
            if (IsCacheActive() == false)
                return false;

            expiredTime = SetDefaultExpireTimeByEnv(expiredTime);

            var cacheService = InstanceCacheService();

            bool isOK = cacheService.Set(key, value, expiredTime.Value);
            return isOK;
        }

        public static bool SetCacheValue(string key, object value, int expiredMinutes)
        {
            if (IsCacheActive() == false)
                return false;

            TimeSpan expiredTime = new TimeSpan(0, expiredMinutes, 0);

            var cacheService = InstanceCacheService();

            bool isOK = cacheService.Set(key, value, expiredTime);
            return isOK;
        }


        public static void FlushAll()
        {
            if (IsCacheActive() == false)
                return;

            var cacheService = InstanceCacheService();
            cacheService.FlushAll();
        }



        public static void RemoveCacheByKey(string key)
        {
            if (IsCacheActive() == false)
                return;

            var cacheService = InstanceCacheService();
            cacheService.Flush(key);
        }



        public static void RemoveKeywordCache(string keyword)
        {
            if (IsCacheActive() == false)
                return;

            var cacheService = InstanceCacheService();
            List<string> cacheKeys = cacheService.GetCacheKeysByKeyword(keyword);

            foreach(string key in cacheKeys)
            {
                cacheService.Flush(key);
            }
        }



        private static bool IsCacheActive()
        {
            return true;
        }


        private static BaseCache InstanceCacheService()
        {
            if (!IsCacheActive())
                return null;

            //如果之後要更換成 redisCache 應該改一下這邊就好
            return InstanceCacheByType();
        }
        


        public static BaseCache InstanceCacheByType()
        {
            var cacheType = SettingHelper.CacheSetting.CacheType;

            switch(cacheType)
            {
                case CacheTypeEnum.MemoryCache:
                    return new MemoryCache();

                case CacheTypeEnum.RedisCache:
                    return new RedisCache();

                default:
                    throw new MSCException(ErrorCode.NotImplement);
            }
        }


        public static TimeSpan SetDefaultExpireTimeByEnv(TimeSpan? expiredTime)
        {
            if (expiredTime.HasValue == false)
                expiredTime = _defaultCacheTime;

            return expiredTime.Value;
        }


    }
}
