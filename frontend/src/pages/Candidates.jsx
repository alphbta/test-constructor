import "../styles/tests.css";
import "../styles/candidates.css";
import LogoutButton from "../components/LogoutButton.jsx";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";

import TaskIcon from "../assets/task.svg?react";
import EventIcon from "../assets/event.svg?react";
import CandidatesIcon from "../assets/Candidates.svg?react";

export default function Candidates() {
    const navigate = useNavigate();
    const [candidates, setCandidates] = useState([]);
    const [filteredCandidates, setFilteredCandidates] = useState([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [loading, setLoading] = useState(true);

    // Mock данные для тестирования
    const mockCandidates = [
        { id: 1, name: "Иван", surname: "Иванов", email: "ivan@example.com" },
        { id: 2, name: "Мария", surname: "Петрова", email: "maria@example.com" },
        { id: 3, name: "Петр", surname: "Сидоров", email: "petr@example.com" },
        { id: 4, name: "Анна", surname: "Кузнецова", email: "anna@example.com" },
        { id: 5, name: "Алексей", surname: "Волков", email: "alex@example.com" },
    ];

    useEffect(() => {
        const fetchCandidates = async () => {
            try {
                const token = localStorage.getItem("token");

                // Используем mock данные если нет токена или для тестирования
                if (!token) {
                    setCandidates(mockCandidates);
                    setFilteredCandidates(mockCandidates);
                    setLoading(false);
                    return;
                }

                const response = await fetch(
                    "http://localhost:8080/api/manager/users",
                    {
                        headers: {
                            "Authorization": `Bearer ${token}`
                        }
                    }
                );

                if (response.ok) {
                    const data = await response.json();
                    const users = data.users || [];
                    setCandidates(users);
                    setFilteredCandidates(users);
                } else {
                    // При ошибке используем mock
                    console.warn("Ошибка загрузки кандидатов, используются тестовые данные");
                    setCandidates(mockCandidates);
                    setFilteredCandidates(mockCandidates);
                }
            } catch (error) {
                console.error("Ошибка при загрузке кандидатов:", error);
                // При ошибке сети используем mock
                setCandidates(mockCandidates);
                setFilteredCandidates(mockCandidates);
            } finally {
                setLoading(false);
            }
        };

        fetchCandidates();
    }, [navigate]);

    const handleSearch = (query) => {
        setSearchQuery(query);
        const filtered = candidates.filter((candidate) => {
            const fullName = `${candidate.name} ${candidate.surname}`.toLowerCase();
            return fullName.includes(query.toLowerCase());
        });
        setFilteredCandidates(filtered);
    };

    const handleCandidateClick = (candidateId) => {
        navigate(`/candidates/${candidateId}`);
    };

    return (
        <div className="tests-page">
            <>
                <LogoutButton />
            </>
            <div className="tests-wrapper">
                <div className="tests">
                    {/* Навигационные вкладки */}
                    <div className="tests-tabs">
                        <button
                            className="tab-btn"
                            onClick={() => navigate("/tests")}
                        >
                            <TaskIcon />
                            Тестовые задания
                        </button>
                        <button
                            className="tab-btn"
                            onClick={() => navigate("/events")}
                        >
                            <EventIcon />
                            Мероприятия
                        </button>
                        <button
                            className="tab-btn tab-btn-active"
                            onClick={() => navigate("/candidates")}
                        >
                            <CandidatesIcon />
                            Кандидаты
                        </button>
                    </div>


                    <div className="candidates-search-box">
                        <input
                            type="text"
                            placeholder="Поиск по названию кандидата..."
                            value={searchQuery}
                            onChange={(e) => handleSearch(e.target.value)}
                            className="candidates-search-input"
                        />
                    </div>

                    <div className="candidates-list">
                        {loading ? (
                            <p className="candidates-loading">Загрузка...</p>
                        ) : filteredCandidates.length === 0 ? (
                            <p className="candidates-empty">
                                Кандидатов не найдено
                            </p>
                        ) : (
                            filteredCandidates.map((candidate) => (
                                <div
                                    key={candidate.id}
                                    className="candidate-card"
                                    onClick={() => handleCandidateClick(candidate.id)}
                                >
                                    <div className="candidate-avatar">
                                        {candidate.name.charAt(0)}
                                        {candidate.surname.charAt(0)}
                                    </div>
                                    <div className="candidate-info">
                                        <p className="candidate-name">
                                            {candidate.name} {candidate.surname}
                                        </p>
                                    </div>
                                    <div className="candidate-arrow">
                                        ›
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}