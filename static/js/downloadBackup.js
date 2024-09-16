$(function () {
    $(".downloadBackup").click(function () { 
        // Получаем данные из второй строки (индекс 1)
        var index = $(this).closest("tr").index()  // получаем индекс
        var row = $('#backupsTable tbody tr').eq(index)  // получаем строку
        var alias = row.find('td').eq(0).text()  // получаем столбцы строки
        var date = row.find('td').eq(1).text()  // получаем столбцы строки
        window.location.href = `download/${alias}/${date}`
    })
})