using UnityEditor;
using UnityEngine;


public static class MakePathMeshRenderer
{
    const string MethodName = "GameObject/Tool/By/Path MeshRenderer";
    const string AttachName = "PathMeshRenderer";

    [MenuItem(MethodName, true, 1)]
    public static bool BuildScriptValidation()
    {
        return Selection.activeObject is GameObject;
    }

    [MenuItem(MethodName, false, 1)]
    public static void BuildScript()
    {
        GameObject root = Selection.activeObject as GameObject;

        GameObject _atthch = root;

        if (root.transform.Find(AttachName) != null)
        {
            _atthch = root.transform.Find(AttachName).gameObject;
            Object.DestroyImmediate(_atthch.gameObject);
        }

        MeshRenderer[] list = root.GetComponentsInChildren<MeshRenderer>(true);

        if (list.Length > 0)
        {
            _atthch = new GameObject();
            _atthch.name = AttachName;
            _atthch.transform.SetParent(root.transform);
        }

        for (int i = 0; i < list.Length; i++)
        {
            GameObject _path = new GameObject();
            _path.name = list[i].name;
            _path.transform.SetParent(_atthch.transform);

            Transform _child = list[i].transform.parent;
            while (_child != null && _child.gameObject != root)
            {
                _path.name = _child.name + "/" + _path.name;
                _child = _child.parent;
            }
        }
    }
}
