import processMark from "../../assets/process-mark.svg";

const PersonalData = () => {
    return (
        <div className="process-personal-data">
            <img src={processMark} width="25" height="25"/>
            <p>Я согласен(а) на обработку персональных данных</p>
        </div>
    );
}

export default PersonalData