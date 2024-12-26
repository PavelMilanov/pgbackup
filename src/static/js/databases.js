$('.close').on('click', function () {
    $(this).closest('.modal').hide()
})

$(window).on('click', function (event) {
    if ($(event.target).hasClass('modal')) {
        $(event.target).hide()
    }
})

$('#app').on('click', '#dbMenu-btn', function () {
    $('#databaseModal').show()
})

$('#app').on('click', '#backupMenu-btn', function () {
    $('#backupMenuModal').show()
})

$('#app').on('click', '#deleteDatabase-btn', function () {
    var db = $('deleteDatabase-btn').val()
    $('#modal-deleteDatabase-btn').val(db)
    $('#deleteDatabaseModal').show()
})
