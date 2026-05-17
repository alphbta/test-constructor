import React from 'react';
import timeIcon from "../../assets/time.svg";

export default function TimeBox({ time, setTime }) {
    return (
        <div className="time-box">
            <div className="time-box1">
                <img src={timeIcon} alt="время" />
                <p>Ограничение по времени</p>
            </div>
            <div className="time-box-inner1">
                <div className="time-input-box">
                    <input
                        type="number"
                        min="0"
                        placeholder="0"
                        value={time.hours}
                        onChange={e => setTime({ ...time, hours: parseInt(e.target.value) || 0 })}
                    />
                    <input
                        type="number"
                        min="0"
                        max="59"
                        placeholder="0 минут"
                        value={time.minutes}
                        onChange={e => setTime({ ...time, minutes: Math.min(59, parseInt(e.target.value) || 0) })}
                    />
                    <input
                        type="number"
                        min="0"
                        max="59"
                        placeholder="0"
                        value={time.seconds}
                        onChange={e => setTime({ ...time, seconds: Math.min(59, parseInt(e.target.value) || 0) })}
                    />
                </div>
            </div>
        </div>
    );
}
