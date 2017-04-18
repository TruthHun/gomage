$(function(){
   $(".gomage-tooltips").tooltip();
	$(".gomage-ajax-get").click(function(e){
		e.preventDefault();
		var _this=$(this);
		if(_this.hasClass("disabled")) return false;
		if (_this.hasClass("gomage-confirm")) {
			var lang=$("html").attr("lang"),confirm_text="Are you sure "+_this.text()+"？";
			if (lang=="zh-CN") {
				confirm_text="您确定要"+_this.text()+"吗？"
			}
			
			if(confirm(confirm_text)){
				
				_this.addClass("disabled");
				$.get(_this.attr("href"),function(rt){
					if (rt.status==1) {
						gomage_alert("succ",rt.msg,3000,location.href);
					} else{
						_this.removeClass("disabled")
						gomage_alert("error",rt.msg,3000,"");
					}
				});
			}
		} else{
			$.get(_this.attr("href"),function(rt){
				if (rt.status==1) {
					gomage_alert("succ",rt.msg,3000,location.href);
				} else{
					_this.removeClass("disabled")
					gomage_alert("error",rt.msg,3000,"");
				}
			});
		}
	});
	


	$("#ModalFriend form input[type=reset]").click(function(){
		$("#ModalFriend").modal("hide");
	});
	
	
	$(".gomage-ajax-form [type=submit]").click(function(e){
		e.preventDefault();
		var form=$(this).parents("form"),method=form.attr("method"),action=form.attr("action"),data=form.serialize(),_url=form.attr("data-url");
		var require=form.find("[required=required]"),l=require.length;
		$.each(require, function() {    
			if (!$(this).val()){
				$(this).focus();
				return false;
			}else{
				l--
			}
		});
        if (!_url || _url==undefined){
            _url=location.href;
        }
		if (l>0) return false;
		if (method=="post") {
			$.post(action,data,function(rt){
				if (rt.status==1) {

                    gomage_alert("success",rt.msg,2000,_url);
				} else{
                    gomage_alert("error",rt.msg,3000,"");
				}
				});
		} else{
			$.get(action,data,function(rt){
				if (rt.status==1) {
                    gomage_alert("success",rt.msg,2000,_url);
				} else{
                    gomage_alert("error",rt.msg,3000,"");
				}
			});
		}
	});
	
	
	
	//内容变更
	$(".gomage-change-update").change(function(){
		var _this=$(this),_url=_this.attr("data-url"),field=_this.attr("name"),value=_this.val();
		$.get(_url,{field:field,value:value},function(rt){
			if (rt.status==1) {
                gomage_alert("success",rt.msg,2000,"")
			} else{
                gomage_alert("error",rt.msg,2000,"")
			}
		});
	});
	
	//site option ,change reload a new link
	$("#SiteOption").change(function(){
		var _this=$(this),sid=_this.val(),_url=_this.attr("data-url");
		location.href=_url+"&sid="+sid
	});

	//创建样式
	$("#gomage-add-style [type=range]").change(function () {
		$("#gomage-quality").val($(this).val());
    });
	//水印位置的选择
	$("#gomage-add-style .gomage-watermark-position span").click(function () {
		$(this).addClass("active").siblings("span").removeClass("active");
		$("#gomage-watermark-position").val($(this).text());
    });
	$("#gomage-add-style select[name=Method]").change(function () {
		var val=$(this).val();
        $("#gomage-add-style input[name=Width]").attr("readonly",null);
        $("#gomage-add-style input[name=Height]").attr("readonly",null);
		if (val==21){
			$("#gomage-add-style input[name=Height]").val("0").attr("readonly","readonly");
		}
        if (val==22){
            $("#gomage-add-style input[name=Width]").val("0").attr("readonly","readonly");
        }
    });


	//原图
	$(".gomage-preview-original").click(function (e) {
		e.preventDefault();
		var _this=$(this),src=_this.attr("href")+"?t="+new Date();
		$(".gomage-preview-box img").attr("src",src);
    });
    //刷新预览
    ///admin/preview/?width=825&height=316&method=11&waterpath=&waterposition=9&zoom=true&ext=jpeg
    $(".gomage-preview-refresh").click(function (e) {
        e.preventDefault();
        var _this=$(this),src=_this.attr("href"),form=$("form"),
            width=form.find("[name=Width]").val(),
            height=form.find("[name=Height]").val(),
            method=form.find("[name=Method]").val(),
            waterpath=form.find("[name=Watermark]").val(),
            waterposition=form.find("[name=WatermarkPosition]").val(),
            zoom=form.find("[name=IsZoom]").val(),
            ext=form.find("[name=Ext]").val();
            top=form.find("[name=Top]").val();
            left=form.find("[name=Left]").val();
            right=form.find("[name=Right]").val();
            bottom=form.find("[name=Bottom]").val();
        src=src+"?width="+width+"&height="+height+"&method="+method+"&waterpath="+waterpath+"&waterposition="+waterposition+"&zoom="+zoom+"&ext="+ext+"&t="+new Date();
        if (_this.hasClass("gomage-preview-window")){
            window.open(src);
        }else{
            $(".gomage-preview-box img").attr("src",src);
        }
    });

	
	//cls：success/error
	//msg:message
	//timeout:超时刷新和跳转时间
	//url:有url链接的话，跳转url链接
	function gomage_alert(cls,msg,timeout,url){
		if(timeout>0){
			t=timeout
		}else{
			t=3000
		}
		if(cls=="error"||cls=="danger"){
			cls="danger";
		}else{
			cls="success";
		}
		html='<div class="alert alert-'+cls+' alert-dismissable gomage-alert"><button type="button" class="close" data-dismiss="alert" aria-hidden="true">&times;</button>'+msg+'</div>';
		$("body").append(html);
		$(".gomage-alert").fadeIn();
		setTimeout(function(){
			$(".gomage-alert").fadeOut();
			$(".gomage-alert").remove();
		},t);
		if(url){
			setTimeout(function(){
				location.href=url
			},t-500);
		}
	}
	
	//图片样式导入
	$(".btn-import").click(function(){
		$("#gomage-import").trigger("click");
	});
	//判断样式格式是否是json
	$("#gomage-import").change(function(){
		var _this=$(this),file=_this.val();
		if (file) {
			$("#gomage-form-import").submit();
		}
	});
	//iframe加载后处理
	$("#target").load(function(){
		var data = $(window.frames['target'].document.body).find("pre").html();
		var obj=eval('(' + data + ')');
		if (obj.status==1) {
			gomage_alert("success",obj.msg,2500,location.href);
		} else{
			gomage_alert("danger",obj.msg,2500,"");
		}
	});


    $(".gomage-tooltip").tooltip();

});



