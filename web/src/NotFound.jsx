import { createEffect } from 'solid-js'
import { useNavigate } from "@solidjs/router"

function NotFound() {
    const navigate = useNavigate()

    createEffect(() => {
        setTimeout(() => navigate("/", { replace: true }), 2000)
    })

    return (
        <div>
            <h1>404 - Страница не найдена</h1>
            <p>Перенаправляем на главную через 2 секунды...</p>
        </div>
    )
}

export default NotFound