$(function () {
    $("#backupRun").change(function () { 
        const collapseElement = document.getElementById('collapseBackupSchedule')
        new bootstrap.Collapse(collapseElement, {
            toggle: true
        })
    })
})