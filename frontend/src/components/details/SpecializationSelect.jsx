import React from 'react';

export default function SpecializationSelect({ specializations, selected, onChange }) {
    return (
        <div className="specialization-select">
            <select value={selected || ''} onChange={e => onChange(e.target.value)}>
                <option value="">Выберите специализацию</option>
                {specializations.length === 0 && <option disabled>Нет специализаций</option>}
                {specializations.map(spec => (
                    <option key={spec.id} value={spec.id}>{spec.name}</option>
                ))}
            </select>
        </div>
    );
}
