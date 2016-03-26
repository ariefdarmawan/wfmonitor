function msValue(obj, returnType) {
    if (obj.data("kendoMultiSelect") == undefined) return;
    var value = obj.data("kendoMultiSelect").value();
    var ret = value;
    if (arguments.length > 1 && arguments[1].toLowerCase() == "string") {
        ret = value.join("|");
    }
    return ret;
}

function toArray(objs) {
    return $.map(objs, function (obj, idx) { return obj; });
}

function loadJs(scriptPath) {
    var oScript = $("<scr" + "ipt" + "></scr" + "ipt>");
    oScript.attr("type", "text/javascript").attr("src", scriptPath);
    $("#script_dynamic").append(oScript);
}

function devide(v1, v2) {
    var ret = 0;
    if (v2 == 0) {
        ret = 0;
    }
    else {
        ret = v1 / v2;
    }
    return ret;
}

function date2time(time)
{
    return kendo.format("{0:HH:mm}", time);
}

function time2decimal(time, nextDay) {
    var times = time.split(":");
    if (times[0] == "00" && nextDay==1) times[0] = "24";
    return parseFloat(times[0]) + (times[1] / 60.00);
}

function decimal2time(time)
{
    var hours = Math.floor(time * 10 / 10);
    var mins = time - hours;
    return hours.toString() + ":" + mins;
}

function input2timePicker(objects) {
    $.each(objects, function (idx, obj) {
        var jobj = $(obj);
        if (jobj.data("kendoTimePicker") == undefined) {
            jobj.kendoTimePicker({format:"HH:mm", interval:30});
        }
    });
}

function getObjProperties(obj) {
    var keys = [];
    for (var key in obj) {
        keys.push(key);
    }
    return keys;
}

function getObjProperty(obj, prop, def) {
    if (def == undefined) def = "";
    if (obj.hasOwnProperty(prop))
        return obj[prop];
    else return def;
}

function cbToogle(obj, selectorTxt) {
    var cbxs = $(selectorTxt);
    var checked = obj.prop("checked");
    cbxs.prop("checked", checked);
}


function gridDelete(deleteProcessUrl, deletedCheckboxes, fnDeleteSuccess, fnNoDelete) {
    var DeletedIds = $.map(deletedCheckboxes,function(obj,idx){
        return $(obj).val();
    });
    if (DeletedIds.length > 0) {
        if (!confirm("Are you sure you want to delete selected record(s) ?")) return;
        executeOnServer(viewModel, deleteProcessUrl + encodeURIComponent(ko.mapping.toJSON(DeletedIds)), fnDeleteSuccess);
    }
    else {
        if (typeof fnNoDelete == "function") fnNoDelete();
    }
}

function dlgHide(ow) {
    if (ow.data("kendoWindow")) {
        var kw = coalesce(ow.data("kendoWindow"), ow.kendoWindow({modal:false}).data("kendoWindow"));
        kw.close();
    }
}

function dlgShow(ow, title, fnClose) {
    var kw = null;
    if (!ow.data("kendoWindow")) {
        ow.kendoWindow({
            position: {
                top:50
            },
            title: title,
            //visible: false,
            modal: false,
            close: typeof fnClose=="function" ? 
                fnClose :
                function () { }
        });
        kw = ow.data("kendoWindow");
        kw.open();
    }
    else {
        kw = ow.data("kendoWindow");
        kw.open();
    }
    ow.show();
    kw.center();
}

function goto(url)
{
    location.href = url;
}

function json2ObsArray(url, data, arrayObject, fnOk) {
    var koResult = "";
    $.ajax({
        url: url,
        cache: false,
        type: 'POST',
        data: ko.mapping.toJSON(data),
        contentType: "application/json; charset=utf-8",
        success: function (data) {
            arrayObject.removeAll();
            if (data.length > 0) {
                data.forEach(function (x) {
                    x = ko.mapping.fromJS(x);
                    arrayObject.push(x);
                });
            }
            if (typeof fnOk == "function") fnOk(data);
            if (status != undefined) status.value = "OK";
        },
        error: function (error) {
            if (status != undefined) status.value = error.responseText;
            alert("There was an error posting the data to the server: " + error.responseText);
        }
    });

    return koResult;
}

function json2Obs(url, data, observableObject, fnOk, dataProperty) {
    var koResult = "";
    $.ajax({
        url: url,
        cache: false,
        type: 'POST',
        data: ko.mapping.toJSON(data),
        contentType: "application/json; charset=utf-8",
        success: function (data) {
            var dataObj = data;
            if (dataProperty != undefined && data.hasOwnProperty(dataProperty))
            {
                dataObj = data[dataProperty];
            }
            var koObj = null;
            if (observableObject == null || observableObject == undefined) {
                koObj = ko.mapping.fromJS(dataObj);
                koResult = ko.observable(koObj);
            }
            else {
                var koObj = ko.mapping.fromJS(dataObj);
                observableObject(koObj);
            }
            if (typeof fnOk == "function") fnOk(dataObj);
        },
        error: function (error) {
            alert("There was an error posting the data to the server: " + error.responseText);
        }
    });

    return koResult;
}

function assignObsArray(obsArrayObj, data, simpleAssign, addEvent, itemModel) {
    obsArrayObj.removeAll();
    var useSimpleAssign = arguments.length > 2 ? simpleAssign : false;
    if (typeof data != "undefined" && data.length > 0) {
        if (useSimpleAssign) {
            data.forEach(function (x) {
                if (typeof addEvent == "function") addEvent(x);
                obsArrayObj.push(x);
            });
        }
        else {
            data.forEach(function (x) {
                if (itemModel == undefined) {
                    x = ko.mapping.fromJS(x);
                }
                else {
                    x = ko.mapping.fromJS(x, itemModel);
                }
                if (typeof addEvent == "function") addEvent(x);
                obsArrayObj.push(x);
            });
        }
    }
}

function push2ObsArray(obsArrayObj, data) {
    obsArrayObj.removeAll();
    data.forEach(function (x) {
        obsArrayObj.push(x);
    });
}

function ajaxPost(url, data, fnOk, fnNok) {
    $.ajax({
        url: url,
        type: 'POST',
        data: ko.mapping.toJSON(data),
        //data: data, 
        contentType: "application/json; charset=utf-8",
        success: function (data) {
            if (typeof fnOk == "function") fnOk(data);
            koResult = "OK";
        },
        error: function (error) {
            if (typeof fnNok == "function") {
                fnNok(error);
            }
            else {
                alert("There was an error posting the data to the server: " + error.responseText);
            }
        }
    });
}

function coalesce(nullcheckobj, defaultvalue) {
    return nullcheckobj == null || nullcheckobj == undefined ?
        defaultvalue : nullcheckobj;
}

function gridResize(gridObject, newSize) {
    var gridElement = gridObject;
    var dataArea = gridElement.find(".k-grid-content");
    var newHeight = newSize;
    var diff = gridElement.innerHeight() - dataArea.innerHeight();
    gridElement.height(newHeight);
    dataArea.height(newHeight - diff);
}

function componentToHex(c) {
    var hex = c.toString(16);
    return hex.length == 1 ? "0" + hex : hex;
}

function rgbToHex(r, g, b) {
    return "#" + componentToHex(r) + componentToHex(g) + componentToHex(b);
}