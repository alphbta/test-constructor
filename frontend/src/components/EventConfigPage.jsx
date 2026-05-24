import React, { useState, useEffect } from 'react';
import SelectTestsModal from './details/SelectTestsModal';
import CriteriaTable from './details/CriteriaTable';
import TimeBox from './details/TimeBox';
import ShareLinkBox from './details/ShareLinkBox';
import SpecializationSelect from './details/SpecializationSelect';
import { testsAPI } from '../services/api.js';
import '../styles/event-config.css';
import back2 from '../assets/back2.svg';
import { useNavigate } from 'react-router-dom';
import plusIcon from '../assets/plus.svg';
import korzinaIcon from '../assets/korzina.svg';
import massageIcon from '../assets/message.svg';

const allSpecsMock = [
    { id: 1, name: 'Frontend' },
    { id: 2, name: 'Backend' },
];

const DEFAULT_CONFIG = {
    selectedSpec: '',
    criteria: [
        { threshold: 50, message: 'Успешно пройден', extraTests: [] },
        { threshold: 30, message: 'Пройдите дополнительный тест', extraTests: [] },
        { threshold: 25, message: 'Пройдите дополнительный тест', extraTests: [] },
    ],
    failMessage: '',
    time: { hours: 0, minutes: 0, seconds: 0 },
    shareLink: 'https://newforms-novaya-forma-konstruktion',
};

export default function EventConfigPage() {
    const navigate = useNavigate();

    const [tests, setTests] = useState([]);
    const [selectedTestIds, setSelectedTestIds] = useState([]);
    const [selectedTestId, setSelectedTestId] = useState(null);
    const [modalOpen, setModalOpen] = useState(false);
    const [modalTarget, setModalTarget] = useState(null);
    const [modalSelected, setModalSelected] = useState([]);
    const [specializations, setSpecializations] = useState(allSpecsMock);

    const [testConfigs, setTestConfigs] = useState({});

    useEffect(() => {
        const fetchTests = async () => {
            try {
                const response = await testsAPI.getTests();
                const data = response.data;

                let testsArray = [];
                if (Array.isArray(data)) {
                    testsArray = data;
                } else if (data.tests && Array.isArray(data.tests)) {
                    testsArray = data.tests;
                } else if (data.data && Array.isArray(data.data)) {
                    testsArray = data.data;
                } else {
                    console.error('Неизвестная структура ответа testsAPI.getTests():', data);
                }

                const normalized = testsArray.map(test => ({
                    ...test,
                    id: test.test_id || test.id,
                    creator_id: test.creator_id ?? test.creatorId ?? test.CreatorID ?? test.creatorID,
                    title: test.title || test.name || test.description || `Тест ${test.test_id || test.id}`,
                }));

                const userStr = localStorage.getItem('user');
                const currentUserId = userStr ? (JSON.parse(userStr).id) : null;

                const filtered = currentUserId != null
                    ? normalized.filter(t => Number(t.creator_id) === Number(currentUserId))
                    : normalized;

                setTests(filtered);
            } catch (err) {
                console.error('Ошибка загрузки тестов:', err);
                setTests([]);
            }
        };

        fetchTests();
    }, []);

    const getCurrentConfig = () => {
        if (!selectedTestId || !testConfigs[selectedTestId]) {
            return { ...DEFAULT_CONFIG };
        }
        return testConfigs[selectedTestId];
    };

    const updateCurrentConfig = (field, value) => {
        if (!selectedTestId) return;

        setTestConfigs(prev => ({
            ...prev,
            [selectedTestId]: {
                ...(prev[selectedTestId] || { ...DEFAULT_CONFIG }),
                [field]: value,
            }
        }));
    };

    const initTestConfig = (testId) => {
        if (!testConfigs[testId]) {
            setTestConfigs(prev => ({
                ...prev,
                [testId]: { ...DEFAULT_CONFIG }
            }));
        }
    };

    const openModal = (target) => {
        setModalTarget(target);
        if (target === 'main') {
            setModalSelected([...selectedTestIds]);
        } else {
            const currentConfig = getCurrentConfig();
            const extraTests = currentConfig.criteria && currentConfig.criteria[target]
                ? currentConfig.criteria[target].extraTests || []
                : [];
            setModalSelected([...extraTests]);
        }
        setModalOpen(true);
    };

    const handleApplyModal = () => {
        if (modalTarget === 'main') {
            setSelectedTestIds([...modalSelected]);
            modalSelected.forEach(testId => {
                if (!testConfigs[testId]) {
                    initTestConfig(testId);
                }
            });
            if (modalSelected.length > 0 && !selectedTestId) {
                setSelectedTestId(modalSelected[0]);
            }
        } else {
            const currentConfig = getCurrentConfig();
            const newCriteria = currentConfig.criteria.map((row, idx) =>
                idx === modalTarget ? { ...row, extraTests: [...modalSelected] } : row
            );
            updateCurrentConfig('criteria', newCriteria);
        }
        setModalOpen(false);
    };

    const handleCriteriaChange = (idx, newRow) => {
        const currentConfig = getCurrentConfig();
        const newCriteria = currentConfig.criteria.map((row, i) => (i === idx ? newRow : row));
        updateCurrentConfig('criteria', newCriteria);
    };

    const handleAddCriteria = () => {
        const currentConfig = getCurrentConfig();
        const newCriteria = [...currentConfig.criteria, { threshold: 0, message: '', extraTests: [] }];
        updateCurrentConfig('criteria', newCriteria);
    };

    const handleRemoveSelected = (idToRemove) => {
        const newSelected = selectedTestIds.filter(id => id !== idToRemove);
        setSelectedTestIds(newSelected);

        if (selectedTestId === idToRemove) {
            setSelectedTestId(newSelected.length > 0 ? newSelected[0] : null);
        }

        setTestConfigs(prev => {
            const newConfigs = { ...prev };
            delete newConfigs[idToRemove];
            return newConfigs;
        });
    };

    const handleDeleteCriteria = (index) => {
        const currentConfig = getCurrentConfig();
        const newCriteria = currentConfig.criteria.filter((_, i) => i !== index);
        updateCurrentConfig('criteria', newCriteria);
    };

    const handleDeleteTest = (criteriaIndex, testIndex) => {
        const currentConfig = getCurrentConfig();
        const newCriteria = currentConfig.criteria.map((item, i) => {
            if (i !== criteriaIndex) return item;
            return {
                ...item,
                extraTests: item.extraTests.filter((_, idx) => idx !== testIndex)
            };
        });
        updateCurrentConfig('criteria', newCriteria);
    };

    const currentConfig = getCurrentConfig();

    const handleSave = async () => {
        try {
            const userStr = localStorage.getItem('user');
            const currentUserId = userStr ? JSON.parse(userStr).id : null;

            if (!currentUserId) {
                alert('Ошибка: пользователь не авторизован');
                return;
            }

            for (const testId of selectedTestIds) {
                const config = testConfigs[testId];
                if (!config) continue;

                const timeInSeconds = (config.time?.hours || 0) * 3600 +
                    (config.time?.minutes || 0) * 60 +
                    (config.time?.seconds || 0);

                const payload = {
                    event_id: parseInt(eventId) || 1,
                    specialization_id: 1,
                    test_id: testId,
                    success_text: config.failMessage || 'Успешно пройден',
                    fail_text: config.failMessage || 'Не пройден',
                    time_limit: timeInSeconds,
                    threshold: config.criteria[0]?.threshold || 50,
                    extra_threshold: config.criteria.slice(1).map(c => ({
                        threshold: c.threshold,
                        message: c.message,
                        test_id: c.extraTests?.[0] || testId,
                        test_threshold: c.threshold,
                    })) || [],
                };

                console.log('Отправляю на бэкэнд:', payload);

                try {
                    const response = await testsAPI.createEventConfig(payload);
                    console.log('Ответ бэкэнда:', response.data);
                } catch (err) {
                    console.error(`Ошибка при сохранении теста ${testId}:`, err);
                }
            }

            alert('Конфигурация сохранена успешно!');
            navigate('/events');
        } catch (err) {
            console.error('Ошибка сохранения конфигурации:', err);
            alert('Ошибка при сохранении конфигурации');
        }
    };
    return (
        <div className="event-config-page">
            <div className="event-config-sidebar">
                <div className="event-config-header">
                    <button
                        type="button"
                        className="event-config-back-btn"
                        onClick={() => navigate('/events')}
                        aria-label="Вернуться к мероприятиям"
                    >
                        <img src={back2} alt="" className="event-config-back-icon" />
                    </button>
                    <p>Настройка тестов мероприятия</p>
                </div>

                <button className="add-tests-btn" onClick={() => openModal('main')}>
                    <span>Добавить тесты</span>
                    <img src={plusIcon} alt="Добавить" className="add-tests-plus" />
                </button>

                <ul className="event-config-tests-list">
                    {selectedTestIds.map(id => {
                        const test = tests.find(t => t.id === id);
                        const isActive = selectedTestId === id;
                        return test ? (
                            <li
                                key={id}
                                className={`event-config-test-item ${isActive ? 'active-red' : ''}`}
                                onClick={() => setSelectedTestId(id)}
                            >
                                <span className="test-title">{test.title}</span>
                                <button
                                    className="test-delete-btn"
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        handleRemoveSelected(id);
                                    }}
                                    aria-label={`Удалить тест ${test.title}`}
                                    type="button"
                                >
                                    <img src={korzinaIcon} alt="Удалить" />
                                </button>
                            </li>
                        ) : null;
                    })}
                </ul>
            </div>

            <div className="event-config-main">
                <SpecializationSelect
                    specializations={specializations}
                    selected={currentConfig.selectedSpec}
                    onChange={(spec) => updateCurrentConfig('selectedSpec', spec)}
                />
                <div className="criteria-table-title">Критерий прохождения теста</div>
                <CriteriaTable
                    criteria={currentConfig.criteria}
                    onChange={handleCriteriaChange}
                    onAdd={handleAddCriteria}
                    onAddTest={idx => openModal(idx)}
                    onDelete={handleDeleteCriteria}
                    onDeleteTest={handleDeleteTest}
                    testsList={tests}
                />
                <div className="fail-message-block">
                    <div className="fail-message-header">
                        <img src={massageIcon} alt="" style={{ width: '32px', height: '32px' }} />
                        <p className="fail-message-title">Сообщение при провальном прохождении</p>
                    </div>
                    <input
                        type="text"
                        placeholder="Введите текст сообщения при провальном прохождении..."
                        value={currentConfig.failMessage}
                        onChange={e => updateCurrentConfig('failMessage', e.target.value)}
                    />
                </div>
                <TimeBox
                    time={currentConfig.time}
                    setTime={(newTime) => updateCurrentConfig('time', newTime)}
                />
                <ShareLinkBox link={currentConfig.shareLink} />
                <button className="save-btn" onClick={handleSave}>Сохранить</button>
            </div>

            <SelectTestsModal
                open={modalOpen}
                tests={tests}
                selected={modalSelected}
                onSelect={id =>
                    setModalSelected(
                        modalSelected.includes(id)
                            ? modalSelected.filter(i => i !== id)
                            : [...modalSelected, id]
                    )
                }
                onApply={handleApplyModal}
                onClose={() => setModalOpen(false)}
            />
        </div>
    );
}