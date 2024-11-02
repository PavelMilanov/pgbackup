$('.close').on('click', function () {
    $(this).closest('.modal').hide()
})

$(window).on('click', function (event) {
    if ($(event.target).hasClass('modal')) {
        $(event.target).hide()
    }
})

$('#app').on('click', '.btn-primary', function () {
    $('#addScheduleModal').show()
})

$('#app').on('click', '#changeShedule', function () {
    var dbName = $('#scheduleDbName').text()
    var dbTime = $('#scheduleDbTime').text()
    var chedule = $('#changeShedule').val()
    $('#ChangeScheduleFormDB').val(dbName)
    $('#ChangeScheduleFormTime').val(dbTime)
    $('#ChangeScheduleID').val(chedule)
    $('#changeScheduleModal').show()
})

$('#app').on('click', '#deleteSchedule', function () {
    var chedule = $('#deleteSchedule').val()
    $('#DeleteScheduleID').val(chedule)
    $('#deleteScheduleModal').show()
})
