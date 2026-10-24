using System;
using System.Collections.Generic;
using example.invalid;
using System.Linq;
using System.Text;
using UnityEditor;
using UnityEngine;
using UnityEngine.UI;


/// <summary>
/// 生成类使用说明
/// 1.需要生成脚本的游戏对象必须在场景根节点下或者是在Project中的预制体
/// 2.Dic维护了一个类型字典表，字典表顺序为游戏对象脚本的生成顺序
/// 3.需要生成到脚本里面作为变量的游戏对象需要将名称修改为带"_"开头
/// </summary>

public static class BuildLuaScriptFromPrefab
{
    const string CREATELUASCRIPT = "GameObject/为GameObject创建Lua脚本";
    const string CREATELUASCRIPT_Asset = "Assets/为GameObject创建Lua脚本";

    [MenuItem(CREATELUASCRIPT, true, 0)]
    [MenuItem(CREATELUASCRIPT_Asset, true, 62)]
    public static bool SetCurrentProject()
    {
        return Selection.activeObject is GameObject;
    }

    [MenuItem(CREATELUASCRIPT, false, 0)]
    [MenuItem(CREATELUASCRIPT_Asset, false, 62)]
    public static void BuildScript()
    {
        Dictionary<Type, string> m_TypeDic = new Dictionary<Type, string>();        // 类型字典表 key: 可生成类型 value: 生成变量后的前缀

        m_TypeDic.Add(typeof(Toggle), "Tog");
        m_TypeDic.Add(typeof(RawImage), "RImg");
        m_TypeDic.Add(typeof(ScrollRect), "Scrr");
        m_TypeDic.Add(typeof(InputField), "Ipt");
        m_TypeDic.Add(typeof(Slider), "Sld");
        m_TypeDic.Add(typeof(Scrollbar), "Scrb");
        m_TypeDic.Add(typeof(Text), "Txt");
        m_TypeDic.Add(typeof(Button), "Btn");
        m_TypeDic.Add(typeof(Image), "Img");
        m_TypeDic.Add(typeof(RectTransform), "RecTran");
        m_TypeDic.Add(typeof(Transform), "Tran");

        if (Selection.activeGameObject == null)
        {
            Debug.LogWarning("请选择需要生成脚本的对象!");
            return;
        }

        for (int i = 0; i < Selection.gameObjects.Length; i++)
        {
            GameObject selectobj = Selection.gameObjects[i];

            Transform[] transforms = selectobj.GetComponentsInChildren<Transform>(true);
            List<Transform> uiTransList = transforms.ToList().FindAll(m => m.name.StartsWith("_"));

            string definition = "";         // 定义部分
            string scriptContent = "";      // 脚本内容
            string buttonEvent = "";        // 按钮事件内容生成
            string notes = "--- 以上内容为工具脚本自动生成，请不要在这之上做任何修改---"; // 注释部分
            string fileString = "";         // 最后写入文件的内容

            foreach (Transform current in uiTransList)
            {
                Transform transform = current;
                List<string> list = new List<string>();
                list.Add(transform.name);
                while (transform != selectobj.transform)
                {
                    transform = transform.parent;
                    list.Add(transform.name);
                }

                string path = string.Empty;

                for (int j = 0; j < list.Count - 1; j++)
                {
                    path = "/" + list[j] + path;
                }

                path = path.TrimStart('/');
                string currentObjName = string.Empty;           // Lua中参数的名称

                foreach (var item in m_TypeDic)
                {
                    if (current.GetComponent(item.Key))
                    {
                        currentObjName = (m_TypeDic[item.Key] + current.name).Replace(" ", "");

                        definition += string.Format("local {0} = this.transform:Find(\"{1}\"):GetComponent(\"{2}\")\r\n", currentObjName, path, item.Key.Name);

                        if (item.Key == typeof(Button))
                        {
                            buttonEvent += string.Format(@"  
    {0}:AddClickListener
    (
    function()

    end
    )
", currentObjName);
                        }

                        if (item.Key == typeof(Toggle))
                        {
                            buttonEvent += string.Format(@"  
    {0}.onValueChanged:AddListener
    (
    function(isOn)

    end
    )
", currentObjName);
                        }
                        break;
                    }
                }
            }

            string filePath = AssetDatabase.GetAllAssetPaths().ToList().Find(m => m.Contains("Assets/" + AppDefine.CurrentProjectFolderPath + "/LuaScript") && Path.GetFileNameWithoutExtension(m) == selectobj.name);

            if (string.IsNullOrEmpty(filePath))
            {
                scriptContent =
@"

-- 初始化函数，实例化之后调用一次
function Start()

#ButtonEvent#

end

-- 每帧调用一次
function Update()

end

-- 脚本删除时调用，释放变量，释放资源
function OnDestroy()

end  

-- 界面被打开时接收界面参数
function RefreshWindowData(windowData)

end

-- 在Window下按下Esc或在Android 下按下Back
function EscapeWindow()

end
";

                fileString = definition + "\r\n" + notes + scriptContent;
                fileString = fileString.Replace("#ButtonEvent#", buttonEvent);

                FileStream fileStream = new FileStream(Application.dataPath + "/" + AppDefine.CurrentProjectFolderPath + "/LuaScript/" + selectobj.name + ".lua", FileMode.CreateNew);
                StreamWriter streamWriter = new StreamWriter(fileStream, new UTF8Encoding(false));
                streamWriter.Write(fileString);
                streamWriter.Flush();
                streamWriter.Close();
                fileStream.Close();

                Debug.Log(string.Format("创建脚本{0}/{1}/LuaScript/{2}.lua 成功!", Application.dataPath, AppDefine.CurrentProjectFolderPath, selectobj.name));
            }
            else
            {
                FileStream fileStream_Open = new FileStream(Path.GetFullPath(filePath), FileMode.Open);
                StreamReader streamWriter = new StreamReader(fileStream_Open);
                string oldScriptContent = streamWriter.ReadToEnd();
                streamWriter.Close();
                fileStream_Open.Close();

                scriptContent = oldScriptContent.Substring(oldScriptContent.IndexOf(notes) + notes.Length, oldScriptContent.Length - oldScriptContent.IndexOf(notes) - notes.Length);

                fileString = definition + "\r\n" + notes + scriptContent;

                FileStream fileStream_Create = new FileStream(Path.GetFullPath(filePath), FileMode.Create);
                StreamWriter streamWriter2 = new StreamWriter(fileStream_Create, new UTF8Encoding(false));
                streamWriter2.Write(fileString);
                streamWriter2.Flush();
                streamWriter2.Close();
                fileStream_Create.Close();
                Debug.Log("更新脚本 " + selectobj.name + ".lua 成功!");
            }
             
            AssetDatabase.Refresh();
            AssetDatabase.SaveAssets();
        }
    }
}
