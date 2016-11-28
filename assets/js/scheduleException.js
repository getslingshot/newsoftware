'use strict';

var exceptionCache = {};

function removeExceptionSuccess(data, text, xhr) {
    var id = data._id;
    $.notify({
        message: 'Exception removed!'
    }, {
        type: 'success'
    });

    if (!id) {
        return;
    }

    $('#row-' + id).remove();
}

function removeExceptionFailure(error) {
    $.notify({
        message: 'Error on remove exception, cant process request'
    }, {
        type: 'danger'
    });
    console.log('on remove exception', error);
}

function removeException(event) {
    var id = $(this).parent().parent().data('id');
    if (!id) {
        return;
    }
    if (confirm('Are you sure you want to delete this item?')) {
        $.ajax({
            method: 'DELETE',
            url: '/delete-exception?id=' + id,
            contentType: 'application/json; charset=utf-8',
            xhrFields: {
                withCredentials: true
            },
            success: removeExceptionSuccess,
            error: removeExceptionFailure,
        });
    }
}


function wholeDay(event) {
    var wholeDayCheckbox = $(this);
    var row = wholeDayCheckbox.closest('.row');
    if (wholeDayCheckbox.is(':checked')) {
        row.find('.time-entry:gt(0)').remove();
        var entryContainer = row.find('.time-entry');
        entryContainer.find('select').filter(':lt(4)').find('option:first').prop('selected', true);
        entryContainer.find('select').filter(':gt(2)').find('option:last').prop('selected', true);
    } else {
        row.find('.time-entry').find('select option:first').prop('selected', true);
    }
}

$('input[name=endAt]').closest('.form-group').toggle();

function dateRange(event) {
    $(this).parent().prev().toggle();
}

function parseException(exceptionBox) {
    // Hack for now
    var applyDaytime = exceptionBox.find('input[name=apply-daytime]');
    var applyAfterHours = exceptionBox.find('input[name=apply-afterhours]');
    var startAt = exceptionBox.find('input[name=startAt]');
    var endAt = exceptionBox.find('input[name=endAt]');
    var exceptionMatrix = exceptionBox.find('.exceptions-row');
    var d = new Date(startAt.val());
    d.setMinutes(0);
    d.setHours(0);
    d.setMilliseconds(0);
    var e = new Date(endAt.val());
    e.setMinutes(0);
    e.setHours(0);
    e.setMilliseconds(0);
    var dayTime = applyDaytime.is(':checked');
    var afterhours = applyAfterHours.is(':checked');
    var entryCount = exceptionMatrix.find('.time-entry').size();

    var data = {
        startAt: d,
        endAt: endAt.val() ? e : d,
        daytime: applyDaytime.is(':checked'),
        afterhours: applyAfterHours.is(':checked'),
        hours: [],
    };

    for (var i = 1; i <= entryCount; i++) {
        var raw = {
            openHour: exceptionMatrix.find('.time-entry:nth-child(' + i + ') select[name="open-hours"]').val(),
            openMinute: exceptionMatrix.find('.time-entry:nth-child(' + i + ') select[name="open-minutes"]').val(),
            openAmPm: exceptionMatrix.find('.time-entry:nth-child(' + i + ') select[name="open-am-pm"]').val(),

            closeHour: exceptionMatrix.find('.time-entry:nth-child(' + i + ') select[name="close-hours"]').val(),
            closeMinute: exceptionMatrix.find('.time-entry:nth-child(' + i + ') select[name="close-minutes"]').val(),
            closeAmPm: exceptionMatrix.find('.time-entry:nth-child(' + i + ') select[name="close-am-pm"]').val(),
        }
        data.hours.push({
            opens: convert(raw.openHour, raw.openMinute, raw.openAmPm),
            closes: convert(raw.closeHour, raw.closeMinute, raw.closeAmPm)
        });
    }

    data.note = exceptionMatrix.find('.note').val();


    return data;
}


function submitException(event) {
    var ref = $(this);
    var exceptionBox = $(this).parent();
    var data = parseException(exceptionBox);
    $.ajax({
        method: 'POST',
        url: '/exception',
        data: JSON.stringify(data),
        success: function(data) {
            $.notify({
                message: 'Exception saved'
            }, {
                type: 'success'
            });
            var expTemplate = $('#existing-exception-template').html();
            var count = $('.existing-exception').size() + 1;
            // TODO: return all the exception data from the newly created exception or build it from edit form
            data.num = count;
            data.id = data['_id'];
            data.range = moment(data.startAt).format('M/D') + ' - ' + moment(data.endAt).format('M/D');
            var template = Mustache.render(expTemplate, data);
            $('.existing-exceptions').prepend(template);
            $('.existing-exception').each(function(i) {
                var lbl = $(this).find('.h4:first');
                lbl.text(lbl.text().split('#')[0] + '# ' + (i + 1));
            });
            var existingException = $('.existing-exception:first');
            existingException.find('.block-content').append(ref.parent().clone());
            existingException.find('.rectangle-exception').hide();;
            $('.nex-exception').find('input, textarea').val('');
            exceptionCache[existingException.data('id')] = existingException.find('.rectangle-exception').clone();
        },
        error: displayErrors
    });
}

function updateException(event) {
    var existingException = $(this).closest('.existing-exception');
    var exceptionBox = $(this).closest('.rectangle-exception');
    var data = parseException(exceptionBox);
    $.ajax({
        method: 'PUT',
        url: '/update-exception?id=' + existingException.data('id'),
        data: JSON.stringify(data),
        success: function(data) {
            $.notify({
                message: 'Exception Updated'
            }, {
                type: 'success'
            });
            var updatedRange = moment(data.startAt).format('M/D') + ' - ' + moment(data.endAt).format('M/D');
            existingException.find('.h4:last').text(updatedRange);
            exceptionCache[existingException.data('id')] = existingException.find('.rectangle-exception').clone();
        },
        error: displayErrors
    });
}

function toggleException() {
    return $(this).parent().find('.rectangle-exception').slideToggle();
}

function unsetWholeDay(event) {
    $(this).closest('.rectangle-exception').find('.whole-day').prop('checked', false);
}

function revertException() {
    var exception = $(this).closest('.existing-exception');
    if (!exception.data('id')) {
        exception = $(this).closest('.new-exception');
    }
    exception.find('.rectangle-exception').replaceWith(exceptionCache[exception.data('id')]);
    exceptionCache[exception.data('id')] = exception.find('.rectangle-exception').clone();
    $('.date-picker').datepicker();
}

$(function() {
    $('.scheduling-logic-page')
        .on('click', '.new-exception .submit-exception', submitException)
        .on('change', '.whole-day', wholeDay)
        .on('change', '.date-range', dateRange)
        .on('click', '.remove-exception', removeException)
        .on('click', '.edit-exception', toggleException)
        .on('click', '.rectangle-exception .add-hours', unsetWholeDay)
        .on('click', '.rectangle-exception .cancel-btn', revertException)
        .on('click', '.existing-exception .submit-exception', updateException);
    $('.existing-exception, .new-exception').each(function() {
        var excep = $(this);
        exceptionCache[excep.data('id')] = $(this).find('.rectangle-exception').clone();
        excep.find('.date-range input:checked').parent().parent().prev().show();
    })
});
