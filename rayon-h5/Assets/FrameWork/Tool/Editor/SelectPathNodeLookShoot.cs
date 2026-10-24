using UnityEditor;
using UnityEngine;


public static class SelectPathNodeLookShoot
{
    const string MethodName = "GameObject/Tool/By/SelectPath NodeLookShoot";

    [MenuItem(MethodName, true, 1)]
    public static bool BuildScriptValidation()
    {
        return Selection.activeObject is GameObject;
    }

    [MenuItem(MethodName, false, 1)]
    public static void BuildScript()
    {
        GameObject root = Selection.activeObject as GameObject;

        GameObject _atthch = new GameObject();
        _atthch.name = "NodeLookShoot";

        GameObject _path = new GameObject();
        _path.name = root.name;
        _path.transform.SetParent(_atthch.transform);

        Transform _child = root.transform.parent;
        while (_child != null)
        {
            if (_child.parent == null)
            {
                _atthch.transform.SetParent(_child.transform);
                _atthch.transform.localPosition = Vector3.zero;
                _atthch.transform.localRotation = Quaternion.identity;
            }
            else
            { 
                _path.name = _child.name + "/" + _path.name;
            }
            _child = _child.parent;
        }
    }
}
