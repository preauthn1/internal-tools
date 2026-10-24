using MSC.Models.CoreLibs;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Reflection;
using System.Text;

namespace MSC.Core.Helper
{

    public static class SQLHelper
    {

        public static string GetInsertCommand<T>(string tableName, List<string> ignoreFields = null) where T : class
        {
            var query = string.Format("INSERT INTO {0} ({1}) VALUES({2})",
                tableName,
                GetFieldList<T>(ignoreFields),
                GetFieldList<T>(ignoreFields , "@"));

            return query;
        }


        public static string ReplaceToTotalCountCommand(string sql)
        {
            if (sql.Contains("select") && sql.Contains("from"))
                throw new Exception("Please use upper case.");

            return System.Text.RegularExpressions.Regex.Replace(sql, "SELECT(.*)FROM", string.Format("SELECT Count(*) FROM"));
        }


        public static string ConcatPagingCommand(string sql, string sortColumn, string descOrder = "desc")
        {
            sql += @" ORDER BY {0} {1} 
                     OFFSET @skipRows Rows 
                     FETCH Next @count Rows only";

            string descOrNot = descOrder == "desc" ? "desc" : "";
            sql = string.Format(sql, sortColumn, descOrNot);
            return sql;
        }
        


        public static string GetUpdateCommand<T>(string tableName)
        {
            var pkFields = GetPrimaryKeyFields<T>();
            return GetUpdateCommand<T>(tableName, null, pkFields);
        }



        private static string GetUpdateCommand<T>(string tableName, List<string> ignoreFields, List<string> pkFields)
        {
            var updateFieldString = GetUpdateFieldString<T>(ignoreFields, pkFields);

            //無可奈何 先這樣寫 (如果有model有定義多組Key、會移除名稱為Id的那一個)
            if (pkFields != null && pkFields.Count() > 1)
                pkFields.RemoveAll(x => x == "Id");

            var primaryKeyString = GetPrimaryKeyString(pkFields);

            var command = $"UPDATE {tableName} SET {updateFieldString} WHERE {primaryKeyString}";
            return command;
        }


        public static string GetInsertOrUpdateCommandByIgnoreId<T>(string tableName) where T : class
        {
            var fields = GetFieldList<T>(null);
            var fieldPatterns = GetFieldList<T>(null , "@");

            var fieldsStr = string.Join(",", fields.Select(fieldName => $"@{fieldName} AS {fieldName}"));

            var joinStr = string.Join(" AND ", GetPrimaryKeyFields<T>().Select(fieldName => $"T.{fieldName} = S.{fieldName}"));

            var updateStr = string.Join(",", fieldPatterns.Select(fieldName => $"T.{fieldName} = S.{fieldName}"));

            string sql = $@"
                MERGE INTO {tableName} as T
                USING (SELECT {fieldsStr}) AS S ON {joinStr}
                WHEN MATCHED THEN
                    UPDATE 
        	        SET 
        		        {updateStr}
                WHEN NOT MATCHED THEN
                    INSERT({string.Join(",", fieldPatterns)}) 
                    VALUES({string.Join(",", fieldPatterns.Select(e => $"S.{e}"))});";

            return sql;
        }



        private static string GetFieldList<T>(List<string> ignoreFields , string prefix = "") where T : class
        {
            var sb = new StringBuilder();
            var filterProperties = GetPropertiesAndFilterIgnore<T>(ignoreFields);

            bool isAutoKeyTable = IsNotAutoKey<T>() == false;

            foreach (var item in filterProperties.OrderBy(x => x.Name))
            {
                if (IsICollection(item.PropertyType))
                {
                    continue;
                }

                if (isAutoKeyTable && IsSkipColumn(item))
                {
                    continue;
                }


                if(prefix == "@")
                {
                    sb.Append($" @{item.Name},");
                }
                else
                {
                    sb.Append($" [{item.Name}],");
                }
            }

            var query = sb.ToString();
            return query.Remove(query.Length - 1);
        }



        private static List<PropertyInfo> GetPropertiesAndFilterIgnore<T>(List<string> ignoreFields)
        {
            var properties = typeof(T).GetProperties();
            List<PropertyInfo> filterProperties = properties.ToList();
            if (ignoreFields != null && ignoreFields.Count > 0)
            {
                filterProperties = properties
                    .Where(a => !ignoreFields.Contains(a.Name))
                    .ToList();
            }

            return filterProperties;
        }


        private static bool IsSkipColumn(PropertyInfo info)
        {
            var attr = info.GetCustomAttribute<KeyAttribute>();

            if (attr == null)
                return false;

            return true;
        }


        private static bool IsNotAutoKey<T>() where T : class
        {
            var keys = GetPrimaryKeyFields<T>();

            int numbersOfKeys = keys?.Count() ?? 0;

            if (numbersOfKeys >= 2)
                return true;

            if (numbersOfKeys == 0)
                return true;

            if (IsNotAutoKeyAttributeExist<T>())
                return true;

            return false;
        }

        private static bool IsNotAutoKeyAttributeExist<T>() where T : class
        {          
            var attr = (IsNotAutoKeyTable)typeof(T).GetCustomAttribute(typeof(IsNotAutoKeyTable), false);

            if (attr == null)
                return false;

            return true;
        }



        private static bool IsICollection(Type type)
        {
            if (type == null || !type.IsGenericType)
            {
                return false;
            }

            return typeof(ICollection<>).IsAssignableFrom(type
                .GetGenericTypeDefinition());

        }
        

        private static string GetUpdateFieldString<T>(List<string> ignoreFields, List<string> pkFields)
        {
            var filterProperties = GetPropertiesAndFilterIgnore<T>(ignoreFields).OrderBy(a => a.Name);
            var array = (from item in filterProperties
                         where !IsICollection(item.PropertyType) &&
                               !pkFields.Contains(item.Name)
                         select item.Name + "=@" + item.Name).ToArray();
            return string.Join(",", array);
        }


        private static string GetPrimaryKeyString(IEnumerable<string> pkFields)
        {
            var array = pkFields.Select(a => a + "=@" + a).ToArray();
            return string.Join(" AND ", array);
        }


        public static List<string> GetPrimaryKeyFields<T>()
        {
            var result = new List<string>();

            foreach (var prop in typeof(T).GetProperties())
            {
                var attr = (KeyAttribute)prop.GetCustomAttribute
                    (typeof(KeyAttribute), false);

                if (attr != null)
                {
                    result.Add(prop.Name);
                }
            }

            return result;
        }



    }








}
