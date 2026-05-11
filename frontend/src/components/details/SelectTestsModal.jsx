import React from 'react';
import './modal.css';

export default function SelectTestsModal({ open, tests, selected, onSelect, onApply, onClose }) {
    if (!open) return null;
    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <h3>Выберите тесты</h3>
                <ul className="modal-tests-list">
                    {tests.map(test => (
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
                <div className="modal-btns">
                    <button onClick={onApply}>Добавить выбранные</button>
                    <button onClick={onClose}>Закрыть</button>
                </div>
            </div>
        </div>
    );
}
