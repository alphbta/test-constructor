import React, { useState } from 'react';
import './modal.css';

export default function SelectTestsModal({ open, tests, selected, onSelect, onApply, onClose }) {
    const [searchTerm, setSearchTerm] = useState('');

    if (!open) return null;

    const filteredTests = tests.filter(test =>
        test.title.toLowerCase().includes(searchTerm.toLowerCase())
    );

    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <h3>
                    Выберите тесты
                    <button className="modal-close-btn" onClick={onClose}>✕</button>
                </h3>

                <div className="modal-search">
                    <input
                        type="text"
                        placeholder="Введите название теста"
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </div>
                <div className="modal-footer">
                <ul className="modal-tests-list">
                    {filteredTests.map(test => (
                        <li key={test.id}>
                            <label>
                                <input
                                    type="checkbox"
                                    checked={selected.includes(test.id)}
                                    onChange={() => onSelect(test.id)}
                                />
                                {test.title}
                            </label>
                        </li>
                    ))}
                </ul>
                </div>
                <div className="modal-btns">
                    <button onClick={onApply}>Сохранить</button>
                </div>
            </div>
        </div>
    );
}