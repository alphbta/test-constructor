import "../styles/MyTestStudent.css";
import LogoutButton from "../components/LogoutButton.jsx";
import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";

export default function MyTestStudent() {
    const navigate = useNavigate();
    const [tests, setTests] = useState([]);
    const [openMenuId, setOpenMenuId] = useState(null);
    const menuRefs = useRef({});

    useEffect(() => {
        const userRaw = localStorage.getItem("user");
        let userEmail = null;

        try {
            if (userRaw) {
                const user = JSON.parse(userRaw);
                userEmail = user.email || user.username || user.login || null;
            }
        } catch (e) {
            console.error("Не удалось распарсить user", e);
        }

        if (!userEmail) {
            setTests([]);
            return;
        }

        const key = `savedTests_${userEmail}`;
        const savedTests = JSON.parse(localStorage.getItem(key)) || [];
        setTests(savedTests);
    }, []);


    const toggleMenu = (id, e) => {
        if (e) e.stopPropagation();
        setOpenMenuId(openMenuId === id ? null : id);
    };

    useEffect(() => {
        const handleClickOutside = (e) => {
            let clickedInsideMenu = false;
            Object.values(menuRefs.current).forEach((ref) => {
                if (ref && ref.contains(e.target)) {
                    clickedInsideMenu = true;
                }
            });
            if (!clickedInsideMenu) {
                setOpenMenuId(null);
            }
        };

        document.addEventListener("mousedown", handleClickOutside);
        return () =>
            document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    return (
        <div className="tests-page">
            <div
                className="test-page"
                style={{ position: "absolute", left: "1430px", top: "0px" }}
            >
                <LogoutButton />
            </div>
            <div className="create-wrapper2">
                <div className="test">
                    <h2>Мои тесты</h2>
                    <div className="tests-line"></div>

                    {tests.length === 0 ? (
                        <p className="mytests-empty">
                            Вы ещё не проходили ни одного теста.
                        </p>
                    ) : (
                        <div className="mytests-list">
                            {tests.map((t, index) => (
                                <div
                                    key={`${t.id}-${index}`}
                                    className="mytests-card"
                                >
                                    <h3 className="mytests-card-title">
                                        {t.title}
                                    </h3>

                                    <p className="mytests-card-message">
                                        {t.message}
                                    </p>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );

}