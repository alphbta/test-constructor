import ManualInputQuestion from "./questions/ManualInputQuestion";
import SingleChoiceQuestion from "./questions/SingleChoiceQuestion";
import MultipleChoiceQuestion from "./questions/MultipleChoiceQuestion";
import MatchingQuestion from "./questions/MatchingQuestion";
import OrderingQuestion from "./questions/OrderingQuestion";

function SortableQuestion({ question, updateQuestion, deleteQuestion, onAddQuestion }) {
    switch (question.type) {
        case "shortText":
            return <ManualInputQuestion
                question={question}
                updateQuestion={updateQuestion}
                deleteQuestion={deleteQuestion}
                onAddQuestion={onAddQuestion}
            />;
        case "singleChoice":
            return <SingleChoiceQuestion
                question={question}
                updateQuestion={updateQuestion}
                deleteQuestion={deleteQuestion}
                onAddQuestion={onAddQuestion}
            />;
        case "multipleChoice":
            return <MultipleChoiceQuestion
                question={question}
                updateQuestion={updateQuestion}
                deleteQuestion={deleteQuestion}
                onAddQuestion={onAddQuestion}
            />;
        case "matching":
            return <MatchingQuestion
                question={question}
                updateQuestion={updateQuestion}
                deleteQuestion={deleteQuestion}
                onAddQuestion={onAddQuestion}
            />;
        case "ordering":
            return <OrderingQuestion
                question={question}
                updateQuestion={updateQuestion}
                deleteQuestion={deleteQuestion}
                onAddQuestion={onAddQuestion}
            />;
        default:
            return <div>Неизвестный тип вопроса</div>;
    }
}

export default SortableQuestion;