$(function () {
    $("#createBackup").click(function () {
        var selectedDatabase = $('#dbname').val()
        var selectedRun = $('#backupRun').val()
        var selectedComment = $('#backupComment').val()
        var selectedCount = $('#backupScheduleCount').val()
        var selectedTime = $('#backupScheduleTime').val()
        var selectedCron = $('#backupScheduleCron').val()
        $.ajax({
            url: "http://localhost:8080/backups/create",
            type: "POST",
            dataType: "json",
            contentType: "application/json",
            data: JSON.stringify({
                db: selectedDatabase,
                run: selectedRun,
                comment: selectedComment,
                count: selectedCount,
                time: selectedTime,
                cron: selectedCron
            }),
            success: function (response) {
                if (response.error) {
                    $("#backupErrorText").text(response.error) // вставляем текст в элемент по id="backupEror"
                    $("#backupError").modal("show") // вызываем элемент по id="backupEror"
                } else {
                    var index = $('#backupsTable tr:eq(1) th:eq(0)').text() // Получаем индекс последней строки

                    // генерируем td с кнопками как в разметке
                    var downloadButton = $('<button>', {
                        type: 'button',
                        id: "downloadBackup",
                        class: 'btn btn-success btn-sm',
                        text: 'Скачать',
                    })
                    var deleteButton = $('<button>', {
                        type: 'button',
                        class: 'btn btn-danger btn-sm',
                        text: 'Удалить',
                    })

                    var newRow = $("<tr>")
                        .append($("<th>").text(parseInt(index) + 1))
                        .append($("<td>").text(response.message["Alias"]))
                        .append($("<td>").text(response.message["Date"]))
                        .append($("<td>").text(response.message["Size"]))
                        .append($("<td>").text(response.message["LeadTime"]))
                        .append($("<td>").text(response.message["Status"]))
                        .append($("<td>").text(response.message["Schedule"]["Run"]))
                        .append($("<td>").text(response.message["Comment"]))
                        .append($("<td>").append(downloadButton).append(deleteButton))

                    // Добавляем новую строку в таблицу
                    $("#backupsTable tbody").append(newRow)
                }
            },
            error: function (error) {
                console.error("Ошибка:", error)
            }
        })
    })
})
