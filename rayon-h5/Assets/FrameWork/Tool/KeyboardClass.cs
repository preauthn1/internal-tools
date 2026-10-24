using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using UnityEngine.EventSystems;
using System.Runtime.InteropServices;
using System;

public class KeyboardClass : MonoBehaviour, IPointerClickHandler
{

	[DllImport("__Internal")]
	private static extern void FocusHandleAction (string _name, string _str);

	public void ReceiveInputData(string value) {
		gameObject.GetComponent<InputField>().text = value;
    }

    public void OnPointerClick(PointerEventData _data)
    {
		#if UNITY_WEBGL
		try{
			FocusHandleAction (gameObject.name, gameObject.GetComponent<InputField> ().text);
		}
		catch(Exception error){print(error);}
        #endif
    }
}
