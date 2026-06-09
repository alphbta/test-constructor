import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { testsAPI } from "../services/api";
import "../styles/TestPreviewPage.css";
import LogoutButton from "../components/LogoutButton.jsx";

export default function TestPreviewPage() {
    const { test_link } = useParams();
    const navigate = useNavigate();

    const [test, setTest] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchTest = async () => {
            try {
                setLoading(true);
                setError(null);
                
                // Загружаем тест с бэкенда
                const response = await testsAPI.startAttempt(test_link);
                
                if (response.data) {
                    // Сохраняем тест в localStorage для дальнейшего использования
                    const key = `shared_test_${test_link}`;
                    localStorage.setItem(key, JSON.stringify(response.data));
                    setTest(response.data);
                } else {
                    setError("Тест не найден");
                }
            } catch (err) {
                console.error("Ошибка загрузки теста:", err);
                setError("Не удалось загрузить тест. Проверьте ссылку.");
            } finally {
                setLoading(false);
            }
        };

        if (test_link) {
            fetchTest();
        }
    }, [test_link]);

    const handleStartTest = () => {
        // Переходим на страницу решения теста
        navigate(`/test/${test_link}`);
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

    if (loading) {
        return (
            <div className="test-preview-page">
                <div className="test-preview-loading">
                    <p>Загрузка теста...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="test-preview-page">
                <div className="test-preview-error">
                    <p>{error}</p>
                    <button 
                        className="test-preview-back-btn"
                        onClick={() => navigate("/")}
                    >
                        Вернуться на главную
                    </button>
                </div>
            </div>
        );
    }

    if (!test) {
        return (
            <div className="test-preview-page">
                <div className="test-preview-error">
                    <p>Тест не найден</p>
                </div>
            </div>
        );
    }

    return (
        <div className="test-preview-page">
            <div className="test-preview-header">
                <LogoutButton />
            </div>

            <div className="test-preview-wrapper">
                <div className="test-preview-container">
                    {/* Заголовок */}
                    <div className="preview-title-section">
                        <h1 className="preview-title">
                            Перед началом тестирования
                        </h1>
                        <p className="preview-subtitle">
                            Ознакомьтесь с информацией о тесте и приступите к тестированию.
                        </p>
                    </div>

                    {/* Карточка теста */}
                    <div className="preview-test-card">
                        {/* Левая часть с информацией */}
                        <div className="preview-card-left">
                            <div className="preview-card-icon">
                                <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
                                    <circle cx="24" cy="24" r="20" fill="currentColor" opacity="0.2"/>
                                    <path d="M24 14C18.48 14 14 18.48 14 24C14 29.52 18.48 34 24 34C29.52 34 34 29.52 34 24C34 18.48 29.52 14 24 14ZM24 31C19.59 31 16 27.41 16 24C16 20.59 19.59 17 24 17C28.41 17 32 20.59 32 24C32 27.41 28.41 31 24 31ZM23 20H25V26H23V20Z" fill="currentColor"/>
                                </svg>
                            </div>
                            
                            <div className="preview-card-info">
                                <div className="preview-info-section">
                                    <span className="preview-info-label">НАЗВАНИЕ</span>
                                    <h2 className="preview-info-title">
                                        {test.title || test.name || "Название теста"}
                                    </h2>
                                </div>

                                {test.description && (
                                    <div className="preview-info-section">
                                        <span className="preview-info-label">ОПИСАНИЕ</span>
                                        <p className="preview-info-description">
                                            {test.description}
                                        </p>
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* Разделитель */}
                        <div className="preview-card-divider"></div>

                        {/* Правая часть с временем */}
                        <div className="preview-card-right">
                            <div className="preview-time-section">
                                <div className="preview-time-icon">
                                    <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
                                        <circle cx="20" cy="20" r="16" fill="currentColor" opacity="0.2"/>
                                        <path d="M20 8C13.37 8 8 13.37 8 20C8 26.63 13.37 32 20 32C26.63 32 32 26.63 32 20C32 13.37 26.63 8 20 8ZM20 29C14.48 29 10 24.52 10 20C10 15.48 14.48 11 20 11C25.52 11 30 15.48 30 20C30 24.52 25.52 29 20 29ZM20.5 14H19V21L25.2 24.5L26 23.16L20.5 20.25V14Z" fill="currentColor"/>
                                    </svg>
                                </div>
                                <div className="preview-time-content">
                                    <span className="preview-time-label">ВРЕМЯ ПРОХОЖДЕНИЯ</span>
                                    <p className="preview-time-value">
                                        {formatTime(test.completetime || test.time_limit || 0)}
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Дополнительная информация */}
                    {test.questions && test.questions.length > 0 && (
                        <div className="preview-additional-info">
                            <div className="info-item">
                                <span className="info-icon">❓</span>
                                <span className="info-text">
                                    Количество вопросов: <strong>{test.questions.length}</strong>
                                </span>
                            </div>
                            {test.maxScore && (
                                <div className="info-item">
                                    <span className="info-icon">⭐</span>
                                    <span className="info-text">
                                        Максимальный балл: <strong>{test.maxScore}</strong>
                                    </span>
                                </div>
                            )}
                        </div>
                    )}

                    {/* Информационное сообщение */}
                    <div className="preview-notice">
                        <span className="notice-icon">ℹ️</span>
                        <p className="notice-text">
                            После начала тестирования отсчёт времени начнётся автоматически.
                        </p>
                    </div>

                    {/* Кнопка для начала теста */}
                    <div className="preview-button-section">
                        <button 
                            className="preview-start-btn"
                            onClick={handleStartTest}
                        >
                            Перейти к тестированию
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
