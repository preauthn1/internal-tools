/*
    标记当前正在编辑的游戏，标记其他游戏的状态
 */

using Boo.Lang;
using System;
using System.Reflection;
using UnityEditor;
using UnityEditor.VersionControl;
using UnityEngine;
using ProjectWindowItemCallback = UnityEditor.EditorApplication.ProjectWindowItemCallback;

[InitializeOnLoad]
public class RainbowFoldersBrowserIcons
{
    private const float LARGE_ICON_SIZE = 64f;

    private static Func<bool> _isCollabEnabled;
    private static Func<bool> _isVcsEnabled;
    private static ProjectWindowItemCallback _drawCollabOverlay;
    private static ProjectWindowItemCallback _drawVcsOverlay;
     
    // Ctors 
    static RainbowFoldersBrowserIcons()
    {
        EditorApplication.projectWindowItemOnGUI += ReplaceFolderIcon;

        var assembly = typeof(EditorApplication).Assembly;
        InitCollabDelegates(assembly);
        InitVcsDelegates(assembly); 
    }

    // Delegates 
    private static void ReplaceFolderIcon(string guid, Rect rect)
    {
        var path = AssetDatabase.GUIDToAssetPath(guid);

        if (!AssetDatabase.IsValidFolder(path))
            return;
         
        if (path != AppDefine.CurrentProjectPath)
            return;

        var isSmall = IsIconSmall(ref rect);

        var texture = AssetDatabase.LoadAssetAtPath<Texture2D>("Assets/FrameWork/ProjectManager/Img/FolderIcon" + (isSmall ? "_S" : "") + ".png");

        if (texture == null) return;

        DrawCustomIcon(guid, rect, texture, isSmall);
    }

    // Helpers 
    private static void InitVcsDelegates(Assembly assembly)
    {
        try
        {
            _isVcsEnabled = () => Provider.isActive;

            var vcsHookType = assembly.GetType("UnityEditorInternal.VersionControl.ProjectHooks");
            var vcsHook = vcsHookType.GetMethod("OnProjectWindowItem", BindingFlags.Static | BindingFlags.Public);
            _drawVcsOverlay = (ProjectWindowItemCallback)Delegate.CreateDelegate(typeof(ProjectWindowItemCallback), vcsHook);
        }
        catch (SystemException ex)
        {
            if (!(ex is NullReferenceException) && !(ex is ArgumentNullException)) throw;
            _isVcsEnabled = () => false;
        }
    }

    private static void InitCollabDelegates(Assembly assembly)
    {
        try
        {
            var collabAccessType = assembly.GetType("UnityEditor.Web.CollabAccess");
            var collabAccessInstance = collabAccessType.GetProperty("Instance", BindingFlags.Static | BindingFlags.Public).GetValue(null, null);
            var collabAccessMethod = collabAccessInstance.GetType().GetMethod("IsServiceEnabled", BindingFlags.Instance | BindingFlags.Public);
            _isCollabEnabled = (Func<bool>)Delegate.CreateDelegate(typeof(Func<bool>), collabAccessInstance, collabAccessMethod);

            var collabHookType = assembly.GetType("UnityEditor.Collaboration.CollabProjectHook");
            var collabHook = collabHookType.GetMethod("OnProjectWindowItemIconOverlay", BindingFlags.Static | BindingFlags.Public);
            _drawCollabOverlay = (ProjectWindowItemCallback)Delegate.CreateDelegate(typeof(ProjectWindowItemCallback), collabHook);
        }
        catch (SystemException ex)
        {
            if (!(ex is NullReferenceException) && !(ex is ArgumentNullException)) throw;
            _isCollabEnabled = () => false;
        }
    }

    private static void DrawCustomIcon(string guid, Rect rect, Texture texture, bool isSmall)
    {
        if (rect.width > LARGE_ICON_SIZE)
        {
            // center the icon if it is zoomed
            var offset = (rect.width - LARGE_ICON_SIZE) / 2f;
            rect = new Rect(rect.x + offset, rect.y + offset, LARGE_ICON_SIZE, LARGE_ICON_SIZE);
        }
        else
        {
            if (isSmall && !IsTreeView(rect))
                rect = new Rect(rect.x + 3, rect.y, rect.width, rect.height);
        }

        if (_isCollabEnabled())
        {
            GUI.DrawTexture(rect, texture);
            _drawCollabOverlay(guid, rect);
        }
        else if (_isVcsEnabled())
        {
            var iconRect = (!isSmall) ? rect : new Rect(rect.x + 7, rect.y, rect.width, rect.height);
            GUI.DrawTexture(iconRect, texture);
            _drawVcsOverlay(guid, rect);
        }
        else
        {
            GUI.DrawTexture(rect, texture);
        }
    }

    private static bool IsIconSmall(ref Rect rect)
    {
        var isSmall = rect.width > rect.height;

        if (isSmall)
            rect.width = rect.height;
        else
            rect.height = rect.width;

        return isSmall;
    }

    private static bool IsTreeView(Rect rect)
    {
        return (rect.x - 16) % 14 == 0;
    }
}