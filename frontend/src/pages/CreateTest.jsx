import { useState, useEffect } from "react";
import {
    DndContext,
    closestCenter,
    KeyboardSensor,
    PointerSensor,
    useSensor,
    useSensors,
} from "@dnd-kit/core";
import {
    arrayMove,
    SortableContext,
    sortableKeyboardCoordinates,
    verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useNavigate, useLocation } from "react-router-dom";
import "../styles/createTest.css";

function SortableQuestion({ question, updateQuestion, deleteQuestion }) {
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

    return (
        <div ref={setNodeRef} style={style} className="question-block">
            <div className="q-header">
                <span>
                    <span {...attributes} {...listeners} className="drag-handle">
                        ⋮⋮
                    </span>{" "}
                    {question.order}. Введите вопрос
                </span>
                <div className="q-icons">
                    <span onClick={() => deleteQuestion(question.id)}>🗑️</span>
                </div>
            </div>

            <input
                className="q-input"
                placeholder="Текст вопроса..."
                value={question.text}
                onChange={(e) =>
                    updateQuestion(question.id, "text", e.target.value)
                }
            />

            {question.type === "yesNo" && (
                <div className="q-options">
                    <label>
                        <input type="radio" name={`yesNo-${question.id}`} /> Да
                    </label>
                    <label>
                        <input type="radio" name={`yesNo-${question.id}`} /> Нет
                    </label>
                </div>
            )}

            {question.type === "multipleChoice" && (
                <div className="q-options">
                    {question.options.map((option, idx) => (
                        <label key={idx}>
                            <input type="checkbox" />{" "}
                            <input
                                type="text"
                                placeholder={`Вариант ${idx + 1}`}
                                value={option}
                                onChange={(e) => {
                                    const newOptions = [...question.options];
                                    newOptions[idx] = e.target.value;
                                    updateQuestion(question.id, "options", newOptions);
                                }}
                            />
                        </label>
                    ))}
                    <button
                        className="add-option-btn"
                        onClick={() =>
                            updateQuestion(question.id, "options", [
                                ...question.options,
                                "",
                            ])
                        }
                    >
                        + Добавить вариант
                    </button>
                </div>
            )}

            {question.type === "singleChoice" && (
                <div className="q-options">
                    {question.options.map((option, idx) => (
                        <label key={idx}>
                            <input type="radio" name={`single-${question.id}`} />{" "}
                            <input
                                type="text"
                                placeholder={`Вариант ${idx + 1}`}
                                value={option}
                                onChange={(e) => {
                                    const newOptions = [...question.options];
                                    newOptions[idx] = e.target.value;
                                    updateQuestion(question.id, "options", newOptions);
                                }}
                            />
                        </label>
                    ))}
                    <button
                        className="add-option-btn"
                        onClick={() =>
                            updateQuestion(question.id, "options", [
                                ...question.options,
                                "",
                            ])
                        }
                    >
                        + Добавить вариант
                    </button>
                </div>
            )}

            {question.type === "shortText" && (
                <div className="q-options">
                    <input
                        type="text"
                        className="short-text-input"
                        placeholder="Короткий ответ..."
                        disabled
                    />
                </div>
            )}

            {question.type === "longText" && (
                <>
                    <textarea
                        className="q-textarea"
                        placeholder="Длинный ответ..."
                        disabled
                    />
                    <div className="q-mode">
                        <span>Режим проверки</span>
                        <label>
                            <input
                                type="radio"
                                name={`mode-${question.id}`}
                                checked={question.checkMode === "auto"}
                                onChange={() =>
                                    updateQuestion(question.id, "checkMode", "auto")
                                }
                            />{" "}
                            Автоматическая
                        </label>
                        <label>
                            <input
                                type="radio"
                                name={`mode-${question.id}`}
                                checked={question.checkMode === "manual"}
                                onChange={() =>
                                    updateQuestion(question.id, "checkMode", "manual")
                                }
                            />{" "}
                            Ручная
                        </label>
                    </div>
                </>
            )}

            <div className="q-score">
                Максимальный балл:{" "}
                <input
                    type="number"
                    className="score-input"
                    value={question.maxScore}
                    onChange={(e) =>
                        updateQuestion(question.id, "maxScore", parseInt(e.target.value) || 0)
                    }
                />
            </div>
        </div>
    );
}

function useAppSensors() {
    const pointerSensor = useSensor(PointerSensor);
    const keyboardSensor = useSensor(KeyboardSensor, {
        coordinateGetter: sortableKeyboardCoordinates,
    });

    return useSensors(pointerSensor, keyboardSensor);
}

export default function CreateTest() {
    const navigate = useNavigate();
    const location = useLocation();

    const isEditing = location.state?.editing || false;
    const editingTest = location.state?.test || null;

    const [title, setTitle] = useState(isEditing ? editingTest.title : "");
    const [description, setDescription] = useState(isEditing ? editingTest.description : "");
    const [questions, setQuestions] = useState(
        isEditing ? editingTest.questions.map((q, idx) => ({
            id: `q-${idx}-${Date.now()}`,
            order: idx + 1,
            type: q.type,
            text: q.text,
            options: q.options || (q.type === "yesNo" ? ["Да", "Нет"] : ["", ""]),
            maxScore: q.maxScore || 10,
            checkMode: q.checkMode || (q.type === "longText" ? "manual" : "auto"),
        })) : [
            {
                id: "1",
                order: 1,
                type: "yesNo",
                text: "Пример вопроса Да/Нет?",
                options: ["Да", "Нет"],
                maxScore: 15,
                checkMode: "auto",
            },
        ]
    );

    const sensors = useAppSensors();

    const addQuestion = (type) => {
        const newQuestion = {
            id: Date.now().toString(),
            order: questions.length + 1,
            type,
            text: "",
            options: type === "yesNo" ? ["Да", "Нет"] : ["", ""],
            maxScore: 10,
            checkMode: type === "longText" ? "manual" : "auto",
        };
        setQuestions([...questions, newQuestion]);
    };

    const updateQuestion = (id, field, value) => {
        setQuestions(
            questions.map((q) =>
                q.id === id ? { ...q, [field]: value } : q
            )
        );
    };

    const deleteQuestion = (id) => {
        setQuestions(questions.filter((q) => q.id !== id));
    };

    const handleDragEnd = (event) => {
        const { active, over } = event;

        if (over && active.id !== over.id) {
            setQuestions((items) => {
                const oldIndex = items.findIndex((item) => item.id === active.id);
                const newIndex = items.findIndex((item) => item.id === over.id);
                const newItems = arrayMove(items, oldIndex, newIndex);

                return newItems.map((item, idx) => ({
                    ...item,
                    order: idx + 1,
                }));
            });
        }
    };

    const handleSave = () => {
        if (!title.trim()) {
            alert("Введите название теста!");
            return;
        }

        const testData = {
            id: isEditing ? editingTest.id : Date.now().toString(),
            title: title.trim(),
            description: description.trim(),
            createdAt: isEditing ? editingTest.createdAt : new Date().toISOString(),
            questions: questions.map((q) => ({
                type: q.type,
                text: q.text,
                options: q.options,
                maxScore: q.maxScore,
                checkMode: q.checkMode,
            })),
        };

        const existingTests = JSON.parse(localStorage.getItem("savedTests")) || [];

        let updatedTests;
        if (isEditing) {
            updatedTests = existingTests.map(test =>
                test.id === editingTest.id ? testData : test
            );
        } else {
            updatedTests = [...existingTests, testData];
        }

        localStorage.setItem("savedTests", JSON.stringify(updatedTests));

        alert(isEditing ? "Тест успешно обновлён!" : "Тест успешно сохранён!");

        navigate("/tests");
    };

    const questionTypes = [
        { key: "shortText", label: "Ввод короткого текста" },
        { key: "longText", label: "Ввод длинного текста" },
        { key: "singleChoice", label: "Одиночный выбор" },
        { key: "yesNo", label: "Да / Нет" },
        { key: "multipleChoice", label: "Множественный выбор" },
        { key: "", label: "Выпадающий список"},
        { key: "", label: "Выбор картинки"},
    ];

    return (
        <div className="create-wrapper">
            {/* Левая часть */}
            <div className="create-left">
                <input
                    className="test-title"
                    placeholder="Название теста"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                />
                <input
                    className="test-desc"
                    placeholder="Описание теста"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                />

                <DndContext
                    sensors={sensors}
                    collisionDetection={closestCenter}
                    onDragEnd={handleDragEnd}
                >
                    <SortableContext
                        items={questions.map((q) => q.id)}
                        strategy={verticalListSortingStrategy}
                    >
                        {questions.map((question) => (
                            <SortableQuestion
                                key={question.id}
                                question={question}
                                updateQuestion={updateQuestion}
                                deleteQuestion={deleteQuestion}
                            />
                        ))}
                    </SortableContext>
                </DndContext>

            </div>

            {/* Правая панель */}
            <div className="create-right">
                <button className="save-btn" onClick={handleSave}> {isEditing ? "Сохранить" : "Сохранить"} </button>

                <h3>Поля теста</h3>

                <div className="right-section">
                    <p>Добавить новый вопрос ▼</p>
                    {questionTypes.map((type) => (
                        <button
                            key={type.key}
                            className="right-btn"
                            onClick={() => addQuestion(type.key)}
                        >
                            {type.label}
                        </button>
                    ))}
                </div>

                <div className="right-section">
                    <p>Установить время</p>
                    <div className="time-box">
                        <input type="number" min="0" placeholder="0" /> ч
                        <input type="number" min="0" max="59" placeholder="0" /> м
                        <input type="number" min="0" max="59" placeholder="0" /> с
                    </div>
                </div>
            </div>
        </div>
    );
}