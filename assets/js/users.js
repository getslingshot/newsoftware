function success( /*data*/ ) {
    jQuery.notify({
        message: 'User updated'
    }, {
        type: 'success'
    })

}

function fail(error) {
    console.log('error', error);

    jQuery.notify({
        message: 'Error while processing your request'
    }, {
        type: 'danger'
    })

}

function ajax(options) {
    $.ajax({
        method: options.method,
        url: options.url,
        contentType: 'application/json; charset=utf-8',
        xhrFields: {
            withCredentials: true
        },
        data: JSON.stringify(options.data),
        success: options.success || success,
        error: options.fail || fail
    });
}

var suspendUsers = $('.suspend-user');
var removeUsers = $('.remove-user');

function suspend(event) {
    var elem = $(event.currentTarget);
    var row = elem.closest('tr');
    var suspended = row.hasClass('suspended');
    suspended ? row.removeClass('suspended') : row.addClass('suspended');
    row.find('.suspend-user').attr('data-original-title', !suspended ? 'Unsuspend User' : 'Suspend User');
    var id = elem.attr('data-id');
    var options = {
        method: 'PUT',
        url: '/users?id=' + id,
        data: {
            suspended: suspended,
        },
    };

    ajax(options);
}

function remove(event) {
    var id = $(event.currentTarget).attr('data-id');
    if (confirm('Are you sure you want to delete this item?')) {
        var options = {
            method: 'DELETE',
            url: '/users?id=' + id,
            success: function(data) {
                var id = JSON.parse(data)._id
                if (!id) {
                    return;
                }
                $('#' + id).remove();
            }
        };

        ajax(options);
    }
}

$(document).ready(function onReady() {
    suspendUsers.on('click', suspend);
    removeUsers.on('click', remove);
});
