import "../styles/MyTestStudent.css";
import LogoutButton from "../components/LogoutButton.jsx";
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { internAPI } from "../services/api.js";

export default function StudentHome() {
    const navigate = useNavigate();
    const [activeTab, setActiveTab] = useState("applications");
    const [applications, setApplications] = useState([]);
    const [completedTests, setCompletedTests] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                setError(null);
                const token = localStorage.getItem("token");
                if (!token) {
                    navigate("/login");
                    return;
                }

                const response = await internAPI.getUserEvents();

                const applicationsData = response.data.map(event => ({
                    event_id: event.event_id,
                    name: `Мероприятие ${event.event_id}`,
                    start_date: "2024-01-15",
                    end_date: "2024-01-20"
                }));

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
                    },
                    {
                        attempt_id: 103,
                        test_title: "JavaScript базовый",
                        result_text: "Успешно пройден",
                        score: 92,
                        max_score: 100,
                        passed: true
                    }
                ];

                setApplications(applicationsData);
                setCompletedTests(mockCompleted);
            } catch (error) {
                console.error("Ошибка при загрузке данных:", error);
                setError("Не удалось загрузить данные");
                setApplications([]);
                setCompletedTests([]);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [navigate]);

    const handleGoToEvent = (eventId) => {
        navigate(`/myTestStudent?eventId=${eventId}`);
    };

    const formatDate = (dateString) => {
        const options = { year: 'numeric', month: 'long', day: 'numeric' };
        return new Date(dateString).toLocaleDateString('ru-RU', options);
    };

    const getStatusBadge = (test) => {
        if (test.passed) {
            return { text: "Пройден", class: "status-passed" };
        }
        return { text: "Не пройден", class: "status-failed" };
    };

    if (loading) {
        return (
            <div className="tests-page">
                <div className="test-page" style={{ position: "absolute", left: "1430px", top: "0px" }}>
                    <LogoutButton />
                </div>
                <div className="create-wrapper2">
                    <div className="test">
                        <p className="mytests-loading">Загрузка данных...</p>
                    </div>
                </div>
            </div>
        );
    }

    if (loading) {
        return (
            <div className="tests-page">
                <div className="test-page" style={{ position: "absolute", left: "1430px", top: "0px" }}>
                    <LogoutButton />
                </div>
                <div className="create-wrapper2">
                    <div className="test">
                        <p className="mytests-loading">Загрузка данных...</p>
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
                    {/* НАВИГАЦИОННЫЕ ВКЛАДКИ */}
                    <div className="tests-tabs">
                        <button
                            className={`tab-btn ${activeTab === "applications" ? "tab-btn-active" : ""}`}
                            onClick={() => setActiveTab("applications")}
                        >
                            Заявки
                        </button>
                        <button
                            className={`tab-btn ${activeTab === "tests" ? "tab-btn-active" : ""}`}
                            onClick={() => setActiveTab("tests")}
                        >
                            Тестовые задания
                        </button>
                    </div>


                    {activeTab === "applications" && (
                        <div>
                            {applications.length === 0 ? (
                                <p className="mytests-empty">
                                    У вас нет заявок на мероприятия.
                                </p>
                            ) : (
                                <div className="mytests-list">
                                    {applications.map((app) => (
                                        <div key={app.event_id} className="mytests-card">
                                            <div className="mytests-card-header">
                                                <div>
                                                    <h4 style={{ margin: '0 0 8px 0', color: '#2F4156', fontSize: '16px', fontWeight: '600' }}>
                                                        {app.name}
                                                    </h4>
                                                    <p style={{ margin: '0', color: '#2F4156', fontSize: '13px', opacity: '0.8' }}>
                                                        {formatDate(app.start_date)} | {formatDate(app.end_date)}
                                                    </p>
                                                </div>
                                                <button
                                                    className="mytests-start-btn"
                                                    onClick={() => handleGoToEvent(app.event_id)}
                                                    style={{ whiteSpace: 'nowrap', marginLeft: '16px' }}
                                                >
                                                    Перейти
                                                </button>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    )}

                    {activeTab === "tests" && (
                        <div>
                            <h2>Тестовые задания</h2>
                            {completedTests.length === 0 ? (
                                <p className="mytests-empty">
                                    Вы еще не прошли ни одного теста.
                                </p>
                            ) : (
                                <div className="mytests-list-completed">
                                    {completedTests.map((test) => {
                                        const statusBadge = getStatusBadge(test);
                                        return (
                                            <div key={test.attempt_id} className="mytests-card-completed">
                                                <div className="mytests-completed-info">
                                                    <h4 className="mytests-completed-title-card">
                                                        {test.test_title}
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
                            )}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
