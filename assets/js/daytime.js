'use strict';

$.ajaxSetup({
    contentType: 'application/json; charset=utf-8',
    xhrFields: {
        withCredentials: true
    },
    dataType: 'json',
    beforeSend(xhr) {
        xhr.withCredentials = true;
    }
});

//TODO: Figure out a better way to clear out validation errors
$(document).ajaxSuccess(function() {
    $('#validation-errors').hide();
});

//Cache unmodified section
var origSchedule;

var scheduleSubmitBtn = $('#submit-daytime');

function convert(hour, minute, ampm) {
    var add = 0;
    if (ampm === 'PM' && hour !== '12') {
        add = 720; // 12 hours * 60 minutes, only for PM values
    }
    return (parseInt(hour, 10) * 60) + parseInt(minute, 10) + add;
}

var days = {
    'Sunday': 0,
    'Monday': 1,
    'Tuesday': 2,
    'Wednesday': 3,
    'Thursday': 4,
    'Friday': 5,
    'Saturday': 6,
};

function displayErrors(error) {
    if (!error || !error.responseJSON || !error.responseJSON.errors) {
        console.log('on error', error);
        $.notify({
            message: 'Cant process request!'
        }, {
            type: 'danger'
        });
        return;
    }
    var errors = error.responseJSON.errors;
    if (!Object.keys(errors).length) {
        return;
    }
    var liTags = Object.keys(errors).map(function(key) {
        return '<li>' + key + ': ' + errors[key] + '</li>'
    });
    $('#validation-errors').show(function() {
        return liTags.length > 0;
    }).find('ul').html(liTags);
}

function parseScheduleRow(scope) {
    var row = $(scope).closest('tr');
    var day = row.find('.weekday-selected').text();
    var entryCount = row.find('.time-entry').size() / 2;
    var entries = [];
    for (var i = 1; i <= entryCount; i++) {
        entries.push({
            day: day,
            daytime: true,
            openHour: row.find('.time-entry:nth-child(' + i + ') .open-hour-selected').val(),
            openMinute: row.find('.time-entry:nth-child(' + i + ') .open-minute-selected').val(),
            openAmpm: row.find('.time-entry:nth-child(' + i + ') .open-ampm-selected').val(),
            closeHour: row.find('.time-entry:nth-child(' + i + ') .close-hour-selected').val(),
            closeMinute: row.find('.time-entry:nth-child(' + i + ') .close-minute-selected').val(),
            closeAmpm: row.find('.time-entry:nth-child(' + i + ') .close-ampm-selected').val(),
        });
    }
    return entries;
}

function parseRawData(day) {
    return {
        day: days[day.day],
        opens: convert(day.openHour, day.openMinute, day.openAmpm),
        closes: convert(day.closeHour, day.closeMinute, day.closeAmpm),
    };
}

function submitDaytime( /*event*/ ) {
    var store = [];
    // For each day selected, let's get information about each table row:
    // weekday
    // open hour, open minute, open ampm
    // close hour, close minute, close ampm
    // Then parse raw data and finally, populate `store` array with entries.
    $('.day-selected:checked').map(function() {
        var self = this;
        var entries = parseScheduleRow(this);
        entries = entries.map(parseRawData);
        entries = entries.map(function(entry) {
            entry.note = $(self).find('.note').val();
            return entry;
        });
        store = store.concat(entries);
    });
    if (!store.length) {
        return;
    }

    var options = {
        url: '/schedule',
        data: store,
        method: 'POST',
    };

    origSchedule = $('#daytime-schedule').html();

    // Post data
    $.ajax({
        method: 'POST',
        url: '/schedule',
        data: JSON.stringify(options.data),
        success: function success( /*data*/ ) {
            $.notify({
                message: 'Schedule saved'
            }, {
                type: 'success'
            });
        },
        error: options.fail || displayErrors
    });
}

function setOptions(jqElem) {
    var start = jqElem.data('start');
    var end = jqElem.data('end');
    var i = start;
    var options = '';
    for (i; i <= end; i++) {
        var val = i <= 9 ? '0' + i : i;
        options += '<option value="' + i + '">' + val + '</option>';
    }
    jqElem.html(options);
}

function addHours(elem) {
    var row = elem.closest('tr');
    row = !row.is('tr') ? elem.closest('.row') : row;
    row.find(':has(.time-entry)').each(function() {
        $(this).append($(this).find('.time-entry:last').clone());
        $(this).find('.time-entry:last').find('select').prop('selectedIndex', 0);
    });
}

function mergeScheduleRows() {
    Object.keys(days).forEach(function(day) {
        var duplicateDays = $('.schedule-row').find('td:contains(' + day + ')');
        if (duplicateDays.size() > 1) {
            duplicateDays.not(':first').each(function() {
                var row = $(this).parent();
                row.find(':has(.time-entry)').each(function(i) {
                    duplicateDays.first().parent().find(':has(.time-entry)').eq(i).append($(this).find('.time-entry:last').clone());
                });
                row.remove();
            });
        }
    })
}


function editNote(event) {
    event.preventDefault();
    var triggerElem = $(event.currentTarget);
    var modal = $('#notes-modal');
    var updateNote = false;
    var noteForm = modal.find('textarea');
    noteForm.val(triggerElem.next('.note').val());
    modal.off('hide.bs.modal click').on('hide.bs.modal', function(event) {
        if (updateNote) {
            var hiddenNote = triggerElem.next('.note')
            hiddenNote.val(noteForm.val());
            var event = new Event('change');
            hiddenNote[0].dispatchEvent(event);
        }
    }).on('click', '.finish-btn', function() {
        updateNote = true;
        modal.modal('toggle');
    });
    modal.modal('toggle');
}

$(function() {
    scheduleSubmitBtn.on('click', submitDaytime);
    $('.schedule-note-btn, .exception-note-btn, .calls-note-btn').on('click', editNote);
    $('.range-time-select').each(function() {
        setOptions($(this));
    });
    $('.date-picker').datepicker();
    $('.new-exception, .existing-exception, #daytime-schedule').on('click', '.add-hours', function() {
        addHours($(this));
    }).on('click', '.remove-time-entry', function() {
        var entry = $(this).parent();
        entry.add(entry.parent().next().find('.time-entry').eq(entry.index())).remove();
    });
    //Reset secedule html section
    $('.daytime-cancel').click(function() {
        $('#daytime-schedule').html(origSchedule);
    });
    // $('#daytime-schedule .time-entry:first-of-type .remove-time-entry').css('visibility', 'hidden');
});
// Manually merge schedule rows that are duplicated days since we don't have a way to do that server side right now
mergeScheduleRows();
origSchedule = $('#daytime-schedule').html();
