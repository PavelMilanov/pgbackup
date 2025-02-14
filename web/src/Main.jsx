function Main() {
    return (
        <div class="container" id="app">
            <h2>Панель управления</h2>
            <div class="dashboard">
                <div class="card">
                    <h3>Общая статистика</h3>
                    <div class="stat">.count.Total </div>
                    <div>Статистика дампов</div>
                    <div class="chart">
                        <canvas id="backupsChart"></canvas>
                    </div>
                </div>
                <div class="card">
                    <h3>Использование хранилища</h3>
                    <div class="stat">.storage.Used / .storage.Total </div>
                    <div>Использовано / Всего</div>
                    <div class="chart">
                        <canvas id="storageChart"></canvas>
                    </div>
                </div>
                <div class="card">
                    <h3>Состояние системы</h3>
                    <div id="systemStatusWidget" class="system-status-widget">
                        <div class="status-indicator">
                            <div class="status-icon status-green">
                            </div>
                            <div class="status-details">
                                <div class="status-name">
                                    Всего расписаний
                                </div>
                                <div class="status-value">
                                    .schedules_count
                                </div>
                            </div>
                        </div>
                        <div class="status-indicator">
                            <div class="status-icon status-green">
                            </div>
                            <div class="status-details">
                                <div class="status-name">
                                    Всего дампов
                                </div>
                                <div class="status-value">
                                    .backups_count
                                </div>
                            </div>
                        </div>
                        <div class="status-indicator">
                            <div class="status-icon status-green">
                            </div>
                            <div class="status-details">
                                <div class="status-name">
                                    Нагрузка системы
                                </div>
                                <div class="status-value">
                                    .system.CPU % CPU, .system.RAM % RAM
                                </div>
                            </div>
                        </div>
                        <div class="status-indicator">
                            <div class="status-icon status-green">
                            </div>
                            <div class="status-details">
                                <div class="status-name">
                                    Хранилие
                                </div>
                                <div class="status-value">
                                    .system.Storage свободно
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="card">
                <h3>Последние бэкапы</h3>
                <table>
                    <thead>
                        <tr>
                            <th>База данных</th>
                            <th>Дата</th>
                            <th>Размер</th>
                            <th>Время</th>
                            <th>Статус</th>
                            <th>Действия</th>
                        </tr>
                    </thead>
                    <tbody>
                        {/* {{ range .backups }}
                        <tr>
                            <td>{{ .DatabaseAlias }}</td>
                            <td>{{ .Date }}</td>
                            <td>{{ .Size }}</td>
                            <td>{{ .LeadTime }}</td>
                            <td>
                                {{ if .Status }}
                                <span class="status status-success">
                                    Успешно
                                </span>
                                {{ else}}
                                <span class="status status-error">
                                    Ошибка
                                </span>
                                {{ end }}
                            </td>
                            <td>
                                <form class="btn-form btn btn-primary" action="/databases/backup/download" method="get">
                                    <button type="submit" name="ID" value="{{.ID}}" class="btn btn-primary">скачать</button>
                                </form>
                            </td>
                        </tr>
                        {{ end }} */}
                    </tbody>
                </table>
            </div> 
        </div>
    )
}

export default Main
