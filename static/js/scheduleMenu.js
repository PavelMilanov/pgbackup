$(function () {
    $("#backupRun").change(function () { 
        // var selectedValue = $(this).val()
        const collapseElement = document.getElementById('collapseBackupSchedule');
        new bootstrap.Collapse(collapseElement, {
            toggle: true
        })
    })
})