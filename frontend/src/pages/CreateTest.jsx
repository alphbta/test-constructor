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

    const [passingCriteria, setPassingCriteria] = useState(
        isEditing ? editingTest.passingCriteria || {
            type: "percentage",
            percentage: 75,
            points: 0,
        } : {
            type: "percentage",
            percentage: 75,
            points: 0,
        }
    );

    const [resultMessages, setResultMessages] = useState(
        isEditing ? editingTest.resultMessages || {
            success: "",
            failure: ""
        } : {
            success: "",
            failure: ""
        }
    );

    const [questions, setQuestions] = useState(
        isEditing ? editingTest.questions.map((q, idx) => ({
            id: `q-${idx}-${Date.now()}`,
            order: idx + 1,
            type: q.type,
            text: q.text,
            ...(q.type === "shortText" && {
                correctAnswers: q.correctAnswers || [""],
                caseSensitive: q.caseSensitive || false
            }),
            ...(q.type === "singleChoice" && {
                options: q.options || [{ text: "", isCorrect: false }]
            }),
            ...(q.type === "multipleChoice" && {
                options: q.options || [{ text: "", isCorrect: false }],
                scoringType: q.scoringType || "allOrNothing"
            }),
            ...(q.type === "matching" && {
                rows: q.rows || [{ option: "", answer: "" }]
            }),
            ...(q.type === "ordering" && {
                items: q.items || [{ text: "" }]
            }),
            maxScore: q.maxScore || 15,
        })) : [
            {
                id: "1",
                order: 1,
                type: "shortText",
                text: "",
                correctAnswers: [""],
                caseSensitive: false,
                maxScore: 15,
            },
        ]
    );

    const calculateTotalPoints = () => {
        return questions.reduce((sum, q) => sum + (q.maxScore || 0), 0);
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
            passingCriteria: {
                ...passingCriteria,
                totalPoints: calculateTotalPoints()
            },
            resultMessages,
            questions: questions.map((q) => ({
                type: q.type,
                text: q.text,
                maxScore: q.maxScore,
                ...(q.type === "shortText" && {
                    correctAnswers: q.correctAnswers,
                    caseSensitive: q.caseSensitive
                }),
                ...(q.type === "singleChoice" && {
                    options: q.options
                }),
                ...(q.type === "multipleChoice" && {
                    options: q.options,
                    scoringType: q.scoringType
                }),
                ...(q.type === "matching" && {
                    rows: q.rows
                }),
                ...(q.type === "ordering" && {
                    items: q.items
                }),
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
        { key: "shortText", label: "Задания на ручной ввод" },
        { key: "singleChoice", label: "Одиночный выбор" },
        { key: "multipleChoice", label: "Множественный выбор" },
        { key: "matching", label: "На соотношение"},
        { key: "ordering", label: "На расположение в правильном порядке"},
    ];

    return (
        <div className="tests-page">
            <div className="test-page" style={{ position: 'absolute', left: '1430px', top: '0px' }}>
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
                    onChange={(e) => setDescription(e.target.value)}
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
                <button className="save-btn" onClick={handleSave}>
                    {isEditing ? "Сохранить" : "Сохранить"}
                </button>

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
                <div className="time-box">
                    <p>Установить время</p>
                    <div className="time-inputs-box">
                        <input type="number" min="0" placeholder="0 ч"/>
                        <input type="number" min="0" max="59" placeholder="0 м" />
                        <input type="number" min="0" max="59" placeholder="0 с" />
                    </div>
                </div>
            </div>
        </div>
        </div>
    );
}