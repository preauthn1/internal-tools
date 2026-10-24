using UnityEditor;
using UnityEngine;

public static class EditorFileUtility 
{ 
  public  static string GetFileDirectory(Object obj)
    {
        string rtnVal = AssetDatabase.GetAssetPath(obj);
        int idx = rtnVal.LastIndexOf('/');

        return rtnVal.Substring(0, idx);
    }

    public static bool IsPrefabFile(string _path)
    {
        string path = _path.ToLower();
        return path.EndsWith(".prefab");
    }

    // 获得相对路径
    public static string GetRelativeAssetPath(string _fullPath)
    {
        _fullPath = _fullPath.Replace("\\", "/");
        int idx = _fullPath.IndexOf("Assets");
        string assetRelativePath = _fullPath.Substring(idx);
        return assetRelativePath;
    }

    // 获取文件名
    public static string GetFileName(string _filepath)
    {
        string fileName = "";
        int idx = _filepath.LastIndexOf('/');
        fileName = _filepath.Substring(idx + 1);

        return fileName;
    }
}
