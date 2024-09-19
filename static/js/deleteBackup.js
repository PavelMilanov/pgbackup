$(function () {
  $(".deleteBackup").click(function () {
    // Получаем данные из второй строки (индекс 1)
    var index = $(this).closest("tr").index(); // получаем индекс
    var row = $("#backupsTable tbody tr").eq(index) // получаем строку
    var alias = row.find("td").eq(0).text() // получаем столбцы строки
    var date = row.find("td").eq(1).text() // получаем столбцы строки
    $.ajax({
      url: `http://localhost:8080/backups/delete/${alias}/${date}`,
      type: "POST",
      dataType: "json",
      contentType: "application/json",
      success: function (response) {
        if (response.error) {
          $("#backupErrorText").text(response.error) // вставляем текст в элемент по id="backupEror"
          $("#backupError").modal("show") // вызываем элемент по id="backupEror"
        }
      },
      error: function (error) {
        console.error("Ошибка:", error)
      },
    })
    row.remove()
  })
})
