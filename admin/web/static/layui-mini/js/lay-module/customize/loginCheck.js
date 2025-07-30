layui.define(["jquery", "layer"], function (exports) {
    let $ = layui.$;
    let layer = layui.layer;

    // 从localStorage中获取accessToken
    let accessToken = localStorage.getItem('accessToken');
    let userId = localStorage.getItem('userId');

    // 全局校验登录态
    $.ajaxSetup({
        headers: {"Authorization": "Bearer " + accessToken},
        // timeout: 5000,
        // dataType: "json",
        complete: function (xhr) {
            if (xhr.status !== 401 && xhr.status !== 200) {
                console.log('xhr', xhr)
                layer.msg(xhr.statusText, {icon: 2});
            }
            if (xhr.status === 401 || xhr.getResponseHeader("sessionstatus") === "timeout") {
                top.location.replace('/login');
            }
        }
    });

    // 校验登录
    function checkLoginStatus() {
        $.ajax({
            url: "/admin/check-login-status",
            type: "POST",
            // success: function(response) {
            //     console.log(response);
            // },
            error: function (xhr, status, error) {
                if (xhr.status === 401) {
                    window.location.replace('/login');
                }
            }
        });
    }

    checkLoginStatus();

    let loginCheck = {
        getUserId: function () {
            return parseInt(userId);
        },

        checkLoginStatus: function () {
            checkLoginStatus();
        }
    }

    exports('loginCheck', loginCheck);
});