import {useNavigate} from "react-router-dom";

const LogoutButton = () => {
    const nav = useNavigate();

    const logoutUser = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        nav('/login');
    };

    return (
        <button onClick={logoutUser}>Выйти</button>
    );
}

export default LogoutButton;