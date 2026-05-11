import React, { useState, useEffect } from 'react';
import SelectTestsModal from './details/SelectTestsModal';
import CriteriaTable from './details/CriteriaTable';
import TimeBox from './details/TimeBox';
import ShareLinkBox from './details/ShareLinkBox';
import SpecializationSelect from './details/SpecializationSelect';
import '../styles/event-config.css';

// Моки для специализаций (замени на загрузку из API, если нужно)
const allSpecsMock = [
    { id: 1, name: 'Frontend' },
    { id: 2, name: 'Backend' },
];

export default function EventConfigPage() {
    // Состояния
    const [tests, setTests] = useState([]);
    const [selectedTests, setSelectedTests] = useState([]);
    const [criteria, setCriteria] = useState([
        { threshold: 50, message: 'Успешно пройден', extraTests: [] },
        { threshold: 30, message: 'Пройдите дополнительный тест', extraTests: [] },
        { threshold: 25, message: 'Пройдите дополнительный тест', extraTests: [] },
    ]);
    const [modalOpen, setModalOpen] = useState(false);
    const [modalTarget, setModalTarget] = useState(null); // null | 'main' | index
    const [modalSelected, setModalSelected] = useState([]);
    const [specializations, setSpecializations] = useState(allSpecsMock);
    const [selectedSpec, setSelectedSpec] = useState('');
    const [failMessage, setFailMessage] = useState('');
    const [time, setTime] = useState({ hours: 0, minutes: 0, seconds: 0 });
    const [shareLink, setShareLink] = useState('https://newforms-novaya-forma-konstruktion');

    // Загрузка тестов из API при монтировании
    useEffect(() => {
        // Пример через fetch:
        fetch('/api/tests') // замени на свой путь к API
            .then(res => res.json())
            .then(data => setTests(data))
            .catch(err => {
                setTests([]); // если ошибка, пусть будет пусто
                console.error('Ошибка загрузки тестов:', err);
            });
    }, []);

    // Модалка для выбора тестов
    const openModal = (target) => {
        setModalTarget(target);
        if (target === 'main') setModalSelected(selectedTests);
        else setModalSelected(criteria[target].extraTests || []);
        setModalOpen(true);
    };
    const handleApplyModal = () => {
        if (modalTarget === 'main') setSelectedTests(modalSelected);
        else {
            setCriteria(criteria.map((row, idx) =>
                idx === modalTarget ? { ...row, extraTests: modalSelected } : row
            ));
        }
        setModalOpen(false);
    };
    const handleCriteriaChange = (idx, newRow) => {
        setCriteria(criteria.map((row, i) => (i === idx ? newRow : row)));
    };
    const handleAddCriteria = () => {
        setCriteria([...criteria, { threshold: 0, message: '', extraTests: [] }]);
    };

    return (
        <div className="event-config-page">
            {/* Левая панель */}
            <div className="event-config-sidebar">
                <button className="add-tests-btn" onClick={() => openModal('main')}>Добавить тесты</button>
                <ul className="event-config-tests-list">
                    {selectedTests.map(id => {
                        const test = tests.find(t => t.id === id);
                        return test ? (
                            <li key={id} className="event-config-test-item">
                                {test.title}
                            </li>
                        ) : null;
                    })}
                </ul>
            </div>
            {/* Центральная панель */}
            <div className="event-config-main">
                <SpecializationSelect
                    specializations={specializations}
                    selected={selectedSpec}
                    onChange={setSelectedSpec}
                />
                <div className="criteria-table-title">Критерий прохождения теста</div>
                <CriteriaTable
                    criteria={criteria}
                    onChange={handleCriteriaChange}
                    onAdd={handleAddCriteria}
                    onAddTest={idx => openModal(idx)}
                />
                <div className="fail-message-block">
                    <p>Сообщение при провальном прохождении</p>
                    <input
                        type="text"
                        placeholder="Введите текст сообщения при провальном прохождении..."
                        value={failMessage}
                        onChange={e => setFailMessage(e.target.value)}
                    />
                </div>
                <TimeBox time={time} setTime={setTime} />
                <ShareLinkBox link={shareLink} />
                <button className="save-btn">Сохранить</button>
            </div>
            {/* Модальное окно */}
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
