/**
 * 现代化Cron表达式生成器
 * 基于Layui的现代化界面设计
 * @author: Your Name
 * @license: MIT
 */

layui.define(['jquery', 'layer', 'form', 'element'], function (exports) {
    "use strict";

    var $ = layui.jquery;
    var layer = layui.layer;
    var form = layui.form;
    var element = layui.element;
    var MOD_NAME = 'cron';
    var MOD_INDEX = 'layui-' + MOD_NAME + '-index';
    var THIS = 'layui-this';

    layui.link(layui.cache.base + 'cron/cron.css');

    // 默认配置
    var defaults = {
        value: function () {
            return '* * * * * ?';
        },
        theme: '#1E9FFF',
        lang: 'cn',
        btns: [], // 移除确定按钮
        position: [],
        zIndex: null,
        elem: null, // 绑定的触发元素
        done: null,
        ready: null,
        change: null
    };

    // 多语言支持
    var locales = {
        cn: {
            seconds: '秒',
            minutes: '分',
            hours: '时',
            days: '日',
            months: '月',
            weeks: '周',
            years: '年',
            per: '间隔',
            assign: '指定',
            run: '运行',
            confirm: '确定',
            cancel: '取消',
            recentTimes: '最近运行时间',
            placeholder: '请选择或输入表达式',
            dayNotSet: '不指定',
            lastDay: '本月最后一天',
            dayWork: '工作日',
            lastWorkDay: '本月最后一个工作日',
            lastWeek: '本月最后一个星期',
            specifyWeek: '第#周的星期',
            Sun: '日',
            Mon: '一',
            Tue: '二',
            Wed: '三',
            Thu: '四',
            Fri: '五',
            Sat: '六'
        },
        en: {
            seconds: 'Seconds',
            minutes: 'Minutes',
            hours: 'Hours',
            days: 'Days',
            months: 'Months',
            weeks: 'Weeks',
            years: 'Years',
            per: 'Per',
            assign: 'Assign',
            run: 'Run',
            confirm: 'Confirm',
            cancel: 'Cancel',
            recentTimes: 'Recent Run Times',
            placeholder: 'Please select or enter expression',
            dayNotSet: 'Not Specified',
            lastDay: 'Last Day of Month',
            dayWork: 'Weekday',
            lastWorkDay: 'Last Weekday of Month',
            lastWeek: 'Last Week of Month',
            specifyWeek: 'Week # of Month',
            Sun: 'Sun',
            Mon: 'Mon',
            Tue: 'Tue',
            Wed: 'Wed',
            Thu: 'Thu',
            Fri: 'Fri',
            Sat: 'Sat'
        }
    };

    // 构造函数
    var Class = function (options) {
        var that = this;
        that.index = ++modernCron.index;
        if (options.elem) {
            options.elem = $(options.elem)
        }
        that.config = $.extend({}, defaults, options);
        that.render();
    };

    // 公共初始化逻辑
    Class.prototype._initCommon = function () {
        var that = this;
        var options = that.config;
        var lang = locales[options.lang] || locales['cn'];

        // 创建容器
        var cronContainer = $(
            '<div class="layui-modern-cron">' +
            '  <div class="layui-tab layui-tab-brief">' +
            '    <ul class="layui-tab-title">' +
            '      <li class="layui-this" data-type="second">' + lang.seconds + '</li>' +
            '      <li data-type="minute">' + lang.minutes + '</li>' +
            '      <li data-type="hour">' + lang.hours + '</li>' +
            '      <li data-type="day">' + lang.days + '</li>' +
            '      <li data-type="month">' + lang.months + '</li>' +
            '      <li data-type="week">' + lang.weeks + '</li>' +
            '      <li data-type="year">' + lang.years + '</li>' +
            '    </ul>' +
            '    <div class="layui-tab-content" style="padding-bottom: 5px">' +
            '      <div class="layui-tab-item layui-show"></div>' +
            '      <div class="layui-tab-item"></div>' +
            '      <div class="layui-tab-item"></div>' +
            '      <div class="layui-tab-item"></div>' +
            '      <div class="layui-tab-item"></div>' +
            '      <div class="layui-tab-item"></div>' +
            '      <div class="layui-tab-item"></div>' +
            '    </div>' +
            '    <div class="modern-cron-title">' + lang.recentTimes + '</div>' +
            '    <div class="modern-cron-recent-times"></div>' +
            '    <input type="text" class="layui-input modern-cron-input" placeholder="' + lang.placeholder + '">' +
            '  </div>' +
            // '  <div class="modern-cron-box">' +
            // '    <div class="modern-cron-expression">' +
            // '      <input type="text" class="layui-input modern-cron-input" placeholder="' + lang.placeholder + '">' +
            // '    </div>' +
            // '  </div>' +
            // '  <div class="modern-cron-box">' +
            // '    <div class="modern-cron-title">' + lang.recentTimes + '</div>' +
            // '    <div class="modern-cron-recent-times"></div>' +
            // '  </div>' +
            // '  <div class="modern-cron-footer-btns"></div>' +
            '</div>'
        );

        // 添加按钮
        var btns = [];
        if (options.btns.indexOf('run') !== -1) {
            btns.push('<span class="modern-cron-btns-run">' + lang.run + '</span>');
        }
        if (options.btns.indexOf('confirm') !== -1) {
            btns.push('<span class="modern-cron-btns-confirm">' + lang.confirm + '</span>');
        }

        cronContainer.find('.modern-cron-footer-btns').html(btns.join(''));

        // 设置样式
        cronContainer.attr('id', 'layui-modern-cron-' + that.index);
        cronContainer.attr(MOD_INDEX, that.index);

        // 添加到页面
        $('body').append(cronContainer);

        that.cronContainer = cronContainer;
        that.elemHeader = cronContainer.find('.layui-tab-title');
        that.elemCont = cronContainer.find('.layui-tab-content');
        that.elemInput = cronContainer.find('.modern-cron-input');
        that.elemRecent = cronContainer.find('.modern-cron-recent-times');

        // 为tab内容添加data-type属性，方便后续查找
        that.elemCont.find('.layui-tab-item').each(function (index) {
            var types = ['second', 'minute', 'hour', 'day', 'month', 'week', 'year'];
            $(this).attr('data-type', types[index]);
        });

        return lang;
    };

    // 核心方法
    Class.prototype.render = function () {
        var that = this;
        var options = that.config;
        var lang = that._initCommon();

        // 添加到页面并初始隐藏
        that.cronContainer.css('display', 'none');

        // 绑定触发元素点击事件
        if (options.elem) {
            // 先解绑之前的点击事件，避免重复绑定
            $(options.elem).off('click');
            $(options.elem).on('click', function () {
                // 确保控件容器存在且未被移除
                let expression = options.value();
                if (that.cronContainer.closest('body').length > 0) {
                    that.cronContainer.toggle();
                    that.position(options.elem);

                    // 表达式回显到各tab
                    that.setValue(expression);
                } else {
                    // 如果控件已被移除，重新初始化
                    that.init();
                    that.cronContainer.show();
                    that.position(options.elem);

                    // 表达式回显到各tab
                    that.setValue(expression);
                }
            });
        }

        // 点击外部关闭控件
        $(document).on('click', function (event) {
            var $target = $(event.target);
            // 点击目标不在控件内且不在触发元素内时关闭
            if (!that.cronContainer.is($target) && that.cronContainer.has($target).length === 0
                && !(options.elem && $(options.elem).is($target)) && !(options.elem && $(options.elem).has($target).length > 0)) {
                that.cronContainer.hide();
            }
        });

        // 定位
        that.position(options.elem);

        // 初始化事件
        that.events();

        // 初始化时加载所有tab内容
        that.loadAllTabContent();

        // 回调
        typeof options.ready === 'function' && options.ready(that);

        return that;
    };

    // 加载所有tab内容
    Class.prototype.loadAllTabContent = function () {
        let tabNames = ['second', 'minute', 'hour', 'day', 'month', 'week', 'year']
        tabNames.forEach((name, index) => this.loadTabContent(name, index));
    };

    // 重新初始化方法
    Class.prototype.init = function () {
        var that = this;
        var options = that.config;
        var lang = that._initCommon();

        // 重新绑定事件
        that.events();

        // 初始化时加载所有tab内容
        that.loadAllTabContent();
    };

    // 设置值
    // 事件绑定
    Class.prototype.events = function () {
        var that = this;
        var options = that.config;

        // tab切换
        that.elemHeader.on('click', 'li', function () {
            var type = $(this).data('type');
            var index = $(this).index();

            $(this).addClass(THIS).siblings().removeClass(THIS);
            that.elemCont.find('.layui-tab-item').eq(index).addClass('layui-show').siblings().removeClass('layui-show');

            // 检查tab内容是否已经加载过，如果没有则加载
            var $tabItem = that.elemCont.find('.layui-tab-item').eq(index);
            if ($tabItem.children().length === 0) {
                // 加载对应内容
                that.loadTabContent(type, index);
            }

            // 添加tab切换时自动计算最近运行时间的功能
            // 生成表达式
            var expression = that.generateExpression();

            // 验证表达式格式是否正确
            if (that.isValidExpression(expression)) {
                // 显示最近运行时间
                that.displayRecentTimes(expression);

                // 更新表达式输入框
                that.elemInput.val(expression);

                // 自动触发确定逻辑
                var value = that.getValue();
                typeof options.done === 'function' && options.done(value, that);
            } else {
                // 表达式不符合规范，显示报错信息
                that.elemRecent.html('<div class="time-error">表达式格式不正确</div>');
            }
        });

        // 按钮事件
        that.cronContainer.find('.modern-cron-footer-btns').on('click', 'span', function () {
            var type = $(this).attr('class').indexOf('run') !== -1 ? 'run' : 'confirm';

            if (type === 'run') {
                that.run();
            } else if (type === 'confirm') {
                that.confirm();
            }
        });

        // 输入框事件
        that.elemInput.on('input', function () {
            var value = $(this).val();

            // 验证表达式格式是否正确
            if (that.isValidExpression(value)) {
                // 解析表达式并更新各tab
                that.parseExpression(value);
                typeof options.change === 'function' && options.change(value, that);

                // 触发输入框的change事件，以便外部监听
                that.elemInput.trigger('change');
            }
        });

        // 单选框变化事件，处理复选框的启用/禁用
        that.cronContainer.on('click', 'input[type=radio]', function () {
            var radioName = $(this).attr('name');
            var radioValue = $(this).val();
            var type = radioName.replace('_type', '');

            // 如果是assign类型的单选框，则启用对应的复选框
            if (radioValue === 'assign') {
                that.cronContainer.find('input[name=' + type + '_assign]').prop('disabled', false);
            } else {
                // 否则禁用对应的复选框
                that.cronContainer.find('input[name=' + type + '_assign]').prop('disabled', true);
            }

            // 生成表达式
            var expression = that.generateExpression();

            // 验证表达式格式是否正确
            if (that.isValidExpression(expression)) {
                // 显示最近运行时间
                that.displayRecentTimes(expression);

                // 更新表达式输入框
                that.elemInput.val(expression);

                // 自动触发确定逻辑
                var value = that.getValue();
                typeof options.done === 'function' && options.done(value, that);
            } else {
                // 表达式不符合规范，显示报错信息
                that.elemRecent.html('<div class="time-error">表达式格式不正确</div>');
            }

            // 触发输入框的change事件，以便外部监听
            that.elemInput.trigger('change');
        });

        // 复选框变化事件
        that.cronContainer.on('change', 'input[type=checkbox]', function () {
            // 更新表达式输入框
            that.updateExpressionInput();

            // 触发输入框的change事件，以便外部监听
            that.elemInput.trigger('change');

            // 添加选项变化时自动计算最近运行时间的功能
            // 生成表达式
            var expression = that.generateExpression();

            // 验证表达式格式是否正确
            if (that.isValidExpression(expression)) {
                // 显示最近运行时间
                that.displayRecentTimes(expression);

                // 更新表达式输入框
                that.elemInput.val(expression);

                // 自动触发确定逻辑
                var value = that.getValue();
                typeof options.done === 'function' && options.done(value, that);
            } else {
                // 表达式不符合规范，显示报错信息
                that.elemRecent.html('<div class="time-error">表达式格式不正确</div>');
            }
        });

        // 数字输入框变化事件
        that.cronContainer.on('input', 'input[type=number]', function () {
            // 更新表达式输入框
            that.updateExpressionInput();

            // 触发输入框的change事件，以便外部监听
            that.elemInput.trigger('change');
        });

        // 数字输入框失去焦点事件
        that.cronContainer.on('blur', 'input[type=number]', function () {
            // 获取输入框所属的tab类型
            var inputName = $(this).attr('name');
            var type = '';
            if (inputName) {
                // 从inputName中提取type，例如second_range_start -> second
                type = inputName.split('_')[0];
            }

            // 检查对应的单选框是否被选中
            var isChecked = false;
            if (type) {
                // 获取该tab下所有单选框
                var $radios = that.cronContainer.find('input[name=' + type + '_type]');
                // 检查是否有单选框被选中
                isChecked = $radios.is(':checked');
            }

            // 只有当对应的单选框被选中时，才执行后续逻辑
            if (isChecked) {
                // 生成表达式
                var expression = that.generateExpression();

                // 验证表达式格式是否正确
                if (that.isValidExpression(expression)) {
                    // 显示最近运行时间
                    that.displayRecentTimes(expression);

                    // 自动触发确定逻辑
                    var value = that.getValue();
                    typeof options.done === 'function' && options.done(value, that);
                }
            }
        });

        // 下拉框变化事件
        that.cronContainer.on('change', 'select', function () {
            // 更新表达式输入框
            that.updateExpressionInput();

            // 触发输入框的change事件，以便外部监听
            that.elemInput.trigger('change');
        });

        // 下拉框选中事件
        that.cronContainer.on('change', 'select', function () {
            // 获取下拉框所属的tab类型
            var selectName = $(this).attr('name');
            var type = '';
            if (selectName) {
                // 从selectName中提取type，例如second_range_unit -> second
                type = selectName.split('_')[0];
            }

            // 检查对应的单选框是否被选中
            var isChecked = false;
            if (type) {
                // 获取该tab下所有单选框
                var $radios = that.cronContainer.find('input[name=' + type + '_type]');
                // 检查是否有单选框被选中
                isChecked = $radios.is(':checked');
            }

            // 只有当对应的单选框被选中时，才执行后续逻辑
            if (isChecked) {
                // 生成表达式
                var expression = that.generateExpression();

                // 验证表达式格式是否正确
                if (that.isValidExpression(expression)) {
                    // 显示最近运行时间
                    that.displayRecentTimes(expression);

                    // 自动触发确定逻辑
                    var value = that.getValue();
                    typeof options.done === 'function' && options.done(value, that);
                }
            }
        });
    };

    // 定位算法
    Class.prototype.position = function (triggerElement) {
        var that = this,
            options = that.config,
            cronWidth = that.cronContainer.outerWidth(), // 控件的宽度
            cronHeight = that.cronContainer.outerHeight(), // 控件的高度
            // 滚动条高度
            scrollArea = function (type) {
                type = type ? "scrollLeft" : "scrollTop";
                return document.body[type] | document.documentElement[type];
            },
            winArea = function (type) {
                return document.documentElement[type ? "clientWidth" : "clientHeight"];
            },
            left = 0,
            top = 0,
            spacing = 2; // 间距2px

        // 如果提供了触发元素，则基于该元素定位
        const triggerRect = triggerElement[0].getBoundingClientRect();

        // 计算触发元素四周的可用位置（下/上/左/右），优先级按此顺序
        // 当控件出现在触发元素左边的时候，控件的右边距离触发元素的左边2px
        // 当控件出现在触发元素右边的时候，控件的左边距离触发元素的右边2px
        const positions = [
            {type: 'bottom', top: triggerRect.bottom + spacing, left: triggerRect.left},
            {type: 'top', top: triggerRect.top - cronHeight - spacing, left: triggerRect.left},
            {type: 'right', top: triggerRect.top, left: triggerRect.right + spacing},
            {type: 'left', top: triggerRect.top, left: triggerRect.left - cronWidth - spacing}
        ];

        // 优先选择下方位置，如果空间不足再按顺序选择其他位置
        for (const pos of positions) {
            if (pos.left >= 0 && pos.top >= 0 &&
                pos.left + cronWidth <= winArea("width") &&
                pos.top + cronHeight <= winArea()) {
                left = pos.left;
                top = pos.top;
                break;
            }
        }

        // 最终调整（防止完全超出视口）
        // 特殊处理左侧定位：确保控件在触发元素左侧时保持正确距离
        if (left === 0 && triggerRect.left - cronWidth - spacing > 0) {
            left = triggerRect.left - cronWidth - spacing;
        } else {
            // 原有的边界保护逻辑
            left = Math.max(0, Math.min(left, winArea("width") - cronWidth));
        }
        top = Math.max(0, Math.min(top, winArea() - cronHeight));

        that.cronContainer.css({
            left: left + (options.position === "fixed" ? 0 : scrollArea(1)),
            top: top + (options.position === "fixed" ? 0 : scrollArea())
        });
    };

    // 加载tab内容
    Class.prototype.loadTabContent = function (type, index) {
        var that = this;
        var content = that.getTabContent(type);
        that.elemCont.find('.layui-tab-item').eq(index).html(content);

        // 重新渲染表单
        form.render();
    };

    // 获取tab内容
    Class.prototype.getTabContent = function (type) {
        var lang = locales[this.config.lang] || locales.cn;
        var html = '';

        // 秒、分、时的tab内容
        if (['second', 'minute', 'hour'].includes(type)) {
            var maxVal = type === 'hour' ? 23 : 59;
            var desc = type === 'second' ? '秒' : type === 'minute' ? '分钟' : '小时'
            html = `
                <div class="modern-cron-form">
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="all" checked> 每${desc} 允许的通配符[, - * /]
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="range">周期 从
                        <input type="number" class="modern-cron-input" name="${type}_range_start" value="1" min="0" max="${maxVal}"> - 
                        <input type="number" class="modern-cron-input" name="${type}_range_end" value="2" min="0" max="${maxVal}">
                        ${desc}
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="per" title="${lang.per}">从
                        <input type="number" class="modern-cron-input" name="${type}_per_start" value="0" min="0" max="${maxVal}">${desc}开始,每 
                        <input type="number" class="modern-cron-input" name="${type}_per_interval" value="1" min="1" max="${maxVal}">${desc}执行一次
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="assign" title="${lang.assign}">指定
                        <div class="modern-cron-assign">
                            ${this.generateCheckboxOptions(type, maxVal)}
                        </div>
                    </div>
                </div>
            `;
        } else if (type === 'day') {
            html = `
                <div class="modern-cron-form">
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="all" title="${lang.all}" checked>每天 允许的通配符[, - * / L W]
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="no" title="${lang.dayNotSet}">不指定
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="range" title="${lang.range}">周期 从
                        <input type="number" class="modern-cron-input" name="${type}_range_start" value="1" min="1" max="31"> - 
                        <input type="number" class="modern-cron-input" name="${type}_range_end" value="2" min="1" max="31">日
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="per" title="${lang.per}">从
                        <input type="number" class="modern-cron-input" name="${type}_per_start" value="1" min="1" max="31">日开始,每
                        <input type="number" class="modern-cron-input" name="${type}_per_interval" value="1" min="1" max="31">天执行一次
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="work" title="${lang.dayWork}">每月
                        <input type="number" class="modern-cron-input" name="${type}_work_day" value="1" min="1" max="31">号最近的那个工作日
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="last" title="${lang.lastDay}">本月最后一天
                    </div>

                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="assign" title="${lang.assign}">指定
                        <div class="modern-cron-assign">
                            ${this.generateCheckboxOptions(type, 31)}
                        </div>
                    </div>
                </div>
            `;
        } else if (type === 'month') {
            html = `
                <div class="modern-cron-form">
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="all" title="${lang.all}" checked>每月 允许的通配符[, - * /]
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="no">不指定
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="range" title="${lang.range}">周期 从
                        <input type="number" class="modern-cron-input" name="${type}_range_start" value="1" min="1" max="12"> - 
                        <input type="number" class="modern-cron-input" name="${type}_range_end" value="2" min="1" max="12">月
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="per" title="${lang.per}">从
                        <input type="number" class="modern-cron-input" name="${type}_per_start" value="1" min="1" max="12">月开始,每
                        <input type="number" class="modern-cron-input" name="${type}_per_interval" value="1" min="1" max="12">月执行一次
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="assign" title="${lang.assign}">指定
                        <div class="modern-cron-assign">
                            ${this.generateCheckboxOptions(type, 12)}
                        </div>
                    </div>
                </div>
            `;
        } else if (type === 'week') {
            html = `
                <div class="modern-cron-form">
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="all" title="${lang.all}">每周 允许的通配符[, - * / L #]
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="no" checked>不指定
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="range" title="${lang.range}">周期 从星期
                        <input type="number" class="modern-cron-input" name="${type}_range_start" value="1" min="1" max="7"> - 
                        <input type="number" class="modern-cron-input" name="${type}_range_end" value="2" min="1" max="7">
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="specifyWeek" title="${lang.specifyWeek}">第
                        <input type="number" class="modern-cron-input" name="${type}_specify_week_num" value="1" min="1" max="4">周的星期
                        <select name="${type}_specify_week" class="modern-cron-input">
                            <option value="1">${lang.Sun}</option>
                            <option value="2">${lang.Mon}</option>
                            <option value="3">${lang.Tue}</option>
                            <option value="4">${lang.Wed}</option>
                            <option value="5">${lang.Thu}</option>
                            <option value="6">${lang.Fri}</option>
                            <option value="7">${lang.Sat}</option>
                        </select>
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="lastWeek" title="${lang.lastWeek}">本月最后一个星期
                        <select name="${type}_last_week" class="modern-cron-input">
                            <option value="1">${lang.Sun}</option>
                            <option value="2">${lang.Mon}</option>
                            <option value="3">${lang.Tue}</option>
                            <option value="4">${lang.Wed}</option>
                            <option value="5">${lang.Thu}</option>
                            <option value="6">${lang.Fri}</option>
                            <option value="7">${lang.Sat}</option>
                        </select>
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="assign" title="${lang.assign}">指定
                        <div class="modern-cron-assign">
                            ${this.generateCheckboxOptions(type, 7)}
                        </div>
                    </div>
                </div>
            `;
        } else if (type === 'year') {
            const currentYear = new Date().getFullYear();
            html = `
                <div class="modern-cron-form">
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="no" checked>不指定 允许的通配符[, - * /] 非必填
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="all">每年
                    </div>
                    <div style="margin-left: 5px;">
                        <input type="radio" name="${type}_type" value="range" title="${lang.range}">周期从
                        <input type="number" class="modern-cron-input" name="${type}_range_start" value="${currentYear}" min="2000" max="2100"> - 
                        <input type="number" class="modern-cron-input" name="${type}_range_end" value="${currentYear + 1}" min="2000" max="2100">年
                    </div>
                </div>
            `;
        }

        return html;
    };

    // 生成复选框选项
    Class.prototype.generateCheckboxOptions = function (type, maxVal) {
        var html = '';
        var max = type === 'year' ? Math.min(maxVal, 2100) : maxVal;

        // 定义每行显示的复选框数量
        var itemsPerLine = 10; // 默认每行10个
        if (type === 'hour' || type === 'month') {
            itemsPerLine = 6; // 小时和月每行6个
        } else if (type === 'week') {
            itemsPerLine = 7; // 周每行7个
        }

        // 对于年份，我们生成一个合理的范围
        if (type !== 'year') {
            // 对于其他类型，生成从1到max的选项
            for (var i = ['second', 'minute', 'hour'].includes(type) ? 0 : 1; i <= max; i++) {
                // 个位数前补0
                var displayValue = i < 10 ? '0' + i : i;
                // 默认选中第一个复选框
                html += `<input type="checkbox" name="${type}_assign" value="${i}" style="margin-left: 5px" disabled>${displayValue}`;

                // 每行结束后添加换行符
                var num = ['second', 'minute', 'hour'].includes(type) ? i + 1 : i
                if (num > 0 && num % itemsPerLine === 0) {
                    html += '<br>';
                }
            }
        }

        return html;
    }

    // 运行
    Class.prototype.run = function () {
        var that = this;
        var value = that.getValue();

        // 解析表达式并生成最近运行时间
        var times = that.generateRecentTimes(value);

        // 显示最近运行时间
        var html = '';
        for (var i = 0; i < times.length; i++) {
            html += `<div class="time-item"><span class="time-icon">⏰</span><span class="time-value">${times[i]}</span></div>`;
        }

        that.elemRecent.html(html);
    };

    // 生成最近运行时间示例
    Class.prototype.generateRecentTimes = function (expression) {
        var that = this;
        var times = [];
        var now = new Date();

        // 解析表达式
        var parts = expression.split(' ');
        if (parts.length < 6) return times;

        var second = parts[0];
        var minute = parts[1];
        var hour = parts[2];
        var day = parts[3];
        var month = parts[4];
        var week = parts[5];
        var year = parts[6] || '';

        // 解析各个部分的值
        var parsePart = function (part, min, max) {
            var result = [];

            // 处理通配符
            if (part === '*') {
                for (var i = min; i <= max; i++) {
                    result.push(i);
                }
                return result;
            }

            // 处理范围
            if (part.indexOf('-') > -1 && part.indexOf('/') === -1) {
                var range = part.split('-');
                var start = parseInt(range[0]);
                var end = parseInt(range[1]);
                for (var i = start; i <= end; i++) {
                    result.push(i);
                }
                return result;
            }

            // 处理间隔
            if (part.indexOf('/') > -1) {
                var intervalParts = part.split('/');
                var start = intervalParts[0] === '*' ? min : parseInt(intervalParts[0]);
                var step = parseInt(intervalParts[1]);
                for (var i = start; i <= max; i += step) {
                    result.push(i);
                }
                return result;
            }

            // 处理具体值
            if (part.indexOf(',') > -1) {
                var values = part.split(',');
                for (var i = 0; i < values.length; i++) {
                    result.push(parseInt(values[i]));
                }
            } else {
                result.push(parseInt(part));
            }

            return result;
        };

        // 解析各个部分
        var seconds = parsePart(second, 0, 59);
        var minutes = parsePart(minute, 0, 59);
        var hours = parsePart(hour, 0, 23);
        var days = parsePart(day, 1, 31);
        var months = parsePart(month, 1, 12);

        // 生成最近5次运行时间
        var currentDate = new Date(now);
        for (var i = 0; i < 5; i++) {
            // 查找下一个运行时间
            var found = false;
            var attempts = 0;

            while (!found && attempts < 1000) { // 限制尝试次数以避免无限循环
                attempts++;

                // 检查当前时间是否匹配
                var currentSecond = currentDate.getSeconds();
                var currentMinute = currentDate.getMinutes();
                var currentHour = currentDate.getHours();
                var currentDay = currentDate.getDate();
                var currentMonth = currentDate.getMonth() + 1;

                if (seconds.indexOf(currentSecond) !== -1 &&
                    minutes.indexOf(currentMinute) !== -1 &&
                    hours.indexOf(currentHour) !== -1 &&
                    days.indexOf(currentDay) !== -1 &&
                    months.indexOf(currentMonth) !== -1) {
                    times.push(new Date(currentDate));
                    found = true;
                }

                // 增加一秒继续查找
                currentDate.setSeconds(currentDate.getSeconds() + 1);
            }

            // 如果找到了时间，继续查找下一个
            if (found) {
                currentDate.setSeconds(currentDate.getSeconds() + 1);
            }
        }

        // 格式化时间
        var formattedTimes = [];
        for (var i = 0; i < times.length; i++) {
            formattedTimes.push(times[i].toLocaleString());
        }

        return formattedTimes;
    };

    // 显示最近运行时间
    Class.prototype.displayRecentTimes = function (expression) {
        var that = this;
        var times = that.generateRecentTimes(expression);
        var html = '';
        for (var i = 0; i < times.length; i++) {
            html += `<div class="time-item"><span class="time-icon">⏰</span><span class="time-value">${times[i]}</span></div>`;
        }
        that.elemRecent.html(html);
    };

    // 解析表达式
    Class.prototype.parseExpression = function (value) {
        var that = this;
        var parts = value.split(' ');
        if (parts.length < 6) return;

        // 设置各个tab的值
        that.setTabValues(parts);

        // 显示最近运行时间
        that.displayRecentTimes(value);

        // 更新表达式输入框的值，但不触发input事件
        that.elemInput.val(value);
    };

    // 验证表达式格式是否正确
    Class.prototype.isValidExpression = function (value) {
        var that = this;

        // 检查表达式是否为空
        if (!value) {
            return false;
        }

        // 分割表达式的各个部分
        var parts = value.trim().split(' ');

        // 检查是否包含6或7个部分
        if (parts.length < 6 || parts.length > 7) {
            return false;
        }

        // 定义各部分的验证规则
        // 秒: 0-59
        if (!that.isValidCronField(parts[0], 0, 59)) return false;
        // 分: 0-59
        if (!that.isValidCronField(parts[1], 0, 59)) return false;
        // 小时: 0-23
        if (!that.isValidCronField(parts[2], 0, 23)) return false;
        // 日: 1-31
        if (!that.isValidCronField(parts[3], 1, 31)) return false;
        // 月: 1-12
        if (!that.isValidCronField(parts[4], 1, 12)) return false;
        // 周: 0-7 (0和7都表示周日)
        if (!that.isValidCronField(parts[5], 0, 7)) return false;
        // 年: 1970-2099 (可选)
        if (parts[6] && !that.isValidCronField(parts[6], 1970, 2099)) return false;

        return true;
    };

    // 验证cron表达式的单个字段
    Class.prototype.isValidCronField = function (field, min, max) {
        // 处理特殊字符: * , - / ? L W #
        // * 表示所有值
        if (field === '*') return true;

        // ? 只能出现在日和周字段
        if (field === '?') return true;

        // 处理逗号分隔的列表
        if (field.includes(',')) {
            var values = field.split(',');
            for (var i = 0; i < values.length; i++) {
                if (!this.isValidCronValue(values[i], min, max)) return false;
            }
            return true;
        }

        // 处理范围
        if (field.includes('-')) {
            var rangeParts = field.split('-');
            if (rangeParts.length !== 2) return false;

            var start = parseInt(rangeParts[0], 10);
            var end = parseInt(rangeParts[1], 10);

            // 检查是否为有效数字
            if (isNaN(start) || isNaN(end)) return false;

            // 检查范围是否在允许的范围内
            if (start < min || end > max || start > end) return false;

            return true;
        }

        // 处理步进值 (如 */5 或 10/2)
        if (field.includes('/')) {
            var stepParts = field.split('/');
            if (stepParts.length !== 2) return false;

            var step = parseInt(stepParts[1], 10);

            // 检查步进值是否为有效数字
            if (isNaN(step) || step <= 0) return false;

            // 如果是 */n 格式
            if (stepParts[0] === '*') return true;

            // 如果是 m/n 格式，检查m是否有效
            var baseValue = parseInt(stepParts[0], 10);
            if (isNaN(baseValue) || baseValue < min) return false;

            return true;
        }

        // 处理单个值
        return this.isValidCronValue(field, min, max);
    };

    // 验证cron表达式的单个值
    Class.prototype.isValidCronValue = function (value, min, max) {
        // 处理 L, W, # 等特殊字符
        if (value === 'L' || value === 'W' || value.includes('#')) return true;

        // 检查是否为有效数字
        var num = parseInt(value, 10);
        if (isNaN(num)) return false;

        // 检查是否在允许范围内
        if (num < min || num > max) return false;

        return true;
    };

    // 更新表达式输入框
    Class.prototype.updateExpressionInput = function () {
        var that = this;
        var expression = that.generateExpression();
        that.elemInput.val(expression);

        // 显示最近运行时间
        that.displayRecentTimes(expression);

        // 触发change事件
        typeof that.config.change === 'function' && that.config.change(expression, that);
    };

    // 生成表达式
    Class.prototype.generateExpression = function () {
        var that = this;
        var values = [];

        // 获取各个tab的值
        values.push(that.getTabValue('second'));
        values.push(that.getTabValue('minute'));
        values.push(that.getTabValue('hour'));
        values.push(that.getTabValue('day'));
        values.push(that.getTabValue('month'));
        values.push(that.getTabValue('week'));

        // 如果有年份
        var yearValue = that.getTabValue('year');
        if (yearValue && yearValue !== '?') {
            values.push(yearValue);
        }

        var expression = values.join(' ');

        return expression;
    };

    // 设置tab的值
    Class.prototype.setTabValues = function (parts) {
        var that = this;

        // 设置秒
        that.setSingleTabValue('second', parts[0]);

        // 设置分
        that.setSingleTabValue('minute', parts[1]);

        // 设置时
        that.setSingleTabValue('hour', parts[2]);

        // 设置日
        that.setSingleTabValue('day', parts[3]);

        // 设置月
        that.setSingleTabValue('month', parts[4]);

        // 设置周
        that.setSingleTabValue('week', parts[5]);

        // 设置年（如果有）
        if (parts[6]) {
            that.setSingleTabValue('year', parts[6]);
        }
    };

    // 设置单个tab的值
    Class.prototype.setSingleTabValue = function (type, value) {
        var that = this;
        var $tabContent = that.elemCont.find('.layui-tab-item[data-type="' + type + '"]');

        // 如果找不到对应tab内容，直接返回
        if ($tabContent.length === 0) {
            // 尝试通过索引来查找
            var tabIndex = that.getTabIndexByType(type);
            if (tabIndex !== -1) {
                $tabContent = that.elemCont.find('.layui-tab-item').eq(tabIndex);
            } else {
                return;
            }
        }

        // 根据值的类型设置选项
        if (value === '*') {
            // 每个都选择
            $tabContent.find('input[name="' + type + '_type"][value="all"]').prop('checked', true);
        } else if (value.indexOf('-') > -1) {
            // 范围选择
            var range = value.split('-');
            $tabContent.find('input[name="' + type + '_type"][value="range"]').prop('checked', true);
            $tabContent.find('input[name="' + type + '_range_start"]').val(range[0]);
            $tabContent.find('input[name="' + type + '_range_end"]').val(range[1]);
        } else if (value.indexOf('/') > -1) {
            // 间隔选择
            var per = value.split('/');
            $tabContent.find('input[name="' + type + '_type"][value="per"]').prop('checked', true);
            $tabContent.find('input[name="' + type + '_per_start"]').val(per[0]);
            $tabContent.find('input[name="' + type + '_per_interval"]').val(per[1]);
        } else if (value.indexOf(',') > -1 || /^\d+$/.test(value)) {
            // 指定选择
            $tabContent.find('input[name="' + type + '_type"][value="assign"]').prop('checked', true);
            var values = value.split(',');

            // 启用复选框
            $tabContent.find('input[name="' + type + '_assign"]').prop('disabled', false);

            // 选中指定的值
            for (var i = 0; i < values.length; i++) {
                $tabContent.find('input[name="' + type + '_assign"][value="' + values[i] + '"]').prop('checked', true);
            }
        }

        // 重新渲染表单
        form.render();
    };

    // 根据tab类型获取索引
    Class.prototype.getTabIndexByType = function (type) {
        var typeMap = {
            'second': 0,
            'minute': 1,
            'hour': 2,
            'day': 3,
            'month': 4,
            'week': 5,
            'year': 6
        };

        return typeMap[type] !== undefined ? typeMap[type] : -1;
    };

    // 获取tab的值
    Class.prototype.getTabValue = function (type) {
        var that = this;
        var $tabContent = that.elemCont.find('.layui-tab-item[data-type="' + type + '"]');

        // 如果找不到对应tab内容，直接返回
        if ($tabContent.length === 0) {
            // 尝试通过索引来查找
            var tabIndex = that.getTabIndexByType(type);
            if (tabIndex !== -1) {
                $tabContent = that.elemCont.find('.layui-tab-item').eq(tabIndex);
            } else {
                return '*';
            }
        }

        // 获取选中的类型
        var selectedType = $tabContent.find('input[name="' + type + '_type"]:checked').val();

        switch (selectedType) {
            case 'all':
                return '*';
            case 'range':
                var start = $tabContent.find('input[name="' + type + '_range_start"]').val();
                var end = $tabContent.find('input[name="' + type + '_range_end"]').val();
                return start + '-' + end;
            case 'per':
                var start = $tabContent.find('input[name="' + type + '_per_start"]').val();
                var interval = $tabContent.find('input[name="' + type + '_per_interval"]').val();
                return start + '/' + interval;
            case 'assign':
                var values = [];
                $tabContent.find('input[name="' + type + '_assign"]:checked').each(function () {
                    values.push(this.value);
                });
                return values.length > 0 ? values.join(',') : '*';
            case 'no':
                return '?';
            case 'last':
                return 'L';
            case 'work':
                var day = $tabContent.find('input[name="' + type + '_work_day"]').val();
                return day + 'W';
            case 'specifyWeek':
                var week = $tabContent.find('select[name="' + type + '_specify_week"]').val();
                var num = $tabContent.find('input[name="' + type + '_specify_week_num"]').val();
                return num + '#' + week;
            case 'lastWeek':
                var week = $tabContent.find('select[name="' + type + '_last_week"]').val();
                return week + 'L';
            default:
                return '*';
        }
    };

    // 确认
    Class.prototype.confirm = function () {
        var that = this;
        var options = that.config;
        var value = that.getValue();

        // 显示最近运行时间
        that.displayRecentTimes(value);

        typeof options.done === 'function' && options.done(value, that);

        // 隐藏控件而不是移除
        that.cronContainer.hide();
    };

    // 获取值
    Class.prototype.getValue = function () {
        var that = this;
        // 根据各tab的选择结果生成表达式
        return that.generateExpression();
    };

    // 设置值
    Class.prototype.setValue = function (value) {
        var that = this;
        // 解析表达式并设置各tab的选项
        that.parseExpression(value);

        // 更新输入框的值
        that.elemInput.val(value);

        // 显示最近运行时间
        that.displayRecentTimes(value);

        return that;
    };

    // 获取tab索引
    Class.prototype.getTabIndexByType = function (type) {
        var typeMap = {
            'second': 0,
            'minute': 1,
            'hour': 2,
            'day': 3,
            'month': 4,
            'week': 5,
            'year': 6
        };

        return typeMap[type] !== undefined ? typeMap[type] : -1;
    };

    // 根据tab类型获取jQuery对象
    Class.prototype.getTabContentByType = function (type) {
        var that = this;
        var $tabContent = that.elemCont.find('.layui-tab-item[data-type="' + type + '"]');

        // 如果找不到对应tab内容，尝试通过索引来查找
        if ($tabContent.length === 0) {
            var tabIndex = that.getTabIndexByType(type);
            if (tabIndex !== -1) {
                $tabContent = that.elemCont.find('.layui-tab-item').eq(tabIndex);
            }
        }

        return $tabContent;
    };

    // 设置单个tab的值
    Class.prototype.setSingleTabValue = function (type, value) {
        var that = this;
        var $tabContent = that.getTabContentByType(type);

        // 如果找不到对应tab内容，直接返回
        if ($tabContent.length === 0) {
            return;
        }

        // 根据值的类型设置选项
        if (value === '*') {
            // 每个都选择
            $tabContent.find('input[name="' + type + '_type"][value="all"]').prop('checked', true);
        } else if (value === '?') {
            // 不指定
            $tabContent.find('input[name="' + type + '_type"][value="no"]').prop('checked', true);
        } else if (value === 'L') {
            // 最后一天/最后一周
            $tabContent.find('input[name="' + type + '_type"][value="last"]').prop('checked', true);

        } else if (value.indexOf('-') > -1) {
            // 范围选择
            var range = value.split('-');
            $tabContent.find('input[name="' + type + '_type"][value="range"]').prop('checked', true);
            $tabContent.find('input[name="' + type + '_range_start"]').val(range[0]);
            $tabContent.find('input[name="' + type + '_range_end"]').val(range[1]);
        } else if (value.indexOf('/') > -1) {
            // 间隔选择
            var per = value.split('/');
            $tabContent.find('input[name="' + type + '_type"][value="per"]').prop('checked', true);
            $tabContent.find('input[name="' + type + '_per_start"]').val(per[0]);
            $tabContent.find('input[name="' + type + '_per_interval"]').val(per[1]);
        } else if (value.indexOf('W') > -1) {
            // 工作日
            var day = value.replace('W', '');
            $tabContent.find('input[name="' + type + '_type"][value="work"]').prop('checked', true);
            $tabContent.find('input[name="' + type + '_work_day"]').val(day);
        } else if (value.indexOf('#') > -1) {
            // 第#周的星期
            var parts = value.split('#');
            var num = parts[0];
            var week = parts[1];
            $tabContent.find('input[name="' + type + '_type"][value="specifyWeek"]').prop('checked', true);
            $tabContent.find('input[name="' + type + '_specify_week_num"]').val(num);
            $tabContent.find('select[name="' + type + '_specify_week"]').val(week);
        } else if (value.indexOf(',') > -1 || /^\d+$/.test(value)) {
            // 指定选择
            $tabContent.find('input[name="' + type + '_type"][value="assign"]').prop('checked', true);
            var values = value.split(',');

            // 启用复选框
            $tabContent.find('input[name="' + type + '_assign"]').prop('disabled', false);

            // 选中指定的值
            for (var i = 0; i < values.length; i++) {
                $tabContent.find('input[name="' + type + '_assign"][value="' + values[i] + '"]').prop('checked', true);
            }
        }

        // 重新渲染表单
        form.render();
    };

    // 移除
    Class.prototype.remove = function () {
        var that = this;
        that.cronContainer.remove();
        delete modernCron.cache[that.index];
    };

    // modernCron对象定义
    var modernCron = {
        index: layui[MOD_NAME] ? (layui[MOD_NAME].index + 10000) : 0,
        cache: {},

        // 设置全局配置
        set: function (options) {
            var that = this;
            that.config = $.extend({}, that.config, options);
            return that;
        },

        // 主入口
        render: function (options) {
            var that = this;

            // 创建实例
            var inst = new Class(options);
            return inst;
        }
    };

    // 保存到layui对象中
    exports(MOD_NAME, modernCron);
})