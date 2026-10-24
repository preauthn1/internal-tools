using Lck.Utility.Extensions.Helper;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using target.Utility.Extensions;
using Lck.Utility.Extensions.Cache;
using target.Utility;

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


        /// <summary>
        /// expiredTime 帶入null，會有預設過期時間
        /// </summary>
        public static bool SetCacheValue(string key, object value, TimeSpan? expiredTime)
        {
            if (IsCacheActive() == false)
                return false;

            expiredTime = GetExpireTimeBySetting(expiredTime);

            var cacheService = InstanceCacheService();

            bool isOK = cacheService.Set(key, value, GetExpireTimeBySetting(expiredTime));
            return isOK;
        }



        public static List<string> GetCacheKeysByKeyword(string keyword)
        {
            if (IsCacheActive() == false)
                return new List<string>();

            var cacheService = InstanceCacheService();
            return cacheService.GetCacheKeysByKeyword(keyword); ;
        }



        public static void RemoveAll()
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


        public static async Task RemoveCacheAsync(string key)
        {
            if (IsCacheActive() == false)
                return;

            var cacheService = InstanceCacheService();
            cacheService.FlushAsync(key);
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
            //根據專案的config設置 決定要用哪一種
            bool isRemoteCache = SettingHelper.IsRemoteCache;

            if (isRemoteCache)
                return new RedisCache();

            return new PGSMemoryCache();
        }


        public static TimeSpan GetExpireTimeBySetting(TimeSpan? expiredTime)
        {
            if (expiredTime.HasValue == false)
                expiredTime = _defaultCacheTime;

            return expiredTime.Value;
        }


    }
}
