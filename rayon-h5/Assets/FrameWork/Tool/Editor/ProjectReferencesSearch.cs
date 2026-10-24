using System.Collections;
using System.Collections.Generic;
using example.invalid;
using UnityEditor;
using UnityEditor.SceneManagement;
using UnityEngine;
using UnityEngine.UI;

/// <summary>
/// 2018.1.24 因为项目添加了可视化的Window->项目依赖 菜单，这个脚本暂时没有使用的必要
/// 
/// </summary>
public class ProjectReferencesSearch
{

    /// <summary>
    /// 检查打包项目是否引用了其他项目的资源
    /// </summary>      
    // [MenuItem("Assets/资源优化/查找项目外引用", priority = 20)]
    static private void AssetsCheck()
    {
        try
        {
            if (string.IsNullOrEmpty(AppDefine.CurrentProjectPath))
            {
                EditorUtility.DisplayDialog("错误", "没有设置当前项目，右键->设置为当前项目可设置当前项目。", "确定");
                return;
            }

            FindAllFontReferences();        // 字体单独处理
            FindOtherAllReferences();
        }
        catch (System.Exception ex)
        {
            throw ex;
        }
        finally
        {
            EditorUtility.ClearProgressBar();
        }
    }

    // 除字体之外的引用
    static private void FindOtherAllReferences()
    {
        string[] types = new string[] { "*.prefab", "*.unity" };

        for (int m = 0; m < types.Length; m++)
        {
            string[] paths = Directory.GetFiles(AppDefine.CurrentProjectPath, types[m], SearchOption.AllDirectories);

            string fileName = string.Empty;

            for (int n = 0; n < paths.Length; n++)
            {
                fileName = EditorFileUtility.GetFileName(paths[n].Replace(@"\", @"/"));

                EditorUtility.DisplayProgressBar("检查资源引用", string.Format("查找文件{0}的资源引用...", fileName), 1f * n / paths.Length);
                string[] dependpaths = AssetDatabase.GetDependencies(new string[] { paths[n] });
                foreach (string item in dependpaths)
                {
                    if (!item.Contains(AppDefine.CurrentProjectFolderPath) && !item.EndsWith(".cs") && !item.EndsWith(".ttf") && !item.Contains("/Gizmos/"))
                    {
                        Debug.LogFormat("<color=#3FFF20>{0}</color>引用了项目外部资源<color=#3FFF20>{1}</color>", fileName, item);
                    }
                }
            }

            EditorUtility.ClearProgressBar();
        }
    }

    static private void FindAllFontReferences()
    {
        string[] paths = Directory.GetFiles(AppDefine.CurrentProjectPath, "*.prefab", SearchOption.AllDirectories);

        int sencesCount = EditorBuildSettings.scenes.Length;

        try
        {
            for (int i = 0; i < paths.Length; i++)
            {
                GameObject prefabObj = AssetDatabase.LoadAssetAtPath(EditorFileUtility.GetRelativeAssetPath(paths[i]), typeof(GameObject)) as GameObject;

                EditorUtility.DisplayProgressBar("查找sprite引用", string.Format("查找预制体{0}的字体引用...", prefabObj.name), 1f * i / paths.Length);

                List<Text> imageList = new List<Text>();

                (prefabObj as GameObject).GetComponentsInChildren<Text>(true, imageList);  // 获取所有Text组件

                for (int j = 0; j < imageList.Count; j++)
                {
                    string fontPath = AssetDatabase.GetAssetPath(imageList[j].font);
                    string fileName = EditorFileUtility.GetFileName(paths[i].Replace(@"\", @"/"));

                    if (!fontPath.Contains("Library/") && !fontPath.Contains(AppDefine.CurrentProjectFolderPath) && !string.IsNullOrEmpty(fontPath))
                    {
                        Debug.LogFormat("<color=#3FFF20>{0}</color>引用了项目外部资源<color=#3FFF20>{1}</color>", fileName, fontPath);
                        break;
                    }
                }
            }

            string[] scenePaths = Directory.GetFiles(AppDefine.CurrentProjectPath, "*.unity", SearchOption.AllDirectories);

            if (scenePaths.Length > 0)
            {
                //打开场景
                EditorSceneManager.OpenScene(scenePaths[0]);

                //获取场景中的所有游戏对象
                GameObject[] gos = (GameObject[])GameObject.FindObjectsOfType(typeof(GameObject));

                foreach (GameObject go in gos)
                {
                    //判断GameObject是否为一个Prefab的引用
                    if (PrefabUtility.GetPrefabType(go) != PrefabType.PrefabInstance)
                    {
                        List<Text> textList = new List<Text>();

                        (go as GameObject).GetComponentsInChildren<Text>(true, textList);  // 获取所有Text组件

                        for (int j = 0; j < textList.Count; j++)
                        {
                            string fontPath = AssetDatabase.GetAssetPath(textList[j].font);
                            string fileName = EditorFileUtility.GetFileName(scenePaths[0].Replace(@"\", @"/"));

                            if (!fontPath.Contains("Library/") && !fontPath.Contains(AppDefine.CurrentProjectFolderPath) && !string.IsNullOrEmpty(fontPath))
                            {
                                Debug.LogFormat("<color=#3FFF20>{0}</color>引用了项目外部资源<color=#3FFF20>{1}</color>", fileName, fontPath);
                                return;
                            }
                        }
                    }
                }
            }
        }
        catch (System.Exception ex)
        {
            throw ex;
        }
        finally
        {
            EditorUtility.ClearProgressBar();
        }
    }

    // [MenuItem("Assets/资源优化/外部字体指向项目字体", priority = 21)]
    static private void ReplaceToProjectFont()
    {
        string[] paths = Directory.GetFiles(AppDefine.CurrentProjectPath, "*.prefab", SearchOption.AllDirectories);

        int sencesCount = EditorBuildSettings.scenes.Length;

        try
        {
            for (int i = 0; i < paths.Length; i++)
            {
                GameObject prefabObj = AssetDatabase.LoadAssetAtPath(EditorFileUtility.GetRelativeAssetPath(paths[i]), typeof(GameObject)) as GameObject;

                EditorUtility.DisplayProgressBar("替换字体", string.Format("替换预制体{0}的字体...", prefabObj.name), 1f * i / paths.Length);

                List<Text> imageList = new List<Text>();

                (prefabObj as GameObject).GetComponentsInChildren<Text>(true, imageList);  // 获取所有Text组件

                for (int j = 0; j < imageList.Count; j++)
                {
                    string fontPath = AssetDatabase.GetAssetPath(imageList[j].font);
                    string projectFontPath = "Assets/" + AppDefine.CurrentProjectFolderPath + "/Font/" + EditorFileUtility.GetFileName(fontPath);

                    if (!fontPath.Contains("Library/") && !fontPath.Contains(AppDefine.CurrentProjectFolderPath) && !string.IsNullOrEmpty(fontPath))
                    {
                        FileUtility.CopyFile(fontPath, projectFontPath, true);

                        AssetDatabase.Refresh();
                        Font myProjectFont = AssetDatabase.LoadAssetAtPath(EditorFileUtility.GetRelativeAssetPath(projectFontPath), typeof(Font)) as Font;
                        imageList[j].font = myProjectFont;

                        EditorUtility.SetDirty(prefabObj);
                        AssetDatabase.SaveAssets();
                    }
                }
            }

            string[] scenePaths = Directory.GetFiles(AppDefine.CurrentProjectPath, "*.unity", SearchOption.AllDirectories);

            if (scenePaths.Length > 0)
            {
                //打开场景
                UnityEngine.SceneManagement.Scene scene = EditorSceneManager.OpenScene(scenePaths[0]);

                //获取场景中的所有游戏对象
                GameObject[] gos = (GameObject[])GameObject.FindObjectsOfType(typeof(GameObject));

                foreach (GameObject go in gos)
                {
                    //判断GameObject是否为一个Prefab的引用
                    if (PrefabUtility.GetPrefabType(go) != PrefabType.PrefabInstance)
                    {
                        List<Text> textList = new List<Text>();

                        (go as GameObject).GetComponentsInChildren<Text>(true, textList);  // 获取所有Text组件

                        for (int j = 0; j < textList.Count; j++)
                        {
                            string fontPath = AssetDatabase.GetAssetPath(textList[j].font);
                            string projectFontPath = "Assets/" + AppDefine.CurrentProjectFolderPath + "/Font/" + EditorFileUtility.GetFileName(fontPath);

                            if (!fontPath.Contains("Library/") && !fontPath.Contains(AppDefine.CurrentProjectFolderPath) && !string.IsNullOrEmpty(fontPath))
                            {
                                FileUtility.CopyFile(fontPath, projectFontPath, true);

                                AssetDatabase.Refresh();
                                Font myProjectFont = AssetDatabase.LoadAssetAtPath(EditorFileUtility.GetRelativeAssetPath(projectFontPath), typeof(Font)) as Font;
                                textList[j].font = myProjectFont;

                                EditorSceneManager.MarkSceneDirty(scene);
                                EditorSceneManager.SaveScene(scene);
                            }
                        }
                    }
                }
            }
        }
        catch (System.Exception ex)
        {
            throw ex;
        }
        finally
        {
            EditorUtility.ClearProgressBar();
        }
    }
}
