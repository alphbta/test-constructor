export default function TimeBox({ time, setTime }) {
    return (
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
                            hours: parseInt(e.target.value) || 0,
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
                                parseInt(e.target.value) || 0
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
                                parseInt(e.target.value) || 0
                            ),
                        })
                    }
                />
            </div>
        </div>
    );
}
