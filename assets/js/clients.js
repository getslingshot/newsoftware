var form = $('#client-form');
var submit = $('#submit');

function cleanInputs(formId) {
    $('#' + formId + ' :input').each(function eachInput() {
        $(this).val('');
    });
}

function getFieldValues(formId) {
    var values = {};

    $('#' + formId + ' :input').each(function eachInput() {
        values[this.name] = $(this).val();
    });

    return values;
}

function createError(response) {
    if (!response.responseText) {
        return;
    }

    var errors = JSON.parse(response.responseText);

    Object.keys(errors).map(function eachError(key) {
        var selector = $('#' + key);
        if (!selector) {
            return;
        }

        selector.parent().addClass('has-error');
        selector.siblings('.help-block').text(errors[key]);
    });
}

function createSuccess(data) {
    var json = JSON.parse(data);
    var tr = '<tr><td>{{companyName}}</td><td>{{address}}, {{city}}, {{state}}</td><td>{{lastName}}, {{firstName}}</td><td>{{phone}}</td><td>{{email}}</td><td></td></tr>';

    Object.keys(json).map(function eachClient(key) {
        tr = tr.replace('{{' + key + '}}', json[key]);
    })

    $('.table').append(tr);
    cleanInputs('client-form');
    $('#modal-create-client').modal('hide');
}

function ajax(data) {
    $.ajax({
        method: 'POST',
        url: '/clients',
        contentType: 'application/json; charset=utf-8',
        xhrFields: {
            withCredentials: true
        },
        data: JSON.stringify(data),
        success: createSuccess,
        error: createError
    });
}

$(document).ready(function onReady() {
    submit.on('click', function submitForm(event) {
        return ajax(getFieldValues('client-form'));
    })
});

jQuery( document ).ready(function() {
        
        var back =jQuery(".prev");
        var next = jQuery(".next");
        var steps = jQuery(".step");
        
        next.bind("click", function() { 
            jQuery.each( steps, function( i ) {
                if (!jQuery(steps[i]).hasClass('current') && !jQuery(steps[i]).hasClass('done')) {
                    jQuery(steps[i]).addClass('current');
                    jQuery(steps[i - 1]).removeClass('current').addClass('done');
                    return false;
                }
            })      
        });
        back.bind("click", function() { 
            jQuery.each( steps, function( i ) {
                if (jQuery(steps[i]).hasClass('done') && jQuery(steps[i + 1]).hasClass('current')) {
                    jQuery(steps[i + 1]).removeClass('current');
                    jQuery(steps[i]).removeClass('done').addClass('current');
                    return false;
                }
            })      
        });

    })