using Microsoft.AspNetCore.Html;
using Microsoft.AspNetCore.Mvc.Rendering;
using MSC.Models.CultureText;
using MSC.Models.ViewModels.Base;
using MSC.Utility.Helper;
using MSC.Utility.Helper.HttpContext;
using System;
using System.Collections.Generic;


namespace MSC.Utility.Helper
{
    public static class MyHTMLHelpers
    {
        public static IHtmlContent RenderTextArea(this IHtmlHelper htmlHelper, string labelText, string vue_model, string v_if = "")
        {
            string v_if_tag = string.IsNullOrEmpty(v_if) ? "" : $"v-if='{v_if}'";
            return new HtmlString($@"
                <div {v_if_tag} >
                    <div class='form-label'>
                        <label>{labelText}</label>
                    </div>

                    <div>
                        <textarea class='form-control'  v-model='{vue_model}' placeHolder=''></textarea>
                    </div>
                </div>
             ");
        }


        public static IHtmlContent RenderInput(this IHtmlHelper htmlHelper, string labelText, string vue_model , string v_if = "")
        {
            string v_if_tag = string.IsNullOrEmpty(v_if) ? "" : $"v-if='{v_if}'";
            return new HtmlString($@"
                <div {v_if_tag} >
                    <div class='form-label'>
                        <label>{labelText}</label>
                    </div>

                    <div>
                        <input type='text' class='form-control'  v-model='{vue_model}' placeHolder=''/>
                    </div>
                </div>
             ");
        }

        public static IHtmlContent RenderDisabledInput(this IHtmlHelper htmlHelper, string labelText, string vue_model)
        {
            return new HtmlString($@"
                <div>
                    <div class='form-label'>
                        <label>{labelText}</label>
                    </div>

                    <div>
                        <input type='text' class='form-control'  v-model='{vue_model}' placeHolder='' disabled='disabled'/>
                    </div>
                </div>
             ");
        }

        public static IHtmlContent RenderDropdown(this IHtmlHelper htmlHelper, string labelText, string vue_model, string dropdownSource ,  string v_if = "")
        {
            string vue_expression_text = "{{item.text}}";
            string v_if_tag = string.IsNullOrEmpty(v_if) ? "" : $"v-if='{v_if}'";
            return new HtmlString($@"
                    <div {v_if_tag} >
                      <div class='form-label'>
                           <label>{labelText}</label>
                      </div>
                      <div>
                            <select class='form-control' v-model='{vue_model}'>
                                <option value=''>{t.Text(i18nKeys.Common_Label_Please_Select)} </option>
                                <option v-for='(item, index) in {dropdownSource}' :value='item.value'>
                                    {vue_expression_text}
                                </option>
                            </select>
                       </div>
                    </div>");
        }




    }
}