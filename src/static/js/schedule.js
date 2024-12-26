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
    $('#changeScheduleModal').show()
})

$('#app').on('click', '#deleteSchedule', function () {
    var chedule = $('#deleteSchedule').val()
    $('#DeleteScheduleID').val(chedule)
    $('#deleteScheduleModal').show()
})
