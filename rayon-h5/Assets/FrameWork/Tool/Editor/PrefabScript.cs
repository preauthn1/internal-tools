using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using UnityEditor;
using UnityEngine;

namespace Assets.FrameWork.Tool.Editor
{
    public static class PrefabScript
    {
        [MenuItem("Tools/批量保存预制体", priority = 10)]
        private static void BatchSavePrefab()
        {
            //获取选中的gameobject对象  
            GameObject[] selectedsGameobject = Selection.gameObjects;

            GameObject prefab = PrefabUtility.FindPrefabRoot(selectedsGameobject[0]);

            for (int i = 0; i < selectedsGameobject.Length; i++)
            {
                GameObject obj = selectedsGameobject[i];

                UnityEngine.Object newsPref = PrefabUtility.GetPrefabObject(obj);

                //判断选择的物体，是否为预设  
                if (PrefabUtility.GetPrefabType(obj) == PrefabType.PrefabInstance || PrefabUtility.GetPrefabType(obj) == PrefabType.DisconnectedPrefabInstance)
                {

                    UnityEngine.Object parentObject = PrefabUtility.GetPrefabParent(obj);
                    //获取路径  
                    string path = AssetDatabase.GetAssetPath(parentObject);
                    Debug.Log("path:" + path);
                    //替换预设  
                    PrefabUtility.ReplacePrefab(obj, parentObject, ReplacePrefabOptions.ConnectToPrefab);
                    //刷新  
                    AssetDatabase.Refresh();
                }
            }
        }
    }
}
