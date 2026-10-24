using System;
using System.Collections.Generic;
using UnityEngine;

public class PrefsUtility
{
    static Dictionary<string, object> PrefsDic = new Dictionary<string, object>();

    public static void Init()
    {
        string TestPrefs = FileUtility.ReadFileTextByFullName("PrefsData/TestPrefs.text");
        ParseData(TestPrefs);
    }

    public static void ParseData(string TestPrefs)
    {
        if (!string.IsNullOrEmpty(TestPrefs))
            PrefsDic = MiniJSON.Json.Deserialize(TestPrefs) as Dictionary<string, object>;
    }

    public static void DeleteKey(string key)
    {
        if (PrefsDic.ContainsKey(key))
            PrefsDic.Remove(key);
    }

    public static float GetFloat(string key, float defaultVal = 0)
    {
#if UNITY_EDITOR || UNITY_STANDALONE_WIN
        return Convert.ToSingle(GetValue(key, defaultVal));
#else
        return PlayerPrefs.GetFloat(key, defaultVal);        
#endif
    }

    public static int GetInt(string key, int defaultVal = 0)
    {
#if UNITY_EDITOR || UNITY_STANDALONE_WIN
        return Convert.ToInt32(GetValue(key, defaultVal));
#else
        return PlayerPrefs.GetInt(key, defaultVal);        
#endif
    }

    public static string GetString(string key, string defaultVal = "")
    {
#if UNITY_EDITOR || UNITY_STANDALONE_WIN
        return Convert.ToString(GetValue(key, defaultVal));
#else
        return PlayerPrefs.GetString(key, defaultVal); 
#endif
    }

    public static void SetFloat(string key, float value)
    {
#if UNITY_EDITOR || UNITY_STANDALONE_WIN
        SaveValue(key, value);
#else
        PlayerPrefs.SetFloat(key, value);
#endif
    }

    public static void SetInt(string key, int value)
    {
#if UNITY_EDITOR || UNITY_STANDALONE_WIN

        SaveValue(key, value);
#else
        PlayerPrefs.SetInt(key, value);
#endif
    }

    public static void SetString(string key, string value)
    {
#if UNITY_EDITOR || UNITY_STANDALONE_WIN
        SaveValue(key, value);
#else
        PlayerPrefs.SetString(key, value);
#endif 
    }

    static object GetValue(string key, object defaultVal)
    {
        object valObj;
        PrefsDic.TryGetValue(key, out valObj);

        if (valObj == null)
        {
            SaveValue(key, defaultVal);
            return defaultVal;
        }
        else
            return valObj;
    }

    static void SaveValue(string key, object value)
    {
        if (PrefsDic.ContainsKey(key))
            PrefsDic[key] = value;
        else
            PrefsDic.Add(key, value);

        string config = MiniJSON.Json.Serialize(PrefsDic);
        FileUtility.WriteTextToFile("PrefsData/TestPrefs.text", config);
    }
}
