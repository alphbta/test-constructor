import "../styles/MyTestStudent.css";
import LogoutButton from "../components/LogoutButton.jsx";
import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { testsAPI } from "../services/api.js";
import notebookIcon from "../assets/bloknot.svg";
import timeIcon from "../assets/time2.svg";

export default function MyTestStudent() {
    const navigate = useNavigate();
    const [availableTests, setAvailableTests] = useState([]);
    const [completedTests, setCompletedTests] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchTests = async () => {
            try {
                setLoading(true);
                setError(null);
                const token = localStorage.getItem("token");
                if (!token) {
                    navigate("/login");
                    return;
                }

                // ПУСТЫШКИ ДЛЯ ТЕСТИРОВАНИЯ
                const mockAvailable = [
                    {
                        config_id: 1,
                        test_id: 1,
                        test_link: "demo-test-1",
                        title: "Тест по базам данных",
                        description: "Проверка знаний SQL-запросов, проектирования баз данных и принципов нормализации данных.",
                        time_limit: 1800,
                        status: "available"
                    },
                    {
                        config_id: 2,
                        test_id: 2,
                        test_link: "demo-test-2",
                        title: "Тест по алгоритмам и структурам данных",
                        description: "Проверка знаний основных алгоритмов, структур данных и принципов их применения при решении задач.",
                        time_limit: 1500,
                        status: "available"
                    },
                    {
                        config_id: 3,
                        test_id: 3,
                        test_link: "demo-test-3",
                        title: "Тест по основам Python",
                        description: "Тест предназначен для оценки знаний синтаксиса Python, работы с функциями, коллекциями и основными конструкциями языка.",
                        time_limit: 1200,
                        status: "available"
                    }
                ];

                const mockCompleted = [
                    {
                        attempt_id: 101,
                        test_title: "Введение в веб-разработку",
                        result_text: "Успешно пройден",
                        score: 85,
                        max_score: 100,
                        passed: true
                    },
                    {
                        attempt_id: 102,
                        test_title: "Основы CSS и HTML",
                        result_text: "Не пройден",
                        score: 45,
                        max_score: 100,
                        passed: false
                    }
                ];

                setAvailableTests(mockAvailable);
                setCompletedTests(mockCompleted);
            } catch (error) {
                console.error("Ошибка при загрузке тестов:", error);
                setError("Не удалось загрузить тесты");
                setAvailableTests([]);
                setCompletedTests([]);
            } finally {
                setLoading(false);
            }
        };

        fetchTests();
    }, [navigate]);

    const handleStartTest = (testLink) => {
        navigate(`/test-preview/${testLink}`);
    };

    const formatTime = (seconds) => {
        if (!seconds) return "Не ограничено";
        const minutes = Math.floor(seconds / 60);
        const hours = Math.floor(minutes / 60);

        if (hours > 0) {
            return `${hours} ч ${minutes % 60} мин`;
        } else if (minutes > 0) {
            return `${minutes} мин`;
        } else {
            return `${seconds} сек`;
        }
    };

    const getStatusBadge = (test) => {
        if (test.passed) {
            return { text: "Пройден", class: "status-passed" };
        }
        return { text: " Не пройден", class: "status-failed" };
    };

    if (loading) {
        return (
            <div className="tests-page">
                <div className="test-page" style={{ position: "absolute", left: "1430px", top: "0px" }}>
                    <LogoutButton />
                </div>
                <div className="create-wrapper2">
                    <div className="test">
                        <p className="mytests-loading">Загрузка тестов...</p>
                    </div>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="tests-page">
                <div className="test-page" style={{ position: "absolute", left: "1430px", top: "0px" }}>
                    <LogoutButton />
                </div>
                <div className="create-wrapper2">
                    <div className="test">
                        <p className="mytests-error">{error}</p>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="tests-page">
            <div className="test-page" style={{ position: "absolute", left: "1430px", top: "0px" }}>
                <LogoutButton />
            </div>
            <div className="create-wrapper2">
                <div className="test">
                    {/* ДОСТУПНЫЕ ТЕСТЫ */}
                    <h2>Доступные тесты</h2>
                    <div className="tests-line"></div>

                    {availableTests.length === 0 ? (
                        <p className="mytests-empty">
                            Нет доступных тестов.
                        </p>
                    ) : (
                        <div className="mytests-list">
                            {availableTests.map((test, index) => (
                                <div key={`available-${test.config_id || index}`} className="mytests-card-new">
                                    <div className="mytests-card-left">
                                        <div className="mytests-card-icon">
                                            <img src={notebookIcon} alt="test" style={{ width: '48px', height: '48px' }} />
                                        </div>
                                        <div className="mytests-card-info">
                                            <p className="mytests-label">Название теста</p>
                                            <h3 className="mytests-card-title">
                                                {test.title || test.test_title || "Название теста"}
                                            </h3>
                                            <p className="mytests-label">Описание</p>
                                            {test.description && (
                                                <p className="mytests-card-description">
                                                    {test.description}
                                                </p>
                                            )}
                                        </div>
                                    </div>

                                    <div className="mytests-card-divider"></div>

                                    <div className="mytests-card-right">
                                        <div className="mytests-time-section">
                                            <p className="mytests-time-label">Время прохождения</p>
                                            <div className="mytests-time-content">
                                                <div className="mytests-time-icon">
                                                    <img src={timeIcon} alt="time" style={{ width: '36px', height: '36px' }} />
                                                </div>
                                                <p className="mytests-time-value">
                                                    {formatTime(test.time_limit)}
                                                </p>
                                            </div>
                                        </div>
                                        <button
                                            className="mytests-start-btn"
                                            onClick={() => handleStartTest(test.test_link)}
                                        >
                                            Перейти к тестированию
                                        </button>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}

                    {/* РАЗДЕЛИТЕЛЬ */}
                    {completedTests.length > 0 && availableTests.length > 0 && (
                        <div className="mytests-separator">
                            <div className="separator-line"></div>
                            <span className="separator-text">Пройденные тесты</span>
                            <div className="separator-line"></div>
                        </div>
                    )}

                    {/* ПРОЙДЕННЫЕ ТЕСТЫ */}
                    {completedTests.length > 0 && (
                        <div className="mytests-completed">
                            <h3 className="mytests-completed-title">Пройденные тесты</h3>
                            <div className="mytests-list-completed">
                                {completedTests.map((test, index) => {
                                    const statusBadge = getStatusBadge(test);
                                    return (
                                        <div key={`completed-${test.attempt_id || index}`} className="mytests-card-completed">
                                            <div className="mytests-completed-info">
                                                <h4 className="mytests-completed-title-card">
                                                    {test.title || test.test_title || "Название теста"}
                                                </h4>
                                                <span className={`mytests-status-badge ${statusBadge.class}`}>
                                                    {statusBadge.text}
                                                </span>
                                            </div>
                                            <div className="mytests-completed-stats">
                                                {test.score !== undefined && test.max_score !== undefined && (
                                                    <div className="stat-item">
                                                        <span className="stat-label">Баллы:</span>
                                                        <span className="stat-value">{test.score}/{test.max_score}</span>
                                                    </div>
                                                )}
                                                {test.result_text && (
                                                    <div className="stat-item">
                                                        <span className="stat-label">Результат:</span>
                                                        <span className="stat-value">{test.result_text}</span>
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    );
                                })}
                            </div>
                        </div>
                    )}

                    {availableTests.length === 0 && completedTests.length === 0 && (
                        <p className="mytests-empty">
                            Нет тестов для отображения.
                        </p>
                    )}
                </div>
            </div>
        </div>
    );
}