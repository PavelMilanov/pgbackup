import { useNavigate } from '@solidjs/router'

function Logout() {
    const navigate = useNavigate()

    function logout() {
        localStorage.removeItem("token")
        navigate("/login", { replace: true })
    }

    return (
        <div class="modal" id="logoutModal">
            <div class="modal-content">
                <span id="logoutModal" class="close">&times;</span>
                <h2>Уверены, что хотите выйти?</h2>
                <button class="btn btn-primary" onclick={logout}>Выйти</button>
            </div>
        </div>
    )
}

export default Logout