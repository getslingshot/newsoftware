(function() {
    var CONDITIONS = [
        "contains",
        "range",
        "equal",
        "greater_than",
        "less_than",
    ];
    var weekDays = [
        'Monday',
        'Tuesday',
        'Wednesday',
        'Thursday',
        'Friday',
        'Saturday',
        'Sunday',
    ];
    var emptyHours = {
        opens: 0,
        closes: 0,
        openHour: 0,
        openMinute: 0,
        openAmPm: 'AM',
        closeHour: 0,
        closeMinute: 0,
        closeAmPm: 'AM',
    };
    var emptyException = {
        '_id': null,
        day: 0,
        logic: [{
            field: 'personal_info',
            condition: 'contains',
            value: ''
        }],
        then: 'schedule',
        scheduleDays: [], //Day numbers
        hours: [$.extend({}, emptyHours)],
        note: ''
    };
    var callOnDays = {};
    weekDays.forEach(function(day) {
        var onDay = {
            forDay: {},
            nextWeek: {
                forDay: {}
            },
        };
        weekDays.forEach(function(forDay) {
            onDay.forDay[forDay] = {
                hours: [$.extend({}, emptyHours)],
                note: '',
                '_id': null,
                checked: false,
                week: 1,
            };
            onDay.nextWeek.forDay[forDay] = {
                hours: [$.extend({}, emptyHours)],
                note: '',
                '_id': null,
                checked: false,
                week: 2
            };
            var exp = $.extend(true, {}, emptyException);
            exp.day = days[day];
            onDay.exceptions = [exp];
        });
        callOnDays[day] = onDay;
    });
    mergeCalls(CALLS_DATA);
    mergeConditions(CONDITIONS_DATA);
    var hoursEntryComponent = Vue.component('hours-notes', {
        template: "#hours-notes-template",
        delimiters: ['{', '}'],
        props: ['hoursEntry'],
        methods: {
            padZero: function(optionStr) {
                return optionStr <= 9 ? '0' + optionStr : optionStr;
            },
            convertOpen: function(hours) {
                hours.opens = convert(hours.openHour, hours.openMinute, hours.openAmPm);
            },
            convertClose: function(hours) {
                hours.closes = convert(hours.closeHour, hours.closeMinute, hours.closeAmPm);;
            },
            addHours: function(hours) {
                hours.push($.extend(true, {}, emptyHours));
            },
            removeHours: function(hours, index) {
                hours.splice(index, 1);
            },
            updateNote: function(noteObj, event) {
                noteObj.note = event.currentTarget.value;
            },
        }
    });
    var ah = new Vue({
        delimiters: ['{', '}'],
        el: "#afterHours",
        data: {
            callOnDays: callOnDays,
            cachedCallOnDays: $.extend(true, {}, callOnDays),
            conditions: CONDITIONS,
            weekDays: weekDays,
        },
        methods: {
            cacheCallOnDays: function() {
                this.cachedCallOnDays = $.extend(true, {}, this.callOnDays);
            },
            revertCallOnDays: function() {
                this.callOnDays = $.extend(true, {}, this.cachedCallOnDays);
            },
            submitCall: function(onDayLbl, forDayLbl, call) {
                call.onDay = days[onDayLbl];
                call.forDay = days[forDayLbl];
                this.cacheCallOnDays();
                if (!call['_id']) {
                    return createCall(call).then(function(data) {
                        call['_id'] = data['_id'];
                    });
                }
                updateCall(call);
            },
            submitNextWeek: function(callDayLbl, nextWeek) {
                var checked = [];
                var unchecked = [];
                for (day in nextWeek) {
                    if (nextWeek.hasOwnProperty(day)) {
                        var call = nextWeek[day];
                        if (call.checked) {
                            var call = nextWeek[day];
                            this.submitCall(callDayLbl, day, call);
                        } else if (call['_id']) {
                            deleteCall(call);
                            call.hours = [$.extend({}, emptyHours)];
                            call.checked = false;
                            call['_id'] = null;
                        }
                    }
                }
            },
            toggleCall: function(call, event) {
                if (!event.target.checked && call['_id']) {
                    if (confirm('Are you sure you want to delete this item?')) {
                        deleteCall(call);
                        call.hours = [$.extend({}, emptyHours)];
                        call.checked = false;
                        call['_id'] = null;
                    } else {
                        event.target.checked = !event.target.checked;
                        call.checked = true;
                        setTimeout(function() {
                            $(event.target).closest('.row').next('.collapse').collapse('show');
                        }, 500);
                    }
                }
                call.checked = event.target.checked;
                this.cacheCallOnDays();
            },
            toggleExceptionAndIf: function(logic, i, event) {
                if (event.target.checked) {
                    logic.push($.extend(true, {}, emptyException.logic[0]));
                } else {
                    logic.splice(i + 1, 1);
                }
            },
            fmtCondition: function(c) {
                var words = c.split('_');
                return words.map(function(w) {
                    return w.charAt(0).toUpperCase() + w.slice(1);
                }).join(' ');
            },
            castScheduleDay: function(dArr, i) {
                if (dArr[i]) {
                    dArr[i] = Number(dArr[i]);
                }
            },
            submitException: function(exceptions, exception) {
                this.cacheCallOnDays();
                if (!exception['_id']) {
                    var exp = $.extend(true, {}, emptyException);
                    exceptions.unshift(exp);
                    return createCondition(exception).then(function(data) {
                        exception['_id'] = data['_id'];
                    });
                }
                updateCondition(exception);
            },
            removeCondition: function(exceptions, i) {
                if (confirm('Are you sure you want to delete this item?')) {
                    deleteCondition(exceptions[i]);
                    exceptions.splice(i, 1);
                }
            },
        },
        filters: {
            uppercase: function(str) {
                return str.toUpperCase();
            }
        }
    });

    function createCall(call) {
        return $.ajax({
            method: 'POST',
            url: '/create-call',
            data: JSON.stringify(call),
            success: function(data) {
                $.notify({
                    message: 'Call saved'
                }, {
                    type: 'success'
                });
            },
            error: displayErrors
        });
    }

    function getCalls() {
        return $.ajax({
            method: 'GET',
            url: '/get-calls',
            success: function(data) {
                return mergeCalls(data);
            },
            error: displayErrors
        });
    }

    function updateCall(call) {
        return $.ajax({
            method: 'PUT',
            url: '/update-calls?id=' + call['_id'],
            data: JSON.stringify(call),
            success: function(data) {
                $.notify({
                    message: 'Call Updated'
                }, {
                    type: 'success'
                });
            },
            error: displayErrors
        });
    }

    function deleteCall(call) {
        return $.ajax({
            method: 'DELETE',
            url: '/delete-calls?id=' + call['_id'],
            data: JSON.stringify(call),
            success: function(data) {
                $.notify({
                    message: 'Call Deleted'
                }, {
                    type: 'success'
                });
            },
            error: displayErrors
        });
    }

    function mergeCalls(calls) {
        var daysArr = Object.keys(days);
        var weekOneCalls = calls.filter(function(call) {
            return call.week === 1;
        });
        var weekTwoCalls = calls.filter(function(call) {
            return call.week === 2;
        });
        weekOneCalls.forEach(function(callEntry) {
            var entry = callOnDays[daysArr[callEntry.onDay]].forDay[daysArr[callEntry.forDay]];
            entry.note = callEntry.note;
            entry.hours = callEntry.hours;
            entry['_id'] = callEntry['_id'];
            entry.checked = true;
        });
        weekTwoCalls.forEach(function(callEntry) {
            var entry = callOnDays[daysArr[callEntry.onDay]].nextWeek.forDay[daysArr[callEntry.forDay]];
            callOnDays[daysArr[callEntry.onDay]].nextWeek.checked = true;
            entry.note = callEntry.note;
            entry.hours = callEntry.hours;
            entry['_id'] = callEntry['_id'];
            entry.checked = true;
        });
    }

    function mergeConditions(conditions) {
        var daysArr = Object.keys(days);
        conditions.forEach(function(condition) {
            callOnDays[daysArr[condition.day]].exceptions.push(condition);
        });
    }

    function createCondition(condition) {
        return $.ajax({
            method: 'POST',
            url: '/create-condition',
            data: JSON.stringify(condition),
            success: function(data) {
                condition = data;
                $.notify({
                    message: 'Condition saved'
                }, {
                    type: 'success'
                });
            },
            error: displayErrors
        });
    }

    function getConditions() {
        return $.ajax({
            method: 'GET',
            url: '/get-Conditions',
            success: function(data) {
                return mergeConditions(data);
            },
            error: displayErrors
        });
    }

    function updateCondition(condition) {
        return $.ajax({
            method: 'PUT',
            url: '/update-conditions?id=' + condition['_id'],
            data: JSON.stringify(condition),
            success: function(data) {
                $.notify({
                    message: 'Condition Updated'
                }, {
                    type: 'success'
                });
            },
            error: displayErrors
        });
    }

    function deleteCondition(condition) {
        return $.ajax({
            method: 'DELETE',
            url: '/delete-condition?id=' + condition['_id'],
            data: JSON.stringify(condition),
            success: function(data) {
                $.notify({
                    message: 'Condition Deleted'
                }, {
                    type: 'success'
                });
            },
            error: displayErrors
        });
    }
})()
