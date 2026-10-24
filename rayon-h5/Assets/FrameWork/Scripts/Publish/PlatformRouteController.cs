using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using example.invalid;
using System.Linq;
using XLua;

public class PlatformRouteController
{
    [CSharpCallLua]
    public delegate Dictionary<string, object> GetAll();

    private static string s_ServerIP = "http://0.0.0.0";

    public static string GetCreateIP()
    {
        string _path = s_ServerIP + "/Home/Create";
        Debug.Log(_path);
        return _path;
    }

    public static string GetEditorIP(string AppKey)
    {
        string _path = s_ServerIP + "/Home/Update?AppKey=" + AppKey;
        Debug.Log(_path);
        return _path;
    }

    public static List<BuildInfo> GetRouteList()
    {
        string filePath = Application.dataPath + "/Hall/LuaScript/Com/PlatformRoute.lua";
        FileStream fileStream_Open = new FileStream(Path.GetFullPath(filePath), FileMode.Open);
        StreamReader streamWriter = new StreamReader(fileStream_Open);
        string oldScriptContent = streamWriter.ReadToEnd();
        streamWriter.Close();
        fileStream_Open.Close();

        LuaEnv luaenv = new LuaEnv();

        luaenv.DoString(oldScriptContent);

        GetAll _GetAll = luaenv.Global.GetInPath<GetAll>("GetAll");

        Dictionary<string, object> _table = _GetAll(); //constructor

        List<BuildInfo> _list = new List<BuildInfo>();

        foreach (XLua.LuaTable lt in _table.Values)
        {
            BuildInfo _info = new BuildInfo();
            _info.Index = lt.Get<System.Int32>("Index");
            _info.CompanyName = lt.Get<string>("CompanyName");
            _info.ProductName = lt.Get<string>("ProductName");
            _info.BundleID = lt.Get<string>("BundleID");
            _info.PlatformID = lt.Get<System.Int32>("PlatformID");
            _info.SiteID = lt.Get<System.Int32>("SiteID");
            _info.AppName = lt.Get<string>("AppName");
            if (lt.ContainsKey("Branch"))
                _info.Branch = lt.Get<System.Int32>("Branch");
            _info.ServerOption = lt.Get<System.Int32>("ServerOption");
            _info.AppleID = lt.Get<string>("AppleID");
            _info.StaticSiteID = lt.Get<System.Int32>("StaticSiteID");
            if (lt.ContainsKey("FlagShare"))
                _info.FlagShare = (lt.Get<System.Int32>("FlagShare") == 1) ? true : false;

            if (lt.ContainsKey("NotExport"))
                _info.NotExport = lt.Get<System.Int32>("NotExport");

            _info.IconName = lt.Get<string>("IconName");
            

            _list.Add(_info);
        }
        _list.Sort((x, y) => {
            int compare = x.Index.CompareTo(y.Index);
            if (compare != 0)
                return compare;

            compare = x.ServerOption.CompareTo(y.ServerOption);
            if (compare != 0)
                return compare;

            compare = x.Branch.CompareTo(y.Branch);
            if (compare != 0)
                return compare;

            compare = x.StaticSiteID.CompareTo(y.StaticSiteID);
            if (compare != 0)
                return compare;

            return x.AppName.CompareTo(y.AppName);
        });

        return _list;
    }

    private static void Save(string defineCode)
    {
        // 寫進lua
        string definition = "";         // 定义部分
        string scriptContent = "";      // 脚本内容
        string fileString = "";         // 最后写入文件的内容
        string notes = "--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---"; // 注释部分
        string filePath = Application.dataPath + "/Hall/LuaScript/Com/PlatformRoute.lua";

        FileStream fileStream_Open = new FileStream(Path.GetFullPath(filePath), FileMode.Open);
        StreamReader streamWriter = new StreamReader(fileStream_Open);
        string oldScriptContent = streamWriter.ReadToEnd();
        streamWriter.Close();
        fileStream_Open.Close();

        scriptContent = oldScriptContent.Substring(oldScriptContent.IndexOf(notes) + notes.Length, oldScriptContent.Length - oldScriptContent.IndexOf(notes) - notes.Length);

        definition = defineCode.Replace("\\\"", "\"");
        definition = definition.Replace("\\n", "\n");
        definition = definition.Replace("\\r", "");
        definition = definition.Replace("\\t", "");

        if (definition.Substring(0, 1) == "\"")
            definition = definition.Substring(1, definition.Length - 1);

        if (definition.Substring(definition.Length - 1, 1) == "\"")
            definition = definition.Substring(0, definition.Length - 1);

        //Debug.Log(definition);
        fileString = definition + "\r\n" + notes + scriptContent;

        FileStream fileStream_Create = new FileStream(Path.GetFullPath(filePath), FileMode.Create);
        StreamWriter streamWriter2 = new StreamWriter(fileStream_Create, new System.Text.UTF8Encoding(false));
        streamWriter2.Write(fileString);
        streamWriter2.Flush();
        streamWriter2.Close();
        fileStream_Create.Close();
    }

    public static void Refresh(System.Action _callback)
    {
        string _url = s_ServerIP + "/api/pr/unity?pl=";

#if UNITY_IOS
            _url += 1;
#else
        _url += 2;
#endif
        Debug.Log(_url);
        WWW www = new WWW(_url);

        bool _check = false;
        while (www.isDone == false)
        {
            if (www.isDone == true)
            {
                _check = true;
                if (string.IsNullOrEmpty(www.error) == false)
                {
                    Debug.Log(www.error);
                }
                else
                {
                    Debug.Log(www.text);
                    Save(www.text);
                    if (_callback != null)
                        _callback();
                }
                break;
            }
        }
        if (www.isDone == true && _check == false)
        {
            if (string.IsNullOrEmpty(www.error) == false)
            {
                Debug.Log(www.error);
            }
            else
            {
                Debug.Log(www.text);
                Save(www.text);
                if (_callback != null)
                    _callback();
            }
        }
    }

    public static Dictionary<string, JiuLiaoInfo> GetJiuLiaoList()
    {
        string filePath = Application.dataPath + "/Hall/LuaScript/Com/PlatformRoute.lua";
        FileStream fileStream_Open = new FileStream(Path.GetFullPath(filePath), FileMode.Open);
        StreamReader streamWriter = new StreamReader(fileStream_Open);
        string oldScriptContent = streamWriter.ReadToEnd();
        streamWriter.Close();
        fileStream_Open.Close();

        LuaEnv luaenv = new LuaEnv();

        luaenv.DoString(oldScriptContent);
        
        GetAll _GetAll = luaenv.Global.GetInPath<GetAll>("GetAppInfoList");

        Dictionary<string, object> _table = _GetAll(); //constructor

        Dictionary<string, JiuLiaoInfo> _list = new Dictionary<string, JiuLiaoInfo>();

        foreach (XLua.LuaTable lt in _table.Values)
        {
            string _key = ((KeyValuePair<string, object>)_table.First(x => x.Value == lt)).Key;

            JiuLiaoInfo _info = new JiuLiaoInfo();
            _info.JiuLiaoId = lt.Get<string>("JiuLiaoId");
            _info.JiuLiaoSecrety = lt.Get<string>("JiuLiaoSecrety");
            _info.JiuLiaoName = lt.Get<string>("JiuLiaoName");
            _info.DownloadUrl = lt.Get<string>("DownloadUrl");
            _list.Add(_key, _info);
        }

        return _list;
    }
}
