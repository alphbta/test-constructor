import "../styles/tests.css";
import LogoutButton from "../components/LogoutButton.jsx";
import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import EditIcon from "../assets/edit.svg?react";
import ShareIcon from "../assets/share.svg?react";
import StatisticsIcon from "../assets/statistics.svg?react";
import DeleteIcon from "../assets/close.svg?react";


export default function Tests() {
    const navigate = useNavigate();
    const [tests, setTests] = useState([]);
    const [openMenuId, setOpenMenuId] = useState(null);
    const menuRefs = useRef({});

    useEffect(() => {
        const savedTests = JSON.parse(localStorage.getItem("savedTests")) || [];
        setTests(savedTests);
    }, []);

    const toggleMenu = (id, e) => {
        if (e) e.stopPropagation();
        setOpenMenuId(openMenuId === id ? null : id);
    };

    useEffect(() => {
        const handleClickOutside = (e) => {
            let clickedInsideMenu = false;

            Object.values(menuRefs.current).forEach(ref => {
                if (ref && ref.contains(e.target)) {
                    clickedInsideMenu = true;
                }
            });

            if (!clickedInsideMenu) {
                setOpenMenuId(null);
            }
        };

        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    const editTest = (test) => {
        navigate("/create", { state: { editing: true, test } });
        setOpenMenuId(null);
    };

    const deleteTest = (id) => {
        if (window.confirm("Удалить этот тест?")) {
            const updatedTests = tests.filter(test => test.id !== id);
            setTests(updatedTests);
            localStorage.setItem("savedTests", JSON.stringify(updatedTests));
            setOpenMenuId(null);
        }
    };
    const shareTest = () => {}
    const closeTest = (id) => {
        const updatedTests = tests.map(test =>
            test.id === id ? { ...test, isClosed: true } : test
        );
        setTests(updatedTests);
        localStorage.setItem("savedTests", JSON.stringify(updatedTests));
        setOpenMenuId(null);
        alert("Тест закрыт (деактивирован)");
    };

    return (
        <div className="tests-page">
            <>
                <LogoutButton />
            </>
        <div className="tests-wrapper">
            <div className="tests-left">
                    <h2>Мои тесты</h2>
                    <div className="tests-line"></div>

                {tests.length === 0 ? (
                    <div className="no-tests">
                        Пока нет тестов. Создайте первый тест →
                    </div>
                ) : (
                    <div className="tests-grid">
                        {tests.map((test) => (
                            <div key={test.id} className="test-card"
                                 style={{
                                     zIndex: openMenuId === test.id ? 100 : 1
                                 }}
                            >
                                <div
                                    className="test-menu-container"
                                    ref={el => menuRefs.current[test.id] = el}
                                >
                                    <button
                                        className="dots-btn"
                                        onClick={(e) => toggleMenu(test.id, e)}
                                    >
                                        ⋮
                                    </button>

                                    {openMenuId === test.id && (
                                        <div className="dropdown-menu">
                                            <button className="menu-item" onClick={() => editTest(test)}>
                                                <EditIcon className="menu-icon" />
                                                <span>Редактировать</span>
                                            </button>
                                            <button className="menu-item share" onClick={() => shareTest(test.id)}>
                                                <ShareIcon className="menu-icon" />
                                                <span>Поделиться</span>
                                            </button>
                                            {/* тут поменять с клос на нормальынй */}
                                            <button className="menu-item" onClick={() => closeTest(test.id)}>
                                                <StatisticsIcon className="menu-icon" />
                                                <span>Статистика</span>
                                            </button>
                                            <button className="menu-item" onClick={() => deleteTest(test.id)}>
                                                <DeleteIcon className="menu-icon" />
                                                <span>Закрыть тест</span>
                                            </button>
                                            {/*<div className="menu-divider"></div>*/}
                                            {/*<button className="menu-item delete" onClick={() => deleteTest(test.id)}>*/}
                                            {/*    <span className="menu-icon">🗑️</span>*/}
                                            {/*    <span>Удалить</span>*/}
                                            {/*</button>*/}
                                        </div>
                                    )}
                                </div>
                                    <span className="test-titles">
                                        {test.title.length > 15
                                        ? `${test.title.substring(0, 15)}...`
                                        : test.title
                                    }</span>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            <div className="tests-right">
                <button className="create-test-btn" onClick={() => navigate("/create")}>
                    Создать тест
                </button>
            </div>
        </div>
        </div>
    );
}