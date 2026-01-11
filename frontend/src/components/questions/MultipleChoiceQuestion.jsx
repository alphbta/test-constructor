import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import DeleteIcon from "../../assets/delete.svg?react";
import DeleteIconSub from "../../assets/delete_sub.svg?react";

function MultipleChoiceQuestion({ question, updateQuestion, deleteQuestion }) {
    const {
        attributes,
        listeners,
        setNodeRef,
        transform,
        transition,
        isDragging,
    } = useSortable({ id: question.id });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.5 : 1,
    };

    const addOption = () => {
        const newOptions = [...question.options, { text: "", isCorrect: false }];
        updateQuestion(question.id, "options", newOptions);
    };

    const updateOption = (index, field, value) => {
        const newOptions = [...question.options];
        newOptions[index][field] = value;
        updateQuestion(question.id, "options", newOptions);
    };

    const deleteOption = (index) => {
        const newOptions = question.options.filter((_, i) => i !== index);
        updateQuestion(question.id, "options", newOptions);
    };

    return (
        <div ref={setNodeRef} style={style} className="question-block multiple-choice">
            <div className = "">
                <div className = "question-up">
                    <span {...attributes} {...listeners} className="drag-handle">
                        <div style={{lineHeight: '0.2'}}>
                            <div>···</div>
                            <div style={{marginTop: '2px'}}>···</div>
                        </div>
                    </span>

                    <div className="q-icons">
                            <span onClick={() => deleteQuestion(question.id)}>
                                <DeleteIcon style={{ width: '24px', height: '24px' }}/>
                            </span>
                    </div>
                </div>
                <div className="q-header">
                <span>

                    {question.order}. <input
                    className="q-text-input"
                    placeholder="Введите текст вопроса..."
                    value={question.text}
                    onChange={(e) => updateQuestion(question.id, "text", e.target.value)}
                />
                </span>
            </div>
            </div>
            <div className="options-list">
                {question.options?.map((option, index) => (
                    <div key={index} className="options-row">
                        <label className="option-label">
                            <input
                                type="checkbox"
                                checked={option.isCorrect}
                                onChange={(e) => updateOption(index, "isCorrect", e.target.checked)}
                            />

                            <input
                                type="text"
                                className="answer-input"
                                placeholder="Введите вариант..."
                                value={option.text}
                                onChange={(e) => updateOption(index, "text", e.target.value)}
                            />
                        </label>
                        <button
                            className="delete-answer-btn"
                            onClick={() => deleteOption(index)}
                        >
                            <DeleteIconSub  />
                        </button>
                    </div>
                ))}
            </div>

            <button className="add-answer-btn" onClick={addOption}>
                + Добавить вариант
            </button>

            <div className="settings">
                <div className="setting-title">Настройки</div>
                <div className="setting-row">
                    <span>Зависимость баллов от % верных ответов</span>
                </div>
                <div className="setting-row">
                    <label>
                        <input
                            type="checkbox"
                            name={`scoring-${question.id}`}
                            checked={question.scoringType === "allOrNothing"}
                            onChange={() => updateQuestion(question.id, "scoringType", "allOrNothing")}
                        />
                        Только 100% или 0
                    </label>
                </div>
                <div className="setting-row">
                    <label>
                        <input
                            type="checkbox"
                            name={`scoring-${question.id}`}
                            checked={question.scoringType === "partial" || false}
                            onChange={() => updateQuestion(question.id, "scoringType", "partial")}
                        />
                        Баллы за частично верные ответы
                    </label>
                </div>
            </div>

            <div className="score-section">
                Баллы за верный ответ :{" "}
                <span className="score-container">
                    <input type="number" className="score-input"
                           value={question.maxScore || 0}
                           onChange={(e) => updateQuestion(question.id, "maxScore", parseInt(e.target.value) || 0)}
                    />{" "}
                    б
                </span>
            </div>
        </div>
    );
}

export default MultipleChoiceQuestion;