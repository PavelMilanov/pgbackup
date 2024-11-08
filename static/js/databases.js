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
    var dbId = $('#backupMenu-btn').val()
    $('#createBackup-btn').val(dbId)
    $('#backupMenuModal').show()
})

$('#app').on('click', '#deleteDatabase-btn', function () {
    var dbId = $('#deleteDatabase-btn').val()
    $('#DeleteDatabaseID').val(dbId)
    $('#deleteDatabaseModal').show()
})