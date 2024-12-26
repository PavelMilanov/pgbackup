window.onload = function () {
    if (window.jQuery) {
        console.log('jQuery is loaded')
    }
    else {
        console.log('jQuery is not loaded')
    }
}

$(document).ready(function () {
    const ctx = $('#storageChart')[0].getContext('2d')

    const used = parseFloat($('#storageChart').data('used'))
    const total = parseFloat($('#storageChart').data('total'))
    const free = total - used

    new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: ['Использовано', 'Свободно'],
            datasets: [{
                data: [used, free],
                backgroundColor: ['#1e3c72', '#e0e0e0']
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false
        }
    })
})

$(document).ready(function () { 
    const ctx = document.getElementById('backupsChart').getContext('2d')

    const success = parseFloat($('#backupsChart').data('completed'))
    const error = parseFloat($('#backupsChart').data('failed'))

    new Chart(ctx, {
        type: 'bar',
        data: {
            labels: ['Успешные', 'Неуспешные'],
            datasets: [{
                label: 'Успешные',
                data: [success, error],
                backgroundColor: [
                    '#d4edda', // Цвет для успешных бэкапов
                    '#f8d7da'   // Цвет для неуспешных бэкапов
                ],
                borderColor: [
                    '#155724',
                    '#721c24'
                ],
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                y: {
                    beginAtZero: true
                }
            },
            responsive: true,
            plugins: {
                legend: {
                    position: 'top',
                },
            }
        }
    })
})
