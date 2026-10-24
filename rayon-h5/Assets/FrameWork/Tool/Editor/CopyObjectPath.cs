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
/// 选中游戏对象生成路径信息 
/// </summary>

public class CopyObjectPath
{
    const string CREATELUASCRIPT = "Tools/复制路径信息";
     
    [MenuItem(CREATELUASCRIPT, true, 62)]
    public static bool SetCurrentProject()
    {
        return Selection.activeObject is GameObject;
    }

    [MenuItem(CREATELUASCRIPT, false, 63)]
    public static void BuildScript()
    {
        Dictionary<Type, string> m_TypeDic = new Dictionary<Type, string>();        // 类型字典表 key: 可生成类型 value: 生成变量后的前缀

        m_TypeDic.Add(typeof(Toggle), "tog");
        m_TypeDic.Add(typeof(ScrollRect), "scrr");
        m_TypeDic.Add(typeof(InputField), "ipt");
        m_TypeDic.Add(typeof(Slider), "sld");
        m_TypeDic.Add(typeof(Scrollbar), "scrb");
        m_TypeDic.Add(typeof(Text), "txt");
        m_TypeDic.Add(typeof(Button), "btn");
        m_TypeDic.Add(typeof(Image), "img");
        m_TypeDic.Add(typeof(RawImage), "rImg");
        m_TypeDic.Add(typeof(RectTransform), "recTran");
        m_TypeDic.Add(typeof(Transform), "tran");


        List<Type> rootScriptType = new List<Type>();        // 类型字典表 key: 可生成类型 value: 生成变量后的前缀

        rootScriptType.Add(typeof(LuaWindowBase));
        rootScriptType.Add(typeof(LuaBehaviour));
        rootScriptType.Add(typeof(LuaGameBehaviour));

        if (Selection.activeGameObject == null)
        {
            Debug.LogWarning("请选择需要复制路径的游戏对象!");
            return;
        }

        string definition = "";         // 定义部分  

        for (int i = 0; i < Selection.gameObjects.Length; i++)
        {
            GameObject selectobj = Selection.gameObjects[i];
             
            Transform transform = selectobj.transform;
            List<string> list = new List<string>();
            list.Add(transform.name);

            bool isBreak = false;
            while (!isBreak)
            {
                if (transform == transform.root) break;

                rootScriptType.ForEach((type) => { if (transform.GetComponent(type)) isBreak = true; });

                if (isBreak) break;

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
                if (selectobj.transform.GetComponent(item.Key))
                {
                    currentObjName = (m_TypeDic[item.Key] + "_" + selectobj.transform.name).Replace(" ", "");

                    definition += string.Format("local {0} = this.transform:Find(\"{1}\"):GetComponent(\"{2}\")\r\n", currentObjName, path, item.Key.Name);
                    break;
                }
            } 
        }


        TextEditor te = new TextEditor();//很强大的文本工具  
        te.text = definition;
        te.OnFocus();
        te.Copy();

        Debug.Log(definition);
    }
}
