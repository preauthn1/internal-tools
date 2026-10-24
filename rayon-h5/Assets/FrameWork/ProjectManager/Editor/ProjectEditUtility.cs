using System.Collections;
using System.Collections.Generic;
using example.invalid;
using UnityEditor;
using UnityEngine;

public class ProjectEditUtility
{
    public const string CONFIGFILENAME = "GameConfig.asset";

    public static GameConfig GetGameConfig()
    {
        Object[] selection = Selection.GetFiltered(typeof(GameConfig), SelectionMode.Assets);

        if (selection.Length > 0 && selection[0] != null)
            return (GameConfig)selection[0];

        return null;
    }

    public static void CreateGameConfig()
    {
        GameConfig asset = ScriptableObject.CreateInstance<GameConfig>();

        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        if (path == "")
        { 
            path = "Assets";
        }
        else if (Path.GetExtension(path) != "")
        {
            path = path.Replace(Path.GetFileName(AssetDatabase.GetAssetPath(Selection.activeObject)), "");
        }

        string assetPathAndName = AssetDatabase.GenerateUniqueAssetPath(path + "/" + CONFIGFILENAME);

        AssetDatabase.CreateAsset(asset, assetPathAndName);

        AssetDatabase.SaveAssets();
        EditorUtility.FocusProjectWindow();
        Selection.activeObject = asset;
    }

    public static GameConfig GetCurrentProjectConfig()
    {
        if (string.IsNullOrEmpty(AppDefine.CurrentProjectPath))
            return null;

        return AssetDatabase.LoadAssetAtPath<GameConfig>(AppDefine.CurrentProjectPath + "/" + CONFIGFILENAME);
    }

    /// 标准文件夹 
    static string[] GameSubFolders = new string[] { "Animation", "UI", "LuaScript", "Prefab", "Audio", "Particle", "Font", "PbFile" };

    public static void CreateGameTemplateForlders()
    {
        string[] assetGUIDArray = Selection.assetGUIDs;

        if (assetGUIDArray.Length == 1)
        {
            string assetPath = AssetDatabase.GUIDToAssetPath(assetGUIDArray[0]);

            if (assetPath == "Assets")
            {
                string assetPathAndName = AssetDatabase.GenerateUniqueAssetPath(assetPath + "/Game/NewGame");
                AssetDatabase.CreateFolder(assetPath + "/Game", assetPathAndName.Replace("Assets/Game/", string.Empty));

                foreach (string item in GameSubFolders)
                    AssetDatabase.CreateFolder(assetPathAndName, item);

                Selection.activeObject = AssetDatabase.LoadAssetAtPath(assetPathAndName, typeof(DefaultAsset));
            }
        }
    }
}
