using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class LuaConst
{
    /// <summary>
    /// lua逻辑代码目录
    /// </summary>
    public static string LuaDir = Application.dataPath + "/{0}/LuaScript/";

    public static string LuaABDir
    {
        get
        {
            return Application.dataPath + "/" + AppDefine.CurrentProjectFolderPath + "/LuaScript/";
        }
    }
}