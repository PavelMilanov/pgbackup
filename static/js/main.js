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
    new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: ['Использовано', 'Свободно'],
            datasets: [{
                data: [used, total],
                backgroundColor: ['#1e3c72', '#e0e0e0']
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false
        }
    })
})