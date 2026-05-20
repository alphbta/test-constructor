import React, { useState, useRef, useEffect } from 'react';

export default function SpecializationSelect({ specializations, selected, onChange }) {
    const [isOpen, setIsOpen] = useState(false);
    const containerRef = useRef(null);
    const listRef = useRef(null);

    const selectedSpec = specializations.find(s => String(s.id) === String(selected));
    const selectedName = selectedSpec ? selectedSpec.name : 'Выберите специализацию';

    // Закрытие при клике вне компонента — используем мгновенное закрытие
    useEffect(() => {
        const handleClickOutside = (e) => {
            if (containerRef.current && !containerRef.current.contains(e.target)) {
                closeInstant();
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    // Функция плавного открытия
    const openAnimated = () => {
        if (!listRef.current) {
            setIsOpen(true);
            return;
        }

        setIsOpen(true); // добавляем класс .open (если вы используете его для стилей)
        // Убедимся, что у списка есть transition для анимации открытия
        listRef.current.style.transition = 'max-height 0.35s ease, opacity 0.25s ease';
        // Установим maxHeight в следующем кадре, чтобы сработал переход
        requestAnimationFrame(() => {
            listRef.current.style.maxHeight = `${listRef.current.scrollHeight}px`;
            listRef.current.style.opacity = '1';
        });
    };

    // Функция мгновенного (без анимации) закрытия
    const closeInstant = () => {
        if (!listRef.current) {
            setIsOpen(false);
            return;
        }

        // Отключаем transition, чтобы изменение высоты было мгновенным
        listRef.current.style.transition = 'none';
        listRef.current.style.maxHeight = '0px';
        listRef.current.style.opacity = '0';
        // Убираем класс open (чтобы border-radius и пр. стали исходными)
        setIsOpen(false);

        // Восстанавливаем transition для следующего открытия (делаем это в следующем тике)
        // Используем setTimeout 0, чтобы браузер применил изменения
        setTimeout(() => {
            if (listRef.current) {
                listRef.current.style.transition = '';
            }
        }, 0);
    };

    const toggleOpen = () => {
        if (isOpen) {
            closeInstant();
        } else {
            openAnimated();
        }
    };

    const handleSelect = (specId) => {
        onChange(String(specId));
        // После выбора закрываем мгновенно по вашему желанию
        closeInstant();
    };

    // На ресайзе пересчитываем maxHeight если открыт (чтобы высота соответствовала содержимому)
    useEffect(() => {
        const handleResize = () => {
            if (isOpen && listRef.current) {
                listRef.current.style.maxHeight = `${listRef.current.scrollHeight}px`;
            }
        };
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, [isOpen]);

    // При монтировании присвоим начальные стили (на случай SSR/HMR)
    useEffect(() => {
        if (listRef.current) {
            listRef.current.style.maxHeight = '0px';
            listRef.current.style.opacity = '0';
            listRef.current.style.overflow = 'hidden';
        }
    }, []);

    return (
        <div
            className={`specialization-select ${isOpen ? 'open' : ''}`}
            ref={containerRef}
        >
            <button
                type="button"
                className="specialization-select-button"
                onClick={toggleOpen}
                aria-expanded={isOpen}
                aria-haspopup="listbox"
            >
                <span className="specialization-select-value">{selectedName}</span>
                <div className="select-arrow" aria-hidden="true" />
            </button>

            <div
                className="specialization-select-list"
                ref={listRef}
                role="listbox"
            >
                <div className="specialization-select-options">
                    {specializations.length === 0 ? (
                        <div className="specialization-option disabled">Нет специализаций</div>
                    ) : (
                        specializations.map(spec => (
                            <button
                                key={spec.id}
                                type="button"
                                className={`specialization-option ${String(spec.id) === String(selected) ? 'selected' : ''}`}
                                onClick={() => handleSelect(spec.id)}
                                role="option"
                                aria-selected={String(spec.id) === String(selected)}
                            >
                                {spec.name}
                            </button>
                        ))
                    )}
                </div>
            </div>
        </div>
    );
}
