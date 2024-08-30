$(function () {
    $("#createBackup").click(function () {
        var selectedDatabase = $('#dbname').val(); // получаем значение тега select dbanme
        var selectedSchedule = $('#schedule').val(); // получаем значение тега select schedule
        $.ajax({
            url: "http://localhost:8080/backups/create",
            type: "POST",
            dataType: "json",
            contentType: "application/json",
            data: JSON.stringify({ db: selectedDatabase, schedule: selectedSchedule }),
            success: function (response) {
                if (response.error) {
                    $("#backupErrorText").text(response.error) // вставляем текст в элемент по id="backupEror"
                    $("#backupError").modal("show") // вызываем элемент по id="backupEror"
                } else {
                    var index = $('#backupsTable tr:eq(1) th:eq(0)').text() // Получаем индекс последней строки
                    var newRow = $("<tr>")
                        .append($("<th>").text(parseInt(index) + 1))
                        .append($("<td>").text(response.message["Alias"]))
                        .append($("<td>").text(response.message["Date"]))
                        .append($("<td>").text(response.message["Size"]))
                        .append($("<td>").text(response.message["LeadTime"]))
                        .append($("<td>").text(response.message["Status"]))
                        .append($("<td>").text(response.message["Run"]))

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