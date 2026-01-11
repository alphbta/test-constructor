import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import "../styles/StatisticsTest.css";
import LogoutButton from "../components/LogoutButton.jsx";
import BackIcon from "../assets/back.svg?react";

export default function StatisticsTest() {
    const { testId } = useParams();
    const navigate = useNavigate();
    const [attempts, setAttempts] = useState([]);
    const [selectedAttempt, setSelectedAttempt] = useState(null);

    useEffect(() => {
        if (!testId) return;
        try {
            const key = `attempts_${testId}`;
            const raw = localStorage.getItem(key);
            const list = raw ? JSON.parse(raw) : [];
            setAttempts(Array.isArray(list) ? list : []);
        } catch (e) {
            console.error("Не удалось загрузить попытки теста", e);
            setAttempts([]);
        }
    }, [testId]);

    const handleBack = () => {
        navigate("/tests");
    };

    const handleOpenDetails = (attempt) => {
        setSelectedAttempt(attempt);
    };

    const handleCloseDetails = () => {
        setSelectedAttempt(null);
    };

    return (
        <div className="tests-page">
            <div
                className="test-page"
                style={{ position: "absolute", left: "1430px", top: "0px" }}
            >
                <LogoutButton />
            </div>
            <div className="create-wrapper2">
                <div className="test2">
                    <div className="stat-top-bar2">
                        <button className="stat-back-btn2" onClick={handleBack}>
                            <BackIcon />
                        </button>
                        <h1>Статистика теста</h1>

                    </div>
                    <div className="tests-line"></div>
                    {attempts.length === 0 ? (
                        <p className="stat-empty">
                            По этому тесту ещё нет попыток.
                        </p>
                    ) : (
                        <div className="stat-attempts-table">
                            <table>
                                <thead>
                                <tr>
                                    <th>Участник</th>
                                    <th>Результат</th>
                                    <th>Время прохождения</th>
                                    <th>Подробная статистика</th>
                                </tr>
                                </thead>
                                <tbody>
                                {attempts.map((a) => (
                                    <tr key={a.id}>
                                        <td className="stat-cell-name">{a.userName}</td>
                                        <td className="stat-cell-score">
                                            {a.passed ? "Пройден" : "Не пройден"}
                                        </td>

                                        <td className="stat-cell-time">
                                            {a.durationMinutes != null
                                                ? `${a.durationMinutes} мин`
                                                : ""}
                                        </td>
                                        <td className="stat-cell-button">
                                            <button
                                                className="stat-details-btn"
                                                onClick={() => handleOpenDetails(a)}
                                            >
                                                Открыть подробную статистику
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                                </tbody>
                            </table>
                        </div>

                    )}

                    {selectedAttempt && (
                        <div className="stat-modal-overlay">
                            <div className="stat-modal">
                                <h3>Подробная статистика</h3>

                                <div className="stat-details-user">
                                    <p>
                                        <strong>Участник:</strong>{" "}
                                        {selectedAttempt.userName}
                                    </p>
                                    <p>
                                        <strong>Почта:</strong>{" "}
                                        {selectedAttempt.userEmail}
                                    </p>
                                    <p>
                                        <strong>Результат теста:</strong>{" "}
                                        {selectedAttempt.score}/{selectedAttempt.totalMax}
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
                                        {(selectedAttempt.perQuestion || []).map((q) => (
                                            <tr key={q.questionIndex}>
                                                <td>{q.questionText}</td>
                                                <td>
                                                    {q.score}/{q.maxScore}
                                                </td>
                                            </tr>
                                        ))}
                                        </tbody>
                                    </table>
                                </div>

                                <button
                                    className="stat-hide-btn"
                                    onClick={handleCloseDetails}
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
