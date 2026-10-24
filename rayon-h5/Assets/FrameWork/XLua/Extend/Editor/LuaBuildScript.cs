using UnityEngine;
using UnityEditor;
using example.invalid;

public class LuaBuildScript
{
    private static string LuaPathOfPrimitive = Application.dataPath + "/" + AppDefine.CurrentProjectFolderPath + "/Primitives/";

    // [MenuItem("XLua/Copy Lua  files to Primitive", false, 52)]
    public static void CopyToPrimitiveLuaPath()
    {
        ClearAllLuaFiles();
        CopyLuaBytesFiles(LuaConst.LuaABDir, LuaPathOfPrimitive);
        AssetDatabase.Refresh(ImportAssetOptions.ForceUpdate); 
    }

    // [MenuItem("XLua/Clear Lua  files", false, 53)]
    public static void ClearAllCopyedLuaFiles()
    {
        ClearAllLuaFiles();
        AssetDatabase.Refresh();
        Debug.Log("Clear lua files over");
    }

    static void CopyLuaBytesFiles(string sourceDir, string destDir, bool appendExt = true, string searchPattern = "*.lua", SearchOption option = SearchOption.AllDirectories)
    {
        if (!Directory.Exists(sourceDir))
        {
            return;
        }

        string[] files = Directory.GetFiles(sourceDir, searchPattern, option);
        int len = sourceDir.Length;

        if (sourceDir[len - 1] == '/' || sourceDir[len - 1] == '\\')
        {
            --len;
        }

        for (int i = 0; i < files.Length; i++)
        {
            string str = files[i].Remove(0, len);
            string dest = destDir + "/" + str;
            if (appendExt) dest += ".bytes";
            string dir = Path.GetDirectoryName(dest);
            Directory.CreateDirectory(dir);
            File.Copy(files[i], dest, true);
        }
    }

    public static void ClearAllLuaFiles()
    {
        if (Directory.Exists(LuaPathOfPrimitive))
        {
            Directory.Delete(LuaPathOfPrimitive, true);
        }
    }
}
