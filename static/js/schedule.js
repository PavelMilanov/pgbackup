$('.close').on('click', function () {
    $(this).closest('.modal').hide()
})

$(window).on('click', function (event) {
    if ($(event.target).hasClass('modal')) {
        $(event.target).hide()
    }
})

function showAddScheduleModal() {
    $('#addScheduleModal').show()
}

function deleteScheduleModal() {
    var chedule = $('#deleteSchedule').val()
    $('#DeleteScheduleID').val(chedule)
    $('#deleteScheduleModal').show()
}

function showChangeScheduleModal() {
    var dbName = $('#scheduleDbName').text()
    var dbTime = $('#scheduleDbTime').text()
    var chedule = $('#changeShedule').val()
    $('#ChangeScheduleFormDB').val(dbName)
    $('#ChangeScheduleFormTime').val(dbTime)
    $('#ChangeScheduleID').val(chedule)
    $('#changeScheduleModal').show()
}

$('#app').on('click', '.btn-primary', showAddScheduleModal)
$('#app').on('click', '#changeShedule', showChangeScheduleModal)
$('#app').on('click', '#deleteSchedule', deleteScheduleModal)
