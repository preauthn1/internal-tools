using UnityEngine;
using UnityEditor;
using example.invalid;
using System.Collections.Generic;
using System.Linq;

public class UISpriteBanding
{
    private static string ASSET_ROOT = Application.dataPath + "/" + AppDefine.CurrentProjectFolderPath + "/UI";
    private static string EXPORT_ROOT = Application.dataPath + "/" + AppDefine.CurrentProjectFolderPath + "/UISpritePrefab";

    /// <summary>
    /// 需要导出的UI图片精灵的前缀(不区分大小写)--填写小写字母
    /// </summary>
    static List<string> exportSpritePrefixList = new List<string>() { "sprite_", };

    /// <summary>
    /// 需要导出的文件后缀
    /// </summary>
    static List<string> exportFileExtensionList = new List<string>() { ".png", ".jpg", };

    /// <summary>
    /// 是否需要导出到prefab
    /// </summary>
    /// <param name="spriteName">精灵名称</param>
    /// <returns></returns>
    static bool IsNeedExportSpriteToPrefab(string spriteName)
    {
        spriteName = spriteName.ToLower();
        foreach (string item in exportSpritePrefixList)
        {
            if (spriteName.StartsWith(item))
            {
                return true;
            }
        }
        return false;
    }


    [MenuItem("Tools/刷新Sprite读取映射", false, 201)]
    public static void CoverExportAllUISptitesToPrefab()
    {
        ExportUISpritesToPrefab();
    }

    /// <summary>
    /// 是否是导出的文件导出
    /// </summary>
    /// <param name="fileName"></param>
    /// <returns></returns>
    static bool IsNeedExprotImageFile(string fileName)
    {
        foreach (string item in exportFileExtensionList)
        {
            if (fileName.EndsWith(item))
            {
                return true;
            }
        }
        return false;
    }

    static void ExportUISpritesToPrefab()
    {
        // 找出所有需要导出的文件路径
        IEnumerable<string> exportFileFullNames = Directory.GetFiles(ASSET_ROOT, "*.*", SearchOption.AllDirectories).Where(t => IsNeedExprotImageFile(t));
        m_HandleTotalCount = exportFileFullNames.Count<string>();
        m_HandleIndex = 0;
        // 创建预制件保存路径
        string prefabRootPath = EXPORT_ROOT;
        if (!Directory.Exists(prefabRootPath))
        {
            Directory.CreateDirectory(prefabRootPath);
        }
        // 循环导出
        foreach (string fileFullName in exportFileFullNames)
        {
            string relativePath = fileFullName.Substring(fileFullName.LastIndexOf("Assets"));
            EditorUtility.DisplayProgressBar(ProgressBarInfo, string.Format("Handle export file : {0}", relativePath), (float)m_HandleIndex / m_HandleTotalCount);
            ExprotImageFileToUISpritePrefabs(relativePath);
            m_HandleIndex++;
        }
        EditorUtility.ClearProgressBar();
        AssetDatabase.Refresh(ImportAssetOptions.ForceUpdate);
        AssetDatabase.SaveAssets();
    }

    const string ProgressBarInfo = "Export sprite to Prefab";
    static int m_HandleIndex = 0;
    static int m_HandleTotalCount = 1;
    /// <summary>
    /// 导出指定图片文件的UI精灵到预制件
    /// </summary>
    /// <param name="relativePath">资源的相对路径</param>
    private static void ExprotImageFileToUISpritePrefabs(string relativePath)
    {
        foreach (Object spriteObject in AssetDatabase.LoadAllAssetsAtPath(relativePath))
        {
            Sprite sprite = spriteObject as Sprite;
            if (sprite != null)
            {
                if (!IsNeedExportSpriteToPrefab(sprite.name))// 需要导出的前缀判断
                    continue;

                string prefabPath = FileUtil.GetProjectRelativePath(string.Format("{0}/{1}.prefab", EXPORT_ROOT, sprite.name));
                // 如果已经存在了，不重复创建                
                GameObject asset = AssetDatabase.LoadAssetAtPath<GameObject>(prefabPath);
                if (asset == null)
                {
                    EditorUtility.DisplayProgressBar(ProgressBarInfo, string.Format("Handle export sprite to prefab: {0}", sprite.name), (float)m_HandleIndex / m_HandleTotalCount);
                    GameObject go = new GameObject(sprite.name);
                    go.AddComponent<SpriteRenderer>().sprite = sprite;
                    PrefabUtility.CreatePrefab(prefabPath, go);
                    GameObject.DestroyImmediate(go);
                }
                else
                {
                    asset.GetComponent<SpriteRenderer>().sprite = sprite;
                }
            }
        }
    }
}
