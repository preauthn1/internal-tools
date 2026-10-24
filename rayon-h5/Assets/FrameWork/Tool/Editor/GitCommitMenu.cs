using System.Collections;
using System.Collections.Generic;
using example.invalid;
using UnityEditor;
using UnityEditor.SceneManagement;
using UnityEngine;
using UnityEngine.UI;

public class GitCommitMenu
{
    const int priorityIndex = 1;

    [MenuItem("Assets/GitCommit/Reset and Pull", priority = priorityIndex)]
    static private void GitPull()
    {
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        GitCommit.Pull(path);
    }
    [MenuItem("Assets/GitCommit/Reset", priority = priorityIndex + 1)]
    static private void GitReset()
    {
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        GitCommit.Reset(path);
    }
    const string GITCHECKLOG = "Assets/GitCommit/Log/";
    [MenuItem(GITCHECKLOG + "Last", priority = priorityIndex + 2)]
    static private void GitLogLast()
    {
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        GitCommit.Show(path);
    }
    [MenuItem(GITCHECKLOG + "ALL", priority = priorityIndex + 2)]
    static private void GitLogALL()
    {
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        GitCommit.Log(path);
    }
    const string GITCHECKOUT = "Assets/GitCommit/Reset and Checkout/";
    const string BranchName3 = "UI3.0";
    [MenuItem(GITCHECKOUT + BranchName3, false, priority = priorityIndex + 3)]
    static private void GitCheckout3()
    {
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        GitCommit.Checkout(path, BranchName3);
    }
    [MenuItem(GITCHECKOUT + BranchName3, true, priorityIndex + 3)]
    static private bool SelectCheckout3()
    {
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        if (path.IndexOf("Game") < 0 && path.IndexOf("Hall") < 0)
            return false;
        return true;
    }

    const string BranchName4 = "UI4.0";
    [MenuItem(GITCHECKOUT + BranchName4, false, priority = priorityIndex + 3)]
    static private void GitCheckout4()
    {
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        GitCommit.Checkout(path, BranchName4);
    }
    [MenuItem(GITCHECKOUT + BranchName4, true, priorityIndex + 3)]
    static private bool SelectCheckout4()
    {
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        if (path.IndexOf("Game") < 0 && path.IndexOf("Hall") < 0)
            return false;
        return true;
    }

    const string BranchName5 = "UI5.0";
    [MenuItem(GITCHECKOUT + BranchName5, false, priorityIndex + 3)]
    static private void GitCheckout5()
    {
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        GitCommit.Checkout(path, BranchName5);
    }
    [MenuItem(GITCHECKOUT + BranchName5, true, priorityIndex + 3)]
    static private bool SelectCheckout5()
    {
        string path = AssetDatabase.GetAssetPath(Selection.activeObject);
        if (path.IndexOf("Game") < 0 && path.IndexOf("Hall") < 0)
            return false;
        return true;
    }
}
