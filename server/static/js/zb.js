var phone_num = "";
var name = "";
var from_class = "";
var to_campus = "";
var to_period = "";
var to_periods = {};

function getClassesByStudentNameAndPhone() {
  $.ajax({
        url: '/api/get-classes-by-name-and-phone-num/' + name + '/' + phone_num,
        type: 'GET',
        dataType: 'json',
        error : function () {
            alert("获取班级失败");
        },
        success: function (data) {
            var classesObj = $('#classes');
            classesObj.html('');
            classesObj.selectmenu("refresh");

            if (!data.success) {
                alert("获取时段信息失败:" + data.err_msg);
            } else {
                $.each(data.classes, function(index, value){
                    classesObj.append('<option value="' + value + '">' + value + '</option>');
                });
                if (data.classes.length > 0) {
                    classesObj.val(data.classes[0]).change();
                }
            }
        }
  });
}

function getClassesByStudentNameAndPhone() {

  $.ajax({
        url: '/api/get-classes-by-name-and-phone-num/' + name + '/' + phone_num,
        type: 'GET',
        dataType: 'json',
        error : function () {
            alert("获取班级失败");
        },
        success: function (data) {
            var classesObj = $('#classes');
            classesObj.html('');
            classesObj.selectmenu("refresh");

            if (!data.success) {
                alert("获取班级失败:" + data.err_msg);
            } else {
                $.each(data.classes, function(index, value){
                    classesObj.append('<option value="' + value + '">' + value + '</option>');
                });
                if (data.classes.length > 0) {
                    classesObj.val(data.classes[0]).change();
                }
            }
        }
  });
}

function getAvailablePeriodsByClass() {
  $.ajax({
        url: '/api/get-available-periods/' + from_class,
        type: 'GET',
        dataType: 'json',
        error : function () {
            alert("获取时段失败");
        },
        success: function (data) {
            var toCampusesObj = $('#to_campuses');
            toCampusesObj.html('');
            toCampusesObj.selectmenu("refresh");

            if (!data.success) {
                alert("获取息失败:" + data.err_msg);
            } else {
                var firstCampus = "";
                to_periods = {};

                $.each(data.campus_periods, function(campus, periods){
                    toCampusesObj.append('<option value="' + campus + '">' + campus + '</option>');
                    if (firstCampus == "") {
                        firstCampus = campus;
                    }
                    to_periods[campus] = periods; 
                });
                if (firstCampus != "") {
                    toCampusesObj.val(firstCampus).change();
                }
            }
        }
  });
}


$(document).ready(function () {
    //alert("document ready.");

    // Page 1 events.
    $('#phone_num').on("input", function () {
        phone_num = $(this).val();
        if (phone_num != "") {
            $('#queryBtn').removeClass("ui-state-disabled");
        } else {
            $('#queryBtn').addClass("ui-state-disabled");
        }
    });

    $('#students').change(function () {
        name = $(this).find('option:selected').val();
        getClassesByStudentNameAndPhone();       
    });
    
    $('#classes').change(function () {
        from_class = $(this).find('option:selected').val();
        getAvailablePeriodsByClass();
    });

    $('#to_campuses').change(function () {
        var toPeriodsObj = $('#to_periods');
        toPeriodsObj.html('');
        toPeriodsObj.selectmenu("refresh");

        to_campus = $(this).find('option:selected').val();
        $.each(to_periods[to_campus], function(index, period) {
            toPeriodsObj.append('<option value="' + period + '">' + period + '</option>');        
        });        
        toPeriodsObj.val(to_periods[to_campus][0]).change();
        to_period = to_periods[to_campus][0];
    });

    $('#to_periods').change(function () {
        to_period = $(this).find('option:selected').val();
    }); 

    $('#queryBtn').click(function () {
        phone_num = $('#phone_num').val();

        $.ajax({
            type: "GET",
            url: "/api/get-names-by-phone-num/" + phone_num,
            error: function (XMLHttpRequest, textStatus, errorThrown) {
		    console.log("/api/get-names-by-phone-num/" + phone_num + " failed");
            },
            success: function (data) {
                var studentsObj = $('#students');
                studentsObj.html('');
		studentsObj.selectmenu("refresh");

		if (data.success) {
			console.log(data.names)
	                $.each(data.names, function(index, name) {
                          studentsObj.append('<option value="' + name + '">' + name + '</option>');
			});
                        if (data.names.length > 0) {
		            studentsObj.val(data.names[0]).change();
                        }
		} else {
       			console.log(data.err_msg);
			alert("查询失败: " + data.err_msg + "\n" + "请更换其它可能预留的联系电话\n" + "如果依旧不能找到绑定的学生，请至报名处修改联系电话.");
		}

            },
            dataType: "json"
        });
    });

    $('#submitBtn').click(function () {
	if ((name == "") || (phone_num == "") || (from_class == "") || (to_campus == "") || (to_period == "")) {
		alert("信息缺失，请确认后再提交");
	        return;
	}

        postData = {name: name, phone_num: phone_num, from_class: from_class, to_campus: to_campus, to_period: to_period};
        console.log(postData);

        $.ajax({
            type: "POST",
            url: "/api/request",
            data: JSON.stringify(postData),
            error: function (XMLHttpRequest, textStatus, errorThrown) {
                    console.log("/api/request" + " failed");
            },
            success: function (data) {
                if (data.success) {
                        console.log(data.request)
                        alert("提交成功\n" + "姓名: " + name + "," + "电话: " + phone_num + "\n" + "当前班级: " + from_class + "\n" + "转入校区: " + to_campus + "," + "时间段: " + to_period + "\n" + "具体结果学校老师会以电话告知，请耐心等待，谢谢.");
                } else {
                        alert("提交失败: " + data.err_msg);
                }
            },
            dataType: "json"
        });

    });

});

$(document).on("pageinit","#page1",function(){
});

$(document).on("pagebeforeshow","#page1",function(){
});
