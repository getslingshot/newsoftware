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

var roleSelector = $('#change-role');

function changeRole(event) {
    var id = $('#role').attr('data-id');
    var value = $('#role').val();
    var options = {
        method: 'PUT',
        url: '/update-user?id=' + id,
        data: {
            roleId: value,
        },
    };

    ajax(options);
}

$(document).ready(function onReady() {
    roleSelector.on('click', changeRole);
});