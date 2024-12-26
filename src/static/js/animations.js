function showNotification(message, type = 'info') {
    // info, errror, success
    $('.notification').remove();
    const notification = $('<div>')
        .addClass(`notification ${type}`)
        .text(message)
        .appendTo('body')

    setTimeout(() => {
        notification.fadeOut(300, function () {
            $(this).remove();
        })
    }, 3000)
}