using System;
using System.Collections.Generic;
using example.invalid;
using UnityEngine;
using XLua;

public class LuaManager : Singleton<LuaManager>
{
    /// <summary>
    /// 上次GC 的时间
    /// </summary>
    internal static float lastGCTime = 0;
    /// <summary>
    /// GC 的时间间隔(s)
    /// </summary>
    internal const float GCInterval = 1;

    // 将使用到的Lua脚本缓存起来，避免每次取值造成的gc
    Dictionary<string, Dictionary<string, string>> CacheLuaScriptDic;

    #region LuaEnv

    /// <summary>
    /// 管理器时候已经就绪
    /// </summary>
    public bool IsReady
    {
        get;
        protected set;
    }

    /// <summary>
    /// 重启管理器
    /// </summary>
    public void ReLaunch()
    {
        ShutDownLuaEnv();
        Launch();
    }

    /// <summary>
    /// 启动管理器
    /// </summary>
    public void Launch()
    {
        if (IsReady) // 如果真正启动或者已经启动则返回 
            return;

        LuaEnv.AddBuildin("pb", XLua.LuaDLL.Lua.LoadLuaProfobuf);

        IsReady = false;

        CacheLuaScriptDic = new Dictionary<string, Dictionary<string, string>>();
        LuaEnv.AddLoader((ref string filename) =>
        {
            string projectNumber = AppDefine.HallNumber;

            //第一位包含/的表示为游戏脚本加载，否则为大厅
            if (filename.IndexOf("/") == 0)
            {
                filename = filename.Substring(1);
                string[] parms = filename.Split('/');
                projectNumber = filename.Substring(0, filename.IndexOf("/"));       //  第一个为项目名
                filename = filename.Substring(filename.IndexOf("/"));               //  后面部分为Lua文件名
            }
            return LoadCustomLuaFile(filename, projectNumber);
        });

        LuaEnv.DoString("require 'StartUp'", "LuaManager");

        IsReady = true;
    }

    protected override void OnDestroy()
    {
        base.OnDestroy();
        LuaEnv = null;
        ShutDownLuaEnv();
    }

    /// <summary>
    /// 关闭管理器
    /// </summary>
    protected void ShutDownLuaEnv()
    {
        StopAllCoroutines();

        CallMethod("StartUp", "Shutdown");

        if (LuaEnv != null)
        {
            LuaEnv.Global.Dispose();
            LuaEnv.Tick();
            LuaEnv.GC();
            LuaEnv = new LuaEnv();
        }

        IsReady = false;
    }

    #endregion

    /// <summary>
    /// 从缓存字典里加载
    /// </summary>
    /// <param name="fileName"></param>
    /// <returns></returns>
    byte[] LoadFormAssetBundle(string fileName, string projectName)
    {
        string loadFileName = string.Empty;
        if (fileName.EndsWith(".lua"))
        {
            loadFileName = fileName;
        }
        else
        {
            loadFileName = fileName + ".lua";
        }
        int index = loadFileName.LastIndexOf('/');
        if (index != -1)
        {
            loadFileName = loadFileName.Substring(index + 1);
        }

        TextAsset asset = AssetBundleManager.Instance.LoadAsset<TextAsset>(projectName, loadFileName, false);

        if (asset != null)
            return asset.bytes;
        else
            return null;
    }

    /// <summary>
    /// 加载自定义的Lua文件
    /// </summary>
    /// <param name="luaFileName"></param>
    /// <returns></returns>
    public byte[] LoadCustomLuaFile(string filename, string projectName)
    {
        // 从缓存中加载
        byte[] result = LoadFormAssetBundle(filename, projectName);
        if (result != null)
            return result;

        string loadFileName = string.Empty;
        if (filename.EndsWith(".lua"))
        {
            loadFileName = filename;
        }
        else
        {
            loadFileName = filename + ".lua";
        }


#if UNITY_EDITOR
        // 从Lua原始目录加载
        string fileFullName = string.Format(LuaConst.LuaDir + "/" + loadFileName, (projectName != AppDefine.HallNumber ? "Game/" : "") + projectName).Replace("//", "/");
        if (File.Exists(fileFullName))
        {
            return File.ReadAllBytes(fileFullName);
        }
        else
        {
            Debug.LogErrorFormat("Don't exsit lua file:[{0}]", fileFullName);
            return null;
        }
#endif

        return null;
    }

    public string LoadCustomLuaFileString(string filename, string projectName)
    {
        if (!CacheLuaScriptDic.ContainsKey(projectName))
            CacheLuaScriptDic.Add(projectName, new Dictionary<string, string>());

        if (!CacheLuaScriptDic[projectName].ContainsKey(filename))
        {
            byte[] luaContent = LuaManager.Instance.LoadCustomLuaFile(filename, projectName);
            string luaString = Utility.BytesToUTF8String(luaContent);

            if (!string.IsNullOrEmpty(luaString))
                CacheLuaScriptDic[projectName].Add(filename, luaString);
        }

        return CacheLuaScriptDic[projectName][filename];
    }

    public void ClearProjectCacheLuaScript(string projectName)
    {
        if (CacheLuaScriptDic.ContainsKey(projectName))
        {
            CacheLuaScriptDic.Remove(projectName);
        }
    }

    private static LuaEnv m_LuaEnv = new LuaEnv();

    public static LuaEnv LuaEnv
    {
        get { return m_LuaEnv; }
        private set { m_LuaEnv = value; }
    }

    void Update()
    {
        if (LuaEnv != null)
        {
            if (Time.time - lastGCTime > GCInterval)
            {
                if (LuaEnv != null)
                {
                    LuaEnv.Tick();
                    lastGCTime = Time.time;
                }
            }
        }
        if (null != UpdateEvent)
        {
            UpdateEvent(Time.deltaTime);
        }
    }

    public static event System.Action<float> UpdateEvent = null;

    /// <summary>
    /// 执行方法
    /// </summary>
    /// <param name="moduleName"></param>
    /// <param name="methodName"></param>
    /// <param name="args"></param>
    /// <returns></returns>
    public static object[] CallMethod(string moduleName, string methodName, params object[] args)
    {
        if (Instance != null)
        {
            LuaFunction func = LuaEnv.Global.GetInPath<LuaFunction>(string.Format("{0}.{1}", moduleName, methodName));
            if (func != null)
            {
                return func.Call(args);
            }
        }
        return null;
    }
}