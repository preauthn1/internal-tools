using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using System.Runtime.Serialization.Formatters.Binary;
using example.invalid;
using StackExchange.Redis;
using Lck.Utility.Extensions;
using Newtonsoft.Json;
using MSC.Models.Models;
using MSC.Utility;
using example.invalid;

namespace Lck.Utility.Extensions.Cache
{
    /// <summary>
    /// 採用Redis進行Cache機制的基礎類別
    /// </summary>
    public class RedisCache : BaseCache
    {
        private static Lazy<ConnectionMultiplexer> Connection;

        private static EndPoint _endPoint { get; set; }


        private int _dbNumber;


        private IDatabase redisDB
        {
            get
            {
                return Connection.Value.GetDatabase(_dbNumber);
            }
        }

        private IServer redisServer
        {
            get
            {
                return Connection.Value.GetServer(_endPoint);
            }
        }

        public RedisCache()
        {
            InitialRedis();
        }

        private void InitialRedis()
        {
            //若已經初始化的話 就不再重新連線
            if (Connection == null)
            {
                string configSetting = SettingHelper.CacheSetting.RedisConn;

                if (string.IsNullOrWhiteSpace(configSetting))
                    throw new MSCException(ErrorCode.BusinessLogicDeny_SettingDataNotCompleted);

                int timeoutMilliseconds = 10000;

                ConfigurationOptions redisConfig = ConfigurationOptions.Parse(configSetting);
                redisConfig.ConnectTimeout = timeoutMilliseconds;
                redisConfig.SyncTimeout = timeoutMilliseconds;
                redisConfig.AllowAdmin = true;

                Connection = new Lazy<ConnectionMultiplexer>(() => ConnectionMultiplexer.Connect(redisConfig)); 
                _endPoint = redisConfig.EndPoints.FirstOrDefault();
            }

            _dbNumber = SettingHelper.CacheSetting.RedisDBNumber;
        }


        public override T Get<T>(string key)
        {
            try
            {
                if (typeof(T) == typeof(byte[]))
                    return (dynamic)redisDB.StringGet(key);

                string jsonText = redisDB.StringGet(key);

                if (typeof(T) == typeof(string))
                    return (dynamic)jsonText;

                return DeserializeObject<T>(jsonText);
            }
            catch (Exception ex)
            {
                return default(T);
            }
        }


        public override bool Set(string key, object value, TimeSpan expiredTime)
        {
            try
            {
                string jsonText = null;
                if (value != null)
                {
                    Type myValueType = value.GetType();

                    if (myValueType == typeof(byte[]))
                    {
                        return SetValue(key, (byte[])value, expiredTime);
                    }

                    if (myValueType == typeof(string))
                    {
                        jsonText = value.ToString();
                    }
                    else
                    {
                        jsonText = JsonConvert.SerializeObject(value);
                    }
                }

                return SetValue(key, jsonText, expiredTime);

            }
            catch (Exception ex)
            {
                return false;
            }
        }


        private bool SetValue(string key, RedisValue jsonValue, TimeSpan? expiredTime = null)
        {
            if (expiredTime == null)
            {
                return redisDB.StringSet(key, jsonValue);
            }

            return redisDB.StringSet(key, jsonValue, expiredTime.Value);
        }


        /// <summary>
        /// 關鍵字找出特定keys
        /// </summary>
        public override List<string> GetCacheKeysByKeyword(string keyword)
        {
            var matchKeys = redisServer.Keys(_dbNumber, pattern: $"*{keyword}*");
            var result = matchKeys.Select(z => z.ToString()).ToList();
            return result;
        }


        /// <summary>
        /// 清除所有Cache
        /// </summary>
        public override void FlushAll()
        {
            redisServer.FlushDatabase(_dbNumber);
        }

        /// <summary>
        /// 清除鍵值為key的特定Cache
        /// </summary>
        public override void Flush(string key)
        {
            redisDB.KeyDelete(key);
        }


        protected static T DeserializeObject<T>(string jsonText)
        {
            if (typeof(T) == typeof(string))
                return (dynamic)jsonText;

            return JsonConvert.DeserializeObject<T>(jsonText);
        }
    }
}