$(function () {
    
})

$('.close').on('click', function () {
    $(this).closest('.modal').hide();
})

$(window).on('click', function (event) {
    if ($(event.target).hasClass('modal')) {
        $(event.target).hide();
    }
})

function showAddDatabaseModal() {
    $('#databaseModal').show()
}
$('#app').on('click', '.btn-primary', showAddDatabaseModal)