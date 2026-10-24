using Microsoft.Extensions.Caching.Memory;
using System;
using System.Collections.Generic;
using System.Linq;

namespace Lck.Utility.Extensions.Cache
{
    internal class MemoryCache : BaseCache
    {
        public static Microsoft.Extensions.Caching.Memory.MemoryCache _memoryCache { get; set; }

        public static HashSet<string> _allCacheKeys { get; set; }



        static MemoryCache()
        {
            _memoryCache = new Microsoft.Extensions.Caching.Memory.MemoryCache(new MemoryCacheOptions()); 
            _allCacheKeys = new HashSet<string>();
        }

        private Microsoft.Extensions.Caching.Memory.MemoryCache InstanceCache()
        {
            return _memoryCache;
        }

        /// <summary>
        /// 取得Cache中的物件
        /// </summary>
        public override T Get<T>(string key)
        {
            try
            {
                var cacheInstance = InstanceCache();
                return cacheInstance.Get<T>(key);
            }
            catch
            {
                return default(T);
            }
        }


        public override bool Set(string key, object value, TimeSpan expiredTime)
        {
            var cacheInstance = InstanceCache();

            //手動的方式 紀錄所有cache key (清除時要使用)
            _allCacheKeys.Add(key);

            cacheInstance.Set(key , value , expiredTime);
            return true;
        }



        /// <summary>
        /// 關鍵字找出特定keys
        /// </summary>
        public override List<string> GetCacheKeysByKeyword(string keyword)
        {
            var keys = _allCacheKeys.Where(x => x.Contains(keyword)).ToList();
            return keys;
        }


        /// <summary>
        /// 清除所有Cache
        /// </summary>
        public override void FlushAll()
        {
            var cacheInstance = InstanceCache();

            foreach (string key in _allCacheKeys)
            {
                Flush(key);
            }

            //清空還原
            _allCacheKeys = new HashSet<string>();
            return;
        }


        /// <summary>
        /// 清除鍵值為key的特定Cache
        /// </summary>
        public override void Flush(string key)
        {
            try
            {
                var cacheInstance = InstanceCache();
                cacheInstance.Remove(key);
            }
            catch
            {

            }
        }

    }
}