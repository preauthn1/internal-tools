using UnityEngine;
using UnityEditor;
using example.invalid;
using System.Collections.Generic;
using UnityEngine.UI;
using UnityEditor.SceneManagement;
using UnityEngine.SceneManagement;

/// <summary>
/// 打包时候对子项目依赖大厅的字体做处理，减少打包文件的大小
/// </summary>
public class FontBuildScript
{
    private static string PathOfBackup = Path.GetFullPath(Application.dataPath + "/../BuildBackup/");

    public static bool HallHasFont(string FontName)
    {
        if (FontName.Contains("font_fzycjt") ||
            FontName.Contains("MicrosoftYaHei"))
        {
            return true;
        }
        return false;
    }

    public static void Backup()
    {
#if !UNITY_WEBGL
        if (AppDefine.CurrentProjectFolder == AppDefine.HallNumber)
            return;
#endif

        ClearAllBackupFiles();

        if (!Directory.Exists(PathOfBackup))
            Directory.CreateDirectory(PathOfBackup);

        DirectoryInfo rootFolder = new DirectoryInfo(AppDefine.CurrentProjectPath);

        if (rootFolder.Exists)
        {
            FileInfo[] prefabInfos = rootFolder.GetFiles("*.prefab", SearchOption.AllDirectories);

            for (int i = 0; i < prefabInfos.Length; i++)
            {
                GameObject prefabObj = AssetDatabase.LoadAssetAtPath(EditorFileUtility.GetRelativeAssetPath(prefabInfos[i].FullName), typeof(GameObject)) as GameObject;

                Text[] textArray = prefabObj.GetComponentsInChildren<Text>(true);

                if (textArray.Length != 0)
                {
                    //FileUtility.CopyFile(prefabInfos[i].FullName, prefabInfos[i].FullName.Replace(@"\Assets\", @"\BuildBackup\"));
                    FileUtility.CopyFile(prefabInfos[i].FullName, prefabInfos[i].FullName.Replace(@"\", @"/").Replace(@"/Assets/", @"/BuildBackup/"));

                    for (int j = 0; j < textArray.Length; j++)
                    {
                        Text text = textArray[j];
                        if (text.font != null)
                        {
                            if (HallHasFont(text.font.name) == true)
                            {
                                BindingFont bindingFont = text.gameObject.GetComponent<BindingFont>();
                                if (bindingFont == null)
                                    bindingFont = text.gameObject.AddComponent<BindingFont>();

                                bindingFont.SetFontName(text.font.name);
                                text.font = null;
                            }
                        }
                    }
#if UNITY_WEBGL
                    SetWebGLText(textArray);
                    InputField[] inputFieldArray = prefabObj.GetComponentsInChildren<InputField>(true);
                    SetWebGLInputField(inputFieldArray);
#endif
                    EditorUtility.SetDirty(prefabObj);
                    AssetDatabase.SaveAssets();
                }
            }

            FileInfo[] sceneInfos = rootFolder.GetFiles("*.unity", SearchOption.AllDirectories);
            for (int i = 0; i < sceneInfos.Length; i++)
            {
                //FileUtility.CopyFile(sceneInfos[i].FullName, sceneInfos[i].FullName.Replace(@"\Assets\", @"\BuildBackup\"));
                FileUtility.CopyFile(sceneInfos[i].FullName, sceneInfos[i].FullName.Replace(@"\", @"/").Replace(@"/Assets/", @"/BuildBackup/"));

                //打开场景
                Scene scene = EditorSceneManager.OpenScene(sceneInfos[i].FullName);

                //获取场景中的所有游戏对象
                GameObject[] gos = (GameObject[])GameObject.FindObjectsOfType(typeof(GameObject));

                foreach (GameObject go in gos)
                {
                    Text[] textArray = go.GetComponentsInChildren<Text>(true);  // 获取所有Text组件
                    
                    for (int j = 0; j < textArray.Length; j++)
                    {
                        Text text = textArray[j];

                        if (text.font != null)
                        {
                            if (HallHasFont(text.font.name) == true)
                            {
                                // 当预制体与场景中放置的预制体字体有冲突的时候，将场景中的预制体实例font置空无效，这里将场景中的Prefab实例和Prefab断开联系
                                if (PrefabUtility.GetPrefabType(text.gameObject) == PrefabType.PrefabInstance)
                                    PrefabUtility.DisconnectPrefabInstance(text.gameObject);

                                BindingFont bindingFont = text.gameObject.GetComponent<BindingFont>();
                                if (bindingFont == null)
                                    bindingFont = text.gameObject.AddComponent<BindingFont>();

                                bindingFont.SetFontName(text.font.name);

                                text.font = null;
                            }
                        }
                    }
#if UNITY_WEBGL
                    SetWebGLText(textArray);
                    InputField[] inputFieldArray = go.GetComponentsInChildren<InputField>(true);// 获取所有InputField组件
                    SetWebGLInputField(inputFieldArray);
#endif
                }
                EditorSceneManager.SaveScene(scene);
            }
        }
        AssetDatabase.Refresh(ImportAssetOptions.ForceUpdate);
    }

    private static void SetWebGLInputField(InputField[] inputFieldArray)
    {
        for (int i = 0; i < inputFieldArray.Length; i++)
        {
            InputField inputField = inputFieldArray[i];
            if (inputField.GetComponent<KeyboardClass>() == null)
                inputField.gameObject.AddComponent<KeyboardClass>();
        }
    }

    private static void SetWebGLText(Text[] textArray)
    {
        for (int i = 0; i < textArray.Length; i++)
        {
            //设定为其他的Font，因WebGL无法显示Arial
            if (textArray[i].font == null || textArray[i].font.name == "Arial")
                textArray[i].font = AssetDatabase.LoadAssetAtPath<Font>("Assets/Hall/Font/MicrosoftYaHei.ttf");
        }
    }



    public static void Restore()
    {
#if !UNITY_WEBGL
        if (AppDefine.CurrentProjectFolder == AppDefine.HallNumber)
            return;
#endif
        DirectoryInfo rootFolder = new DirectoryInfo(PathOfBackup);

        if (rootFolder.Exists)
        {
            FileInfo[] prefabInfos = rootFolder.GetFiles("*.prefab", SearchOption.AllDirectories);

            for (int i = 0; i < prefabInfos.Length; i++)
                //FileUtility.CopyFile(prefabInfos[i].FullName, prefabInfos[i].FullName.Replace(@"\BuildBackup\", @"\Assets\"), true);
                FileUtility.CopyFile(prefabInfos[i].FullName, prefabInfos[i].FullName.Replace(@"\", @"/").Replace(@"/BuildBackup/", @"/Assets/"), true);


            FileInfo[] sceneInfos = rootFolder.GetFiles("*.unity", SearchOption.AllDirectories);

            for (int i = 0; i < sceneInfos.Length; i++)
                //FileUtility.CopyFile(sceneInfos[i].FullName, sceneInfos[i].FullName.Replace(@"\BuildBackup\", @"\Assets\"), true);
                FileUtility.CopyFile(sceneInfos[i].FullName, sceneInfos[i].FullName.Replace(@"\", @"/").Replace(@"/BuildBackup/", @"/Assets/"), true);
        }

        AssetDatabase.Refresh(ImportAssetOptions.ForceUpdate);
    }

    static void ClearAllBackupFiles()
    {
        if (Directory.Exists(PathOfBackup))
        {
            Directory.Delete(PathOfBackup, true);
            AssetDatabase.Refresh();
        }
    }
}
