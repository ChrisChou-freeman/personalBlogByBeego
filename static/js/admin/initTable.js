
(function () {
    var requestUrl = null;
    function bindChangePager() {
        $('#idPagination').on('click','a',function () {
            var num = $(this).text();
            init(num);
        })
    }

    function bindSave() {
        $('#idSave').click(function () {
            var hasEditorMod = $("#idEditMode").hasClass('btn-warning')
            if (hasEditorMod){
                alert("请先退出编辑模式");
            }else{
                var postList = [];
                //找到已经编辑过的tr，tr has-edit='true'
                $('#table_tb').find('tr[has-edit="true"]').each(function () {
                    // $(this) => tr
                    var temp = {};
                    var id = $(this).attr('row-id');
                    $(this).children('[edit-enable="true"]').each(function () {
                        // $(this) = > td
                        var name = $(this).attr('name');
                        var origin = $(this).attr('origin');
                        var newVal = $(this).attr('new-val');
                        if (origin != newVal){
                            temp['Id'] = Number(id);
                            temp[name] = newVal;
                        }
                    });
                    if (JSON.stringify(temp) != "{}"){
                        postList.push(temp);
                    }
                })
                console.log(postList)
                $.ajax({
                    url:requestUrl,
                    type: 'POST',
                    data: {'post_list': JSON.stringify(postList)},
                    dataType: 'JSON',
                    success:function (arg) {
                        arg  = JSON.parse(arg)
                        console.log(arg)
                        if(arg.Status){
                            init(1);
                        }else{
                            alert(arg.Message);
                        }
                    }
                })
            }
            
        })
    }

    function bindReverseAll() {
        $('#idReverseAll').click(function () {
            $('#table_tb').find(':checkbox').each(function () {
                // $(this) => checkbox
                if($('#idEditMode').hasClass('btn-warning')) {
                    if($(this).prop('checked')){
                        $(this).prop('checked',false);
                        trOutEditMode($(this).parent().parent());
                    }else{
                        $(this).prop('checked',true);
                        trIntoEditMode($(this).parent().parent());
                    }
                }else{
                    if($(this).prop('checked')){
                        $(this).prop('checked',false);
                    }else{
                        $(this).prop('checked',true);
                    }
                }
            })
        })
    }

    function bindCancelAll() {
        $('#idCancelAll').click(function () {
            $('#table_tb').find(':checked').each(function () {
                // $(this) => checkbox
                if($('#idEditMode').hasClass('btn-warning')){
                    $(this).prop('checked',false);
                    // 退出编辑模式
                    trOutEditMode($(this).parent().parent());
                }else{
                    $(this).prop('checked',false);
                }
            });
        })
    }

    function bindCheckAll() {
        $('#idCheckAll').click(function () {
            $('#table_tb').find(':checkbox').each(function () {
                // $(this)  = checkbox
                if($('#idEditMode').hasClass('btn-warning')){
                    if($(this).prop('checked')){
                        // 当前行已经进入编辑模式了
                    }else{
                        // 进入编辑模式
                        var $currentTr = $(this).parent().parent();
                        trIntoEditMode($currentTr);
                        $(this).prop('checked',true);
                    }
                }else{
                    $(this).prop('checked',true);
                }
            })
        })
    }

    function bindEditMode() {
        $('#idEditMode').click(function () {
            var editing = $(this).hasClass('btn-warning');
            if(editing){
                // 退出编辑模式
                $(this).removeClass('btn-warning');
                $(this).text('进入编辑模式');

                $('#table_tb').find(':checked').each(function () {
                    var $currentTr = $(this).parent().parent();
                    trOutEditMode($currentTr);
                })

            }else{
                // 进入编辑模式
                $(this).addClass('btn-warning');
                $(this).text('退出编辑模式');
                $('#table_tb').find(':checked').each(function () {
                    var $currentTr = $(this).parent().parent();
                    trIntoEditMode($currentTr);
                })
            }
        })
    }

    function bindCheckbox() {
        // $('#table_tb').find(':checkbox').click()
        $('#table_tb').on('click',':checkbox',function () {

            if($('#idEditMode').hasClass('btn-warning')){
                var ck = $(this).prop('checked');
                var $currentTr = $(this).parent().parent();
                if(ck){
                    // 进入编辑模式
                    trIntoEditMode($currentTr);
                }else{
                    // 退出编辑模式
                    trOutEditMode($currentTr)
                }
            }
        })
    }

    function trIntoEditMode($tr) {
        $tr.addClass('success');
        $tr.attr('has-edit','true');
        $tr.children().each(function () {
            // $(this) => td
            var editEnable = $(this).attr('edit-enable');
            var editType = $(this).attr('edit-type');
            if(editEnable=='true'){
                if(editType == 'select'){
                    var globalName = $(this).attr('global-name'); //  "device_status_choices"
                    var origin = $(this).attr('origin'); // 1
                    // 生成select标签
                    var sel = document.createElement('select');
                    sel.className = "form-control";
                    $.each(window[globalName],function(k1,v1){
                        var op = document.createElement('option');
                        op.setAttribute('value',v1[0]);
                        op.innerHTML = v1[1];
                        $(sel).append(op);
                    });
                    $(sel).val(origin);

                    $(this).html(sel);
                    
                }else if(editType == 'input'){
                    // input文本框
                    // *******可以进入编辑模式*******
                    var innerText = $(this).text();
                    var tag = document.createElement('input');
                    tag.className = "form-control";
                    tag.value = innerText;
                    $(this).html(tag);
                }
            }
        })
    }

    function trOutEditMode($tr){
        $tr.removeClass('success');
        $tr.children().each(function () {
            // $(this) => td
            var editEnable = $(this).attr('edit-enable');
            var editType = $(this).attr('edit-type');
            if(editEnable=='true'){
                if (editType == 'select'){
                    // 获取正在编辑的select对象
                    var $select = $(this).children().first();
                    // 获取选中的option的value
                    var newId = $select.val();
                    // 获取选中的option的文本内容
                    var newText = $select[0].selectedOptions[0].innerHTML;
                    // 在td中设置文本内容
                    $(this).html(newText);
                    $(this).attr('new-val',newId);

                }else if(editType == 'input') {
                    // *******可以退出编辑模式*******
                    var $input = $(this).children().first();
                    var inputValue = $input.val();
                    $(this).html(inputValue);
                    $(this).attr('new-val',inputValue);
                }

            }
        })
    }

    String.prototype.format = function (kwargs) {
        var ret = this.replace(/\{(\w+)\}/g,function (km,m) {
            return kwargs[m];
        });
        return ret;
    };

    function init(pager) {
        $.ajax({
            url: requestUrl,
            type: 'GET',
            data: {'pager':pager},
            dataType: 'JSON',
            success:function (result) {
                result = JSON.parse(result)
                initGlobalData(result.GlobalDict);
                initHeader(result.TableConfig);
                initBody(result.TableConfig,result.DataList);
                initPager(result.Pager);
            }
        })

    }
    function initPager(pager){
        $('#idPagination').html(pager);
        $('#idPagination').find("a").prop("href", "#")
    }

    function initHeader(table_config) {
        var tr = document.createElement('tr');
        $.each(table_config,function (k,item) {
            if(item.Display){
                var th = document.createElement('th');
                th.innerHTML = item.Title;
                $(tr).append(th);
            }

        });
        $('#table_th').empty();
        $('#table_th').append(tr);
    }

    function initBody(table_config, data_list){
        $('#table_tb').empty();
        for(var i=0;i<data_list.length;i++){
            var row = data_list[i];
            // row = {'cabinet_num': '12B', 'cabinet_order': '1', 'id': 1},
            var tr = document.createElement('tr');
            tr.setAttribute('row-id',row['Id']);
            $.each(table_config,function (i,colConfig) {
               if(colConfig.Display){
                   var td = document.createElement('td');
                   /* 生成文本信息 */
                   var kwargs = {};
                   text_kwargs = JSON.parse(colConfig.Text.kwargs)
                   $.each(text_kwargs,function (key,value) {

                       if(value.substring(0,2) == '@@'){
                           var globalName = value.substring(2,value.length); // 全局变量的名称
                           var currentId = row[colConfig.Q]; // 获取的数据库中存储的数字类型值
                           var t = getTextFromGlobalById(globalName,currentId);
                           kwargs[key] = t;
                       }
                       else if (value[0] == '@'){
                            kwargs[key] = row[value.substring(1,value.length)]; //cabinet_num
                       }else{
                            kwargs[key] = value;
                       }
                   });
                   var temp = colConfig.Text.content.format(kwargs);
                   td.innerHTML = temp;


                   /* 属性colConfig.attrs = {'edit-enable': 'true','edit-type': 'select'}  */
                    $.each(colConfig.Attrs,function (kk,vv) {
                        if(vv[0] == '@'){
                            td.setAttribute(kk,row[vv.substring(1,vv.length)]);
                        }else{
                            td.setAttribute(kk,vv);
                        }
                    });
                   $(tr).append(td);
               }
            });

            $('#table_tb').append(tr);
        }
    }

    function initGlobalData(global_dict) {
        $.each(global_dict,function (k,v) {
            // k = "device_type_choices"
            // v= [[0,'xx'],[1,'xxx']]
            // device_type_choices = 123;
            window[k] = v;
        })
    }

    function getTextFromGlobalById(globalName,currentId) {
        // globalName = "device_type_choices"
        // currentId = 1
        var ret = null;
        $.each(window[globalName],function (k,item) {
            if(item[0] == currentId){
                ret = item[1];
                return
            }
        });
        return ret;
    }


    jQuery.extend({
        'initTable': function (url) {
            requestUrl = url;
            init();
            bindEditMode();
            bindCheckbox();
            bindCheckAll();
            bindCancelAll();
            bindReverseAll();
            bindSave();
            bindChangePager();
        },
        'changePager': function (num) {
            init(num);
        }
    })
})();


