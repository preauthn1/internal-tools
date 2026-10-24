using System.Collections.Generic;
using example.invalid;
using System.Linq;
using UnityEditor;
using UnityEngine;

public class ProjectMenu : Editor
{
    const string NEWGAMEMENU = "Assets/新建游戏";

    [MenuItem(NEWGAMEMENU, false, 60)]
    public static void CreateGame()
    {
        ProjectEditUtility.CreateGameTemplateForlders();
        ProjectEditUtility.CreateGameConfig();
    }

    [MenuItem(NEWGAMEMENU, true, 60)]
    public static bool CreateGameValidate()
    {
        string[] assetGUIDArray = Selection.assetGUIDs;

        if (assetGUIDArray.Length == 1)
            return AssetDatabase.GUIDToAssetPath(assetGUIDArray[0]) == "Assets";

        return false;
    }

    const string SETCURRENTMENU = "Assets/设置为当前项目";

    // 设置为当前项目目录
    [MenuItem(SETCURRENTMENU, false, 61)]
    public static void SetCurrentProject()
    {
        string[] assetGUIDArray = Selection.assetGUIDs;

        if (assetGUIDArray.Length == 1)
            AppDefine.CurrentProjectPath = AssetDatabase.GUIDToAssetPath(assetGUIDArray[0]);
    }

    [MenuItem(SETCURRENTMENU, true, 61)]
    public static bool SelectProjectFounderValidate()
    {
        string[] assetGUIDArray = Selection.assetGUIDs;

        if (assetGUIDArray.Length == 1)
        {
            string path = AssetDatabase.GUIDToAssetPath(Selection.assetGUIDs[0]);
            return path == "Assets/Hall" || (path.Split('/').Length == 3 && path.Contains("Assets/Game"));
        }

        return false;
    }
}
