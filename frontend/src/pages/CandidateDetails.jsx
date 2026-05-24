import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import "../styles/CandidateDetails.css";
import LogoutButton from "../components/LogoutButton.jsx";
import BackIcon from "../assets/back.svg?react";
import StatisticsIcon from "../assets/statistics2.svg?react";

export default function CandidateDetails() {
    const { candidateId } = useParams();
    const navigate = useNavigate();
    const [candidateData, setCandidateData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [selectedAttempt, setSelectedAttempt] = useState(null);

    const mockCandidatesData = {
        1: {
            first_name: "Иван",
            last_name: "Иванов",
            email: "ivan@example.com",
            attempts: [
                {
                    attempt_id: 1,
                    test_title: "JavaScript основы",
                    event_name: "Квалификация Frontend 2024",
                    is_extra: false,
                    score: 85,
                    max_score: 100,
                    questions: [
                        { text: "Что такое closure?", points_earned: 10, max_points: 10 },
                        { text: "Объясни this в JavaScript", points_earned: 8, max_points: 10 },
                        { text: "Как работает async/await?", points_earned: 9, max_points: 10 },
                        { text: "Что такое Promise?", points_earned: 8, max_points: 10 },
                        { text: "Разница между let и var", points_earned: 10, max_points: 10 },
                        { text: "Стрелочные функции", points_earned: 9, max_points: 10 },
                        { text: "Деструктуризация", points_earned: 8, max_points: 10 },
                        { text: "Spread оператор", points_earned: 9, max_points: 10 },
                    ]
                },
                {
                    attempt_id: 2,
                    test_title: "React продвинутый",
                    event_name: "Квалификация Frontend 2024",
                    is_extra: false,
                    score: 72,
                    max_score: 100,
                    questions: [
                        { text: "Hooks в React", points_earned: 9, max_points: 10 },
                        { text: "useEffect и побочные эффекты", points_earned: 8, max_points: 10 },
                        { text: "Context API", points_earned: 7, max_points: 10 },
                        { text: "Оптимизация производительности", points_earned: 6, max_points: 10 },
                        { text: "Redux паттерны", points_earned: 8, max_points: 10 },
                        { text: "Error Boundaries", points_earned: 7, max_points: 10 },
                        { text: "Suspense и lazy loading", points_earned: 6, max_points: 10 },
                        { text: "Portals и refs", points_earned: 7, max_points: 10 },
                    ]
                },
                {
                    attempt_id: 3,
                    test_title: "TypeScript",
                    event_name: "Квалификация Frontend 2024",
                    is_extra: true,
                    score: 45,
                    max_score: 100,
                    questions: [
                        { text: "Типизация базовых типов", points_earned: 8, max_points: 10 },
                        { text: "Интерфейсы", points_earned: 5, max_points: 10 },
                        { text: "Generics", points_earned: 4, max_points: 10 },
                        { text: "Union и Intersection типы", points_earned: 6, max_points: 10 },
                        { text: "Decorators", points_earned: 3, max_points: 10 },
                        { text: "Advanced типы", points_earned: 4, max_points: 10 },
                        { text: "Utility типы", points_earned: 5, max_points: 10 },
                        { text: "Module система", points_earned: 0, max_points: 10 },
                    ]
                }
            ]
        },
        2: {
            first_name: "Мария",
            last_name: "Петрова",
            email: "maria@example.com",
            attempts: [
                {
                    attempt_id: 1,
                    test_title: "Node.js",
                    event_name: "Квалификация Backend 2024",
                    is_extra: false,
                    score: 90,
                    max_score: 100,
                    questions: [
                        { text: "Event Loop", points_earned: 10, max_points: 10 },
                        { text: "Модули и require", points_earned: 10, max_points: 10 },
                        { text: "Express.js", points_earned: 9, max_points: 10 },
                        { text: "Middleware", points_earned: 10, max_points: 10 },
                        { text: "Async операции", points_earned: 9, max_points: 10 },
                        { text: "Файловая система", points_earned: 9, max_points: 10 },
                        { text: "Streams", points_earned: 8, max_points: 10 },
                        { text: "Безопасность", points_earned: 8, max_points: 10 },
                    ]
                },
                {
                    attempt_id: 2,
                    test_title: "SQL",
                    event_name: "Квалификация Backend 2024",
                    is_extra: false,
                    score: 88,
                    max_score: 100,
                    questions: [
                        { text: "SELECT и WHERE", points_earned: 10, max_points: 10 },
                        { text: "JOINs", points_earned: 10, max_points: 10 },
                        { text: "Aggregation функции", points_earned: 9, max_points: 10 },
                        { text: "Subqueries", points_earned: 9, max_points: 10 },
                        { text: "Indexes", points_earned: 8, max_points: 10 },
                        { text: "Transactions", points_earned: 8, max_points: 10 },
                        { text: "Stored Procedures", points_earned: 8, max_points: 10 },
                        { text: "Query оптимизация", points_earned: 8, max_points: 10 },
                    ]
                }
            ]
        },
        3: {
            first_name: "Петр",
            last_name: "Сидоров",
            email: "petr@example.com",
            attempts: [
                {
                    attempt_id: 1,
                    test_title: "Тестирование ПО",
                    event_name: "Квалификация QA 2024",
                    is_extra: false,
                    score: 92,
                    max_score: 100,
                    questions: [
                        { text: "Виды тестирования", points_earned: 10, max_points: 10 },
                        { text: "Тест кейсы", points_earned: 10, max_points: 10 },
                        { text: "Баг репорты", points_earned: 9, max_points: 10 },
                        { text: "Регрессионное тестирование", points_earned: 10, max_points: 10 },
                        { text: "Автоматизация тестов", points_earned: 9, max_points: 10 },
                        { text: "Selenium", points_earned: 9, max_points: 10 },
                        { text: "API тестирование", points_earned: 9, max_points: 10 },
                        { text: "Performance тестирование", points_earned: 8, max_points: 10 },
                    ]
                }
            ]
        }
    };

    useEffect(() => {
        const fetchCandidateData = async () => {
            try {
                const token = localStorage.getItem("token");

                if (!token) {
                    const mockData = mockCandidatesData[candidateId];
                    if (mockData) {
                        setCandidateData(mockData);
                    }
                    setLoading(false);
                    return;
                }

                const response = await fetch(
                    `http://localhost:8080/api/manager/users/${candidateId}`,
                    {
                        headers: {
                            "Authorization": `Bearer ${token}`
                        }
                    }
                );

                if (response.ok) {
                    const data = await response.json();
                    setCandidateData(data);
                } else {
                    const mockData = mockCandidatesData[candidateId];
                    setCandidateData(mockData || null);
                }
            } catch (error) {
                console.error("Ошибка при загрузке данных:", error);
                const mockData = mockCandidatesData[candidateId];
                setCandidateData(mockData || null);
            } finally {
                setLoading(false);
            }
        };

        fetchCandidateData();
    }, [candidateId, navigate]);

    const handleBack = () => {
        navigate("/candidates");
    };

    const handleOpenStatistics = (attempt) => {
        setSelectedAttempt(attempt);
    };

    const handleCloseStatistics = () => {
        setSelectedAttempt(null);
    };

    if (loading) {
        return (
            <div className="tests-page">
                <LogoutButton />
                <p>Загрузка...</p>
            </div>
        );
    }

    if (!candidateData) {
        return (
            <div className="tests-page">
                <LogoutButton />
                <p>Данные кандидата не найдены</p>
            </div>
        );
    }

    return (
        <div className="tests-page">
            <div className="test-page">
                <LogoutButton />
            </div>
            <div className="create-wrapper2">
                <div className="test2">
                    <div className="candidate-details-top">
                        <button className="stat-back-btn2" onClick={handleBack}>
                            <BackIcon />
                        </button>
                        <div className="candidate-header-info">
                            <h1>{candidateData.first_name} {candidateData.last_name}</h1>
                            <p className="candidate-email">{candidateData.email}</p>
                        </div>
                    </div>

                    <div className="tests-line"></div>

                    <h2 className="candidate-tests-title">Тестовые задания</h2>

                    {/* Таблица с тестами */}
                    <div className="candidate-tests-table">
                        <table>
                            <thead>
                            <tr>
                                <th>Название теста</th>
                                <th>Мероприятие</th>
                                <th>Дополнительный тест</th>
                                <th>Баллы</th>
                                <th>Статистика</th>
                            </tr>
                            </thead>
                            <tbody>
                            {candidateData.attempts && candidateData.attempts.length > 0 ? (
                                candidateData.attempts.map((attempt) => (
                                    <tr key={attempt.attempt_id}>
                                        <td>{attempt.test_title}</td>
                                        <td>{attempt.event_name}</td>
                                        <td className="extra-test-cell">
                                                <span className={attempt.is_extra ? "badge-yes" : "badge-no"}>
                                                    {attempt.is_extra ? "Да" : "Нет"}
                                                </span>
                                        </td>
                                        <td className="score-cell">
                                            {attempt.score}/{attempt.max_score}
                                        </td>
                                        <td className="statistics-btn-cell">
                                            <button
                                                className="candidate-statistics-btn"
                                                onClick={() => handleOpenStatistics(attempt)}
                                            >
                                                <StatisticsIcon />
                                            </button>
                                        </td>
                                    </tr>
                                ))
                            ) : (
                                <tr>
                                    <td colSpan="5" className="no-attempts">
                                        Тестов не найдено
                                    </td>
                                </tr>
                            )}
                            </tbody>
                        </table>
                    </div>

                    {/* Модальное окно статистики */}
                    {selectedAttempt && (
                        <div className="stat-modal-overlay">
                            <div className="stat-modal">
                                <h3>Подробная статистика</h3>

                                <div className="stat-details-user">
                                    <p>
                                        <strong>Участник:</strong>{" "}
                                        {candidateData.first_name} {candidateData.last_name}
                                    </p>
                                    <p>
                                        <strong>Почта:</strong>{" "}
                                        {candidateData.email}
                                    </p>
                                    <p>
                                        <strong>Результат теста:</strong>{" "}
                                        {selectedAttempt.score}/{selectedAttempt.max_score}
                                    </p>
                                </div>

                                <div className="stat-table-wrapper">
                                    <table className="stat-table">
                                        <thead>
                                        <tr>
                                            <th>Вопрос</th>
                                            <th>Баллы</th>
                                        </tr>
                                        </thead>
                                        <tbody>
                                        {(selectedAttempt.questions || []).map((q, idx) => (
                                            <tr key={idx}>
                                                <td>{q.text}</td>
                                                <td>
                                                    {q.points_earned}/{q.max_points}
                                                </td>
                                            </tr>
                                        ))}
                                        </tbody>
                                    </table>
                                </div>

                                <button
                                    className="stat-hide-btn"
                                    onClick={handleCloseStatistics}
                                >
                                    Скрыть подробную статистику
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}