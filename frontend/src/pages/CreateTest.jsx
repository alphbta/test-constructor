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
import { useNavigate, useLocation } from "react-router-dom";
import SortableQuestion from "../components/Question";
import PassingCriteria from "../components/questions/PassingCriteria.jsx";
import ResultMessages from "../components/questions/ResultMessages";
import "../styles/createTest.css";
import LogoutButton from "../components/LogoutButton.jsx";
import { testsAPI } from "../services/api.js";

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


    const [title, setTitle] = useState(
        isEditing ? editingTest?.title || "" : ""
    );

    const [description, setDescription] = useState(
        isEditing ? editingTest?.description || "" : ""
    );


    const [time, setTime] = useState(
        isEditing && editingTest?.complete_time
            ? {
                hours: Math.floor(editingTest.complete_time / 3600),
                minutes: Math.floor(
                    (editingTest.complete_time % 3600) / 60
                ),
                seconds: editingTest.complete_time % 60,
            }
            : {
                hours: 0,
                minutes: 0,
                seconds: 0,
            }
    );

    const [passingCriteria, setPassingCriteria] = useState(
        isEditing
            ? {
                type: editingTest?.is_percentage ? "percentage" : "points",
                percentage: editingTest?.is_percentage
                    ? editingTest.threshold
                    : 75,
                points: editingTest?.is_percentage
                    ? 0
                    : editingTest.threshold,
            }
            : {
                type: "percentage",
                percentage: 75,
                points: 0,
            }
    );

    const [resultMessages, setResultMessages] = useState(
        isEditing
            ? {
                success: editingTest?.success_text || "",
                failure: editingTest?.fail_text || "",
            }
            : {
                success: "",
                failure: "",
            }
    );

    const [questions, setQuestions] = useState(() => {
        if (isEditing && editingTest && Array.isArray(editingTest.questions)) {
            return editingTest.questions.map((q, idx) => {
                const base = {
                    id: `q-${idx}-${Date.now()}`,
                    order: idx + 1,
                    type: q.type || "shortText",
                    text: q.text || "",
                    maxScore: q.maxScore || q.points || 15,
                };

                switch (base.type) {
                    case "shortText":
                        return {
                            ...base,
                            correctAnswers:
                                q.correctAnswers ||
                                q.correct_input ||
                                [""],
                            caseSensitive:
                                q.caseSensitive !== undefined
                                    ? q.caseSensitive
                                    : q.case_sensitive || false,
                        };
                    case "singleChoice":
                        return {
                            ...base,
                            options:
                                q.options ||
                                q.choice ||
                                [{ text: "", isCorrect: false }],
                        };
                    case "multipleChoice":
                        return {
                            ...base,
                            options:
                                q.options ||
                                q.choice ||
                                [{ text: "", isCorrect: false }],
                            scoringType: q.scoringType || "allOrNothing",
                        };
                    case "matching":
                        return {
                            ...base,
                            rows:
                                q.rows ||
                                q.matching ||
                                [{ option: "", answer: "" }],
                        };
                    case "ordering":
                        return {
                            ...base,
                            items:
                                q.items ||
                                q.sequence ||
                                [{ text: "" }],
                        };
                    default:
                        return {
                            ...base,
                            correctAnswers: [""],
                            caseSensitive: false,
                        };
                }
            });
        }

        return [
            {
                id: "1",
                order: 1,
                type: "shortText",
                text: "",
                correctAnswers: [""],
                caseSensitive: false,
                maxScore: 15,
            },
        ];
    });

    const calculateTotalPoints = () => {
        return questions.reduce((sum, q) => sum + (q.maxScore || 0), 0);
    };

    const calculateCompleteTime = () => {
        return time.hours * 3600 + time.minutes * 60 + time.seconds;
    };

    const sensors = useAppSensors();

    const addQuestion = (type) => {
        const baseQuestion = {
            id: Date.now().toString(),
            order: questions.length + 1,
            type,
            text: "",
            maxScore: 15,
        };

        switch (type) {
            case "shortText":
                baseQuestion.correctAnswers = [""];
                baseQuestion.caseSensitive = false;
                break;
            case "singleChoice":
                baseQuestion.options = [{ text: "", isCorrect: false }];
                break;
            case "multipleChoice":
                baseQuestion.options = [{ text: "", isCorrect: false }];
                baseQuestion.scoringType = "allOrNothing";
                break;
            case "matching":
                baseQuestion.rows = [{ option: "", answer: "" }];
                break;
            case "ordering":
                baseQuestion.items = [{ text: "" }];
                break;
        }

        setQuestions([...questions, baseQuestion]);
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
                const oldIndex = items.findIndex(
                    (item) => item.id === active.id
                );
                const newIndex = items.findIndex(
                    (item) => item.id === over.id
                );
                const newItems = arrayMove(items, oldIndex, newIndex);

                return newItems.map((item, idx) => ({
                    ...item,
                    order: idx + 1,
                }));
            });
        }
    };

    const handleSave = async () => {
        if (!title.trim()) {
            alert("Введите название теста!");
            return;
        }

        const testData = {
            title: title.trim(),
            description: description.trim(),
            is_percentage: passingCriteria.type === "percentage",
            fail_text: resultMessages.failure || "",
            success_text: resultMessages.success || "",
            complete_time: calculateCompleteTime() || 3600,
            threshold:
                passingCriteria.type === "percentage"
                    ? passingCriteria.percentage
                    : passingCriteria.points,
            questions: questions.map((q, index) => {
                let options = {};
                let questionType;

                switch (q.type) {
                    case "singleChoice":
                        questionType = "single_choice";
                        options = {
                            choice:
                                q.options?.map((opt) => ({
                                    text: opt.text,
                                    is_true: opt.isCorrect,
                                })) || [],
                        };
                        break;
                    case "multipleChoice":
                        questionType = "multiple_choice";
                        options = {
                            choice:
                                q.options?.map((opt) => ({
                                    text: opt.text,
                                    is_true: opt.isCorrect,
                                })) || [],
                        };
                        break;
                    case "shortText":
                        questionType = "text_input";
                        options = {
                            correct_input: q.correctAnswers || [],
                            case_sensitive: q.caseSensitive || false,
                        };
                        break;
                    case "matching":
                        questionType = "matching";
                        options = {
                            matching:
                                q.rows?.map((row) => ({
                                    left: row.option,
                                    right: row.answer,
                                })) || [],
                        };
                        break;
                    case "ordering":
                        questionType = "correct_order";
                        options = {
                            sequence:
                                q.items?.map((item, itemIdx) => ({
                                    text: item.text,
                                    order: itemIdx + 1,
                                })) || [],
                        };
                        break;
                    default:
                        questionType = "text_input";
                }

                return {
                    text: q.text || "",
                    points: q.maxScore || 0,
                    type: questionType,
                    order_number: index + 1,
                    options: options,
                };
            }),
        };

        console.log(
            "Отправляемые данные на бэкенд:",
            JSON.stringify(testData, null, 2)
        );

        try {
            const token = localStorage.getItem("token");
            if (!token) {
                alert("Требуется авторизация!");
                navigate("/login");
                return;
            }

            console.log("Отправляемые данные на бэкенд:", JSON.stringify(testData, null, 2));

            const response = await testsAPI.createTest(testData);
            const result = response.data;

            console.log("Успешный ответ от сервера:", result);



            alert(
                isEditing
                    ? "Тест успешно обновлен!"
                    : "Тест успешно создан на сервере!"
            );
            navigate("/tests", { replace: true });

        } catch (error) {
            console.error("Ошибка при создании теста:", error);
            alert(
                `Не удалось ${
                    isEditing ? "обновить" : "создать"
                } тест на сервере: ${error.message}\n\nПроверьте консоль для деталей.`
            );
        }
    };

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

    return (
        <div className="tests-page">
            <div
                className="test-page"
                style={{ position: "absolute", left: "1430px", top: "0px" }}
            >
                <LogoutButton />
            </div>
            <div className="create-wrapper">
                <div className="create-left">
                    <input
                        className="test-title"
                        placeholder="Название"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                    />
                    <input
                        className="test-desc"
                        placeholder="Описание теста"
                        value={description}
                        onChange={(e) =>
                            setDescription(e.target.value)
                        }
                    />
                    <div className="tests-line"></div>

                    <PassingCriteria
                        criteria={passingCriteria}
                        updateCriteria={setPassingCriteria}
                        totalPoints={calculateTotalPoints()}
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

                    <ResultMessages
                        messages={resultMessages}
                        updateMessages={setResultMessages}
                    />
                </div>

                <div className="create-right">
                    <div className="create-right-inner">
                        <button
                            className="save-btn"
                            onClick={handleSave}
                        >
                            {isEditing
                                ? "Сохранить изменения"
                                : "Создать тест"}
                        </button>


                        <div className="right-section">
                            <p>Добавить новый вопрос</p>
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
                        <div className="time-box">
                            <p>Установить время</p>
                            <div className="time-inputs-box">
                                <input
                                    type="number"
                                    min="0"
                                    placeholder="0 ч"
                                    value={time.hours}
                                    onChange={(e) =>
                                        setTime({
                                            ...time,
                                            hours:
                                                parseInt(
                                                    e.target.value
                                                ) || 0,
                                        })
                                    }

                                    />
                                <input
                                    type="number"
                                    min="0"
                                    max="59"
                                    placeholder="0 м"
                                    value={time.minutes}
                                    onChange={(e) =>
                                        setTime({
                                            ...time,
                                            minutes: Math.min(
                                                59,
                                                parseInt(
                                                    e.target.value
                                                ) || 0
                                            ),
                                        })
                                    }
                                />
                                <input
                                    type="number"
                                    min="0"
                                    max="59"
                                    placeholder="0 с"
                                    value={time.seconds}
                                    onChange={(e) =>
                                        setTime({
                                            ...time,
                                            seconds: Math.min(
                                                59,
                                                parseInt(
                                                    e.target.value
                                                ) || 0
                                            ),
                                        })
                                    }
                                />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
