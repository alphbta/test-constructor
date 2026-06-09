import { useState } from "react";
import "../styles/questionTypeModal.css";

function QuestionTypeModal({ isOpen, onClose, onSelectType, position }) {
    if (!isOpen) return null;

    const questionTypes = [
        { key: "shortText", label: "Задания на ручной ввод" },
        { key: "singleChoice", label: "Одиночный выбор" },
        { key: "multipleChoice", label: "Множественный выбор" },
        { key: "matching", label: "На соотношение" },
        {
            key: "ordering",
            label: "На расположение в правильном порядке",
        },
    ];

    const handleTypeClick = (typeKey) => {
        onSelectType(typeKey);
        onClose();
    };

    return (
        <div className="modal-overlay" onClick={onClose}>
            <div
                className="modal-content"
                onClick={(e) => e.stopPropagation()}
                style={{
                    position: "absolute",
                    top: position?.top || "auto",
                    left: position?.left || "auto",
                }}
            >
                <div className="modal-header">
                    <h3>Выберите тип вопроса</h3>
                    <button className="modal-close" onClick={onClose}>
                        ×
                    </button>
                </div>
                <div className="modal-body">
                    {questionTypes.map((type) => (
                        <button
                            key={type.key}
                            className="modal-option"
                            onClick={() => handleTypeClick(type.key)}
                        >
                            {type.label}
                        </button>
                    ))}
                </div>
            </div>
        </div>
    );
}

export default QuestionTypeModal;
