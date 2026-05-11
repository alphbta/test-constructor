import React from 'react';

export default function CriteriaTable({ criteria, onChange, onAdd, onAddTest }) {
    return (
        <div className="criteria-table-block">
            <div className="criteria-table-header">
                <div>Нижняя граница</div>
                <div>Сообщение</div>
                <div>Дополнительный тест</div>
            </div>
            {criteria.map((row, idx) => (
                <div className="criteria-table-row" key={idx}>
                    <input
                        type="number"
                        value={row.threshold}
                        onChange={e => onChange(idx, { ...row, threshold: Number(e.target.value) })}
                        className="criteria-input"
                    />
                    <input
                        type="text"
                        value={row.message}
                        onChange={e => onChange(idx, { ...row, message: e.target.value })}
                        className="criteria-input"
                    />
                    <button className="criteria-add-test-btn" onClick={() => onAddTest(idx)}>
                        {row.extraTests && row.extraTests.length > 0
                            ? `Тестов: ${row.extraTests.length}`
                            : 'Добавить тесты'}
                    </button>
                </div>
            ))}
            <button className="criteria-add-btn" onClick={onAdd}>Добавить критерий</button>
        </div>
    );
}
